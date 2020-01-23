package zenity

import (
	"runtime"
	"syscall"
	"unsafe"
)

var (
	messageBox = user32.NewProc("MessageBoxW")
)

func message(kind messageKind, text string, options []Option) (bool, error) {
	opts := optsParse(options)

	var flags, caption uintptr

	switch {
	case kind == questionKind && opts.extra != "":
		flags |= 0x3 // MB_YESNOCANCEL
	case kind == questionKind || opts.extra != "":
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

	if kind == questionKind && opts.defcancel {
		if opts.extra == "" {
			flags |= 0x100 // MB_DEFBUTTON2
		} else {
			flags |= 0x200 // MB_DEFBUTTON3
		}
	}

	if opts.title != "" {
		caption = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(opts.title)))
	}

	if opts.ok != "" || opts.cancel != "" || opts.extra != "" {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		hook, err := hookMessageLabels(kind, opts)
		if hook == 0 {
			return false, err
		}
		defer unhookWindowsHookEx.Call(hook)
	}

	n, _, err := messageBox.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		caption, flags)

	if n == 0 {
		return false, err
	}
	if n == 7 || n == 2 && kind != questionKind { // IDNO
		return false, ErrExtraButton
	}
	if n == 1 || n == 6 { // IDOK, IDYES
		return true, nil
	}
	return false, nil
}

func hookMessageLabels(kind messageKind, opts options) (hook uintptr, err error) {
	tid, _, _ := getCurrentThreadId.Call()
	hook, _, err = setWindowsHookEx.Call(12, // WH_CALLWNDPROCRET
		syscall.NewCallback(func(code int, wparam uintptr, lparam *_CWPRETSTRUCT) uintptr {
			if lparam.Message == 0x0110 { // WM_INITDIALOG
				name := [7]uint16{}
				getClassName.Call(lparam.Wnd, uintptr(unsafe.Pointer(&name)), uintptr(len(name)))
				if syscall.UTF16ToString(name[:]) == "#32770" { // The class for a dialog box
					enumChildWindows.Call(lparam.Wnd,
						syscall.NewCallback(func(wnd, lparam uintptr) uintptr {
							name := [7]uint16{}
							getClassName.Call(wnd, uintptr(unsafe.Pointer(&name)), uintptr(len(name)))
							if syscall.UTF16ToString(name[:]) == "Button" {
								ctl, _, _ := getDlgCtrlID.Call(wnd)
								var text string
								switch ctl {
								case 1, 6: // IDOK, IDYES
									text = opts.ok
								case 2: // IDCANCEL
									if kind == questionKind {
										text = opts.cancel
									} else if opts.extra != "" {
										text = opts.extra
									} else {
										text = opts.ok
									}
								case 7: // IDNO
									text = opts.extra
								}
								if text != "" {
									ptr := syscall.StringToUTF16Ptr(text)
									setWindowText.Call(wnd, uintptr(unsafe.Pointer(ptr)))
								}
							}
							return 1
						}), 0)
				}
			}
			next, _, _ := callNextHookEx.Call(hook, uintptr(code), wparam, uintptr(unsafe.Pointer(lparam)))
			return next
		}), 0, tid)
	return
}
