package god

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGod(t *testing.T) {
	d := NewGod("ls", []string{"-l"})
	assert.Equal(t, "ls", d.name)
	assert.Equal(t, []string{"-l"}, d.args)
	d.Start()
	pid := d.cmd.Process.Pid
	started := d.started
	assert.NotEmpty(t, pid)
	assert.True(t, d.cmd.ProcessState.Exited())
	d.Restart()
	// new pid
	assert.NotEqual(t, pid, d.cmd.Process.Pid)
	assert.True(t, d.cmd.ProcessState.Exited())
	// restored after restart
	assert.False(t, d.stopping)
	// updated started
	assert.NotEqual(t, started, d.started)
}

func TestGodPanic(t *testing.T) {
	d := NewGod("./someprogram", []string{})
	assert.Equal(t, "./someprogram", d.name)
	assert.Panics(t, func() {
		d.Start()
	})
}

func TestGodRestartWhileRunning(t *testing.T) {
	d := NewGod("sleep", []string{"10"})
	MIMIMUM_AGE = 0.1
	go d.Start()
	time.Sleep(200 * time.Millisecond)
	pid := d.cmd.Process.Pid
	go d.Restart()
	time.Sleep(100 * time.Millisecond)
	assert.NotEqual(t, pid, d.cmd.Process.Pid)
	d.stopping = true
	d.Stop()
}


func TestGodRestartKill(t *testing.T) {
	d := NewGod("sleep", []string{"10"})
	MIMIMUM_AGE = 0.1
	go d.Start()
	time.Sleep(200 * time.Millisecond)
	pid := d.cmd.Process.Pid
	d.cmd.Process.Kill()
	time.Sleep(100 * time.Millisecond)
	assert.NotEqual(t, pid, d.cmd.Process.Pid)
	d.stopping = true
	d.Stop()
}
