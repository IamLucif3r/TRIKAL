package util

import "unicode"

func StripUnsafe(s string) string {
	// remove non-printable; keep ASCII (add more rules as needed)
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r <= unicode.MaxASCII && unicode.IsPrint(r) {
			out = append(out, r)
		}
	}
	return string(out)
}
