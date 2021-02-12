// +build !windows,!darwin

package zenity

import (
	"fmt"

	"github.com/ncruces/zenity/internal/zenutil"
)

func notify(text string, options []Option) error {
	opts := applyOptions(options)

	args := []string{"--notification"}

	if text != "" {
		args = append(args, "--text", text, "--no-markup")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.width > 0 {
		args = append(args, "--width", fmt.Sprint(opts.width))
	}
	if opts.height > 0 {
		args = append(args, "--height", fmt.Sprint(opts.height))
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
