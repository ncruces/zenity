package zencmd

import (
	"testing"
)

var markupTests = []struct {
	data string
	want string
}{
	// success cases
	{"", ``},
	{"abc", `abc`},
	{"&lt;", `<`},
	{"&amp;", `&`},
	{"&quot;", `"`},
	{"<i></i>", ``},
	{"<i>abc</i>", `abc`},
	{"<i>&quot;</i>", `"`},
	{"<!--abc-->", ``},

	// failure cases
	{"<", `<`},
	{"<i", `<i`},
	{"<i>", `<i>`},
	{"<i></b>", `<i></b>`},
	{"<i>&amp</i>", `<i>&amp</i>`},
}

func TestStripMarkup(t *testing.T) {
	t.Parallel()
	for _, test := range markupTests {
		if got := StripMarkup(test.data); got != test.want {
			t.Errorf("StripMarkup(%q) = %q; want %q", test.data, got, test.want)
		}
	}
}
