package fsutils

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func FindOrphans(ctx context.Context, dir, prefix, suffix string, progress func(event Event) bool) ([]string, error) {
	orphans := newSliceRef(100)
	err := findOrphans(ctx, dir, prefix, suffix, orphans, progress)
	return orphans.ref, err
}

func findOrphans(ctx context.Context, dir, prefix, suffix string, orphans *sliceRef, progress func(event Event) bool) error {
	progress(Event{
		Type:   EventProgressDir,
		SrcDir: dir,
	})

	all, err := afero.ReadDir(Fs, dir)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		progress(Event{
			Type: EventError,
			Err:  err,
		})
		// we send the error event but we don't stop
		return nil
	}

	// count all elements first
	var totalDirs, totalFiles int
	for _, info := range all {
		if info.IsDir() {
			totalDirs++
			continue
		}
		totalFiles++
	}

	progress(Event{
		Type:            EventTotal,
		TotalFilesInDir: totalFiles,
		TotalDirsInDir:  totalDirs,
	})

	for _, info := range all {
		// check if the context is cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if info.IsDir() {
			err := findOrphans(ctx, filepath.Join(dir, info.Name()), prefix, suffix, orphans, progress)
			if err != nil {
				return err
			}
			continue
		}
		progress(Event{
			Type:        EventProgressFile,
			SrcDir:      dir,
			SrcFilename: info.Name(),
		})

		if prefix != "" && !strings.HasPrefix(info.Name(), prefix) {
			continue
		}
		if suffix != "" && !strings.HasPrefix(info.Name(), suffix) {
			continue
		}
		source := strings.TrimPrefix(info.Name(), prefix)
		source = strings.TrimSuffix(source, suffix)
		if has(source, all) {
			continue
		}
		orphans.append(filepath.Join(dir, info.Name()))
		progress(Event{
			Type:        EventProgressFileProcessed,
			SrcDir:      dir,
			SrcFilename: info.Name(),
		})
	}
	return nil
}

func has(name string, in []fs.FileInfo) bool {
	for _, info := range in {
		if info.Name() == name {
			return true
		}
	}
	return false
}
