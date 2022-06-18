//go:build windows

package win

// https://docs.microsoft.com/en-us/windows/win32/api/commctrl/ns-commctrl-initcommoncontrolsex
type INITCOMMONCONTROLSEX struct {
	Size uint32
	ICC  uint32
}

//sys InitCommonControlsEx(icc *INITCOMMONCONTROLSEX) (ok bool) = comctl32.InitCommonControlsEx
