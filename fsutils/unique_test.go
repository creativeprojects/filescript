package fsutils

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniqueFilename(t *testing.T) {
	Fs = afero.NewMemMapFs()
	require.NoError(t, Fs.MkdirAll("/root/path", DirectoryPermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/a", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/b", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/b-1", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/c-1.txt", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/d.txt", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/d-1.txt", []byte("a"), FilePermission))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/d-2.txt", []byte("a"), FilePermission))

	fixtures := []struct {
		filename string
		unique   string
	}{
		{"", ""},
		{"a", "a"},
		{"/root", "/root-1"},
		{"/root/path", "/root/path-1"},
		{"/root/path/a", "/root/path/a-1"},
		{"/root/path/b", "/root/path/b-2"},
		{"/root/path/c.txt", "/root/path/c.txt"},
		{"/root/path/d.txt", "/root/path/d-3.txt"},
		{"/root/path/b-1", "/root/path/b-2"},
		{"/root/path/c-1.txt", "/root/path/c-2.txt"},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.filename, func(t *testing.T) {
			unique, err := FindUniqueName(fixture.filename)
			require.NoError(t, err)
			assert.Equal(t, fixture.unique, unique)
		})
	}
}
