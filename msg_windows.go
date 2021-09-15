package zenity

import (
	"context"
	"syscall"
	"unsafe"
)

var (
	messageBox   = user32.NewProc("MessageBoxW")
	getDlgCtrlID = user32.NewProc("GetDlgCtrlID")
)

func message(kind messageKind, text string, opts options) error {
	var flags uintptr

	switch {
	case kind == questionKind && opts.extraButton != nil:
		flags |= 0x3 // MB_YESNOCANCEL
	case kind == questionKind:
		flags |= 0x1 // MB_OKCANCEL
	case opts.extraButton != nil:
		flags |= 0x4 // MB_YESNO
	}

	switch opts.icon {
	case ErrorIcon:
		flags |= 0x10 // MB_ICONERROR
	case QuestionIcon:
		flags |= 0x20 // MB_ICONQUESTION
	case WarningIcon:
		flags |= 0x30 // MB_ICONWARNING
	case InfoIcon:
		flags |= 0x40 // MB_ICONINFORMATION
	}

	if kind == questionKind && opts.defaultCancel {
		if opts.extraButton == nil {
			flags |= 0x100 // MB_DEFBUTTON2
		} else {
			flags |= 0x200 // MB_DEFBUTTON3
		}
	}

	defer setup()()

	if opts.ctx != nil || opts.okLabel != nil || opts.cancelLabel != nil || opts.extraButton != nil {
		unhook, err := hookMessageLabels(kind, opts)
		if err != nil {
			return err
		}
		defer unhook()
	}

	var title *uint16
	if opts.title != nil {
		title = syscall.StringToUTF16Ptr(*opts.title)
	}

	s, _, err := messageBox.Call(0, strptr(text), uintptr(unsafe.Pointer(title)), flags)

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return opts.ctx.Err()
	}
	switch s {
	case 1, 6: // IDOK, IDYES
		return nil
	case 2: // IDCANCEL
		return ErrCanceled
	case 7: // IDNO
		return ErrExtraButton
	default:
		return err
	}
}

func hookMessageLabels(kind messageKind, opts options) (unhook context.CancelFunc, err error) {
	return hookDialog(opts.ctx, func(wnd uintptr) {
		enumChildWindows.Call(wnd,
			syscall.NewCallback(hookMessageLabelsCallback),
			uintptr(unsafe.Pointer(&opts)))
	})
}

func hookMessageLabelsCallback(wnd uintptr, lparam *options) uintptr {
	var name [8]uint16
	getClassName.Call(wnd, uintptr(unsafe.Pointer(&name)), uintptr(len(name)))
	if syscall.UTF16ToString(name[:]) == "Button" {
		ctl, _, _ := getDlgCtrlID.Call(wnd)
		var text *string
		switch ctl {
		case 1, 6: // IDOK, IDYES
			text = lparam.okLabel
		case 2: // IDCANCEL
			text = lparam.cancelLabel
		case 7: // IDNO
			text = lparam.extraButton
		}
		if text != nil {
			setWindowText.Call(wnd, strptr(*text))
		}
	}
	return 1
}
