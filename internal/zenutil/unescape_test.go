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
	{`abc\`, "abc"},
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
	{`a\1\c`, "a\001c"},
	{`a\12\c`, "a\012c"},
	{`a\123\c`, "a\123c"},
	{`a\1\b`, "a\001\b"},
	{`a\12\b`, "a\012\b"},
	{`a\123\b`, "a\123\b"},
	{`abc\1`, "abc\001"},
	{`abc\12`, "abc\012"},
	{`abc\123`, "abc\123"},
	{`abc\1234`, "abc\1234"},
	{`abc\0123`, "abc\0123"},
	{`abc\0012`, "abc\0012"},
	{`abc\4567`, "abc\0567"},
	{`abc\5678`, "abc\1678"},
	{`abc\6789`, "abc\06789"},
	{`abc\7890`, "abc\007890"},
	{"ab\xcdef", "ab\xcdef"},
	{"ab\\\xcdef", "ab\xcdef"},
	{"ab\xcd\\ef", "ab\xcdef"},
	{"ab\\0\xcdef", "ab\000\xcdef"},
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
