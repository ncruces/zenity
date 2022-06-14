package zencmd

import "strings"

// Unescape is internal.
func Unescape(s string) string {
	// Apply rules described in:
	// https://docs.gtk.org/glib/func.strescape.html

	const (
		initial = iota
		escape1
		escape2
		escape3
	)
	var oct byte
	var res strings.Builder
	state := initial
	for _, b := range []byte(s) {
		switch state {
		default:
			switch b {
			case '\\':
				state = escape1
			default:
				res.WriteByte(b)
				state = initial
			}

		case escape1:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7':
				oct = b - '0'
				state = escape2
			case 'b':
				res.WriteByte('\b')
				state = initial
			case 'f':
				res.WriteByte('\f')
				state = initial
			case 'n':
				res.WriteByte('\n')
				state = initial
			case 'r':
				res.WriteByte('\r')
				state = initial
			case 't':
				res.WriteByte('\t')
				state = initial
			case 'v':
				res.WriteByte('\v')
				state = initial
			default:
				res.WriteByte(b)
				state = initial
			}

		case escape2:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7':
				oct = oct<<3 | (b - '0')
				state = escape3
			case '\\':
				res.WriteByte(oct)
				state = escape1
			default:
				res.WriteByte(oct)
				res.WriteByte(b)
				state = initial
			}

		case escape3:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7':
				oct = oct<<3 | (b - '0')
				res.WriteByte(oct)
				state = initial
			case '\\':
				res.WriteByte(oct)
				state = escape1
			default:
				res.WriteByte(oct)
				res.WriteByte(b)
				state = initial
			}
		}
	}
	if state == escape2 || state == escape3 {
		res.WriteByte(oct)
	}

	return res.String()
}
