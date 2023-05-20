package fsutils

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(ctx context.Context, filename string, exclude []string, progress func(event Event) bool) error {
	extractTo := filepath.Join(filepath.Dir(filename), filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))]))
	err := Fs.Mkdir(extractTo, 0755)
	if err != nil {
		return err
	}

	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	// Iterate through the files in the archive,
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/") {
			progress(Event{
				Type:   EventProgressDir,
				SrcDir: f.Name,
			})
			err = Fs.Mkdir(filepath.Join(extractTo, f.Name), 0755)
			if err != nil {
				return err
			}
			continue
		}
		progress(Event{
			Type:   EventProgressFile,
			SrcDir: f.Name,
		})
		input, err := f.Open()
		if err != nil {
			return err
		}
		defer input.Close()

		output, err := Fs.OpenFile(filepath.Join(extractTo, f.Name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0655)
		if err != nil {
			return err
		}
		defer output.Close()

		written, err := io.Copy(output, input)
		if err != nil {
			return err
		}
		if written != int64(f.UncompressedSize64) {
			return fmt.Errorf("written %d bytes, expected %d", written, f.UncompressedSize64)
		}
		progress(Event{
			Type:   EventProgressFileProcessed,
			SrcDir: f.Name,
		})
	}
	return nil
}

func isExcluded(filename string, excludes []string) bool {
	parts := strings.Split(filename, "/")
	for _, part := range parts {
		for _, exclude := range excludes {
			if part == exclude {
				return true
			}
		}
	}
	return false
}
