package zenity

import (
	"image/color"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectColor(opts options) (color.Color, error) {
	var data zenutil.Color
	if opts.attach != nil {
		data.Application = opts.attach
	}
	if i, ok := opts.windowIcon.(string); ok {
		data.WindowIcon = i
	}

	var col color.Color
	if opts.color != nil {
		col = opts.color
	} else {
		col = color.White
	}
	r, g, b, _ := col.RGBA()
	data.Color = [3]float32{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
	}

	out, err := zenutil.Run(opts.ctx, "color", data)
	str, err := strResult(opts, out, err)
	if err != nil {
		return nil, err
	}
	return zenutil.ParseColor(str), nil
}
