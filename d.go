package god

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type God struct {
	cmd     *exec.Cmd
	name    string
	args    []string
	started time.Time
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
	go d.Watch()
}

func (d *God) Pid() int {
	return d.cmd.Process.Pid
}

func (d *God) Wait() error {
	return d.cmd.Wait()
}

func (d *God) Watch() {
	restart := false
	for {
		log.Printf("Waiting for command to finish...")
		err := d.Wait()
		if err == nil {
			log.Println("Terminate without error")
			break
		}
		log.Printf("Command finished with error: %v", err)
		if time.Now().Sub(d.started).Seconds() < 2 {
			log.Printf("Program '%s' restart too fast. No restart!", d.name)
		} else {
			restart = true
		}
		break
	}
	if restart {
		d.Restart()
	}
}

func (d *God) Restart() {
	log.Printf("Restart program %s", d.name)
	pid := d.cmd.Process.Pid
	d.cmd.Process.Signal(syscall.SIGTERM)
	d.stopWait(pid)
	d.Start()
}

func (d *God) stopWait(pid int) {
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

func (d *God) Stop() {
	pid := d.cmd.Process.Pid
	d.cmd.Process.Signal(syscall.SIGTERM)
	d.stopWait(pid)
}
