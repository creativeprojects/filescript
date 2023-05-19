package fsutils

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestUnzip(t *testing.T) {
	Fs = afero.NewMemMapFs()
	err := Unzip(context.Background(), "test.zip", func(event Event) bool {
		t.Logf("%+v", event)
		return true
	})
	assert.NoError(t, err)

	assertFileExists(t, "test/zip", true)
	assertFileExists(t, "test/zip/file1.txt", true)
	assertFileExists(t, "test/zip/file2.txt", true)

	assertFileExists(t, "test/zip/dir1", true)
	assertFileExists(t, "test/zip/dir1/sub2", true)
	assertFileExists(t, "test/zip/dir1/sub2/sub3", true)
	assertFileExists(t, "test/zip/dir2", true)
}
