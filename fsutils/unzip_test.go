package fsutils

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestUnzip(t *testing.T) {
	Fs = afero.NewMemMapFs()
	err := Unzip(context.Background(), "test.zip", []string{}, func(event Event) bool {
		t.Logf("%+v", event)
		return true
	})
	assert.NoError(t, err)

	assertDirExists(t, "test/zip", true)
	assertFileExists(t, "test/zip/file1.txt", true)
	assertFileExists(t, "test/zip/file2.txt", true)

	assertDirExists(t, "test/zip/dir1", true)
	assertDirExists(t, "test/zip/dir1/sub2", true)
	assertDirExists(t, "test/zip/dir1/sub2/sub3", true)
	assertDirExists(t, "test/zip/dir2", true)
}

func TestExclusion(t *testing.T) {

}
