package zenity

import (
	"os"

	"github.com/ncruces/zenity/internal/zenutil"
)

func message(kind messageKind, text string, opts options) error {
	var data zenutil.Dialog
	data.Text = text
	data.Options.Timeout = zenutil.Timeout
	if opts.attach != nil {
		data.Application = opts.attach
	}
	if i, ok := opts.windowIcon.(string); ok {
		data.WindowIcon = i
	}

	// dialog is more flexible, alert prettier
	var dialog bool
	if opts.icon != nil { // use if we want to show a specific icon
		dialog = true
	} else if kind == questionKind && opts.cancelLabel == nil { // use for questions with default buttons
		dialog = true
	}

	if dialog {
		data.Operation = "displayDialog"
		data.Options.Title = opts.title
		switch i := opts.icon.(type) {
		case string:
			_, err := os.Stat(i)
			if err != nil {
				return err
			}
			data.IconPath = i
		case DialogIcon:
			data.Options.Icon = i.String()
		}
	} else {
		data.Operation = "displayAlert"
		data.Options.As = kind.String()
		if opts.title != nil {
			data.Text = *opts.title
			data.Options.Message = text
		}
	}

	data.SetButtons(getButtons(dialog, kind == questionKind, opts))

	out, err := zenutil.Run(opts.ctx, "dialog", data)
	_, err = strResult(opts, out, err)
	return err
}
