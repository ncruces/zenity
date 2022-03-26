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
		{"rgba(128,128,128,0)", color.NRGBA{0x80, 0x80, 0x80, 0x00}},
		{"rgba(128,128,128,1)", color.NRGBA{0x80, 0x80, 0x80, 0xff}},
		{"rgba(128,128,128,0.0)", color.NRGBA{0x80, 0x80, 0x80, 0x00}},
		{"rgba(128,128,128,1.0)", color.NRGBA{0x80, 0x80, 0x80, 0xff}},
		{"not a color", nil},
		{"#0", nil},
		{"#00", nil},
		{"#000", color.Black},
		{"#0000", color.Transparent},
		{"#00000", nil},
		{"#000000", color.Black},
		{"#0000000", nil},
		{"#00000000", color.Transparent},
		{"#000000000", nil},
		{"rgb(-1,-1,-1)", nil},
		{"rgb(256,256,256)", nil},
		{"rgb(128,128,128,0.5)", nil},
		{"rgb(127.5,127.5,127.5)", nil},
		{"rgba(127.5,127.5,127.5,0.5)", nil},
		{"rgba(128,128,128)", nil},
	}
	for _, test := range tests {
		c := ParseColor(test.data)
		if !ColorEquals(c, test.want) {
			t.Errorf("ParseColor(%s) = %v; want %v", test.data, c, test.want)
		}
	}
}

func FuzzParseColor(f *testing.F) {
	f.Add("#000")
	f.Add("#000f")
	f.Add("#000000")
	f.Add("#000000ff")
	f.Add("#fff")
	f.Add("#ffff")
	f.Add("#ffffff")
	f.Add("#ffffffff")
	f.Add("#FFF")
	f.Add("#FFFF")
	f.Add("#FFFFFF")
	f.Add("#FFFFFFFF")
	f.Add("#0")
	f.Add("#00")
	f.Add("#000")
	f.Add("#0000")
	f.Add("#00000")
	f.Add("#000000")
	f.Add("#0000000")
	f.Add("#00000000")
	f.Add("#000000000")
	f.Add("#8888")
	f.Add("#80808080")
	f.Add("rgb(-1,-1,-1)")
	f.Add("rgb(128,128,128)")
	f.Add("rgb(256,256,256)")
	f.Add("rgb(128,128,128,0.5)")
	f.Add("rgb(127.5,127.5,127.5)")
	f.Add("rgba(128,128,128)")
	f.Add("rgba(128,128,128,0)")
	f.Add("rgba(128,128,128,1)")
	f.Add("rgba(128,128,128,0.0)")
	f.Add("rgba(128,128,128,0.5)")
	f.Add("rgba(128,128,128,1.0)")
	f.Add("rgba(127.5,127.5,127.5,0.5)")
	f.Add("not a color")
	f.Add("")

	for _, name := range colornames.Names {
		f.Add(name)
	}

	f.Fuzz(func(t *testing.T, s string) {
		ParseColor(s)
	})
}
