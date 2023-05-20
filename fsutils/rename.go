package fsutils

import (
	"context"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

// Rename/move a file from oldpath to newpath. If a file exists at newpath, a dash followed by a number will be added to the filename.
// If the macOS extended attributes file is found, it is also moved ("._*")
func Rename(oldpath, newpath string) (string, error) {
	newpath, _ = FindUniqueName(newpath)

	double := GetAppleDouble(oldpath)
	if found, err := afero.Exists(Fs, double); found && err == nil {
		// also move the macOS hidden file
		newDouble := GetAppleDouble(newpath)
		_ = Fs.Rename(double, newDouble)
	}
	return newpath, Fs.Rename(oldpath, newpath)
}

func MoveAllPerYear(ctx context.Context, dir string, progress func(event Event) bool, dryRun bool) error {
	moveFunc := func(entry fs.FileInfo) error {
		progress(Event{
			Type:        EventProgressFile,
			SrcFilename: entry.Name(),
		})
		if entry.IsDir() {
			return nil
		}
		if strings.HasPrefix(entry.Name(), ".") {
			return nil
		}
		year := strconv.Itoa(entry.ModTime().Year())
		if len(year) != 4 {
			return nil
		}
		orig := filepath.Join(dir, entry.Name())
		moveTo := filepath.Join(dir, year, entry.Name())
		newpath := moveTo
		if dryRun {
			pterm.Info.Printf("would move %q to %q\n", orig, newpath)
		} else {
			err := Fs.MkdirAll(filepath.Join(dir, year), DirectoryPermission)
			if err != nil {
				return err
			}
			newpath, err := Rename(orig, moveTo)
			if err != nil {
				progress(Event{
					Type:        EventError,
					Err:         err,
					SrcFilename: orig,
					DstFilename: moveTo,
				})
			}
			pterm.Debug.Printf("moving %q to %q\n", orig, newpath)
		}
		progress(Event{
			Type:        EventProgressFileProcessed,
			SrcFilename: orig,
			DstFilename: newpath,
		})
		return nil
	}
	return ForEachFile(ctx, dir, moveFunc)
}
