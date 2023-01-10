package zenity

import (
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func notify(text string, opts options) error {
	var data zenutil.Notify
	data.Options.Title = opts.title
	if sub, body, cut := strings.Cut(text, "\n"); cut {
		data.Options.Subtitle = sub
		data.Text = body
	} else {
		data.Text = text
	}

	_, err := zenutil.Run(opts.ctx, "notify", data)
	if err != nil {
		return err
	}
	return nil
}
