package zenity

import (
	"syscall"
	"unsafe"
)

var (
	user32     = syscall.NewLazyDLL("user32.dll")
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

	if opts.title != "" {
		caption = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(opts.title)))
	}

	n, _, err := messageBox.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		caption, flags)

	if n == 0 {
		return false, err
	} else {
		return n == 1 /* IDOK */ || n == 6 /* IDYES */, nil
	}
}
