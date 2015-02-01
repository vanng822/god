package god

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

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
	if time.Now().Sub(d.started).Seconds() < 2 {
		log.Printf("Program '%s' restart too fast. No restart!", d.name)
		return
	}
	d.Restart()
}

func (d *God) Restart() {
	log.Printf("Restart program %s", d.name)
	d.stopping = true
	d.Stop()
	d.stopping = false
	d.Start()
}

func (d *God) Stop() {
	pid := d.cmd.Process.Pid
	d.cmd.Process.Signal(syscall.SIGTERM)
	stopWait(pid)
}

func stopWait(pid int) {
	// wait for process to completely terminated
	if process, err := os.FindProcess(pid); err != nil {
		if _, err := process.Wait(); err != nil {
			for i := 0; i < 50; i++ {
				time.Sleep(5 * time.Millisecond)
				if _, err := os.FindProcess(pid); err != nil {
					break
				}
			}
		}
	}
}
