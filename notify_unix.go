//go:build !windows && !darwin

package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func notify(text string, opts options) error {
	args := []string{"--notification", "--text", text}
	args = appendTitle(args, opts)
	switch opts.icon {
	case ErrorIcon:
		args = append(args, "--window-icon=dialog-error")
	case WarningIcon:
		args = append(args, "--window-icon=dialog-warning")
	case InfoIcon:
		args = append(args, "--window-icon=dialog-information")
	case QuestionIcon:
		args = append(args, "--window-icon=dialog-question")
	case PasswordIcon:
		args = append(args, "--window-icon=dialog-password")
	case NoIcon:
		args = append(args, "--window-icon=dialog")
	}

	_, err := zenutil.Run(opts.ctx, args)
	if err != nil {
		return err
	}
	return nil
}
