package fsutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlural(t *testing.T) {
	fixtures := []struct {
		count  int
		word   string
		result string
	}{
		{0, "", ""},
		{10, "", ""},
		{0, "file", "no file"},
		{1, "file", "1 file"},
		{2, "file", "2 files"},
		{0, "directory", "no directory"},
		{1, "directory", "1 directory"},
		{2, "directory", "2 directories"},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.result, func(t *testing.T) {
			result := Plural(fixture.count, fixture.word)
			assert.Equal(t, fixture.result, result)
		})
	}
}
