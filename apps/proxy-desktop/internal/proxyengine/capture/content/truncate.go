package content

import "unicode/utf8"

// TruncateUTF8 limits s to at most maxBytes UTF-8 without splitting a code point.
func TruncateUTF8(s string, maxBytes int64) (string, bool) {
	if maxBytes <= 0 {
		return "", len(s) > 0
	}
	if int64(len(s)) <= maxBytes {
		return s, false
	}
	end := int(maxBytes)
	for end > 0 && !utf8.RuneStart(s[end]) {
		end--
	}
	if end == 0 {
		return "", true
	}
	return s[:end], true
}
