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
	{"%a %b %e %T %Y", time.ANSIC, "", "Fri Aug  7 06:05:04 2009"},
	{"%a %b %e %T %Z %Y", time.UnixDate, "", "Fri Aug  7 06:05:04 UTC 2009"},
	{"%a %b %d %T %z %Y", time.RubyDate, "E MMM dd HH:mm:ss Z yyyy", "Fri Aug 07 06:05:04 +0000 2009"},
	{"%a, %d %b %Y %T %Z", time.RFC1123, "E, dd MMM yyyy HH:mm:ss zzz", "Fri, 07 Aug 2009 06:05:04 UTC"},
	{"%a, %d %b %Y %T GMT", http.TimeFormat, "E, dd MMM yyyy HH:mm:ss 'GMT'", "Fri, 07 Aug 2009 06:05:04 GMT"},
	{"%Y-%m-%dT%H:%M:%SZ", time.RFC3339[:20], "yyyy-MM-dd'T'HH:mm:ss'Z'", "2009-08-07T06:05:04Z"},
	// Date formats
	{"%v", "_2-Jan-2006", "d-MMM-yyyy", " 7-Aug-2009"},
	{"%F", "2006-01-02", "yyyy-MM-dd", "2009-08-07"},
	{"%D", "01/02/06", "MM/dd/yy", "08/07/09"},
	{"%x", "01/02/06", "MM/dd/yy", "08/07/09"},
	{"%e-%b-%Y", "_2-Jan-2006", "", " 7-Aug-2009"},
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
	{"%-", "%-", "%-", "%-"},
	{"%n", "\n", "\n", "\n"},
	{"%t", "\t", "\t", "\t"},
	{"%q", "", "", "%q"},
	{"%-q", "", "", "%-q"},
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

func TestFormat_ruby(t *testing.T) {
	// https://ruby-doc.org/stdlib-2.6.1/libdoc/date/rdoc/DateTime.html#method-i-strftime
	reference := time.Date(2007, 11, 19, 8, 37, 48, 0, time.FixedZone("", -6*3600))
	tests := []struct {
		format string
		time   string
	}{
		{"Printed on %m/%d/%Y", "Printed on 11/19/2007"},
		{"at %I:%M%p", "at 08:37AM"},
		// Various ISO 8601 formats:
		{"%Y%m%d", "20071119"},                       // Calendar date (basic)
		{"%F", "2007-11-19"},                         // Calendar date (extended)
		{"%Y-%m", "2007-11"},                         // Calendar date, reduced accuracy, specific month
		{"%Y", "2007"},                               // Calendar date, reduced accuracy, specific year
		{"%C", "20"},                                 // Calendar date, reduced accuracy, specific century
		{"%Y%j", "2007323"},                          // Ordinal date (basic)
		{"%Y-%j", "2007-323"},                        // Ordinal date (extended)
		{"%GW%V%u", "2007W471"},                      // Week date (basic)
		{"%G-W%V-%u", "2007-W47-1"},                  // Week date (extended)
		{"%GW%V", "2007W47"},                         // Week date, reduced accuracy, specific week (basic)
		{"%G-W%V", "2007-W47"},                       // Week date, reduced accuracy, specific week (extended)
		{"%H%M%S", "083748"},                         // Local time (basic)
		{"%T", "08:37:48"},                           // Local time (extended)
		{"%H%M", "0837"},                             // Local time, reduced accuracy, specific minute (basic)
		{"%H:%M", "08:37"},                           // Local time, reduced accuracy, specific minute (extended)
		{"%H", "08"},                                 // Local time, reduced accuracy, specific hour
		{"%H%M%S,%L", "083748,000"},                  // Local time with decimal fraction, comma as decimal sign (basic)
		{"%T,%L", "08:37:48,000"},                    // Local time with decimal fraction, comma as decimal sign (extended)
		{"%H%M%S.%L", "083748.000"},                  // Local time with decimal fraction, full stop as decimal sign (basic)
		{"%T.%L", "08:37:48.000"},                    // Local time with decimal fraction, full stop as decimal sign (extended)
		{"%H%M%S%z", "083748-0600"},                  // Local time and the difference from UTC (basic)
		{"%Y%m%dT%H%M%S%z", "20071119T083748-0600"},  // Date and time of day for calendar date (basic)
		{"%Y%jT%H%M%S%z", "2007323T083748-0600"},     // Date and time of day for ordinal date (basic)
		{"%GW%V%uT%H%M%S%z", "2007W471T083748-0600"}, // Date and time of day for week date (basic)
		{"%Y%m%dT%H%M", "20071119T0837"},             // Calendar date and local time (basic)
		{"%FT%R", "2007-11-19T08:37"},                // Calendar date and local time (extended)
		{"%Y%jT%H%MZ", "2007323T0837Z"},              // Ordinal date and UTC of day (basic)
		{"%Y-%jT%RZ", "2007-323T08:37Z"},             // Ordinal date and UTC of day (extended)
		{"%GW%V%uT%H%M%z", "2007W471T0837-0600"},     // Week date and local time and difference from UTC (basic)
		// {"%T%:z", "08:37:48-06:00"},                      // Local time and the difference from UTC (extended)
		// {"%FT%T%:z", "2007-11-19T08:37:48-06:00"},        // Date and time of day for calendar date (extended)
		// {"%Y-%jT%T%:z", "2007-323T08:37:48-06:00"},       // Date and time of day for ordinal date (extended)
		// {"%G-W%V-%uT%T%:z", "2007-W47-1T08:37:48-06:00"}, // Date and time of day for week date (extended)
		// {"%G-W%V-%uT%R%:z", "2007-W47-1T08:37-06:00"},    // Week date and local time and difference from
	}

	for _, test := range tests {
		if got := Format(test.format, reference); got != test.time {
			t.Errorf("Format(%q) = %q, want %q", test.format, got, test.time)
		}
	}
}

func TestFormat_tebeka(t *testing.T) {
	// github.com/tebeka/strftime
	// github.com/hhkbp2/go-strftime
	reference := time.Date(2009, time.November, 10, 23, 1, 2, 3, time.UTC)
	tests := []struct {
		format string
		time   string
	}{
		{"%a", "Tue"},
		{"%A", "Tuesday"},
		{"%b", "Nov"},
		{"%B", "November"},
		{"%c", "Tue Nov 10 23:01:02 2009"}, // we use a different format
		{"%d", "10"},
		{"%H", "23"},
		{"%I", "11"},
		{"%j", "314"},
		{"%m", "11"},
		{"%M", "01"},
		{"%p", "PM"},
		{"%S", "02"},
		{"%U", "45"},
		{"%w", "2"},
		{"%W", "45"},
		{"%x", "11/10/09"},
		{"%X", "23:01:02"},
		{"%y", "09"},
		{"%Y", "2009"},
		{"%Z", "UTC"},
		{"%L", "000"},       // we use a different specifier
		{"%f", "000000"},    // we use a different specifier
		{"%N", "000000003"}, // we use a different specifier

		// Escape
		{"%%%Y", "%2009"},
		{"%3%%", "%3%"},
		{"%3%L", "%3000"},     // we use a different specifier
		{"%3xy%L", "%3xy000"}, // we use a different specifier

		// Embedded
		{"/path/%Y/%m/report", "/path/2009/11/report"},

		// Empty
		{"", ""},
	}

	for _, test := range tests {
		if got := Format(test.format, reference); got != test.time {
			t.Errorf("Format(%q) = %q, want %q", test.format, got, test.time)
		}
	}
}

func TestFormat_fastly(t *testing.T) {
	// github.com/fastly/go-utils/strftime
	timezone, err := time.LoadLocation("MST")
	if err != nil {
		t.Skip("could not load timezone:", err)
	}

	reference := time.Unix(1136239445, 0).In(timezone)

	tests := []struct {
		format string
		time   string
	}{
		{"", ``},

		// invalid formats
		{"%", `%`},
		{"%^", `%^`},
		{"%^ ", `%^ `},
		{"%^ x", `%^ x`},
		{"x%^ x", `x%^ x`},

		// supported locale-invariant formats
		{"%a", `Mon`},
		{"%A", `Monday`},
		{"%b", `Jan`},
		{"%B", `January`},
		{"%C", `20`},
		{"%d", `02`},
		{"%D", `01/02/06`},
		{"%e", ` 2`},
		{"%F", `2006-01-02`},
		{"%G", `2006`},
		{"%g", `06`},
		{"%h", `Jan`},
		{"%H", `15`},
		{"%I", `03`},
		{"%j", `002`},
		{"%k", `15`},
		{"%l", ` 3`},
		{"%m", `01`},
		{"%M", `04`},
		{"%n", "\n"},
		{"%p", `PM`},
		{"%r", `03:04:05 PM`},
		{"%R", `15:04`},
		{"%s", `1136239445`},
		{"%S", `05`},
		{"%t", "\t"},
		{"%T", `15:04:05`},
		{"%u", `1`},
		{"%U", `01`},
		{"%V", `01`},
		{"%w", `1`},
		{"%W", `01`},
		{"%x", `01/02/06`},
		{"%X", `15:04:05`},
		{"%y", `06`},
		{"%Y", `2006`},
		{"%z", `-0700`},
		{"%Z", `MST`},
		{"%%", `%`},

		// supported locale-varying formats
		{"%c", `Mon Jan  2 15:04:05 2006`},
		{"%E", `%E`},
		{"%EF", `%EF`},
		{"%O", `%O`},
		{"%OF", `%OF`},
		{"%P", `pm`},
		{"%+", `Mon Jan  2 15:04:05 MST 2006`},
		{
			"%a|%A|%b|%B|%c|%C|%d|%D|%e|%E|%EF|%F|%G|%g|%h|%H|%I|%j|%k|%l|%m|%M|%O|%OF|%p|%P|%r|%R|%s|%S|%t|%T|%u|%U|%V|%w|%W|%x|%X|%y|%Y|%z|%Z|%%",
			`Mon|Monday|Jan|January|Mon Jan  2 15:04:05 2006|20|02|01/02/06| 2|%E|%EF|2006-01-02|2006|06|Jan|15|03|002|15| 3|01|04|%O|%OF|PM|pm|03:04:05 PM|15:04|1136239445|05|	|15:04:05|1|01|01|1|01|01/02/06|15:04:05|06|2006|-0700|MST|%`,
		},
	}

	for _, test := range tests {
		if got := Format(test.format, reference); got != test.time {
			t.Errorf("Format(%q) = %q, want %q", test.format, got, test.time)
		}
	}
}

func TestFormat_jehiah(t *testing.T) {
	// github.com/jehiah/go-strftime
	reference := time.Unix(1340244776, 0).UTC()
	tests := []struct {
		format string
		time   string
	}{
		{"%Y-%m-%d %H:%M:%S", "2012-06-21 02:12:56"},
		{"aaabbb0123456789%Y", "aaabbb01234567892012"},
		{"%H:%M:%S.%L", "02:12:56.000"}, // jehiah disagrees with Ruby on this one
		{"%0%1%%%2", "%0%1%%2"},
	}

	for _, test := range tests {
		if got := Format(test.format, reference); got != test.time {
			t.Errorf("Format(%q) = %q, want %q", test.format, got, test.time)
		}
	}
}

func TestFormat_lestrrat(t *testing.T) {
	// github.com/lestrrat-go/strftime
	reference := time.Unix(1136239445, 123456789).UTC()
	tests := []struct {
		format string
		time   string
	}{
		{
			`%A %a %B %b %C %c %D %d %e %F %H %h %I %j %k %l %M %m %n %p %R %r %S %T %t %U %u %V %v %W %w %X %x %Y %y %Z %z`,
			"Monday Mon January Jan 20 Mon Jan  2 22:04:05 2006 01/02/06 02  2 2006-01-02 22 Jan 10 002 22 10 04 01 \n PM 22:04 10:04:05 PM 05 22:04:05 \t 01 1 01  2-Jan-2006 01 1 22:04:05 01/02/06 2006 06 UTC +0000",
		},
	}

	for _, test := range tests {
		if got := Format(test.format, reference); got != test.time {
			t.Errorf("Format(%q) = %q, want %q", test.format, got, test.time)
		}
	}
}

func FuzzFormat(f *testing.F) {
	for _, test := range timeTests {
		f.Add(test.format)
	}

	f.Fuzz(func(t *testing.T, fmt string) {
		s := Format(fmt, reference)
		if s == "" && fmt != "" {
			t.Errorf("Format(%q) = %q", fmt, s)
		}
	})
}

func FuzzParse(f *testing.F) {
	for _, test := range timeTests {
		f.Add(test.format, Format(test.format, reference))
	}

	f.Fuzz(func(t *testing.T, format, value string) {
		tm, err := Parse(format, value)
		if tm.IsZero() && err == nil {
			t.Errorf("Parse(%q, %q) = (%v, %v)", format, value, tm, err)
		}
	})
}

func FuzzLayout(f *testing.F) {
	for _, test := range timeTests {
		f.Add(test.format)
	}

	f.Fuzz(func(t *testing.T, format string) {
		layout, err := Layout(format)
		if format != "" && layout == "" && err == nil {
			t.Errorf("Layout(%q) = (%q, %v)", format, layout, err)
		}
	})
}

func FuzzUTS35(f *testing.F) {
	for _, test := range timeTests {
		f.Add(test.format)
	}

	f.Fuzz(func(t *testing.T, format string) {
		pattern, err := UTS35(format)
		if format != "" && pattern == "" && err == nil {
			t.Errorf("UTS35(%q) = (%q, %v)", format, pattern, err)
		}
	})
}
