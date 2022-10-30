package fsutils

import (
	"strconv"
	"strings"
)

func Plural(count int, word string) string {
	if word == "" {
		return ""
	}
	if count <= 0 {
		return "no " + word
	}
	if count == 1 {
		return "1 " + word
	}
	if strings.HasSuffix(word, "y") {
		return strconv.Itoa(count) + " " + strings.TrimSuffix(word, "y") + "ies"
	}
	return strconv.Itoa(count) + " " + word + "s"
}
