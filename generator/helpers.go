package generator

import (
	"strings"
)

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	if s == strings.ToUpper(s) {
		s = strings.ToLower(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func lowercase(s string) string {
	if s == "" {
		return ""
	}
	if s == strings.ToUpper(s) {
		s = strings.ToLower(s)
	}
	return strings.ToLower(s[:1]) + s[1:]
}
