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

func TestExclusions(t *testing.T) {
	excludes := []string{"file1.txt", "sub2"}
	fixtures := []struct {
		filename string
		excluded bool
	}{
		{"file1.txt", true},
		{"dir1/file1.txt", true},
		{"dir1/dir2/file1.txt", true},
		{"sub2", true},
		{"sub2/file1.txt", true},
		{"sub3/sub2", true},
		{"sub3/sub2/file.txt", true},
		{"file2.txt", false},
		{"dir1/file2.txt", false},
		{"dir1/dir2/file2.txt", false},
		{"sub1", false},
		{"sub1/file2.txt", false},
		{"sub3/sub1", false},
		{"sub3/sub1/file.txt", false},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.filename, func(t *testing.T) {
			assert.Equal(t, fixture.excluded, isExcluded(fixture.filename, excludes))
		})
	}
}
