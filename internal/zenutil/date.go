package zenutil

import "strings"

// Strftime is internal.
func Strftime(fmt string) string {
	// https://strftime.org/
	return strftime(fmt, map[byte]string{
		'B': "January",
		'b': "Jan",
		'h': "Jan",
		'm': "01",
		'A': "Monday",
		'a': "Mon",
		'e': "_2",
		'd': "02",
		'j': "002",
		'H': "15",
		'I': "03",
		'M': "04",
		'S': "05",
		'Y': "2006",
		'y': "06",
		'p': "PM",
		'Z': "MST",
		'z': "-0700",
		'L': "000",
		'f': "000000",
		'N': "000000000",

		'+': "Mon Jan _2 15:04:05 MST 2006",
		'c': "Mon Jan _2 15:04:05 2006",
		'F': "2006-01-02",
		'D': "01/02/06",
		'x': "01/02/06",
		'r': "03:04:05 PM",
		'T': "15:04:05",
		'X': "15:04:05",
		'R': "15:04",

		'%': "%",
		't': "\t",
		'n': LineBreak,
	}, map[byte]string{
		'm': "1",
		'd': "2",
		'I': "3",
		'M': "4",
		'S': "5",
	})
}

// StrftimeUTS35 is internal.
func StrftimeUTS35(fmt string) string {
	// https://nsdateformatter.com/
	return strftime(fmt, map[byte]string{
		'B': "MMMM",
		'b': "MMM",
		'h': "MMM",
		'm': "MM",
		'A': "EEEE",
		'a': "E",
		'd': "dd",
		'j': "DDD",
		'H': "HH",
		'I': "hh",
		'M': "mm",
		'S': "ss",
		'Y': "yyyy",
		'y': "yy",
		'G': "YYYY",
		'g': "YY",
		'V': "ww",
		'p': "a",
		'Z': "zzz",
		'z': "Z",
		'L': "SSS",
		'f': "SSSSSS",
		'N': "SSSSSSSSS",

		'+': "E MMM d HH:mm:ss zzz yyyy",
		'c': "E MMM d HH:mm:ss yyyy",
		'F': "yyyy-MM-dd",
		'D': "MM/dd/yy",
		'x': "MM/dd/yy",
		'r': "hh:mm:ss a",
		'T': "HH:mm:ss",
		'X': "HH:mm:ss",
		'R': "HH:mm",

		'%': "%",
		't': "\t",
		'n': LineBreak,
	}, map[byte]string{
		'm': "M",
		'd': "d",
		'j': "D",
		'H': "H",
		'I': "h",
		'M': "m",
		'S': "s",
	})
}

func strftime(fmt string, formats, unpadded map[byte]string) string {
	var res strings.Builder
	res.Grow(len(fmt))

	const (
		initial = iota
		special
		padding
	)

	state := initial
	for _, b := range []byte(fmt) {
		switch state {
		case initial:
			if b == '%' {
				state = special
			} else {
				res.WriteByte(b)
			}

		case special:
			if b == '-' {
				state = padding
				continue
			}
			if s, ok := formats[b]; ok {
				res.WriteString(s)
			} else {
				res.WriteByte('%')
				res.WriteByte(b)
			}
			state = initial

		case padding:
			if s, ok := unpadded[b]; ok {
				res.WriteString(s)
			} else if s, ok := formats[b]; ok {
				res.WriteString(s)
			} else {
				res.WriteString("%-")
				res.WriteByte(b)
			}
			state = initial
		}
	}

	switch state {
	case padding:
		res.WriteString("%-")
	case special:
		res.WriteByte('%')
	}
	return res.String()
}
