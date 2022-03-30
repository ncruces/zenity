package zenutil

import (
	"strings"
	"time"
)

// Strftime is internal.
func Strftime(fmt string, t time.Time) string {
	var res strings.Builder
	writeLit := res.WriteByte
	writeFmt := func(fmt string) (int, error) {
		return res.WriteString(t.Format(fmt))
	}
	strftimeGo(fmt, writeLit, writeFmt)
	return res.String()
}

// StrftimeLayout is internal.
func StrftimeLayout(fmt string) string {
	var res strings.Builder
	strftimeGo(fmt, res.WriteByte, res.WriteString)
	return res.String()
}

func strftimeGo(fmt string, writeLit func(byte) error, writeFmt func(string) (int, error)) {
	// https://strftime.org/
	fmts := map[byte]string{
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
	}

	unpaded := map[byte]string{
		'm': "1",
		'd': "2",
		'I': "3",
		'M': "4",
		'S': "5",
	}

	parser(fmt, fmts, unpaded, writeLit, writeFmt)
}

// StrftimeUTS35 is internal.
func StrftimeUTS35(fmt string) string {
	// https://nsdateformatter.com/
	fmts := map[byte]string{
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
	}

	unpaded := map[byte]string{
		'm': "M",
		'd': "d",
		'j': "D",
		'H': "H",
		'I': "h",
		'M': "m",
		'S': "s",
	}

	const quote = '\''
	var literal bool
	var res strings.Builder

	writeLit := func(b byte) error {
		if b == quote {
			res.WriteByte(quote)
			return res.WriteByte(quote)
		}
		if !literal && ('a' <= b && b <= 'z' || 'A' <= b && b <= 'Z') {
			literal = true
			res.WriteByte(quote)
		}
		return res.WriteByte(b)
	}

	writeFmt := func(s string) (int, error) {
		if literal {
			literal = false
			res.WriteByte(quote)
		}
		return res.WriteString(s)
	}

	parser(fmt, fmts, unpaded, writeLit, writeFmt)
	writeFmt("")

	return res.String()
}

func parser(
	fmt string,
	formats, unpadded map[byte]string,
	writeLit func(byte) error, writeFmt func(string) (int, error)) {

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
				writeLit(b)
			}

		case special:
			if b == '-' {
				state = padding
				continue
			}
			if s, ok := formats[b]; ok {
				writeFmt(s)
			} else {
				writeLit('%')
				writeLit(b)
			}
			state = initial

		case padding:
			if s, ok := unpadded[b]; ok {
				writeFmt(s)
			} else if s, ok := formats[b]; ok {
				writeFmt(s)
			} else {
				writeLit('%')
				writeLit('-')
				writeLit(b)
			}
			state = initial
		}
	}

	switch state {
	case padding:
		writeLit('%')
		fallthrough
	case special:
		writeLit('-')
	}
}
