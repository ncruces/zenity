package zenity_test

import (
	"image/color"

	"github.com/ncruces/zenity"
)

func ExampleSelectColor() {
	zenity.SelectColor(
		zenity.Color(color.RGBA{R: 0x66, G: 0x33, B: 0x99, A: 0xff}))
	// Output:
}

func ExampleSelectColor_palette() {
	zenity.SelectColor(
		zenity.ShowPalette(),
		zenity.Color(color.RGBA{R: 0x66, G: 0x33, B: 0x99, A: 0xff}))
	// Output:
}
