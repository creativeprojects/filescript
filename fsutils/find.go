package fsutils

import (
	"context"
	"io/fs"
	"strings"

	"github.com/spf13/afero"
)

type FileMatcher func(filename string, fileinfo fs.FileInfo) bool

func WithExtension(extension string) FileMatcher {
	return func(filename string, fileinfo fs.FileInfo) bool {
		return strings.HasSuffix(filename, extension)
	}
}

func WithExecutionBit() FileMatcher {
	return func(filename string, fileinfo fs.FileInfo) bool {
		return fileinfo.Mode().Perm()&0111 != 0
	}
}

func FindFiles(ctx context.Context, matcher FileMatcher, root string, found chan string, progress func(event Event) bool) error {
	err := afero.Walk(Fs, root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			progress(Event{
				Type: EventError,
				Err:  err,
			})
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if path == root {
			return nil
		}

		if info.IsDir() {
			progress(Event{
				Type:   EventProgressDir,
				SrcDir: path,
			})
			return nil
		}

		if matcher(path, info) {
			progress(Event{
				Type:        EventProgressFile,
				SrcFilename: path,
			})
			found <- path
			return nil
		}
		return nil
	})
	return err
}
