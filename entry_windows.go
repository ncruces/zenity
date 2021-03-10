// This file was imported from: github.com/gen2brain/dlgs
// Copyright (c) 2017, Milan Nikolic <gen2brain>
// Licensed under the BSD 2-Clause "Simplified" License.

package zenity

import (
	"syscall"
	"unsafe"
)

func entry(text string, opts options) (string, bool, error) {
	var title string
	if opts.title != nil {
		title = *opts.title
	}
	return editBox(title, text, opts.entryText, "ClassEntry", false)
}

func password(opts options) (string, string, bool, error) {
	var title string
	if opts.title != nil {
		title = *opts.title
	}
	pass, ok, err := editBox(title, "Type your password", "", "ClassPassword", true)
	return "", pass, ok, err
}

var (
	gdi32 = syscall.NewLazyDLL("gdi32.dll")

	createWindowExW       = user32.NewProc("CreateWindowExW")
	defWindowProcW        = user32.NewProc("DefWindowProcW")
	destroyWindowW        = user32.NewProc("DestroyWindow")
	dispatchMessage       = user32.NewProc("DispatchMessageW")
	postQuitMessageW      = user32.NewProc("PostQuitMessage")
	registerClassExW      = user32.NewProc("RegisterClassExW")
	unregisterClassW      = user32.NewProc("UnregisterClassW")
	translateMessage      = user32.NewProc("TranslateMessage")
	getWindowTextLengthW  = user32.NewProc("GetWindowTextLengthW")
	getWindowTextW        = user32.NewProc("GetWindowTextW")
	getWindowLongW        = user32.NewProc("GetWindowLongW")
	setWindowLongW        = user32.NewProc("SetWindowLongW")
	getWindowRectW        = user32.NewProc("GetWindowRect")
	setWindowPosW         = user32.NewProc("SetWindowPos")
	showWindowW           = user32.NewProc("ShowWindow")
	updateWindowW         = user32.NewProc("UpdateWindow")
	isDialogMessage       = user32.NewProc("IsDialogMessageW")
	getSystemMetricsW     = user32.NewProc("GetSystemMetrics")
	systemParametersInfoW = user32.NewProc("SystemParametersInfoW")

	createFontIndirectW = gdi32.NewProc("CreateFontIndirectW")

	getModuleHandleW = kernel32.NewProc("GetModuleHandleW")
)

const (
	swShow       = 5
	swShowNormal = 1
	swUseDefault = 0x80000000

	swpNoZOrder = 0x0004
	swpNoSize   = 0x0001

	smCxScreen = 0
	smCyScreen = 1

	wsThickFrame       = 0x00040000
	wsSysMenu          = 0x00080000
	wsBorder           = 0x00800000
	wsCaption          = 0x00C00000
	wsChild            = 0x40000000
	wsVisible          = 0x10000000
	wsMaximizeBox      = 0x00010000
	wsMinimizeBox      = 0x00020000
	wsTabStop          = 0x00010000
	wsGroup            = 0x00020000
	wsOverlappedWindow = 0x00CF0000
	wsExClientEdge     = 0x00000200

	wmCreate     = 0x0001
	wmDestroy    = 0x0002
	wmClose      = 0x0010
	wmCommand    = 0x0111
	wmSetFont    = 0x0030
	wmKeydown    = 0x0100
	wmInitDialog = 0x0110

	esPassword    = 0x0020
	esAutoVScroll = 0x0040
	esAutoHScroll = 0x0080

	dtmFirst         = 0x1000
	dtmGetSystemTime = dtmFirst + 1
	dtmSetSystemTime = dtmFirst + 2

	vkEscape               = 0x1B
	enUpdate               = 0x0400
	bsPushButton           = 0
	colorWindow            = 5
	spiGetNonClientMetrics = 0x0029
	gwlStyle               = -16
	maxPath                = 260
)

// wndClassExW https://msdn.microsoft.com/en-us/library/windows/desktop/ms633577.aspx
type wndClassExW struct {
	size       uint32
	style      uint32
	wndProc    uintptr
	clsExtra   int32
	wndExtra   int32
	instance   syscall.Handle
	icon       syscall.Handle
	cursor     syscall.Handle
	background syscall.Handle
	menuName   *uint16
	className  *uint16
	iconSm     syscall.Handle
}

// https://docs.microsoft.com/en-ie/windows/win32/api/winuser/ns-winuser-msg
type _MSG struct {
	Owner   syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      _POINT
}

// nonClientMetricsW https://msdn.microsoft.com/en-us/library/windows/desktop/ff729175.aspx
type nonClientMetricsW struct {
	cbSize           uint32
	iBorderWidth     int32
	iScrollWidth     int32
	iScrollHeight    int32
	iCaptionWidth    int32
	iCaptionHeight   int32
	lfCaptionFont    logfontW
	iSmCaptionWidth  int32
	iSmCaptionHeight int32
	lfSmCaptionFont  logfontW
	iMenuWidth       int32
	iMenuHeight      int32
	lfMenuFont       logfontW
	lfStatusFont     logfontW
	lfMessageFont    logfontW
}

// logfontW https://msdn.microsoft.com/en-us/library/windows/desktop/dd145037.aspx
type logfontW struct {
	lfHeight         int32
	lfWidth          int32
	lfEscapement     int32
	lfOrientation    int32
	lfWeight         int32
	lfItalic         byte
	lfUnderline      byte
	lfStrikeOut      byte
	lfCharSet        byte
	lfOutPrecision   byte
	lfClipPrecision  byte
	lfQuality        byte
	lfPitchAndFamily byte
	lfFaceName       [32]uint16
}

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-point
type _POINT struct {
	x, y int32
}

// https://docs.microsoft.com/en-us/windows/win32/api/windef/ns-windef-rect
type _RECT struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

func getModuleHandle() (syscall.Handle, error) {
	ret, _, err := getModuleHandleW.Call(uintptr(0))
	if ret == 0 {
		return 0, err
	}

	return syscall.Handle(ret), nil
}

func createWindow(exStyle uint64, className, windowName string, style uint64, x, y, width, height int64,
	parent, menu, instance syscall.Handle) (syscall.Handle, error) {
	ret, _, err := createWindowExW.Call(uintptr(exStyle), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(className))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName))), uintptr(style), uintptr(x), uintptr(y),
		uintptr(width), uintptr(height), uintptr(parent), uintptr(menu), uintptr(instance), uintptr(0))

	if ret == 0 {
		return 0, err
	}

	return syscall.Handle(ret), nil
}

func destroyWindow(hwnd syscall.Handle) error {
	ret, _, err := destroyWindowW.Call(uintptr(hwnd))
	if ret == 0 {
		return err
	}

	return nil
}

func defWindowProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	ret, _, _ := defWindowProcW.Call(uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam))
	return uintptr(ret)
}

func registerClassEx(wcx *wndClassExW) (uint16, error) {
	ret, _, err := registerClassExW.Call(uintptr(unsafe.Pointer(wcx)))

	if ret == 0 {
		return 0, err
	}

	return uint16(ret), nil
}

func unregisterClass(className string, instance syscall.Handle) bool {
	ret, _, _ := unregisterClassW.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(className))), uintptr(instance))

	return ret != 0
}

func postQuitMessage(exitCode int32) {
	postQuitMessageW.Call(uintptr(exitCode))
}

func getWindowTextLength(hwnd syscall.Handle) int {
	ret, _, _ := getWindowTextLengthW.Call(uintptr(hwnd))
	return int(ret)
}

func getWindowText(hwnd syscall.Handle) string {
	textLen := getWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	getWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func systemParametersInfo(uiAction, uiParam uint32, pvParam unsafe.Pointer, fWinIni uint32) bool {
	ret, _, _ := systemParametersInfoW.Call(uintptr(uiAction), uintptr(uiParam), uintptr(pvParam), uintptr(fWinIni), 0, 0)
	return int32(ret) != 0
}

func createFontIndirect(lplf *logfontW) uintptr {
	ret, _, _ := createFontIndirectW.Call(uintptr(unsafe.Pointer(lplf)), 0, 0)
	return uintptr(ret)
}

func getWindowLong(hwnd syscall.Handle, index int32) int32 {
	ret, _, _ := getWindowLongW.Call(uintptr(hwnd), uintptr(index), 0)
	return int32(ret)
}

func setWindowLong(hwnd syscall.Handle, index, value int32) int32 {
	ret, _, _ := setWindowLongW.Call(uintptr(hwnd), uintptr(index), uintptr(value))
	return int32(ret)
}

func getWindowRect(hwnd syscall.Handle, rect *_RECT) bool {
	ret, _, _ := getWindowRectW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(rect)), 0)
	return ret != 0
}

func setWindowPos(hwnd, hwndInsertAfter syscall.Handle, x, y, width, height int32, flags uint32) bool {
	ret, _, _ := setWindowPosW.Call(uintptr(hwnd), uintptr(hwndInsertAfter),
		uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(flags), 0, 0)
	return ret != 0
}

func showWindow(hwnd syscall.Handle, nCmdShow int32) bool {
	ret, _, _ := showWindowW.Call(uintptr(hwnd), uintptr(nCmdShow), 0)
	return ret != 0
}

func updateWindow(hwnd syscall.Handle) bool {
	ret, _, _ := updateWindowW.Call(uintptr(hwnd), 0, 0)
	return ret != 0
}

func getSystemMetrics(nindex int32) int32 {
	ret, _, _ := getSystemMetricsW.Call(uintptr(nindex), 0, 0)
	return int32(ret)
}

func centerWindow(hwnd syscall.Handle) {
	var rc _RECT
	getWindowRect(hwnd, &rc)
	xPos := (getSystemMetrics(smCxScreen) - (rc.right - rc.left)) / 2
	yPos := (getSystemMetrics(smCyScreen) - (rc.bottom - rc.top)) / 2
	setWindowPos(hwnd, 0, xPos, yPos, 0, 0, swpNoZOrder|swpNoSize)
}

func getMessageFont() uintptr {
	var metrics nonClientMetricsW
	metrics.cbSize = uint32(unsafe.Sizeof(metrics))
	systemParametersInfo(spiGetNonClientMetrics, uint32(unsafe.Sizeof(metrics)), unsafe.Pointer(&metrics), 0)
	return createFontIndirect(&metrics.lfMessageFont)
}

func registerClass(className string, instance syscall.Handle, fn interface{}) error {
	var wcx wndClassExW
	wcx.size = uint32(unsafe.Sizeof(wcx))
	wcx.wndProc = syscall.NewCallback(fn)
	wcx.instance = instance
	wcx.background = colorWindow + 1
	wcx.className = syscall.StringToUTF16Ptr(className)

	_, err := registerClassEx(&wcx)
	return err
}

// https://docs.microsoft.com/en-us/windows/win32/winmsg/using-messages-and-message-queues
func messageLoop(hwnd uintptr) error {
	for {
		var msg _MSG
		ret, _, err := getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(ret) == -1 {
			return err
		}
		if ret == 0 {
			return nil
		}

		ret, _, _ = isDialogMessage.Call(hwnd, uintptr(unsafe.Pointer(&msg)), 0)
		if ret == 0 {
			translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}
}

// editBox displays textedit/inputbox dialog.
func editBox(title, text, defaultText, className string, password bool) (string, bool, error) {
	var out string
	notCancledOrClosed := true

	var hwndEdit syscall.Handle

	instance, err := getModuleHandle()
	if err != nil {
		return out, false, err
	}

	fn := func(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
		switch msg {
		case wmClose:
			notCancledOrClosed = false
			destroyWindow(hwnd)
		case wmDestroy:
			postQuitMessage(0)
		case wmKeydown:
			if wparam == vkEscape {
				notCancledOrClosed = false
				destroyWindow(hwnd)
			}
		case wmCommand:
			if wparam == 100 {
				out = getWindowText(hwndEdit)
				destroyWindow(hwnd)
			} else if wparam == 110 {
				notCancledOrClosed = false
				destroyWindow(hwnd)
			}
		default:
			ret := defWindowProc(hwnd, msg, wparam, lparam)
			return ret
		}

		return 0
	}

	err = registerClass(className, instance, fn)
	if err != nil {
		return out, false, err
	}
	defer unregisterClass(className, instance)

	hwnd, _ := createWindow(0, className, title, wsOverlappedWindow, swUseDefault, swUseDefault, 235, 140, 0, 0, instance)
	hwndText, _ := createWindow(0, "STATIC", text, wsChild|wsVisible, 10, 10, 200, 16, hwnd, 0, instance)

	flags := wsBorder | wsChild | wsVisible | wsGroup | wsTabStop | esAutoHScroll
	if password {
		flags |= esPassword
	}
	hwndEdit, _ = createWindow(wsExClientEdge, "EDIT", defaultText, uint64(flags), 10, 30, 200, 24, hwnd, 0, instance)

	hwndOK, _ := createWindow(wsExClientEdge, "BUTTON", "OK", wsChild|wsVisible|bsPushButton|wsGroup|wsTabStop, 10, 65, 90, 24, hwnd, 100, instance)
	hwndCancel, _ := createWindow(wsExClientEdge, "BUTTON", "Cancel", wsChild|wsVisible|bsPushButton|wsGroup|wsTabStop, 120, 65, 90, 24, hwnd, 110, instance)

	setWindowLong(hwnd, gwlStyle, getWindowLong(hwnd, gwlStyle)^wsMinimizeBox)
	setWindowLong(hwnd, gwlStyle, getWindowLong(hwnd, gwlStyle)^wsMaximizeBox)

	font := getMessageFont()
	sendMessage.Call(uintptr(hwndText), wmSetFont, font, 0)
	sendMessage.Call(uintptr(hwndEdit), wmSetFont, font, 0)
	sendMessage.Call(uintptr(hwndOK), wmSetFont, font, 0)
	sendMessage.Call(uintptr(hwndCancel), wmSetFont, font, 0)

	centerWindow(hwnd)

	showWindow(hwnd, swShowNormal)
	updateWindow(hwnd)

	err = messageLoop(uintptr(hwnd))
	if err != nil {
		return out, false, err
	}

	return out, notCancledOrClosed, nil
}
