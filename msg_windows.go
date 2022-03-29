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
		flags |= _MB_YESNOCANCEL
	case kind == questionKind:
		flags |= _MB_OKCANCEL
	case opts.extraButton != nil:
		flags |= _MB_YESNO
	}

	switch opts.icon {
	case ErrorIcon:
		flags |= _MB_ICONERROR
	case QuestionIcon:
		flags |= _MB_ICONQUESTION
	case WarningIcon:
		flags |= _MB_ICONWARNING
	case InfoIcon:
		flags |= _MB_ICONINFORMATION
	case unspecifiedIcon:
		switch kind {
		case errorKind:
			flags |= _MB_ICONERROR
		case questionKind:
			flags |= _MB_ICONQUESTION
		case warningKind:
			flags |= _MB_ICONWARNING
		case infoKind:
			flags |= _MB_ICONINFORMATION
		}
	}

	if kind == questionKind && opts.defaultCancel {
		if opts.extraButton == nil {
			flags |= _MB_DEFBUTTON2
		} else {
			flags |= _MB_DEFBUTTON3
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
	case _IDOK, _IDYES:
		return nil
	case _IDCANCEL:
		return ErrCanceled
	case _IDNO:
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
		case _IDOK, _IDYES:
			text = lparam.okLabel
		case _IDCANCEL:
			text = lparam.cancelLabel
		case _IDNO:
			text = lparam.extraButton
		}
		if text != nil {
			setWindowText.Call(wnd, strptr(*text))
		}
	}
	return 1
}
