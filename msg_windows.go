package zenity

import (
	"context"
	"syscall"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

func message(kind messageKind, text string, opts options) error {
	var flags uint32

	switch {
	case kind == questionKind && opts.extraButton != nil:
		flags |= win.MB_YESNOCANCEL
	case kind == questionKind:
		flags |= win.MB_OKCANCEL
	case opts.extraButton != nil:
		flags |= win.MB_YESNO
	default:
		opts.cancelLabel = opts.okLabel
	}

	switch opts.icon {
	case ErrorIcon:
		flags |= win.MB_ICONERROR
	case QuestionIcon:
		flags |= win.MB_ICONQUESTION
	case WarningIcon:
		flags |= win.MB_ICONWARNING
	case InfoIcon:
		flags |= win.MB_ICONINFORMATION
	case NoIcon:
		//
	default:
		switch kind {
		case errorKind:
			flags |= win.MB_ICONERROR
		case questionKind:
			flags |= win.MB_ICONQUESTION
		case warningKind:
			flags |= win.MB_ICONWARNING
		case infoKind:
			flags |= win.MB_ICONINFORMATION
		}
	}

	if kind == questionKind && opts.defaultCancel {
		if opts.extraButton == nil {
			flags |= win.MB_DEFBUTTON2
		} else {
			flags |= win.MB_DEFBUTTON3
		}
	}

	defer setup()()

	if opts.ctx != nil || opts.okLabel != nil || opts.cancelLabel != nil || opts.extraButton != nil || opts.icon != nil {
		unhook, err := hookMessageDialog(opts)
		if err != nil {
			return err
		}
		defer unhook()
	}

	var title *uint16
	if opts.title != nil {
		title = strptr(*opts.title)
	}

	owner, _ := opts.attach.(win.HWND)
	s, err := win.MessageBox(owner, strptr(text), title, flags)

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return opts.ctx.Err()
	}
	switch s {
	case win.IDOK, win.IDYES:
		return nil
	case win.IDCANCEL:
		return ErrCanceled
	case win.IDNO:
		return ErrExtraButton
	default:
		return err
	}
}

func hookMessageDialog(opts options) (unhook context.CancelFunc, err error) {
	// TODO: use GetDlgItem, SetDlgItemText instead of EnumChildWindows.
	return hookDialog(opts.ctx, opts.windowIcon, nil, func(wnd win.HWND) {
		win.EnumChildWindows(wnd, syscall.NewCallback(hookMessageDialogCallback),
			unsafe.Pointer(&opts))
	})
}

func hookMessageDialogCallback(wnd win.HWND, lparam *options) uintptr {
	ctl := win.GetDlgCtrlID(wnd)

	var text *string
	switch ctl {
	case win.IDOK, win.IDYES:
		text = lparam.okLabel
	case win.IDCANCEL:
		text = lparam.cancelLabel
	case win.IDNO:
		text = lparam.extraButton
	}
	if text != nil {
		win.SetWindowText(wnd, strptr(*text))
	}

	if ctl == win.IDC_STATIC_OK {
		icon := getIcon(lparam.icon)
		if icon.handle != 0 {
			defer icon.delete()
			win.SendMessage(wnd, win.STM_SETICON, uintptr(icon.handle), 0)
		}
	}
	return 1
}
