// +build !windows,!darwin

package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func notify(text string, options []Option) error {
	opts := applyOptions(options)

	args := []string{"--notification"}

	if text != "" {
		args = append(args, "--text", text, "--no-markup")
	}
	if opts.title != nil {
		args = append(args, "--title", *opts.title)
	}
	switch opts.icon {
	case ErrorIcon:
		args = append(args, "--window-icon=error")
	case WarningIcon:
		args = append(args, "--window-icon=warning")
	case InfoIcon:
		args = append(args, "--window-icon=info")
	case QuestionIcon:
		args = append(args, "--window-icon=question")
	}

	_, err := zenutil.Run(opts.ctx, args)
	if err != nil {
		return err
	}
	return nil
}
