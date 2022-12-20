//go:build !windows && !darwin

package zenity

import (
	"strconv"
	"time"

	"github.com/ncruces/zenity/internal/zenutil"
)

func calendar(text string, opts options) (time.Time, error) {
	args := []string{"--calendar", "--text", quoteMarkup(text), "--date-format", zenutil.DateFormat}
	args = appendGeneral(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendWindowIcon(args, opts)
	if opts.time != nil {
		year, month, day := opts.time.Date()
		args = append(args, "--month", strconv.Itoa(int(month)))
		args = append(args, "--day", strconv.Itoa(day))
		args = append(args, "--year", strconv.Itoa(year))
	}

	out, err := zenutil.Run(opts.ctx, args)
	str, err := strResult(opts, out, err)
	if err != nil {
		return time.Time{}, err
	}
	return zenutil.DateParse(str)
}
