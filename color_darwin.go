package zenity

import (
	"image/color"
	"os/exec"

	"github.com/ncruces/zenity/internal/osa"
	"github.com/ncruces/zenity/internal/zen"
)

func SelectColor(options ...Option) (color.Color, error) {
	opts := optsParse(options)

	var data osa.Color
	if opts.color != nil {
		r, g, b, _ := opts.color.RGBA()
		data.Color = []float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff}
	}

	out, err := osa.Run("color", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return zen.ParseColor(string(out)), nil
}
