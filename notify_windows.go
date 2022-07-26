package zenity

import (
	"math/rand"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
	"github.com/ncruces/zenity/internal/zenutil"
)

func notify(text string, opts options) error {
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return opts.ctx.Err()
	}

	var args win.NOTIFYICONDATA
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.ID = rand.Uint32()
	args.Flags = win.NIF_INFO
	args.State = win.NIS_HIDDEN

	info := syscall.StringToUTF16(text)
	copy(args.Info[:len(args.Info)-1], info)

	if opts.title != nil {
		title := syscall.StringToUTF16(*opts.title)
		copy(args.InfoTitle[:len(args.InfoTitle)-1], title)
	}

	switch opts.icon {
	case InfoIcon, QuestionIcon:
		args.InfoFlags |= win.NIIF_INFO
	case WarningIcon:
		args.InfoFlags |= win.NIIF_WARNING
	case ErrorIcon:
		args.InfoFlags |= win.NIIF_ERROR
	default:
		icon, _ := getIcon(opts.icon)
		if icon.handle != 0 {
			defer icon.delete()
			args.Icon = icon.handle
			args.Flags |= win.NIF_ICON
			args.InfoFlags |= win.NIIF_USER
		}
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if !win.ShellNotifyIcon(win.NIM_ADD, &args) {
		return wtsMessage(text, opts)
	}

	major, minor, _ := win.RtlGetNtVersionNumbers()
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

	win.ShellNotifyIcon(win.NIM_DELETE, &args)
	return nil
}

func wtsMessage(text string, opts options) error {
	var flags uint32

	switch opts.icon {
	case ErrorIcon:
		flags |= win.MB_ICONERROR
	case QuestionIcon:
		flags |= win.MB_ICONQUESTION
	case WarningIcon:
		flags |= win.MB_ICONWARNING
	case InfoIcon:
		flags |= win.MB_ICONINFORMATION
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
	return win.WTSSendMessage(
		win.WTS_CURRENT_SERVER_HANDLE, win.WTS_CURRENT_SESSION,
		&ptitle[0], 2*len(ptitle), &ptext[0], 2*len(ptext),
		flags, timeout, &res, false)
}
