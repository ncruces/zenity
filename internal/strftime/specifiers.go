package strftime

import (
	"strconv"
	"time"
)

func (p *parser) goSpecifiers() {
	// https://strftime.org/
	p.specs = map[byte]string{
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
		'P': "pm",
		'Z': "MST",
		'z': "-0700",
		'L': "000",
		'f': "000000",
		'N': "000000000",

		'+': "Mon Jan _2 15:04:05 MST 2006",
		'c': "Mon Jan _2 15:04:05 2006",
		'v': "_2-Jan-2006",
		'F': "2006-01-02",
		'D': "01/02/06",
		'x': "01/02/06",
		'r': "03:04:05 PM",
		'T': "15:04:05",
		'X': "15:04:05",
		'R': "15:04",

		'%': "%",
		't': "\t",
		'n': "\n",
	}

	p.unpadded = map[byte]string{
		'm': "1",
		'd': "2",
		'I': "3",
		'M': "4",
		'S': "5",
	}
}

func (p *parser) uts35Specifiers() {
	// https://nsdateformatter.com/
	p.specs = map[byte]string{
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
		'v': "d-MMM-yyyy",
		'F': "yyyy-MM-dd",
		'D': "MM/dd/yy",
		'x': "MM/dd/yy",
		'r': "hh:mm:ss a",
		'T': "HH:mm:ss",
		'X': "HH:mm:ss",
		'R': "HH:mm",

		'%': "%",
		't': "\t",
		'n': "\n",
	}

	p.unpadded = map[byte]string{
		'm': "M",
		'd': "d",
		'j': "D",
		'H': "H",
		'I': "h",
		'M': "m",
		'S': "s",
	}
}

func weekNumber(t time.Time, pad, monday bool) string {
	day := t.YearDay()
	offset := int(t.Weekday())
	if monday {
		if offset == 0 {
			offset = 6
		} else {
			offset--
		}
	}

	if day < offset {
		if pad {
			return "00"
		} else {
			return "0"
		}
	}

	n := (day-offset)/7 + 1
	if n < 10 && pad {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
