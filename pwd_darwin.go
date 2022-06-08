package zenity

import (
	"os"

	"github.com/ncruces/zenity/internal/zenutil"
)

func password(opts options) (string, string, error) {
	if !opts.username {
		opts.entryText = ""
		opts.hideText = true
		str, err := entry("Password:", opts)
		return "", str, err
	}

	var data zenutil.Password
	data.Separator = zenutil.Separator
	data.Options.Title = opts.title
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
			return "", "", err
		}
		data.IconPath = i
	case DialogIcon:
		data.Options.Icon = i.String()
	}
	data.SetButtons(getButtons(true, true, opts))

	out, err := zenutil.Run(opts.ctx, "pwd", data)
	return pwdResult(zenutil.Separator, opts, out, err)
}
