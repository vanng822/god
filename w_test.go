package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanFileNotExist(t *testing.T) {
	folders, err := scan("./something")
	assert.Regexp(t, "./something: no such file or directory", err.Error())
	assert.Nil(t, folders)
}

func TestScanFile(t *testing.T) {
	folders, err := scan("./README.md")
	assert.Nil(t, err)
	assert.Equal(t, len(folders), 0)
}

func TestScanFolder(t *testing.T) {
	folders, err := scan("./test_program")
	assert.Nil(t, err)
	assert.Equal(t, len(folders), 1)
	assert.Equal(t, "./test_program/test_folder", folders[0])
}
