package zenity

import (
	"context"

	"github.com/ncruces/zenity/internal/win"
)

func message(kind messageKind, text string, opts options) error {
	var flags uint32 = win.MB_SETFOREGROUND

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

	owner, _ := opts.attach.(win.HWND)
	defer setup(owner)()
	unhook, err := hookMessageDialog(opts)
	if err != nil {
		return err
	}
	defer unhook()

	var title *uint16
	if opts.title != nil {
		title = strptr(*opts.title)
	}

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

func hookMessageDialog(opts options) (context.CancelFunc, error) {
	var init func(wnd win.HWND)
	var icon icon
	if opts.okLabel != nil || opts.cancelLabel != nil || opts.extraButton != nil || opts.icon != nil {
		init = func(wnd win.HWND) {
			if opts.okLabel != nil {
				win.SetDlgItemText(wnd, win.IDOK, strptr(quoteAccelerators(*opts.okLabel)))
				win.SetDlgItemText(wnd, win.IDYES, strptr(quoteAccelerators(*opts.okLabel)))
			}
			if opts.cancelLabel != nil {
				win.SetDlgItemText(wnd, win.IDCANCEL, strptr(quoteAccelerators(*opts.cancelLabel)))
			}
			if opts.extraButton != nil {
				win.SetDlgItemText(wnd, win.IDNO, strptr(quoteAccelerators(*opts.extraButton)))
			}

			if icon.handle != 0 {
				ctl, _ := win.GetDlgItem(wnd, win.IDC_STATIC_OK)
				win.SendMessage(ctl, win.STM_SETICON, uintptr(icon.handle), 0)
			}
		}
	}
	icon, err := getIcon(opts.icon)
	if err != nil {
		return nil, err
	}
	unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, init)
	if err != nil {
		icon.delete()
		return nil, err
	}
	return func() {
		icon.delete()
		unhook()
	}, nil
}
