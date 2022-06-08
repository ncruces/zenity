package zenity

import (
	"os"

	"github.com/ncruces/zenity/internal/zenutil"
)

func entry(text string, opts options) (string, error) {
	var data zenutil.Dialog
	data.Text = text
	data.Operation = "displayDialog"
	data.Options.Title = opts.title
	data.Options.Answer = &opts.entryText
	data.Options.Hidden = opts.hideText
	data.Options.Timeout = zenutil.Timeout
	if opts.attach != nil {
		data.Application = opts.attach
	}
	if i, ok := opts.windowIcon.(string); ok {
		data.WindowIcon = i
	}
	switch i := opts.icon.(type) {
	case string:
		_, err := os.Stat(i)
		if err != nil {
			return "", err
		}
		data.IconPath = i
	case DialogIcon:
		data.Options.Icon = i.String()
	}
	data.SetButtons(getButtons(true, true, opts))

	out, err := zenutil.Run(opts.ctx, "dialog", data)
	return strResult(opts, out, err)
}
