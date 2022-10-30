package fsutils

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindOrphans(t *testing.T) {
	Fs = afero.NewMemMapFs()
	require.NoError(t, Fs.MkdirAll("/root/path", 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/a", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/._a", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/path/._b", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/._a", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/root/b", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/._a", []byte("a"), 0777))
	require.NoError(t, afero.WriteFile(Fs, "/._root", []byte("a"), 0777))

	orphans, err := FindOrphans(context.Background(), "/", "._", "", func(event Event) bool {
		assert.NoError(t, event.Err)
		return true
	})
	assert.NoError(t, err)
	expected := []string{"/._a", "/root/._a", "/root/path/._b"}
	assert.ElementsMatch(t, expected, orphans)
}
