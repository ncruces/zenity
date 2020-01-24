package zenity

import (
	"image/color"
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectColor(options ...Option) (color.Color, error) {
	opts := applyOptions(options)

	var data zenutil.Color
	if opts.color != nil {
		n := color.NRGBA64Model.Convert(opts.color).(color.NRGBA64)
		data.Color = []uint16{n.R, n.G, n.B}
	}

	out, err := zenutil.Run("color", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return zenutil.ParseColor(string(out)), nil
}
