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
)

func MessageBox(hwnd HWND, text *uint16, caption *uint16, boxtype uint32) (ret int32, err error) {
	return windows.MessageBox(hwnd, text, caption, boxtype)
}
