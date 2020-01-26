// +build !windows,!darwin

package zenity

import (
	"image/color"
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectColor(options []Option) (color.Color, error) {
	opts := applyOptions(options)

	args := []string{"--color-selection"}

	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.color != nil {
		args = append(args, "--color", zenutil.UnparseColor(opts.color))
	}
	if opts.showPalette {
		args = append(args, "--show-palette")
	}

	out, err := zenutil.Run(args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() != 255 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return zenutil.ParseColor(string(out)), nil
}
