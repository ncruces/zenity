//go:build !windows && !darwin

package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func progress(opts options) (ProgressDialog, error) {
	args := []string{"--progress"}
	args = appendTitle(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendIcon(args, opts)
	if opts.maxValue == 0 {
		opts.maxValue = 100
	}
	if opts.maxValue < 0 {
		args = append(args, "--pulsate")
	}
	if opts.noCancel {
		args = append(args, "--no-cancel")
	}
	if opts.timeRemaining {
		args = append(args, "--time-remaining")
	}
	return zenutil.RunProgress(opts.ctx, opts.maxValue, opts.extraButton, args)
}
