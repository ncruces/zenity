//go:build !windows && !darwin

package zenity

import (
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func password(opts options) (string, string, error) {
	args := []string{"--password"}
	args = appendTitle(args, opts)
	args = appendButtons(args, opts)
	if opts.username {
		args = append(args, "--username")
	}

	out, err := zenutil.Run(opts.ctx, args)
	str, err := strResult(opts, out, err)
	if opts.username {
		if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 255 {
			return "", "", ErrUnsupported
		}
		if split := strings.SplitN(str, "|", 2); err == nil && len(split) == 2 {
			return split[0], split[1], nil
		}
	}
	return "", str, err
}
