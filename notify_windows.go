package zenity

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/ncruces/zenity/internal/zenutil"
)

var (
	shellNotifyIcon = shell32.NewProc("Shell_NotifyIconW")
	wtsSendMessage  = wtsapi32.NewProc("WTSSendMessageW")
)

func notify(text string, options []Option) error {
	opts := applyOptions(options)

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return opts.ctx.Err()
	}

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

	s, _, err := shellNotifyIcon.Call(0 /* NIM_ADD */, uintptr(unsafe.Pointer(&args)))
	if s == 0 {
		if errno, ok := err.(syscall.Errno); ok && errno == 0 {
			return wtsMessage(text, opts)
		}
		return err
	}

	shellNotifyIcon.Call(2 /* NIM_DELETE */, uintptr(unsafe.Pointer(&args)))
	return nil
}

func wtsMessage(text string, opts options) error {
	var flags uintptr

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

	title := opts.title
	if title == "" {
		title = "Notification"
	}

	timeout := zenutil.Timeout
	if timeout == 0 {
		timeout = 10
	}

	ptext := syscall.StringToUTF16(text)
	ptitle := syscall.StringToUTF16(title)

	var res uint32
	s, _, err := wtsSendMessage.Call(
		0,          // WTS_CURRENT_SERVER_HANDLE
		0xffffffff, // WTS_CURRENT_SESSION
		uintptr(unsafe.Pointer(&ptitle[0])), uintptr(2*len(ptitle)),
		uintptr(unsafe.Pointer(&ptext[0])), uintptr(2*len(ptext)),
		flags, uintptr(timeout), uintptr(unsafe.Pointer(&res)), 0)

	if s == 0 {
		return err
	}
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
