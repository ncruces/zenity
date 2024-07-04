//go:build windows

package win

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// Button IDs
	IDOK     = 1
	IDCANCEL = 2
	IDABORT  = 3
	IDRETRY  = 4
	IDIGNORE = 5
	IDYES    = 6
	IDNO     = 7

	// Control IDs
	IDC_STATIC_OK = 20

	// MessageBox types
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
	MB_SETFOREGROUND     = windows.MB_SETFOREGROUND

	// Window messages
	WM_DESTROY     = 0x0002
	WM_CLOSE       = 0x0010
	WM_SETFONT     = 0x0030
	WM_SETICON     = 0x0080
	WM_NCCREATE    = 0x0081
	WM_NCDESTROY   = 0x0082
	WM_INITDIALOG  = 0x0110
	WM_COMMAND     = 0x0111
	WM_SYSCOMMAND  = 0x0112
	WM_DPICHANGED  = 0x02e0
	WM_USER        = 0x0400
	EM_SETSEL      = 0x00b1
	LB_ADDSTRING   = 0x0180
	LB_SETSEL      = 0x0185
	LB_SETCURSEL   = 0x0186
	LB_GETCURSEL   = 0x0188
	LB_GETSELCOUNT = 0x0190
	LB_GETSELITEMS = 0x0191
	MCM_GETCURSEL  = 0x1001
	MCM_SETCURSEL  = 0x1002
	PBM_SETPOS     = WM_USER + 2
	PBM_SETRANGE32 = WM_USER + 6
	PBM_SETMARQUEE = WM_USER + 10
	STM_SETICON    = 0x0170

	// CreateWindow
	CW_USEDEFAULT = -0x80000000

	// Window classes
	PROGRESS_CLASS = "msctls_progress32"
	MONTHCAL_CLASS = "SysMonthCal32"

	// Window styles
	WS_OVERLAPPED       = 0x00000000
	WS_TILED            = 0x00000000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_TABSTOP          = 0x00010000
	WS_GROUP            = 0x00020000
	WS_MINIMIZEBOX      = 0x00020000
	WS_SIZEBOX          = 0x00040000
	WS_THICKFRAME       = 0x00040000
	WS_SYSMENU          = 0x00080000
	WS_HSCROLL          = 0x00100000
	WS_VSCROLL          = 0x00200000
	WS_DLGFRAME         = 0x00400000
	WS_BORDER           = 0x00800000
	WS_CAPTION          = 0x00c00000
	WS_MAXIMIZE         = 0x01000000
	WS_CLIPCHILDREN     = 0x02000000
	WS_CLIPSIBLINGS     = 0x04000000
	WS_DISABLED         = 0x08000000
	WS_VISIBLE          = 0x10000000
	WS_ICONIC           = 0x20000000
	WS_MINIMIZE         = 0x20000000
	WS_CHILD            = 0x40000000
	WS_CHILDWINDOW      = 0x40000000
	WS_POPUP            = 0x80000000
	WS_POPUPWINDOW      = WS_POPUP | WS_BORDER | WS_SYSMENU
	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_TILEDWINDOW      = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX

	// Extended window styles
	WS_EX_LEFT                = 0x00000000
	WS_EX_LTRREADING          = 0x00000000
	WS_EX_RIGHTSCROLLBAR      = 0x00000000
	WS_EX_DLGMODALFRAME       = 0x00000001
	WS_EX_NOPARENTNOTIFY      = 0x00000004
	WS_EX_TOPMOST             = 0x00000008
	WS_EX_ACCEPTFILES         = 0x00000010
	WS_EX_TRANSPARENT         = 0x00000020
	WS_EX_MDICHILD            = 0x00000040
	WS_EX_TOOLWINDOW          = 0x00000080
	WS_EX_WINDOWEDGE          = 0x00000100
	WS_EX_CLIENTEDGE          = 0x00000200
	WS_EX_CONTEXTHELP         = 0x00000400
	WS_EX_RIGHT               = 0x00001000
	WS_EX_RTLREADING          = 0x00002000
	WS_EX_LEFTSCROLLBAR       = 0x00004000
	WS_EX_CONTROLPARENT       = 0x00010000
	WS_EX_STATICEDGE          = 0x00020000
	WS_EX_APPWINDOW           = 0x00040000
	WS_EX_LAYERED             = 0x00080000
	WS_EX_NOINHERITLAYOUT     = 0x00100000
	WS_EX_NOREDIRECTIONBITMAP = 0x00200000
	WS_EX_LAYOUTRTL           = 0x00400000
	WS_EX_COMPOSITED          = 0x02000000
	WS_EX_NOACTIVATE          = 0x08000000
	WS_EX_OVERLAPPEDWINDOW    = WS_EX_WINDOWEDGE | WS_EX_CLIENTEDGE
	WS_EX_PALETTEWINDOW       = WS_EX_WINDOWEDGE | WS_EX_TOOLWINDOW | WS_EX_TOPMOST

	// Button styles
	BS_DEFPUSHBUTTON = 0x0001

	// Static control styles
	SS_NOPREFIX     = 0x0080
	SS_EDITCONTROL  = 0x2000
	SS_WORDELLIPSIS = 0xc000

	// Edic control styles
	ES_PASSWORD    = 0x0020
	ES_AUTOHSCROLL = 0x0080

	// List box control styles
	LBS_NOTIFY      = 0x0001
	LBS_EXTENDEDSEL = 0x0800

	// Month calendar control styles
	MCS_NOTODAY = 0x0010

	// Progress bar control styles
	PBS_SMOOTH  = 0x0001
	PBS_MARQUEE = 0x0008

	// ShowWindow command
	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_SHOWMAXIMIZED   = 3
	SW_MAXIMIZE        = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11

	// SetWindowPos flags
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	SWP_NOREDRAW       = 0x0008
	SWP_NOACTIVATE     = 0x0010
	SWP_DRAWFRAME      = 0x0020
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOREPOSITION   = 0x0200
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOSENDCHANGING = 0x0400
	SWP_DEFERERASE     = 0x2000
	SWP_ASYNCWINDOWPOS = 0x4000

	// Get/SetWindowLong ids
	GWL_WNDPROC   = -4
	GWL_HINSTANCE = -6
	GWL_ID        = -12
	GWL_STYLE     = -16
	GWL_EXSTYLE   = -20
	GWL_USERDATA  = -21

	// Get/SetSystemMetrics ids
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1

	// SystemParametersInfo ids
	SPI_GETNONCLIENTMETRICS = 0x29

	// WM_SYSCOMMAND wParam
	SC_SIZE         = 0xf000
	SC_MOVE         = 0xf010
	SC_MINIMIZE     = 0xf020
	SC_MAXIMIZE     = 0xf030
	SC_NEXTWINDOW   = 0xf040
	SC_PREVWINDOW   = 0xf050
	SC_CLOSE        = 0xf060
	SC_VSCROLL      = 0xf070
	SC_HSCROLL      = 0xf080
	SC_MOUSEMENU    = 0xf090
	SC_KEYMENU      = 0xf100
	SC_RESTORE      = 0xf120
	SC_TASKLIST     = 0xf130
	SC_SCREENSAVE   = 0xf140
	SC_HOTKEY       = 0xf150
	SC_DEFAULT      = 0xf160
	SC_MONITORPOWER = 0xf170
	SC_CONTEXTHELP  = 0xf180

	// SetWindowsHookEx types
	WH_CALLWNDPROCRET = 12

	// System colors
	COLOR_WINDOW = 5

	// DPI awareness
	USER_DEFAULT_SCREEN_DPI                    = 96
	DPI_AWARENESS_CONTEXT_UNAWARE              = ^uintptr(1) + 1
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         = ^uintptr(2) + 1
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    = ^uintptr(3) + 1
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = ^uintptr(4) + 1
	DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED    = ^uintptr(5) + 1

	// LoadIcon resources
	IDI_APPLICATION = 32512
	IDI_ERROR       = 32513
	IDI_HAND        = 32513
	IDI_QUESTION    = 32514
	IDI_WARNING     = 32515
	IDI_EXCLAMATION = 32515
	IDI_ASTERISK    = 32516
	IDI_INFORMATION = 32516
	IDI_WINLOGO     = 32517
	IDI_SHIELD      = 32518

	// LoadResource (image/icon) flags
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

func GetWindowText(wnd HWND) string {
	len, _ := getWindowTextLength(wnd)
	if len == 0 {
		return ""
	}
	buf := make([]uint16, len+1)
	getWindowText(wnd, &buf[0], len+1)
	return windows.UTF16ToString(buf)
}

func SendMessagePointer(wnd HWND, msg uint32, wparam uintptr, lparam unsafe.Pointer) (ret uintptr) {
	r0, _, _ := syscall.SyscallN(procSendMessageW.Addr(), uintptr(wnd), uintptr(msg), uintptr(wparam), uintptr(lparam))
	ret = uintptr(r0)
	return
}

// https://docs.microsoft.com/en-us/windows/win32/winmsg/using-messages-and-message-queues
func MessageLoop(wnd HWND) error {
	msg, err := GlobalAlloc(0, unsafe.Sizeof(MSG{}))
	if err != nil {
		return err
	}
	defer GlobalFree(msg)

	getMessage := procGetMessageW.Addr()
	translateMessage := procTranslateMessage.Addr()
	dispatchMessage := procDispatchMessageW.Addr()
	isDialogMessage := procIsDialogMessageW.Addr()

	for {
		s, _, err := syscall.SyscallN(getMessage, uintptr(msg), 0, 0, 0)
		if int32(s) == -1 {
			return err
		}
		if s == 0 {
			return nil
		}

		s, _, _ = syscall.SyscallN(isDialogMessage, uintptr(wnd), uintptr(msg))
		if s == 0 {
			syscall.SyscallN(translateMessage, uintptr(msg))
			syscall.SyscallN(dispatchMessage, uintptr(msg))
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
	X, Y int32
}

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-rect
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
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

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-nonclientmetricsw
type NONCLIENTMETRICS struct {
	Size            uint32
	BorderWidth     int32
	ScrollWidth     int32
	ScrollHeight    int32
	CaptionWidth    int32
	CaptionHeight   int32
	CaptionFont     LOGFONT
	SmCaptionWidth  int32
	SmCaptionHeight int32
	SmCaptionFont   LOGFONT
	MenuWidth       int32
	MenuHeight      int32
	MenuFont        LOGFONT
	StatusFont      LOGFONT
	MessageFont     LOGFONT
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-cwpretstruct
type CWPRETSTRUCT struct {
	Result  uintptr
	LParam  uintptr
	WParam  uintptr
	Message uint32
	Wnd     HWND
}

//sys CallNextHookEx(hk Handle, code int32, wparam uintptr, lparam unsafe.Pointer) (ret uintptr) = user32.CallNextHookEx
//sys CreateIconFromResourceEx(resBits []byte, icon bool, ver uint32, cx int, cy int, flags int) (ret Handle, err error) = user32.CreateIconFromResourceEx
//sys CreateWindowEx(exStyle uint32, className *uint16, windowName *uint16, style uint32, x int, y int, width int, height int, parent HWND, menu Handle, instance Handle, param unsafe.Pointer) (ret HWND, err error) = user32.CreateWindowExW
//sys DefWindowProc(wnd HWND, msg uint32, wparam uintptr, lparam unsafe.Pointer) (ret uintptr) = user32.DefWindowProcW
//sys DestroyIcon(icon Handle) (err error) = user32.DestroyIcon
//sys DestroyWindow(wnd HWND) (err error) = user32.DestroyWindow
//sys DispatchMessage(msg *MSG) (ret uintptr) = user32.DispatchMessageW
//sys EnableWindow(wnd HWND, enable bool) (ok bool) = user32.EnableWindow
//sys EnumWindows(enumFunc uintptr, lparam unsafe.Pointer) (err error) = user32.EnumChildWindows
//sys GetDlgItem(dlg HWND, dlgItemID int) (ret HWND, err error) = user32.GetDlgItem
//sys GetDpiForWindow(wnd HWND) (ret int, err error) [false] = user32.GetDpiForWindow?
//sys GetMessage(msg *MSG, wnd HWND, msgFilterMin uint32, msgFilterMax uint32) (ret uintptr, err error) [int32(failretval)==-1] = user32.GetMessageW
//sys GetSystemMetrics(index int) (ret int) = user32.GetSystemMetrics
//sys GetWindowDC(wnd HWND) (ret Handle) = user32.GetWindowDC
//sys GetWindowRect(wnd HWND, cmdShow *RECT) (err error) = user32.GetWindowRect
//sys getWindowText(wnd HWND, str *uint16, maxCount int) (ret int, err error) = user32.GetWindowTextW
//sys getWindowTextLength(wnd HWND) (ret int, err error) = user32.GetWindowTextLengthW
//sys IsDialogMessage(wnd HWND, msg *MSG) (ok bool) = user32.IsDialogMessageW
//sys LoadIcon(instance Handle, resource uintptr) (ret Handle, err error) = user32.LoadIconW
//sys PostQuitMessage(exitCode int) = user32.PostQuitMessage
//sys RegisterClassEx(cls *WNDCLASSEX) (err error) = user32.RegisterClassExW
//sys ReleaseDC(wnd HWND, dc Handle) (ok bool) = user32.ReleaseDC
//sys SendMessage(wnd HWND, msg uint32, wparam uintptr, lparam uintptr) (ret uintptr) = user32.SendMessageW
//sys SetDlgItemText(dlg HWND, dlgItemID int, str *uint16) (err error) = user32.SetDlgItemTextW
//sys SetFocus(wnd HWND) (ret HWND, err error) = user32.SetFocus
//sys SetForegroundWindow(wnd HWND) (ok bool) = user32.SetForegroundWindow
//sys SetThreadDpiAwarenessContext(dpiContext uintptr) (ret uintptr, err error) [false] = user32.SetThreadDpiAwarenessContext?
//sys SetWindowLong(wnd HWND, index int, newLong int) (ret int, err error) = user32.SetWindowLongW
//sys SetWindowPos(wnd HWND, wndInsertAfter HWND, x int, y int, cx int, cy int, flags int) (err error) = user32.SetWindowPos
//sys SetWindowsHookEx(idHook int, fn uintptr, mod Handle, threadID uint32) (ret Handle, err error) = user32.SetWindowsHookExW
//sys SetWindowText(wnd HWND, text *uint16) (err error) = user32.SetWindowTextW
//sys ShowWindow(wnd HWND, cmdShow int) (ok bool) = user32.ShowWindow
//sys SystemParametersInfo(action int, uiParam uintptr, pvParam unsafe.Pointer, winIni int) (err error) = user32.SystemParametersInfoW
//sys TranslateMessage(msg *MSG) (ok bool) = user32.TranslateMessage
//sys UnhookWindowsHookEx(hk Handle) (err error) = user32.UnhookWindowsHookEx
//sys UnregisterClass(className *uint16, instance Handle) (err error) = user32.UnregisterClassW
