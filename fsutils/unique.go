package fsutils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

func FindUniqueName(name string) (string, error) {
	if name == "" {
		return "", nil
	}
	ext := filepath.Ext(name)
	orig := strings.TrimSuffix(name, ext)

	pattern, err := regexp.Compile(`-\d+$`)
	match := pattern.FindStringSubmatch(orig)
	if len(match) == 1 {
		orig = strings.TrimSuffix(orig, match[0])
	}

	if err != nil {
		return "", err
	}
	number := 0
	for {
		found, err := afero.Exists(Fs, name)
		if err != nil {
			return name, err
		}
		if !found {
			return name, nil
		}
		number++
		name = fmt.Sprintf("%s-%d%s", orig, number, ext)
	}
}
