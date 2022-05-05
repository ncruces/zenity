//go:build !windows && !darwin

package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func message(kind messageKind, text string, opts options) error {
	args := []string{"--text", text, "--no-markup"}
	switch kind {
	case questionKind:
		args = append(args, "--question")
	case infoKind:
		args = append(args, "--info")
	case warningKind:
		args = append(args, "--warning")
	case errorKind:
		args = append(args, "--error")
	}
	args = appendTitle(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendIcon(args, opts)
	if opts.noWrap {
		args = append(args, "--no-wrap")
	}
	if opts.ellipsize {
		args = append(args, "--ellipsize")
	}
	if opts.defaultCancel {
		args = append(args, "--default-cancel")
	}
	switch opts.icon {
	case ErrorIcon:
		args = append(args, "--icon-name=dialog-error")
	case WarningIcon:
		args = append(args, "--icon-name=dialog-warning")
	case InfoIcon:
		args = append(args, "--icon-name=dialog-information")
	case QuestionIcon:
		args = append(args, "--icon-name=dialog-question")
	case PasswordIcon:
		args = append(args, "--icon-name=dialog-password")
	case NoIcon:
		args = append(args, "--icon-name=")
	}

	out, err := zenutil.Run(opts.ctx, args)
	_, err = strResult(opts, out, err)
	return err
}
