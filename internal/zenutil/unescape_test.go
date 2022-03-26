package zenutil

import (
	"testing"
)

func TestUnescape(t *testing.T) {
	tests := []struct {
		data string
		want string
	}{
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
	}
	for _, tt := range tests {
		if got := Unescape(tt.data); got != tt.want {
			t.Errorf("Unescape(%q) = %q; want %q", tt.data, got, tt.want)
		}
	}
}

func FuzzUnescape(f *testing.F) {
	f.Add(`abc`)
	f.Add(`ab\c`)
	f.Add(`a\bc`)
	f.Add(`abc\f`)
	f.Add(`abc\n`)
	f.Add(`abc\r`)
	f.Add(`abc\t`)
	f.Add(`abc\v`)
	f.Add(`a\1c`)
	f.Add(`a\12c`)
	f.Add(`a\123c`)
	f.Add(`a\1\b`)
	f.Add(`a\12\b`)
	f.Add(`a\123\b`)
	f.Add(`abc\1`)
	f.Add(`abc\12`)
	f.Add(`abc\123`)
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
