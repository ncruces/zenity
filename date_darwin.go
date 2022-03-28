package zenity

import (
	"time"

	"github.com/ncruces/zenity/internal/zenutil"
)

func calendar(text string, opts options) (time.Time, error) {
	var date zenutil.Date

	year, month, day := time.Now().Date()
	if time.January <= opts.month && opts.month <= time.December {
		month = opts.month
	}
	if 1 <= opts.day && opts.day <= 31 {
		day = opts.day
	}
	if opts.year != nil {
		year = *opts.year
	}
	date.Date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix()
	date.OK, date.Cancel, date.Extra = getAlertButtons(opts)
	date.Format = "yyyy-MM-dd"

	if opts.title != nil {
		date.Text = *opts.title
		date.Info = text
	} else {
		date.Text = text
	}

	out, err := zenutil.Run(opts.ctx, "date", date)
	str, err := strResult(opts, out, err)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse("2006-01-02", str)
}
