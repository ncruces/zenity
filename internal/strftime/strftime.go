package strftime

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func Format(fmt string, t time.Time) string {
	var res strings.Builder

	var parser parser
	parser.fmt = fmt
	parser.goSpecifiers()

	parser.writeLit = res.WriteByte

	parser.writeFmt = func(fmt string) error {
		switch fmt {
		case "000", "000000", "000000000":
			res.WriteString(t.Format("." + fmt)[1:])
		default:
			res.WriteString(t.Format(fmt))
		}
		return nil
	}

	parser.fallback = func(spec byte, pad bool) error {
		switch spec {
		default:
			res.WriteByte('%')
			if !pad {
				res.WriteByte('-')
			}
			res.WriteByte(spec)
		case 'C':
			s := t.Format("2006")
			res.WriteString(s[:len(s)-2])
		case 'g':
			y, _ := t.ISOWeek()
			res.WriteString(time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC).Format("06"))
		case 'G':
			y, _ := t.ISOWeek()
			res.WriteString(time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"))
		case 'V':
			_, w := t.ISOWeek()
			if w < 10 && pad {
				res.WriteByte('0')
			}
			res.WriteString(strconv.Itoa(w))
		case 'W':
			res.WriteString(weekNumber(t, pad, true))
		case 'U':
			res.WriteString(weekNumber(t, pad, false))
		case 'w':
			w := int(t.Weekday())
			res.WriteString(strconv.Itoa(w))
		case 'u':
			if w := int(t.Weekday()); w == 0 {
				res.WriteByte('7')
			} else {
				res.WriteString(strconv.Itoa(w))
			}
		case 'k':
			h := t.Hour()
			if h < 10 {
				res.WriteByte(' ')
			}
			res.WriteString(strconv.Itoa(h))
		case 'l':
			h := t.Hour()
			if h == 0 {
				h = 12
			} else if h > 12 {
				h -= 12
			}
			if h < 10 {
				res.WriteByte(' ')
			}
			res.WriteString(strconv.Itoa(h))
		case 's':
			res.WriteString(strconv.FormatInt(t.Unix(), 10))
		}
		return nil
	}

	parser.parse()
	return res.String()
}

func Parse(fmt, value string) (time.Time, error) {
	layout, err := Layout(fmt)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(layout, value)
}

func Layout(fmt string) (string, error) {
	var res strings.Builder

	var parser parser
	parser.fmt = fmt
	parser.goSpecifiers()

	parser.writeLit = func(b byte) error {
		if '0' <= b && b <= '9' {
			return errors.New("strftime: unsupported literal digit: '" + string(b) + "'")
		}
		res.WriteByte(b)
		if b == 'M' || b == 'T' || b == 'm' || b == 'n' {
			cur := res.String()
			switch {
			case strings.HasSuffix(cur, "Jan"):
				return errors.New("strftime: unsupported literal: 'Jan'")
			case strings.HasSuffix(cur, "Mon"):
				return errors.New("strftime: unsupported literal: 'Mon'")
			case strings.HasSuffix(cur, "MST"):
				return errors.New("strftime: unsupported literal: 'MST'")
			case strings.HasSuffix(cur, "PM"):
				return errors.New("strftime: unsupported literal: 'PM'")
			case strings.HasSuffix(cur, "pm"):
				return errors.New("strftime: unsupported literal: 'pm'")
			}
		}
		return nil
	}

	parser.writeFmt = func(fmt string) error {
		switch fmt {
		case "000", "000000", "000000000":
			if cur := res.String(); !(strings.HasSuffix(cur, ".") || strings.HasSuffix(cur, ",")) {
				return errors.New("strftime: unsupported specifier: fractional seconds must follow '.' or ','")
			}
		}
		res.WriteString(fmt)
		return nil
	}

	parser.fallback = func(spec byte, pad bool) error {
		return errors.New("strftime: unsupported specifier: %" + string(spec))
	}

	if err := parser.parse(); err != nil {
		return "", err
	}

	parser.writeFmt("")
	return res.String(), nil
}

func UTS35(fmt string) (string, error) {
	var parser parser
	parser.fmt = fmt
	parser.uts35Specifiers()

	const quote = '\''
	var literal bool
	var res strings.Builder

	parser.writeLit = func(b byte) error {
		if b == quote {
			res.WriteByte(quote)
			res.WriteByte(quote)
			return nil
		}
		if !literal && ('a' <= b && b <= 'z' || 'A' <= b && b <= 'Z') {
			literal = true
			res.WriteByte(quote)
		}
		res.WriteByte(b)
		return nil
	}

	parser.writeFmt = func(fmt string) error {
		if literal {
			literal = false
			res.WriteByte(quote)
		}
		res.WriteString(fmt)
		return nil
	}

	parser.fallback = func(spec byte, pad bool) error {
		return errors.New("strftime: unsupported specifier: %" + string(spec))
	}

	if err := parser.parse(); err != nil {
		return "", err
	}

	parser.writeFmt("")
	return res.String(), nil
}
