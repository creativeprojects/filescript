package fsutils

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

func CountFiles(ctx context.Context, root string, progress func(event Event) bool) (map[string]int, error) {
	pathCount := make(map[string]int, 1000)
	pathCount[root] = 0

	err := afero.Walk(Fs, root, func(path string, d fs.FileInfo, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if path == root {
			return nil
		}

		if d.IsDir() {
			pathCount[path] = 0
		}
		parent := filepath.Dir(path)
		if _, found := pathCount[parent]; !found {
			return nil
		}
		pathCount[parent]++

		progress(Event{
			Type:        EventProgressFile,
			SrcFilename: path,
		})
		return nil
	})

	return pathCount, err
}
