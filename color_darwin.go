package zenity

import (
	"image/color"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectColor(opts options) (color.Color, error) {
	var col color.Color
	if opts.color != nil {
		col = opts.color
	} else {
		col = color.White
	}
	r, g, b, _ := col.RGBA()

	out, err := zenutil.Run(opts.ctx, "color", []float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
	})
	str, err := strResult(opts, out, err)
	if err != nil {
		return nil, err
	}
	return zenutil.ParseColor(str), nil
}
