package zenity

import (
	"time"

	"github.com/ncruces/zenity/internal/zenutil"
)

func calendar(text string, opts options) (t time.Time, err error) {
	var data zenutil.Date

	data.OK, data.Cancel, data.Extra = getAlertButtons(opts)
	data.Format, err = zenutil.DateUTS35()
	if err != nil {
		return
	}
	if opts.time != nil {
		data.Date = ptr(opts.time.Unix())
	}

	if opts.title != nil {
		data.Text = *opts.title
		data.Info = text
	} else {
		data.Text = text
	}
	if i, ok := opts.windowIcon.(string); ok {
		data.WindowIcon = i
	}

	out, err := zenutil.Run(opts.ctx, "date", data)
	str, err := strResult(opts, out, err)
	if err != nil {
		return
	}
	return zenutil.DateParse(str)
}
