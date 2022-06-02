//go:build !windows && !darwin

package zenity

import "github.com/ncruces/zenity/internal/zenutil"

func password(opts options) (string, string, error) {
	args := []string{"--password"}
	args = appendGeneral(args, opts)
	args = appendButtons(args, opts)
	if opts.username {
		args = append(args, "--username")
	}

	out, err := zenutil.Run(opts.ctx, args)
	return pwdResult("|", opts, out, err)
}
