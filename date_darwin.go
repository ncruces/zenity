package zenity

import (
	"time"

	"github.com/ncruces/zenity/internal/zenutil"
)

func calendar(text string, opts options) (time.Time, error) {
	var date zenutil.Date

	date.OK, date.Cancel, date.Extra = getAlertButtons(opts)
	date.Format = zenutil.StrftimeUTS35(zenutil.DateFormat)
	if opts.time != nil {
		date.Date = opts.time.Unix()
	}

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
	layout := zenutil.StrftimeLayout(zenutil.DateFormat)
	return time.Parse(layout, str)
}
