//go:build windows

package win

import "golang.org/x/sys/windows"

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
)

func MessageBox(hwnd HWND, text *uint16, caption *uint16, boxtype uint32) (ret int32, err error) {
	return windows.MessageBox(hwnd, text, caption, boxtype)
}

//sys GetDlgCtrlID(wnd HWND) (ret int) = user32.GetDlgCtrlID
//sys SendMessage(wnd HWND, msg uint32, wparam uintptr, lparam uintptr) (ret uintptr) = user32.SendMessageW
//sys SetWindowText(wnd HWND, text *uint16) (err error) = user32.SetWindowTextW
//sys EnumChildWindows(parent HWND, enumFunc uintptr, lparam uintptr) = user32.EnumChildWindows
//sys EnumWindows(enumFunc uintptr, lparam uintptr) (err error) = user32.EnumChildWindows
