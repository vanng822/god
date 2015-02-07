package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArgsEmpty(t *testing.T) {
	a := Args{}
	assert.Nil(t, a.Parse([]string{}))
	assert.Empty(t, a.pidFile)
	assert.False(t, a.force)
	assert.Empty(t, a.args)
	assert.Empty(t, a.programs)
}

func TestArgsHelp(t *testing.T) {
	a := Args{}
	assert.Nil(t, a.Parse([]string{"--help"}))
	assert.True(t, a.help)
}

func TestArgsOneProgram(t *testing.T) {
	a := Args{}
	
	args := []string{"--pidfile", "test.pid", "--pidclean",
		"-s", "./example/test_bin", "-p", "8080"}

	assert.Nil(t, a.Parse(args))
	assert.Equal(t, "test.pid", a.pidFile)
	assert.True(t, a.force)
	assert.Equal(t, []string{"-p", "8080"}, a.args)
	assert.Equal(t, []string{"./example/test_bin"}, a.programs)
}

func TestArgsMultipleProgram(t *testing.T) {
	a := Args{}
	args := []string{"--pidfile", "test.pid", "--pidclean",
		"-s", "./example/test_bin",
		"-s", "./example/test_bin2", "-p", "8080"}

	assert.Nil(t, a.Parse(args))
	assert.Equal(t, "test.pid", a.pidFile)
	assert.True(t, a.force)
	assert.Equal(t, []string{"-p", "8080"}, a.args)
	assert.Equal(t, []string{"./example/test_bin", "./example/test_bin2"}, a.programs)
}

func TestArgsMultipleProgramArgs(t *testing.T) {
	a := Args{}
	args := []string{"--pidfile", "test.pid", "--pidclean",
		"-s", "sleep", "10",
		"-s", "node", "server.js",
		"-p", "8080"}

	assert.Nil(t, a.Parse(args))
	assert.Equal(t, "test.pid", a.pidFile)
	assert.True(t, a.force)
	assert.Equal(t, []string{"-p", "8080"}, a.args)
	assert.Equal(t, []string{"sleep", "node"}, a.programs)
	assert.Equal(t, [][]string{{"10"}, {"server.js"}}, a.programArgs)
}

func TestArgsInvalidProgram(t *testing.T) {
	a := Args{}
	
	args := []string{"--pidfile", "test.pid", "--pidclean", "-s", "-p", "8080"}
	err := a.Parse(args)
	assert.NotNil(t, err)
	assert.Regexp(t, "Invalid program", err)
}

func TestArgsInvalidNoProgram(t *testing.T) {
	a := Args{}
	
	args := []string{"--pidfile", "test.pid", "--pidclean", "-s"}
	err := a.Parse(args)
	assert.NotNil(t, err)
	assert.Regexp(t, "Invalid program", err)
}

func TestArgsInvalidPidfile(t *testing.T) {
	a := Args{}
	
	args := []string{"--pidfile", "--pidclean"}
	err := a.Parse(args)
	assert.NotNil(t, err)
	assert.Regexp(t, "Invalid pidfile value", err)
}

func TestArgsInvalidNoPidfile(t *testing.T) {
	a := Args{}
	
	args := []string{"--pidclean", "--pidfile"}
	err := a.Parse(args)
	assert.NotNil(t, err)
	assert.Regexp(t, "Invalid pidfile value", err)
}

func TestArgsInvalidIntervalNoValue(t *testing.T) {
	a := Args{}
	args := []string{"--interval", "--pidclean"}
	err := a.Parse(args)
	assert.NotNil(t, err)
	assert.Regexp(t, "Invalid interval value", err)
}

func TestArgsInvalidIntervalSmallValue(t *testing.T) {
	a := Args{}
	args := []string{"--interval", "1"}
	err := a.Parse(args)
	assert.NotNil(t, err)
	assert.Regexp(t, "Minium value for interval is 2 seconds", err)
}

func TestArgsInvalidIntervalNoneIntValue(t *testing.T) {
	a := Args{}
	args := []string{"--interval", "abc123"}
	err := a.Parse(args)
	assert.NotNil(t, err)
	assert.Regexp(t, "Invalid interval value", err)
}

func TestArgsIntervalOK(t *testing.T) {
	a := Args{}
	args := []string{"--interval", "123"}
	err := a.Parse(args)
	assert.Nil(t, err)
	assert.Equal(t, a.interval, 123)
}