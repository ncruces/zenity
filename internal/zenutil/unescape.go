package zenutil

// Unescape is internal.
func Unescape(s string) string {
	// Apply rules described in:
	// https://developer.gnome.org/glib/stable/glib-String-Utility-Functions.html#g-strescape

	const (
		initial = iota
		escape1
		escape2
		escape3
	)
	var oct byte
	var res []byte
	state := initial
	for _, b := range []byte(s) {
		switch state {
		case initial:
			switch b {
			case '\\':
				state = escape1
			default:
				res = append(res, b)
				state = initial
			}

		case escape1:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7':
				oct = b - '0'
				state = escape2
			case 'b':
				res = append(res, '\b')
				state = initial
			case 'f':
				res = append(res, '\f')
				state = initial
			case 'n':
				res = append(res, '\n')
				state = initial
			case 'r':
				res = append(res, '\r')
				state = initial
			case 't':
				res = append(res, '\t')
				state = initial
			case 'v':
				res = append(res, '\v')
				state = initial
			default:
				res = append(res, b)
				state = initial
			}

		case escape2:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7':
				oct = oct<<3 | (b - '0')
				state = escape3
			case '\\':
				res = append(res, oct)
				state = escape1
			default:
				res = append(res, oct, b)
				state = initial
			}

		case escape3:
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7':
				oct = oct<<3 | (b - '0')
				res = append(res, oct)
				state = initial
			case '\\':
				res = append(res, oct)
				state = escape1
			default:
				res = append(res, oct, b)
				state = initial
			}
		}
	}
	if state == escape2 || state == escape3 {
		res = append(res, oct)
	}

	return string(res)
}
