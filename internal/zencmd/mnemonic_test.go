package zencmd

import (
	"strings"
	"testing"
)

var mnemonicTests = []struct {
	data string
	want string
}{
	{"", ""},
	{"abc", "abc"},
	{"_abc", `abc`},
	{"_a_b_c_", "abc"},
	{"a__c", `a_c`},
	{"a___c", `a_c`},
	{"a____c", `a__c`},
}

func TestStripMnemonic(t *testing.T) {
	t.Parallel()
	for _, test := range mnemonicTests {
		if got := StripMnemonic(test.data); got != test.want {
			t.Errorf("StripMnemonic(%q) = %q; want %q", test.data, got, test.want)
		}
	}
}

func FuzzStripMnemonic(f *testing.F) {
	for _, test := range mnemonicTests {
		f.Add(test.data)
	}

	f.Fuzz(func(t *testing.T, s string) {
		q := quoteMnemonic(s)
		u := StripMnemonic(q)
		if s != u {
			t.Errorf("strip(quote(%q)) = strip(%q) = %q", s, q, u)
		}
	})
}

func quoteMnemonic(s string) string {
	return strings.ReplaceAll(s, "_", "__")
}
