package main

import (
	"github.com/vanng822/gopid"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Goz struct {
	gods []*God
	sigc chan os.Signal
}

func NewGoz() *Goz {
	z := new(Goz)
	z.gods = make([]*God, 0)
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

	for i, p := range args.programs {
		pargs := args.programArgs[i]
		pargs = append(pargs, args.args...)
		z.Add(NewGod(p, pargs))
	}

	for _, d := range z.gods {
		d.Start()
	}

	if float64(args.interval) > MIMIMUM_AGE {
		ticker := time.NewTicker(time.Duration(args.interval) * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					z.Restart()
				}
			}
		}()
	}

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
			z.Restart()
			break
		case syscall.SIGTERM:
			z.Stop()
			return
		case os.Kill:
			z.Stop()
			return
		case os.Interrupt:
			z.Stop()
			return
		case syscall.SIGUSR1:
			z.Signal(sig)
		case syscall.SIGUSR2:
			z.Signal(sig)
		default:
			log.Printf("Unhandled signal %v, stop program", sig)
			z.Stop()
			return
		}
	}
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

func (z *Goz) Restart() {
	for _, d := range z.gods {
		d.Restart()
	}
}
