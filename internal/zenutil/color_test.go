package zenutil

import (
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestColor(t *testing.T) {
	tests := []string{
		"chocolate",
		"lime",
		"olive",
		"orange",
		"plum",
		"salmon",
		"tomato",
	}
	eq := func(c1, c2 color.Color) bool {
		r1, g1, b1, a1 := c1.RGBA()
		r2, g2, b2, a2 := c2.RGBA()
		return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
	}
	for _, test := range tests {
		c1 := colornames.Map[test]
		c2 := ParseColor(test)
		c3 := ParseColor(UnparseColor(c2))
		if !eq(c1, c2) {
			t.Errorf("ParseColor(%s) = %v, want %v", test, c2, c1)
		}
		if !eq(c1, c3) {
			t.Errorf("ParseColor(UnparseColor(%s)) = %v, want %v", test, c3, c1)
		}
	}
}
