package fsutils

import (
	"context"
	"io/fs"

	"github.com/spf13/afero"
)

// ForEachFile
func ForEachFile(ctx context.Context, dir string, callback func(entry fs.FileInfo) error) error {
	entries, err := afero.ReadDir(Fs, dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = callback(entry)
		if err != nil {
			return err
		}
	}
	return nil
}
