package fsutils

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

func Unzip(ctx context.Context, filename string, progress func(event Event) bool) error {
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

		output, err := Fs.Create(filepath.Join(extractTo, f.Name))
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
