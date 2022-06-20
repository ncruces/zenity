//go:build windows

package win

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	IDOK     = 1
	IDCANCEL = 2
	IDABORT  = 3
	IDRETRY  = 4
	IDIGNORE = 5
	IDYES    = 6
	IDNO     = 7

	MB_OK                = windows.MB_OK
	MB_OKCANCEL          = windows.MB_OKCANCEL
	MB_ABORTRETRYIGNORE  = windows.MB_ABORTRETRYIGNORE
	MB_YESNOCANCEL       = windows.MB_YESNOCANCEL
	MB_YESNO             = windows.MB_YESNO
	MB_RETRYCANCEL       = windows.MB_RETRYCANCEL
	MB_CANCELTRYCONTINUE = windows.MB_CANCELTRYCONTINUE
	MB_ICONERROR         = windows.MB_ICONERROR
	MB_ICONQUESTION      = windows.MB_ICONQUESTION
	MB_ICONWARNING       = windows.MB_ICONWARNING
	MB_ICONINFORMATION   = windows.MB_ICONINFORMATION
	MB_DEFBUTTON1        = windows.MB_DEFBUTTON1
	MB_DEFBUTTON2        = windows.MB_DEFBUTTON2
	MB_DEFBUTTON3        = windows.MB_DEFBUTTON3

	WM_DESTROY     = 0x0002
	WM_CLOSE       = 0x0010
	WM_SETFONT     = 0x0030
	WM_SETICON     = 0x0080
	WM_NCCREATE    = 0x0081
	WM_NCDESTROY   = 0x0082
	WM_COMMAND     = 0x0111
	WM_SYSCOMMAND  = 0x0112
	WM_DPICHANGED  = 0x02e0
	WM_USER        = 0x0400
	EM_SETSEL      = 0x00b1
	LB_ADDSTRING   = 0x0180
	LB_GETCURSEL   = 0x0188
	LB_GETSELCOUNT = 0x0190
	LB_GETSELITEMS = 0x0191
	MCM_GETCURSEL  = 0x1001
	MCM_SETCURSEL  = 0x1002
	PBM_SETPOS     = WM_USER + 2
	PBM_SETRANGE32 = WM_USER + 6
	PBM_SETMARQUEE = WM_USER + 10
	STM_SETICON    = 0x0170

	USER_DEFAULT_SCREEN_DPI = 96

	DPI_AWARENESS_CONTEXT_UNAWARE              = ^uintptr(1) + 1
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         = ^uintptr(2) + 1
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    = ^uintptr(3) + 1
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = ^uintptr(4) + 1
	DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED    = ^uintptr(5) + 1

	IMAGE_BITMAP = 0
	IMAGE_ICON   = 1
	IMAGE_CURSOR = 2

	LR_DEFAULTCOLOR     = 0x00000000
	LR_MONOCHROME       = 0x00000001
	LR_LOADFROMFILE     = 0x00000010
	LR_LOADTRANSPARENT  = 0x00000020
	LR_DEFAULTSIZE      = 0x00000040
	LR_VGACOLOR         = 0x00000080
	LR_LOADMAP3DCOLORS  = 0x00001000
	LR_CREATEDIBSECTION = 0x00002000
	LR_SHARED           = 0x00008000
)

func MessageBox(hwnd HWND, text *uint16, caption *uint16, boxtype uint32) (ret int32, err error) {
	return windows.MessageBox(hwnd, text, caption, boxtype)
}

func GetWindowThreadProcessId(hwnd HWND, pid *uint32) (tid uint32, err error) {
	return windows.GetWindowThreadProcessId(hwnd, pid)
}

func GetDpiForWindow(wnd HWND) (ret int, err error) {
	if err := procGetDpiForWindow.Find(); err != nil {
		return 0, err
	}
	return getDpiForWindow(wnd), nil
}

func SetThreadDpiAwarenessContext(dpiContext uintptr) (ret uintptr, err error) {
	if err := procSetThreadDpiAwarenessContext.Find(); err != nil {
		return 0, err
	}
	return setThreadDpiAwarenessContext(dpiContext), nil
}

// https://docs.microsoft.com/en-us/windows/win32/winmsg/using-messages-and-message-queues
func MessageLoop(wnd HWND) error {
	getMessage := procGetMessageW.Addr()
	translateMessage := procTranslateMessage.Addr()
	dispatchMessage := procDispatchMessageW.Addr()
	isDialogMessage := procIsDialogMessageW.Addr()

	for {
		var msg MSG
		s, _, err := syscall.Syscall6(getMessage, 4, uintptr(unsafe.Pointer(&msg)), 0, 0, 0, 0, 0)
		if int32(s) == -1 {
			return err
		}
		if s == 0 {
			return nil
		}

		s, _, _ = syscall.Syscall(isDialogMessage, 2, uintptr(wnd), uintptr(unsafe.Pointer(&msg)), 0)
		if s == 0 {
			syscall.Syscall(translateMessage, 1, uintptr(unsafe.Pointer(&msg)), 0, 0)
			syscall.Syscall(dispatchMessage, 1, uintptr(unsafe.Pointer(&msg)), 0, 0)
		}
	}
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type MSG struct {
	Owner   syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
	private uint32
}

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-point
type POINT struct {
	x, y int32
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-wndclassexw
type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   Handle
	Icon       Handle
	Cursor     Handle
	Background Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     Handle
}

//sys DestroyIcon(icon Handle) (err error) = user32.DestroyIcon
//sys DispatchMessage(msg *MSG) (ret uintptr) = user32.DispatchMessageW
//sys EnumChildWindows(parent HWND, enumFunc uintptr, lparam unsafe.Pointer) = user32.EnumChildWindows
//sys EnumWindows(enumFunc uintptr, lparam unsafe.Pointer) (err error) = user32.EnumChildWindows
//sys GetDlgCtrlID(wnd HWND) (ret int) = user32.GetDlgCtrlID
//sys getDpiForWindow(wnd HWND) (ret int) = user32.GetDpiForWindow
//sys GetMessage(msg *MSG, wnd HWND, msgFilterMin uint32, msgFilterMax uint32) (ret uintptr) = user32.GetMessageW
//sys GetWindowDC(wnd HWND) (ret Handle) = user32.GetWindowDC
//sys IsDialogMessage(wnd HWND, msg *MSG) (ok bool) = user32.IsDialogMessageW
//sys LoadIcon(instance Handle, resource uintptr) (ret Handle, err error) = user32.LoadIconW
//sys LoadImage(instance Handle, name *uint16, typ int, cx int, cy int, load int) (ret Handle, err error) = user32.LoadImageW
//sys RegisterClassEx(cls *WNDCLASSEX) (ret uint16, err error) = user32.RegisterClassExW
//sys ReleaseDC(wnd HWND, dc Handle) (ok bool) = user32.ReleaseDC
//sys SendMessage(wnd HWND, msg uint32, wparam uintptr, lparam uintptr) (ret uintptr) = user32.SendMessageW
//sys SetForegroundWindow(wnd HWND) (ok bool) = user32.SetForegroundWindow
//sys setThreadDpiAwarenessContext(dpiContext uintptr) (ret uintptr) = user32.SetThreadDpiAwarenessContext
//sys SetWindowText(wnd HWND, text *uint16) (err error) = user32.SetWindowTextW
//sys TranslateMessage(msg *MSG) (ok bool) = user32.TranslateMessage
//sys UnregisterClass(atom uint16, instance Handle) (err error) = user32.UnregisterClassW
//sys CreateIconFromResource(resBits []byte, resSize int, icon bool, ver uint32) (ret Handle, err error) = user32.CreateIconFromResource
