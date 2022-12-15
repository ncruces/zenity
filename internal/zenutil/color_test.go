package zenutil

import (
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestColorEquals(t *testing.T) {
	t.Parallel()
	if ColorEquals(nil, nil) == false {
		t.Error("ColorEquals(nil, nil) == false")
	}
	if ColorEquals(nil, color.Black) == true {
		t.Error("ColorEquals(nil, color.Black) == true")
	}
	if ColorEquals(color.Black, nil) == true {
		t.Error("ColorEquals(color.Black, nil) == true")
	}
	if ColorEquals(color.Black, color.Black) == false {
		t.Error("ColorEquals(color.Black, color.Black) == false")
	}
	if ColorEquals(color.Black, color.White) == true {
		t.Error("ColorEquals(color.Black, color.White) == true")
	}
	if ColorEquals(color.Black, colornames.Red) == true {
		t.Error("ColorEquals(color.Black, colornames.Red) == true")
	}
	if ColorEquals(color.Black, colornames.Green) == true {
		t.Error("ColorEquals(color.Black, colornames.Green) == true")
	}
	if ColorEquals(color.Black, colornames.Blue) == true {
		t.Error("ColorEquals(color.Black, colornames.Blue) == true")
	}
	if ColorEquals(color.Black, color.Transparent) == true {
		t.Error("ColorEquals(color.Black, color.Transparent) == true")
	}
}

func TestColor_names(t *testing.T) {
	t.Parallel()
	for _, name := range colornames.Names {
		c1 := colornames.Map[name]
		c2 := ParseColor(name)
		c3 := ParseColor(UnparseColor(c1))
		if !ColorEquals(c1, c2) {
			t.Errorf("ParseColor(%q) = %v; want %v", name, c2, c1)
		}
		if !ColorEquals(c1, c3) {
			t.Errorf("ParseColor(UnparseColor(%v)) = %v; want %v", c1, c3, c1)
		}
	}
}

func TestColor_colors(t *testing.T) {
	t.Parallel()
	colors := []color.Color{
		color.Black,
		color.White,
		color.Opaque,
		color.Transparent,
	}
	for _, color := range colors {
		c := ParseColor(UnparseColor(color))
		if !ColorEquals(c, color) {
			t.Errorf("ParseColor(UnparseColor(%v)) = %v; want %v", color, c, color)
		}
	}
}

var colorTests = []struct {
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
	{"rgba(128,128,128,0.5)", color.NRGBA{0x80, 0x80, 0x80, 0x80}},
	{"rgba(128,128,128,1.0)", color.NRGBA{0x80, 0x80, 0x80, 0xff}},
	{"not a color", nil},
	{"", nil},
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

func TestParseColor(t *testing.T) {
	t.Parallel()
	for _, test := range colorTests {
		c := ParseColor(test.data)
		if !ColorEquals(c, test.want) {
			t.Errorf("ParseColor(%q) = %v; want %v", test.data, c, test.want)
		}
	}
}

func FuzzParseColor(f *testing.F) {
	for _, test := range colorTests {
		f.Add(test.data)
	}
	for _, name := range colornames.Names {
		f.Add(name)
	}

	f.Fuzz(func(t *testing.T, s string) {
		ParseColor(s)
	})
}
