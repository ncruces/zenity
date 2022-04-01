package zenity

import (
	"time"

	"github.com/ncruces/go-strftime"
	"github.com/ncruces/zenity/internal/zenutil"
)

func calendar(text string, opts options) (t time.Time, err error) {
	var date zenutil.Date

	date.OK, date.Cancel, date.Extra = getAlertButtons(opts)
	date.Format, err = strftime.UTS35(zenutil.DateFormat)
	if err != nil {
		return
	}
	if opts.time != nil {
		unix := opts.time.Unix()
		date.Date = &unix
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
		return
	}
	return strftime.Parse(zenutil.DateFormat, str)
}
