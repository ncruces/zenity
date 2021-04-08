package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func entry(text string, opts options) (string, bool, error) {
	var data zenutil.Dialog
	data.Text = text
	data.Operation = "displayDialog"
	data.Options.Title = opts.title
	data.Options.Answer = &opts.entryText
	data.Options.Hidden = opts.hideText
	data.Options.Timeout = zenutil.Timeout
	data.Options.Icon = opts.icon.String()
	data.SetButtons(getButtons(true, true, opts))

	out, err := zenutil.Run(opts.ctx, "dialog", data)
	return strResult(opts, out, err)
}

func password(opts options) (string, string, bool, error) {
	opts.hideText = true
	str, ok, err := entry("Password:", opts)
	return "", str, ok, err
}
