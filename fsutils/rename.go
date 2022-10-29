package fsutils

import (
	"errors"
	"os"
	"path/filepath"
)

// Rename/move a file from oldpath to newpath. If a file exists at newpath, a dash followed by a number will be added to the filename.
// If the macOS double file is found, it is also copied ("._*")
func Rename(oldpath, newpath string) (string, error) {
	newpath, _ = FindUniqueName(newpath)

	oldDir := filepath.Dir(oldpath)
	oldFilename := filepath.Base(oldpath)

	double := filepath.Join(oldDir, "._"+oldFilename)
	if fileExists(double) {
		// also move the macOS hidden file
		newDouble := filepath.Join(filepath.Dir(newpath), "._"+filepath.Base(newpath))
		_ = os.Rename(double, newDouble)
	}
	return newpath, os.Rename(oldpath, newpath)
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err == nil || errors.Is(err, os.ErrExist) {
		return true
	}
	return false
}
