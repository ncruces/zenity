package zenity

import (
	"image/color"
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func SelectColor(options ...Option) (color.Color, error) {
	opts := optsParse(options)

	var data zenutil.Color
	if opts.color != nil {
		r, g, b, _ := opts.color.RGBA()
		data.Color = []uint32{r, g, b}
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
