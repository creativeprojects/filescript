package fsutils

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

func Unzip(ctx context.Context, filename string, exclude []string, progress func(event Event) bool) error {
	var err error

	extractTo := filepath.Join(filepath.Dir(filename), filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))]))
	extractTo, err = FindUniqueName(extractTo)
	if err != nil {
		return err
	}
	err = Fs.Mkdir(extractTo, 0755)
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
			dirname := filepath.Join(extractTo, f.Name)
			err = Fs.MkdirAll(dirname, 0755)
			if err != nil {
				return err
			}
			// a directory modification time will be set when the files are extracted
			// by deferring this call, we set the modification time once all files have been extracted
			defer func() {
				err = Fs.Chtimes(dirname, time.Now(), f.Modified)
				if err != nil {
					pterm.Warning.Println(err)
				}
			}()
			continue
		}
		if isExcluded(f.Name, exclude) {
			continue
		}
		progress(Event{
			Type:   EventProgressFile,
			SrcDir: f.Name,
		})
		outputFilename, err := unzipFile(f, extractTo)
		if err != nil {
			return err
		}
		err = Fs.Chtimes(outputFilename, time.Now(), f.Modified)
		if err != nil {
			pterm.Warning.Println(err)
		}
		progress(Event{
			Type:   EventProgressFileProcessed,
			SrcDir: f.Name,
		})
	}
	return nil
}

func unzipFile(f *zip.File, extractTo string) (string, error) {
	input, err := f.Open()
	if err != nil {
		return "", err
	}
	defer input.Close()

	outputFilename, err := FindUniqueName(filepath.Join(extractTo, f.Name))
	if err != nil {
		return "", err
	}
	dirpath := path.Dir(outputFilename)
	if dirpath != "" {
		err = Fs.MkdirAll(dirpath, DirectoryPermission)
		if err != nil {
			return "", err
		}
	}
	output, err := Fs.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0655)
	if err != nil {
		return "", err
	}
	defer output.Close()

	written, err := io.Copy(output, input)
	if err != nil {
		return "", err
	}
	if written != int64(f.UncompressedSize64) {
		return "", fmt.Errorf("written %d bytes, expected %d", written, f.UncompressedSize64)
	}
	return outputFilename, nil
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
