package fsutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppleDoubleFilename(t *testing.T) {
	fixtures := []struct {
		filename string
		dotFile  string
	}{
		{"", ""},
		{"file", "._file"},
		{"file.ext", "._file.ext"},
		{"path/file", "path/._file"},
		{"/root/file", "/root/._file"},
		{"._file", "._file"},
		{"._file.ext", "._file.ext"},
		{"path/._file", "path/._file"},
		{"/root/._file", "/root/._file"},
		{"path/._file.ext", "path/._file.ext"},
		{"/root/._file.ext", "/root/._file.ext"},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.filename, func(t *testing.T) {
			dotFile := GetAppleDouble(fixture.filename)
			assert.Equal(t, fixture.dotFile, dotFile)
		})
	}
}
