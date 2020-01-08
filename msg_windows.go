package zenity

import (
	"runtime"
	"syscall"
	"unsafe"
)

var (
	messageBox = user32.NewProc("MessageBoxW")
)

func Error(text string, options ...Option) (bool, error) {
	return message(0, text, options)
}

func Info(text string, options ...Option) (bool, error) {
	return message(1, text, options)
}

func Question(text string, options ...Option) (bool, error) {
	return message(2, text, options)
}

func Warning(text string, options ...Option) (bool, error) {
	return message(3, text, options)
}

func message(typ int, text string, options []Option) (bool, error) {
	opts := optsParse(options)

	var flags, caption uintptr

	switch {
	case typ == 2 && opts.extra != "":
		flags |= 0x3 // MB_YESNOCANCEL
	case typ == 2 || opts.extra != "":
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

	if typ == 2 && opts.defcancel {
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

		hook, err := hookMessageLabels(typ, opts)
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
	if n == 7 || n == 2 && typ != 2 { // IDNO
		return false, ErrExtraButton
	}
	if n == 1 || n == 6 { // IDOK, IDYES
		return true, nil
	}
	return false, nil
}

func hookMessageLabels(typ int, opts options) (hook uintptr, err error) {
	tid, _, _ := getCurrentThreadId.Call()
	hook, _, err = setWindowsHookEx.Call(12, // WH_CALLWNDPROCRET
		syscall.NewCallback(func(code int, wparam, lparam uintptr) uintptr {
			msg := *(*_CWPRETSTRUCT)(unsafe.Pointer(lparam))
			if msg.Message == 0x0110 { // WM_INITDIALOG
				name := [7]byte{}
				n, _, _ := getClassName.Call(msg.HWnd, uintptr(unsafe.Pointer(&name)), uintptr(len(name)))
				if string(name[:n]) == "#32770" {
					enumChildWindows.Call(msg.HWnd,
						syscall.NewCallback(func(hwnd, lparam uintptr) uintptr {
							name := [7]byte{}
							n, _, _ := getClassName.Call(hwnd, uintptr(unsafe.Pointer(&name)), uintptr(len(name)))
							if string(name[:n]) == "Button" {
								ctl, _, _ := getDlgCtrlID.Call(hwnd)
								var text string
								switch ctl {
								case 1, 6: // IDOK, IDYES
									text = opts.ok
								case 2: // IDCANCEL
									if typ == 2 {
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
									setWindowText.Call(hwnd, uintptr(unsafe.Pointer(ptr)))
								}
							}
							return 1
						}), 0)
				}
			}
			next, _, _ := callNextHookEx.Call(hook, uintptr(code), wparam, lparam)
			return next
		}), 0, tid)
	return
}
