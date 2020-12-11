package zenity

import (
	"image/color"
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectColor(options []Option) (color.Color, error) {
	opts := applyOptions(options)

	var col color.Color
	if opts.color != nil {
		col = opts.color
	} else {
		col = color.White
	}
	r, g, b, _ := col.RGBA()

	out, err := zenutil.Run(opts.ctx, "color", []uint32{r, g, b})
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return zenutil.ParseColor(string(out)), nil
}
