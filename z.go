package god

import (
	"github.com/vanng822/gopid"
	"log"
	"os"
	"os/signal"
	"syscall"
	"runtime"
)

type Goz struct {
	gods []*God
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
	args.Parse()
	
	if args.pidFile != "" {
		gopid.CheckPid(args.pidFile, args.force)
		gopid.CreatePid(args.pidFile)
		defer gopid.CleanPid(args.pidFile)
	}

	log.Println(args.args)

	if len(args.programs) > 0 {
		for _, p := range args.programs {
			z.Add(NewGod(p, args.args))
		}
		// need to handle panic and shut down others
		log.Printf("Number of goroutines %d before start", runtime.NumGoroutine())
		for _, d := range z.gods {
			go d.Start()
		}
		log.Printf("Number of goroutines %d after start", runtime.NumGoroutine())
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		for {
			sig := <-sigc
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
}

func (z *Goz) Stop() {
	for _, d := range z.gods {
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
