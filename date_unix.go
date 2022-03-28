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
	if time.January <= opts.month && opts.month <= time.December {
		args = append(args, "--month", strconv.Itoa(int(opts.month)))
	}
	if 1 <= opts.day && opts.day <= 31 {
		args = append(args, "--day", strconv.Itoa(opts.day))
	}
	if opts.year != nil {
		args = append(args, "--year", strconv.Itoa(*opts.year))
	}

	out, err := zenutil.Run(opts.ctx, args)
	str, err := strResult(opts, out, err)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse("2006-01-02", str)
}
