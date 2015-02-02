package god

import (
	"github.com/vanng822/gopid"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
	args := Args{}

	if err := args.Parse(os.Args[1:]); err != nil || args.help || len(args.programs) == 0 {
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

	// need to handle panic and shut down others
	log.Printf("Number of goroutines %d before start", runtime.NumGoroutine())
	for _, d := range z.gods {
		go d.Start()
	}

	log.Printf("Number of goroutines %d after start", runtime.NumGoroutine())
	z.sigc = make(chan os.Signal, 1)
	signal.Notify(z.sigc, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	for {
		sig := <-z.sigc
		log.Printf("Got signal %v", sig)
		switch sig {
		case syscall.SIGHUP:
			z.Restart()
			break
		case syscall.SIGTERM:
			z.Stop()
			goto programExit
		case os.Kill:
			z.Stop()
			goto programExit
		case os.Interrupt:
			z.Stop()
			goto programExit
		default:
			log.Printf("Unhandled signal %v, stop program", sig)
			z.Stop()
			goto programExit

		}

	}
programExit:
	log.Println("Program exit")
}

func (z *Goz) Stop() {
	for _, d := range z.gods {
		d.stopping = true
		d.Stop()
	}
}

func (z *Goz) Restart() {
	log.Printf("Number of goroutines %d before restart", runtime.NumGoroutine())
	for _, d := range z.gods {
		go d.Restart()
	}
	log.Printf("Number of goroutines %d after restart", runtime.NumGoroutine())
}
