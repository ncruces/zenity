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
	if opts.customIcon != "" {
		_, err := os.Stat(opts.customIcon)
		if err != nil {
			return "", "", err
		}
		data.IconPath = opts.customIcon
	} else {
		data.Options.Icon = opts.icon.String()
	}
	data.SetButtons(getButtons(true, true, opts))

	out, err := zenutil.Run(opts.ctx, "pwd", data)
	return pwdResult(zenutil.Separator, opts, out, err)
}
