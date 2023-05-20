package fsutils

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountFiles(t *testing.T) {
	Fs = afero.NewMemMapFs()
	require.NoError(t, Fs.MkdirAll("/root/path", DirectoryPermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/a", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/._a", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/._b", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/._a", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/b", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/._a", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/._root", []byte("a"), FilePermission))
	require.NoError(t, Fs.MkdirAll("/root/path/1", DirectoryPermission))
	require.NoError(t, Fs.MkdirAll("/root/path/1/2/3", DirectoryPermission))
	require.NoError(t, Fs.MkdirAll("/root/path/2", DirectoryPermission))
	require.NoError(t, Fs.MkdirAll("/root/path/3", DirectoryPermission))

	expected := map[string]int{
		"/":                3,
		"/root":            3,
		"/root/path":       6,
		"/root/path/1":     1,
		"/root/path/1/2":   1,
		"/root/path/1/2/3": 0,
		"/root/path/2":     0,
		"/root/path/3":     0,
	}

	count, err := CountFiles(context.Background(), "/", func(event Event) bool {
		assert.NoError(t, event.Err)
		t.Log(event.SrcFilename)
		return true
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, count)
}
