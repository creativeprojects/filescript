package fsutils

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRename(t *testing.T) {
	fixtures := []struct {
		oldpath          string
		newpath          string
		expected         string
		valid            bool
		withExtendedFile bool
	}{
		{"/root/path/a", "/root/path/b", "/root/path/b", true, true},
		{"/root/path/a", "/root/b", "/root/b-1", true, true},
		{"/root/b", "/root/c", "/root/c", true, false},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.expected, func(t *testing.T) {
			Fs = afero.NewMemMapFs()
			require.NoError(t, Fs.MkdirAll("/root/path", 0777))
			require.NoError(t, afero.WriteFile(Fs, "/root/path/a", []byte("a"), 0777))
			require.NoError(t, afero.WriteFile(Fs, "/root/path/._a", []byte("a"), 0777))
			require.NoError(t, afero.WriteFile(Fs, "/root/path/._b", []byte("a"), 0777))
			require.NoError(t, afero.WriteFile(Fs, "/root/._a", []byte("a"), 0777))
			require.NoError(t, afero.WriteFile(Fs, "/root/b", []byte("a"), 0777))
			require.NoError(t, afero.WriteFile(Fs, "/._a", []byte("a"), 0777))
			require.NoError(t, afero.WriteFile(Fs, "/._root", []byte("a"), 0777))

			result, err := Rename(fixture.oldpath, fixture.newpath)
			if fixture.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, fixture.expected, result)
			found, err := afero.Exists(Fs, GetAppleDouble(fixture.expected))
			assert.NoError(t, err)
			assert.Equal(t, fixture.withExtendedFile, found, "macOS extended file")
		})
	}
}

func TestMoveAllPerYear(t *testing.T) {
	Fs = afero.NewMemMapFs()
	require.NoError(t, Fs.MkdirAll("/root/path", 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/a", []byte("a"), 0777))
	require.NoError(t, Fs.Chtimes("/root/path/a", time.Now(), time.Date(2001, 1, 1, 1, 1, 0, 0, time.Local)))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/._a", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/b", []byte("a"), 0777))
	require.NoError(t, Fs.Chtimes("/root/path/b", time.Now(), time.Date(2002, 1, 1, 1, 1, 0, 0, time.Local)))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/._b", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/c", []byte("a"), 0777))
	require.NoError(t, Fs.Chtimes("/root/path/c", time.Now(), time.Date(2001, 1, 1, 1, 1, 0, 0, time.Local)))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/d", []byte("a"), 0777))
	require.NoError(t, Fs.Chtimes("/root/path/d", time.Now(), time.Date(2010, 1, 1, 1, 1, 0, 0, time.Local)))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/2010/d", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/._a", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/b", []byte("a"), 0777))

	progress := func(event Event) bool {
		assert.NoError(t, event.Err)
		if event.Type == EventProgressFileProcessed {
			t.Logf("moved %q to %q\n", event.SrcFilename, event.DstFilename)
		}
		return true
	}
	err := MoveAllPerYear(context.Background(), "/root/path", progress)
	assert.NoError(t, err)

	assertFileExists(t, "/root/path/2001/a", true)
	assertFileExists(t, "/root/path/2001/._a", true)
	assertFileExists(t, "/root/path/2002/b", true)
	assertFileExists(t, "/root/path/2002/._b", true)
	assertFileExists(t, "/root/path/2001/c", true)
	assertFileExists(t, "/root/path/2001/._c", false)
	assertFileExists(t, "/root/path/2010/d-1", true)
	assertFileExists(t, "/root/path/2010/._d-1", false)
}

func assertFileExists(t *testing.T, filename string, exists bool) {
	found, err := afero.Exists(Fs, filename)
	assert.NoError(t, err)
	assert.Equal(t, exists, found)
}
