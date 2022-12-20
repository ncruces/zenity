package zencmd

import "strings"

// StripMnemonic is internal.
func StripMnemonic(s string) string {
	// Strips mnemonics described in:
	// https: //docs.gtk.org/gtk4/class.Label.html#mnemonics

	var res strings.Builder

	underscore := false
	for _, b := range []byte(s) {
		switch {
		case underscore:
			underscore = false
		case b == '_':
			underscore = true
			continue
		}
		res.WriteByte(b)
	}

	return res.String()
}
