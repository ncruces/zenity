//go:build !windows && !darwin

package zenity

import (
	"strconv"
	"time"

	"github.com/ncruces/zenity/internal/zenutil"
)

func calendar(text string, opts options) (time.Time, error) {
	args := []string{"--calendar", "--text", text, "--date-format=%F"}
	args = appendTitle(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendIcon(args, opts)
	if opts.day != 0 {
		args = append(args, "--day", strconv.Itoa(opts.day))
	}
	if opts.month != 0 {
		args = append(args, "--month", strconv.Itoa(opts.month))
	}
	if opts.year != 0 {
		args = append(args, "--year", strconv.Itoa(opts.year))
	}

	out, err := zenutil.Run(opts.ctx, args)
	str, err := strResult(opts, out, err)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse("2006-01-02", str)
}
