package types

import (
	"regexp"
	"strings"
	"unicode"
)

func NormalizeSpaces(s string) string {
	var multispaceRegex = regexp.MustCompile(`\s+`)
	return multispaceRegex.ReplaceAllString(strings.TrimSpace(s), " ")
}

func HasNewlines(s string) bool {
	return strings.ContainsRune(s, '\n')
}

func HasUnprintableChar(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}
