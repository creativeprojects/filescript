package fsutils

import "github.com/spf13/afero"

func Exists(filename string) (bool, error) {
	return afero.Exists(Fs, filename)
}
