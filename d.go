package god

import (
	"log"
	"os/exec"
	"syscall"
	"time"
)

var MIMIMUM_AGE = 2.0

type God struct {
	cmd      *exec.Cmd
	name     string
	args     []string
	started  time.Time
	stopping bool
}

func NewGod(name string, args []string) *God {
	d := new(God)
	d.name = name
	d.args = args
	return d
}

func (d *God) Start() {
	cmd := exec.Command(d.name, d.args...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	d.cmd = cmd
	d.started = time.Now()
	d.Watch()
}

func (d *God) Watch() {
	if d.cmd == nil {
		panic("You must call Start first")
	}
	log.Printf("Waiting for command to finish...")
	err := d.cmd.Wait()
	if err == nil {
		log.Println("Terminate without error")
		return
	}

	if d.stopping {
		log.Printf("Stopping. Process %s exited with %v", d.name, err)
		return
	}

	log.Printf("Command finished with error: %v", err)
	if time.Now().Sub(d.started).Seconds() < MIMIMUM_AGE {
		log.Printf("Program '%s' restart too fast. No restart!", d.name)
		return
	}
	d.Restart()
}

func (d *God) Restart() {
	if d.cmd == nil {
		panic("You must call Start first")
	}
	log.Printf("Restart program %s", d.name)
	// stopping race between goroutines
	// should figure out how to solve with channel
	d.stopping = true
	d.Stop()
	d.stopping = false
	d.Start()
}

func (d *God) Stop() {
	if d.cmd == nil {
		panic("You must call Start first")
	}
	d.cmd.Process.Signal(syscall.SIGTERM)
	// this is still bad, let see if we can fix with channel one day
	d.waitExited()
}

func (d *God) waitExited() {
	for i := 0; i < 400; i++ {
		if d.cmd.ProcessState != nil {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
}