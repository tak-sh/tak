package stringutils

import (
	"strings"
	"unicode"
)

// Capitalize the first letter of a string.
func Capitalize(str string) string {
	switch len(str) {
	case 0:
		return str
	case 1:
		return string(unicode.ToUpper(rune(str[0])))
	default:
		return strings.Join([]string{string(unicode.ToUpper(rune(str[0]))), str[1:]}, "")
	}
}
