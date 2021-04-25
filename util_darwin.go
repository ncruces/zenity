package zenity

import "github.com/ncruces/zenity/internal/zenutil"

func getButtons(dialog, okcancel bool, opts options) (btns zenutil.DialogButtons) {
	if !okcancel {
		opts.cancelLabel = nil
		opts.defaultCancel = false
	}

	if opts.okLabel != nil || opts.cancelLabel != nil || opts.extraButton != nil || dialog != okcancel {
		if opts.okLabel == nil {
			opts.okLabel = stringPtr("OK")
		}
		if okcancel {
			if opts.cancelLabel == nil {
				opts.cancelLabel = stringPtr("Cancel")
			}
			if opts.extraButton == nil {
				btns.Buttons = []string{*opts.cancelLabel, *opts.okLabel}
				btns.Default = 2
				btns.Cancel = 1
			} else {
				btns.Buttons = []string{*opts.extraButton, *opts.cancelLabel, *opts.okLabel}
				btns.Default = 3
				btns.Cancel = 2
				btns.Extra = 1
			}
		} else {
			if opts.extraButton == nil {
				btns.Buttons = []string{*opts.okLabel}
				btns.Default = 1
			} else {
				btns.Buttons = []string{*opts.extraButton, *opts.okLabel}
				btns.Default = 2
				btns.Extra = 1
			}
		}
	}

	if opts.defaultCancel {
		if btns.Cancel != 0 {
			btns.Default = btns.Cancel
		} else {
			btns.Default = 1
		}
	}
	return
}

func (i DialogIcon) String() string {
	switch i {
	case ErrorIcon:
		return "stop"
	case WarningIcon:
		return "caution"
	case InfoIcon, QuestionIcon:
		return "note"
	default:
		return ""
	}
}

func (k messageKind) String() string {
	switch k {
	case infoKind:
		return "informational"
	case warningKind:
		return "warning"
	case errorKind:
		return "critical"
	default:
		return ""
	}
}
