package zenutil

import (
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestColor_names(t *testing.T) {
	tests := []string{
		"chocolate",
		"lime",
		"olive",
		"orange",
		"plum",
		"salmon",
		"tomato",
	}
	for _, test := range tests {
		c1 := colornames.Map[test]
		c2 := ParseColor(test)
		c3 := ParseColor(UnparseColor(c1))
		if !ColorEquals(c1, c2) {
			t.Errorf("ParseColor(%s) = %v; want %v", test, c2, c1)
		}
		if !ColorEquals(c1, c3) {
			t.Errorf("ParseColor(UnparseColor(%v)) = %v; want %v", c1, c3, c1)
		}
	}
}

func TestColor_colors(t *testing.T) {
	tests := []color.Color{
		color.Black,
		color.White,
		color.Opaque,
		color.Transparent,
	}
	for _, test := range tests {
		c := ParseColor(UnparseColor(test))
		if !ColorEquals(c, test) {
			t.Errorf("ParseColor(UnparseColor(%v)) = %v; want %v", test, c, test)
		}
	}
}

func TestColor_strings(t *testing.T) {
	tests := []struct {
		data string
		want color.Color
	}{
		{"#000", color.Black},
		{"#000f", color.Black},
		{"#000000", color.Black},
		{"#000000ff", color.Black},
		{"#fff", color.White},
		{"#ffff", color.White},
		{"#ffffff", color.White},
		{"#ffffffff", color.White},
		{"#FFF", color.Opaque},
		{"#FFFF", color.Opaque},
		{"#FFFFFF", color.Opaque},
		{"#FFFFFFFF", color.Opaque},
		{"#0000", color.Transparent},
		{"#00000000", color.Transparent},
		{"#8888", color.NRGBA{0x88, 0x88, 0x88, 0x88}},
		{"#80808080", color.NRGBA{0x80, 0x80, 0x80, 0x80}},
		{"rgb(128,128,128)", color.NRGBA{0x80, 0x80, 0x80, 0xff}},
		{"rgba(128,128,128,0.5)", color.NRGBA{0x80, 0x80, 0x80, 0x80}},
		{"rgba(128,128,128,1.0)", color.NRGBA{0x80, 0x80, 0x80, 0xff}},
		{"not a color", nil},
	}
	for _, test := range tests {
		c := ParseColor(test.data)
		if !ColorEquals(c, test.want) {
			t.Errorf("ParseColor(%s) = %v; want %v", test.data, c, test.want)
		}
	}
}
