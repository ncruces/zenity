package zenutil

import (
	"testing"
)

var unescapeTests = []struct {
	data string
	want string
}{
	{``, ""},
	{`abc`, "abc"},
	{`ab\c`, "abc"},
	{`a\bc`, "a\bc"},
	{`abc\f`, "abc\f"},
	{`abc\n`, "abc\n"},
	{`abc\r`, "abc\r"},
	{`abc\t`, "abc\t"},
	{`abc\v`, "abc\v"},
	{`a\1c`, "a\001c"},
	{`a\12c`, "a\012c"},
	{`a\123c`, "a\123c"},
	{`a\1\b`, "a\001\b"},
	{`a\12\b`, "a\012\b"},
	{`a\123\b`, "a\123\b"},
	{`abc\1`, "abc\001"},
	{`abc\12`, "abc\012"},
	{`abc\123`, "abc\123"},
	{`abc\1234`, "abc\1234"},
	{`abc\001`, "abc\001"},
	{`abc\012`, "abc\012"},
	{`abc\123`, "abc\123"},
	{`abc\4`, "abc\004"},
	{`abc\45`, "abc\045"},
	{`abc\456`, "abc\056"},
	{`abc\4567`, "abc\0567"},
}

func TestUnescape(t *testing.T) {
	for _, test := range unescapeTests {
		if got := Unescape(test.data); got != test.want {
			t.Errorf("Unescape(%q) = %q; want %q", test.data, got, test.want)
		}
	}
}

func FuzzUnescape(f *testing.F) {
	for _, test := range unescapeTests {
		f.Add(test.data)
	}

	f.Fuzz(func(t *testing.T, e string) {
		u := Unescape(e)
		switch {
		case u == e:
			return
		case len(u) < len(e):
			return
		}
		t.Errorf("Unescape(%q) = %q", e, u)
	})
}
