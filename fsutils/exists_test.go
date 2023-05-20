package fsutils

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func assertDirExists(t *testing.T, filename string, exists bool) {
	found, err := afero.DirExists(Fs, filename)
	assert.NoError(t, err)
	assert.Equalf(t, exists, found, "dir %q", filename)
}

func assertFileExists(t *testing.T, filename string, exists bool) {
	found, err := afero.Exists(Fs, filename)
	assert.NoError(t, err)
	assert.Equalf(t, exists, found, "file %q", filename)
}
