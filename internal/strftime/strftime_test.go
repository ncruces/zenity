package strftime

import (
	"net/http"
	"testing"
	"time"
)

var reference = time.Date(2009, 8, 7, 6, 5, 4, 300000000, time.UTC)

var timeTests = []struct {
	format string
	layout string
	uts35  string
	time   string
}{
	// Date and time formats
	{"%c", time.ANSIC, "E MMM d HH:mm:ss yyyy", "Fri Aug  7 06:05:04 2009"},
	{"%+", time.UnixDate, "E MMM d HH:mm:ss zzz yyyy", "Fri Aug  7 06:05:04 UTC 2009"},
	{"%FT%TZ", time.RFC3339[:20], "yyyy-MM-dd'T'HH:mm:ss'Z'", "2009-08-07T06:05:04Z"},
	{"%a %b %e %H:%M:%S %Y", time.ANSIC, "", "Fri Aug  7 06:05:04 2009"},
	{"%a %b %e %H:%M:%S %Z %Y", time.UnixDate, "", "Fri Aug  7 06:05:04 UTC 2009"},
	{"%a %b %d %H:%M:%S %z %Y", time.RubyDate, "E MMM dd HH:mm:ss Z yyyy", "Fri Aug 07 06:05:04 +0000 2009"},
	{"%a, %d %b %Y %H:%M:%S %Z", time.RFC1123, "E, dd MMM yyyy HH:mm:ss zzz", "Fri, 07 Aug 2009 06:05:04 UTC"},
	{"%a, %d %b %Y %H:%M:%S GMT", http.TimeFormat, "E, dd MMM yyyy HH:mm:ss 'GMT'", "Fri, 07 Aug 2009 06:05:04 GMT"},
	{"%Y-%m-%dT%H:%M:%SZ", time.RFC3339[:20], "yyyy-MM-dd'T'HH:mm:ss'Z'", "2009-08-07T06:05:04Z"},
	// Date formats
	{"%F", "2006-01-02", "yyyy-MM-dd", "2009-08-07"},
	{"%D", "01/02/06", "MM/dd/yy", "08/07/09"},
	{"%x", "01/02/06", "MM/dd/yy", "08/07/09"},
	{"%Y-%m-%d", "2006-01-02", "yyyy-MM-dd", "2009-08-07"},
	{"%m/%d/%y", "01/02/06", "MM/dd/yy", "08/07/09"},
	// Time formats
	{"%R", "15:04", "HH:mm", "06:05"},
	{"%T", "15:04:05", "HH:mm:ss", "06:05:04"},
	{"%X", "15:04:05", "HH:mm:ss", "06:05:04"},
	{"%r", "03:04:05 PM", "hh:mm:ss a", "06:05:04 AM"},
	{"%H:%M", "15:04", "HH:mm", "06:05"},
	{"%H:%M:%S", "15:04:05", "HH:mm:ss", "06:05:04"},
	{"%I:%M:%S %p", "03:04:05 PM", "hh:mm:ss a", "06:05:04 AM"},
	// Misc
	{"%g", "", "YY", "09"},
	{"%V/%G", "", "ww/YYYY", "32/2009"},
	{"%Cth Century Fox", "", "", "20th Century Fox"},
	{"%-I o'clock", "3 o'clock", "h 'o''clock'", "6 o'clock"},
	{"%-A, the %uth day of the week", "", "", "Friday, the 5th day of the week"},
	{"%Nns since %T", "", "SSSSSSSSS'ns since 'HH:mm:ss", "300000000ns since 06:05:04"},
	{"%-S.%Ls since %R", "5.000s since 15:04", "s.SSS's since 'HH:mm", "4.300s since 06:05"},
	// Parsing
	{"", "", "", ""},
	{"%", "%", "%", "%"},
	{"%-", "-", "-", "-"},
	{"%n", "\n", "\n", "\n"},
	{"%t", "\t", "\t", "\t"},
	{"%q", "", "", ""},
	{"%-q", "", "", ""},
	{"'", "'", "''", "'"},
	{"100%", "", "100%", "100%"},
	{"Monday", "", "'Monday'", "Monday"},
	{"January", "", "'January'", "January"},
	{"MST", "", "'MST'", "MST"},
	{"AM", "AM", "'AM'", "AM"},
	{"am", "am", "'am'", "am"},
	{"PM", "", "'PM'", "PM"},
	{"pm", "", "'pm'", "pm"},
}

func TestFormat(t *testing.T) {
	for _, test := range timeTests {
		if got := Format(test.format, reference); got != test.time {
			t.Errorf("Format(%q) = %q, want %q", test.format, got, test.time)
		}
	}
}

func TestLayout(t *testing.T) {
	for _, test := range timeTests {
		if got, err := Layout(test.format); err != nil && test.layout != "" {
			t.Errorf("Layout(%q) = %v", test.format, err)
		} else if got != test.layout {
			t.Errorf("Layout(%q) = %q, want %q", test.format, got, test.layout)
		}
	}
}

func TestUTS35(t *testing.T) {
	for _, test := range timeTests {
		if got, err := UTS35(test.format); err != nil && test.uts35 != "" {
			t.Errorf("UTS35(%q) = %v", test.format, err)
		} else if got != test.uts35 {
			t.Errorf("UTS35(%q) = %q, want %q", test.format, got, test.uts35)
		}
	}
}
