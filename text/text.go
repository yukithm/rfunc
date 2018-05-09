package text

import "regexp"

// ConvertLineEnding converts line endings of s.
// If eol is empty, it just returns s.
func ConvertLineEnding(s, eol string) string {
	if eol == "" {
		return s
	}

	pat := regexp.MustCompile(`\r\n|\r|\n`)
	return pat.ReplaceAllLiteralString(s, eol)
}
