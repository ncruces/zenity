package zenity

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	messageBox             = user32.NewProc("MessageBoxW")
	getDlgCtrlID           = user32.NewProc("GetDlgCtrlID")
	createIconFromResource = user32.NewProc("CreateIconFromResource")
	destroyIcon            = user32.NewProc("DestroyIcon")
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

	data, err := os.ReadFile(opts.customIcon)
	if err != nil {
		return err
	}
	if !bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
		return fmt.Errorf("%w: unsupported icon", ErrUnsupported)
	}
	opts.iconData = data

	defer setup()()

	if opts.ctx != nil || opts.okLabel != nil || opts.cancelLabel != nil || opts.extraButton != nil || opts.customIcon != "" {
		unhook, err := hookMessageDialog(kind, opts)
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

func hookMessageDialog(kind messageKind, opts options) (unhook context.CancelFunc, err error) {
	return hookDialog(opts.ctx, func(wnd uintptr) {
		enumChildWindows.Call(wnd,
			syscall.NewCallback(hookMessageDialogCallback),
			uintptr(unsafe.Pointer(&opts)))
	})
}

func hookMessageDialogCallback(wnd uintptr, lparam *options) uintptr {
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

	if ctl == 20 /*IDC_STATIC_OK*/ && lparam.customIcon != "" {
		icon, _, _ := createIconFromResource.Call(
			uintptr(unsafe.Pointer(&lparam.iconData[0])),
			uintptr(len(lparam.iconData)),
			1, 0x00030000)
		if icon != 0 {
			sendMessage.Call(wnd, 0x0170 /*STM_SETICON*/, icon, 0)
			destroyIcon.Call(icon)
		}
	}
	return 1
}
