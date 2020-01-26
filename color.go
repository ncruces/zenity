package zenity

import "image/color"

// SelectColor displays the color selection dialog.
//
// Returns nil on cancel.
//
// Valid options: Title, Color, ShowPalette.
func SelectColor(options ...Option) (color.Color, error) {
	return selectColor(options)
}

// Color returns an Option to set the color.
func Color(c color.Color) Option {
	return funcOption(func(o *options) { o.color = c })
}

// ShowPalette returns an Option to show the palette.
func ShowPalette() Option {
	return funcOption(func(o *options) { o.showPalette = true })
}
