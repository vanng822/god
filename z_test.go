package god

import (
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestUsage(t *testing.T) {
	args := os.Args[:]
	defer func() {
		os.Args = args
	}()
	os.Args = []string{""}
	z := NewGoz()
	assert.NotNil(t, z)
	z.Start()
}

func TestStop(t *testing.T) {
	MIMIMUM_AGE = 0.1
	args := os.Args[:]
	defer func() {
		os.Args = args
	}()
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	defer func() {
		z.sigc <- syscall.SIGTERM
		time.Sleep(200 * time.Millisecond)
	}()
	go z.Start()
	time.Sleep(100 * time.Millisecond)
	pid := z.gods[0].cmd.Process.Pid
	assert.Nil(t, z.gods[0].cmd.ProcessState)
	// Prob
	z.Stop()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, z.gods[0].cmd.ProcessState)
}

func TestRestart(t *testing.T) {
	MIMIMUM_AGE = 0.1
	args := os.Args[:]
	defer func() {
		os.Args = args
	}()
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	defer func() {
		z.sigc <- syscall.SIGTERM
		time.Sleep(200 * time.Millisecond)
	}()
	go z.Start()
	time.Sleep(100 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.Nil(t, cmd.ProcessState)
	// Prob
	z.Restart()
	time.Sleep(100 * time.Millisecond)
	assert.NotEqual(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}

func TestSignalSigHup(t *testing.T) {
	MIMIMUM_AGE = 0.1
	args := os.Args[:]
	defer func() {
		os.Args = args
	}()
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	defer func() {
		z.sigc <- syscall.SIGTERM
		time.Sleep(200 * time.Millisecond)
	}()
	go z.Start()
	time.Sleep(200 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- syscall.SIGHUP
	time.Sleep(100 * time.Millisecond)
	assert.NotEqual(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}

func TestSignalKill(t *testing.T) {
	MIMIMUM_AGE = 0.1
	args := os.Args[:]
	defer func() {
		os.Args = args
	}()
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	go z.Start()
	time.Sleep(100 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- os.Kill
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}

func TestSignalInterrupt(t *testing.T) {
	MIMIMUM_AGE = 0.1
	args := os.Args[:]
	defer func() {
		os.Args = args
	}()
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	go z.Start()
	time.Sleep(100 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- os.Interrupt
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}

func TestSignalOther(t *testing.T) {
	MIMIMUM_AGE = 0.1
	args := os.Args[:]
	defer func() {
		os.Args = args
	}()
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	go z.Start()
	time.Sleep(100 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- syscall.SIGALRM
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}
