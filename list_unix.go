//go:build !windows && !darwin

package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func list(text string, items []string, opts options) (string, error) {
	args := []string{"--list", "--column=", "--hide-header", "--text", text}
	args = appendTitle(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendIcon(args, opts)
	args = append(args, items...)

	out, err := zenutil.Run(opts.ctx, args)
	return strResult(opts, out, err)
}

func listMultiple(text string, items []string, opts options) ([]string, error) {
	args := []string{"--list", "--column=", "--hide-header", "--text", text, "--multiple", "--separator", zenutil.Separator}
	args = appendTitle(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendIcon(args, opts)
	args = append(args, items...)

	out, err := zenutil.Run(opts.ctx, args)
	return lstResult(opts, out, err)
}
