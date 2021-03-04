package zenity

import (
	"context"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	messageBox = user32.NewProc("MessageBoxW")
)

func message(kind messageKind, text string, options []Option) (bool, error) {
	opts := applyOptions(options)

	var flags uintptr

	switch {
	case kind == questionKind && opts.extraButton != nil:
		flags |= 0x3 // MB_YESNOCANCEL
	case kind == questionKind || opts.extraButton != nil:
		flags |= 0x1 // MB_OKCANCEL
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

	if opts.ctx != nil || opts.okLabel != nil || opts.cancelLabel != nil || opts.extraButton != nil {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		unhook, err := hookMessageLabels(kind, opts)
		if err != nil {
			return false, err
		}
		defer unhook()
	}

	var title *uint16
	if opts.title != nil {
		title = syscall.StringToUTF16Ptr(*opts.title)
	}

	activate()
	s, _, err := messageBox.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(title)), flags)

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return false, opts.ctx.Err()
	}
	if s == 0 {
		return false, err
	}
	if s == 7 || s == 2 && kind != questionKind { // IDNO
		return false, ErrExtraButton
	}
	if s == 1 || s == 6 { // IDOK, IDYES
		return true, nil
	}
	return false, nil
}

func hookMessageLabels(kind messageKind, opts options) (unhook context.CancelFunc, err error) {
	return hookDialog(opts.ctx, func(wnd uintptr) {
		enumChildWindows.Call(wnd,
			syscall.NewCallback(func(wnd, lparam uintptr) uintptr {
				name := [8]uint16{}
				getClassName.Call(wnd, uintptr(unsafe.Pointer(&name)), uintptr(len(name)))
				if syscall.UTF16ToString(name[:]) == "Button" {
					ctl, _, _ := getDlgCtrlID.Call(wnd)
					var text *string
					switch ctl {
					case 1, 6: // IDOK, IDYES
						text = opts.okLabel
					case 2: // IDCANCEL
						if kind == questionKind {
							text = opts.cancelLabel
						} else if opts.extraButton != nil {
							text = opts.extraButton
						} else {
							text = opts.okLabel
						}
					case 7: // IDNO
						text = opts.extraButton
					}
					if text != nil {
						ptr := syscall.StringToUTF16Ptr(*text)
						setWindowText.Call(wnd, uintptr(unsafe.Pointer(ptr)))
					}
				}
				return 1
			}), 0)
	})
}
