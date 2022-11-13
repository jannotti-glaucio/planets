package validator

import (
	"regexp"
	"strings"
)

// Matches checks if string matches the pattern (pattern is regular expression)
// In case of error return false
func Matches(str, pattern string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}

// ReplacePattern replaces regular expression pattern in string
func ReplacePattern(str, pattern, replace string) string {
	r, _ := regexp.Compile(pattern)
	return r.ReplaceAllString(str, replace)
}

// Reverse returns reversed string
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Contains checks if the string contains the substring.
func Contains(str, substring string) bool {
	return strings.Contains(str, substring)
}
