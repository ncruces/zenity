package zenity

import (
	"runtime"
	"syscall"
	"unsafe"
)

var (
	shellNotifyIcon = shell32.NewProc("Shell_NotifyIconW")
)

func notify(text string, options []Option) error {
	opts := applyOptions(options)

	var args _NOTIFYICONDATA
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.ID = 0x378eb49c    // random
	args.Flags = 0x00000010 // NIF_INFO
	args.State = 0x00000001 // NIS_HIDDEN

	info := syscall.StringToUTF16(text)
	copy(args.Info[:len(args.Info)-1], info)

	title := syscall.StringToUTF16(opts.title)
	copy(args.InfoTitle[:len(args.InfoTitle)-1], title)

	switch opts.icon {
	case InfoIcon:
		args.InfoFlags |= 0x1 // NIIF_INFO
	case WarningIcon:
		args.InfoFlags |= 0x2 // NIIF_WARNING
	case ErrorIcon:
		args.InfoFlags |= 0x3 // NIIF_ERROR
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	n, _, err := shellNotifyIcon.Call(0 /* NIM_ADD */, uintptr(unsafe.Pointer(&args)))
	if n == 0 {
		if errno, ok := err.(syscall.Errno); ok && errno == 0 {
			_, err = Info(text, Title(opts.title), Icon(opts.icon))
		}
		return err
	}

	shellNotifyIcon.Call(2 /* NIM_DELETE */, uintptr(unsafe.Pointer(&args)))
	return nil
}

type _NOTIFYICONDATA struct {
	StructSize      uint32
	Owner           uintptr
	ID              uint32
	Flags           uint32
	CallbackMessage uint32
	Icon            uintptr
	Tip             [128]uint16
	State           uint32
	StateMask       uint32
	Info            [256]uint16
	Version         uint32
	InfoTitle       [64]uint16
	InfoFlags       uint32
	// GuidItem     [16]byte
	// BalloonIcon  uintptr
}
