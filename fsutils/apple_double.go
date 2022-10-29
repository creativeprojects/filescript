package fsutils

import (
	"path/filepath"
	"strings"
)

const (
	AppleDoublePrefix = "._"
)

func GetAppleDouble(filename string) string {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)

	if base == "" || base == "." {
		return ""
	}

	if strings.HasPrefix(base, AppleDoublePrefix) {
		return filename
	}

	return filepath.Join(dir, AppleDoublePrefix+base)
}
