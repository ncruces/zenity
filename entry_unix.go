//go:build !windows && !darwin

package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func entry(text string, opts options) (string, error) {
	args := []string{"--entry", "--text", text}
	args = appendTitle(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendIcon(args, opts)
	if opts.entryText != "" {
		args = append(args, "--entry-text", opts.entryText)
	}
	if opts.hideText {
		args = append(args, "--hide-text")
	}

	out, err := zenutil.Run(opts.ctx, args)
	return strResult(opts, out, err)
}
