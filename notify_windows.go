package zenity

import (
	"math/rand"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/ncruces/zenity/internal/zenutil"
)

var (
	rtlGetNtVersionNumbers = ntdll.NewProc("RtlGetNtVersionNumbers")
	shellNotifyIcon        = shell32.NewProc("Shell_NotifyIconW")
	wtsSendMessage         = wtsapi32.NewProc("WTSSendMessageW")
)

func notify(text string, opts options) error {
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return opts.ctx.Err()
	}

	var args _NOTIFYICONDATA
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.ID = rand.Uint32()
	args.Flags = 0x00000010 // NIF_INFO
	args.State = 0x00000001 // NIS_HIDDEN

	info := syscall.StringToUTF16(text)
	copy(args.Info[:len(args.Info)-1], info)

	if opts.title != nil {
		title := syscall.StringToUTF16(*opts.title)
		copy(args.InfoTitle[:len(args.InfoTitle)-1], title)
	}

	switch opts.icon {
	case InfoIcon, QuestionIcon:
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

	var major, minor, build uint32
	rtlGetNtVersionNumbers.Call(
		uintptr(unsafe.Pointer(&major)),
		uintptr(unsafe.Pointer(&minor)),
		uintptr(unsafe.Pointer(&build)))

	// On Windows 7 (6.1) and lower, wait up to 10 seconds to clean up.
	if major < 6 || major == 6 && minor < 2 {
		if opts.ctx != nil {
			select {
			case <-opts.ctx.Done():
			case <-time.After(10 * time.Second):
			}
		} else {
			time.Sleep(10 * time.Second)
		}
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
	if title == nil {
		title = stringPtr("Notification")
	}

	timeout := zenutil.Timeout
	if timeout == 0 {
		timeout = 10
	}

	ptext := syscall.StringToUTF16(text)
	ptitle := syscall.StringToUTF16(*title)

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

// https://docs.microsoft.com/en-us/windows/win32/api/shellapi/ns-shellapi-notifyicondataw
type _NOTIFYICONDATA struct {
	StructSize      uint32
	Wnd             uintptr
	ID              uint32
	Flags           uint32
	CallbackMessage uint32
	Icon            uintptr
	Tip             [128]uint16 // NOTIFYICONDATAA_V1_SIZE
	State           uint32
	StateMask       uint32
	Info            [256]uint16
	Version         uint32
	InfoTitle       [64]uint16
	InfoFlags       uint32
	// GuidItem     [16]byte       // NOTIFYICONDATAA_V2_SIZE
	// BalloonIcon  syscall.Handle // NOTIFYICONDATAA_V3_SIZE
}
