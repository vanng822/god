package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vanng822/gopid"
)

type Goz struct {
	gods         []*God
	sigc         chan os.Signal
	ticker       *time.Ticker
	args         Args
	gracefulWait time.Duration
}

func NewGoz() *Goz {
	z := new(Goz)
	z.gods = make([]*God, 0)
	z.gracefulWait = 3 * time.Second
	return z
}

func (z *Goz) Add(d *God) {
	z.gods = append(z.gods, d)
}

func (z *Goz) Start() {
	defer func() {
		if r := recover(); r != nil {
			// print error first
			log.Println(r)
			for _, d := range z.gods {
				if !d.started.IsZero() {
					log.Printf("Stopping '%s' before exit", d.name)
					d.Stop()
				}
			}
		}
	}()

	args := Args{}

	z.args = args

	if err := args.Parse(os.Args[1:]); err != nil {
		usage()
		return
	}

	if args.version {
		version()
		return
	}

	if args.help || len(args.programs) == 0 {
		usage()
		return
	}

	if args.pidFile != "" {
		gopid.CheckPid(args.pidFile, args.force)
		gopid.CreatePid(args.pidFile)
		defer gopid.CleanPid(args.pidFile)
	}

	if args.logFile != "" {
		logwritter, err := os.OpenFile(args.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		os.Stdout = logwritter
	}

	if args.logFileErr != "" {
		errlogwritter, err := os.OpenFile(args.logFileErr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		os.Stderr = errlogwritter
	}

	for i, p := range args.programs {
		pargs := args.programArgs[i]
		pargs = append(pargs, args.args...)
		z.Add(NewGod(p, pargs))
	}

	for _, d := range z.gods {
		d.Start(nil)
	}

	if args.fileWatched {
		go func() {
			restartTimer := time.NewTimer(24 * 365 * 10 * time.Hour)
			shouldRestart := make(chan bool)
			go watchDirs(args.fileWatchedDirs, args.fileWatchedExts, shouldRestart)
			for {
				select {
				case <-restartTimer.C:
					z.Restart(args.graceful)
				case restart := <-shouldRestart:
					if restart {
						restartTimer.Reset(time.Duration(1000 * 1000 * 1000 * args.delaySecs))
					}
				}
			}
		}()
	}

	z.startInterval(args.interval)

	z.sigc = make(chan os.Signal, 1)
	// stop
	signals := []os.Signal{os.Kill, os.Interrupt, syscall.SIGTERM}
	// restart
	signals = append(signals, syscall.SIGHUP)
	// generic signal to send to processes
	signals = append(signals, syscall.SIGUSR1, syscall.SIGUSR2)
	signal.Notify(z.sigc, signals...)
	for {
		sig := <-z.sigc
		log.Printf("Got signal %v", sig)
		switch sig {
		case syscall.SIGHUP:
			z.Restart(args.graceful)
			z.startInterval(args.interval)
		case syscall.SIGUSR1:
			z.Signal(sig)
		case syscall.SIGUSR2:
			z.Restart(true)
			z.startInterval(args.interval)
		default:
			log.Println("Stop program")
			z.Stop()
			return
		}
	}
}

func (z *Goz) stopInterval() {
	if z.ticker != nil {
		z.ticker.Stop()
	}
}

func (z *Goz) startInterval(secs int) {
	if float64(secs) <= MIMIMUM_AGE {
		return
	}
	z.stopInterval()
	z.ticker = time.NewTicker(time.Duration(secs) * time.Second)
	go func() {
		for {
			select {
			case <-z.ticker.C:
				z.Restart(z.args.graceful)
			}
		}
	}()
}

func (z *Goz) Signal(s os.Signal) {
	for _, d := range z.gods {
		d.Signal(s)
	}
}

func (z *Goz) Stop() {
	for _, d := range z.gods {
		d.Stop()
	}
}
func waitDone(done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
func (z *Goz) Restart(graceful bool) {

	for _, d := range z.gods {
		if graceful {
			done := make(chan bool)
			d.Restart(done)
			waitDone(done)
			close(done)
			time.Sleep(z.gracefulWait)
		} else {
			d.Restart(nil)
		}
	}
}
