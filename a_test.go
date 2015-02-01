package god

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
