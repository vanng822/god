package main

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUsage(t *testing.T) {
	os.Args = []string{""}
	z := NewGoz()
	assert.NotNil(t, z)
	z.Start()
}

func TestStop(t *testing.T) {
	MIMIMUM_AGE = 0.1
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
	z.Restart(false)
	time.Sleep(100 * time.Millisecond)
	assert.NotEqual(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}

func TestSignalSigHup(t *testing.T) {
	MIMIMUM_AGE = 0.1
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

func TestSignalSigUsr2(t *testing.T) {
	MIMIMUM_AGE = 0.1
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "5", "-s", "sleep", "5"}
	z := NewGoz()
	z.gracefulWait = 100 * time.Millisecond
	defer func() {
		z.sigc <- syscall.SIGTERM
		time.Sleep(200 * time.Millisecond)
	}()
	go z.Start()
	time.Sleep(200 * time.Millisecond)
	cmd := z.gods[0].cmd
	cmd2 := z.gods[1].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- syscall.SIGUSR2
	time.Sleep(600 * time.Millisecond)
	assert.NotEqual(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotEqual(t, cmd2.Process.Pid, z.gods[1].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
	assert.NotNil(t, cmd2.ProcessState)

	// Second time
	cmd = z.gods[0].cmd
	cmd2 = z.gods[1].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- syscall.SIGUSR2
	time.Sleep(600 * time.Millisecond)
	assert.NotEqual(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotEqual(t, cmd2.Process.Pid, z.gods[1].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
	assert.NotNil(t, cmd2.ProcessState)
}

func TestSignalKill(t *testing.T) {
	MIMIMUM_AGE = 0.1
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	go z.Start()
	time.Sleep(200 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- os.Kill
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}

func TestSignalInterrupt(t *testing.T) {
	MIMIMUM_AGE = 0.1
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	go z.Start()
	time.Sleep(200 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- os.Interrupt
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}

func TestSignalOther(t *testing.T) {
	MIMIMUM_AGE = 0.1
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "10"}
	z := NewGoz()
	go z.Start()
	time.Sleep(200 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.Nil(t, cmd.ProcessState)
	z.sigc <- syscall.SIGALRM
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	assert.NotNil(t, cmd.ProcessState)
}

func TestRecoverPanic(t *testing.T) {
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "100", "-s", "./something"}
	z := NewGoz()
	assert.NotPanics(t, func() {
		z.Start()
	})
}

func TestIntervalAutorestart(t *testing.T) {
	MIMIMUM_AGE = 0.1
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "--interval", "2", "-s", "sleep", "1"}
	z := NewGoz()
	defer func() {
		z.sigc <- syscall.SIGTERM
		time.Sleep(200 * time.Millisecond)
	}()
	go z.Start()
	time.Sleep(200 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.NotNil(t, cmd)
	time.Sleep(2200 * time.Millisecond)
	assert.NotEqual(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
}

func TestSignal(t *testing.T) {
	MIMIMUM_AGE = 0.1
	os.Args = []string{"", "--pidfile", "testing.pid", "--pidclean", "-s", "sleep", "5"}
	z := NewGoz()
	defer func() {
		z.sigc <- syscall.SIGTERM
		time.Sleep(200 * time.Millisecond)
	}()
	go z.Start()
	time.Sleep(200 * time.Millisecond)
	cmd := z.gods[0].cmd
	assert.NotNil(t, cmd)
	z.sigc <- syscall.SIGUSR1
	time.Sleep(200 * time.Millisecond)
	assert.NotEqual(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
	cmd = z.gods[0].cmd
	z.sigc <- syscall.SIGUSR2
	time.Sleep(200 * time.Millisecond)
	assert.NotEqual(t, cmd.Process.Pid, z.gods[0].cmd.Process.Pid)
}
