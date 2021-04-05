package zenity

import (
	"math"
	"strconv"
	"syscall"
	"unsafe"
)

func entry(text string, opts options) (string, bool, error) {
	var title string
	if opts.title != nil {
		title = *opts.title
	}
	if opts.okLabel == nil {
		opts.okLabel = stringPtr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = stringPtr("Cancel")
	}
	return editBox(title, text, opts)
}

func password(opts options) (string, string, bool, error) {
	opts.hideText = true
	pass, ok, err := entry("Password:", opts)
	return "", pass, ok, err
}

var (
	registerClassEx      = user32.NewProc("RegisterClassExW")
	unregisterClass      = user32.NewProc("UnregisterClassW")
	createWindowEx       = user32.NewProc("CreateWindowExW")
	destroyWindow        = user32.NewProc("DestroyWindow")
	isDialogMessage      = user32.NewProc("IsDialogMessageW")
	translateMessage     = user32.NewProc("TranslateMessage")
	dispatchMessage      = user32.NewProc("DispatchMessageW")
	postQuitMessage      = user32.NewProc("PostQuitMessage")
	defWindowProc        = user32.NewProc("DefWindowProcW")
	getWindowRect        = user32.NewProc("GetWindowRect")
	setWindowPos         = user32.NewProc("SetWindowPos")
	setFocus             = user32.NewProc("SetFocus")
	showWindow           = user32.NewProc("ShowWindow")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
	getSystemMetrics     = user32.NewProc("GetSystemMetrics")
	getWindowDC          = user32.NewProc("GetWindowDC")
	releaseDC            = user32.NewProc("ReleaseDC")
	getDpiForWindow      = user32.NewProc("GetDpiForWindow")

	deleteObject       = gdi32.NewProc("DeleteObject")
	getDeviceCaps      = gdi32.NewProc("GetDeviceCaps")
	createFontIndirect = gdi32.NewProc("CreateFontIndirectW")
)

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-wndclassexw
type _WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type _MSG struct {
	Owner   syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      _POINT
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-nonclientmetricsw
type _NONCLIENTMETRICS struct {
	Size            uint32
	BorderWidth     int32
	ScrollWidth     int32
	ScrollHeight    int32
	CaptionWidth    int32
	CaptionHeight   int32
	CaptionFont     _LOGFONT
	SmCaptionWidth  int32
	SmCaptionHeight int32
	SmCaptionFont   _LOGFONT
	MenuWidth       int32
	MenuHeight      int32
	MenuFont        _LOGFONT
	StatusFont      _LOGFONT
	MessageFont     _LOGFONT
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-logfontw
type _LOGFONT struct {
	Height         int32
	Width          int32
	Escapement     int32
	Orientation    int32
	Weight         int32
	Italic         byte
	Underline      byte
	StrikeOut      byte
	CharSet        byte
	OutPrecision   byte
	ClipPrecision  byte
	Quality        byte
	PitchAndFamily byte
	FaceName       [32]uint16
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

type dpi uintptr

func getDPI(wnd uintptr) dpi {
	var res uintptr

	if wnd != 0 && getDpiForWindow.Find() == nil {
		res, _, _ = getDpiForWindow.Call(wnd)
	} else if dc, _, _ := getWindowDC.Call(wnd); dc != 0 {
		res, _, _ = getDeviceCaps.Call(dc, 90) // LOGPIXELSY
		releaseDC.Call(0, dc)
	}

	if res == 0 {
		return 96 // USER_DEFAULT_SCREEN_DPI
	}
	return dpi(res)
}

func (d dpi) Scale(dim uintptr) uintptr {
	if d == 0 {
		return dim
	}
	return dim * uintptr(d) / 96 // USER_DEFAULT_SCREEN_DPI
}

type font struct {
	handle uintptr
	face   _LOGFONT
}

func getFont() font {
	var metrics _NONCLIENTMETRICS
	metrics.Size = uint32(unsafe.Sizeof(metrics))
	systemParametersInfo.Call(0x29, // SPI_GETNONCLIENTMETRICS
		unsafe.Sizeof(metrics), uintptr(unsafe.Pointer(&metrics)), 0)
	return font{face: metrics.MessageFont}
}

func (f *font) ForDPI(dpi dpi) uintptr {
	if h := -int32(dpi.Scale(12)); f.handle == 0 || f.face.Height != h {
		f.Delete()
		f.face.Height = h
		f.handle, _, _ = createFontIndirect.Call(uintptr(unsafe.Pointer(&f.face)))
	}
	return f.handle
}

func (f *font) Delete() {
	if f.handle != 0 {
		deleteObject.Call(f.handle)
		f.handle = 0
	}
}

func getWindowTextString(wnd uintptr) string {
	len, _, _ := getWindowTextLength.Call(wnd)
	buf := make([]uint16, len+1)
	getWindowText.Call(wnd, uintptr(unsafe.Pointer(&buf[0])), len+1)
	return syscall.UTF16ToString(buf)
}

func centerWindow(wnd uintptr) {
	getMetric := func(i uintptr) int32 {
		ret, _, _ := getSystemMetrics.Call(i)
		return int32(ret)
	}

	var rect _RECT
	getWindowRect.Call(wnd, uintptr(unsafe.Pointer(&rect)))
	x := (getMetric(0 /* SM_CXSCREEN */) - (rect.right - rect.left)) / 2
	y := (getMetric(1 /* SM_CYSCREEN */) - (rect.bottom - rect.top)) / 2
	setWindowPos.Call(wnd, 0, uintptr(x), uintptr(y), 0, 0, 0x5) // SWP_NOZORDER|SWP_NOSIZE
}

func registerClass(instance, proc uintptr) (uintptr, error) {
	name := "WC_" + strconv.FormatUint(uint64(proc), 16)

	var wcx _WNDCLASSEX
	wcx.Size = uint32(unsafe.Sizeof(wcx))
	wcx.WndProc = proc
	wcx.Instance = instance
	wcx.Background = 5 // COLOR_WINDOW
	wcx.ClassName = syscall.StringToUTF16Ptr(name)

	ret, _, err := registerClassEx.Call(uintptr(unsafe.Pointer(&wcx)))
	return ret, err
}

// https://docs.microsoft.com/en-us/windows/win32/winmsg/using-messages-and-message-queues
func messageLoop(wnd uintptr) error {
	getMessage := getMessage.Addr()
	isDialogMessage := isDialogMessage.Addr()
	translateMessage := translateMessage.Addr()
	dispatchMessage := dispatchMessage.Addr()

	for {
		var msg _MSG
		ret, _, err := syscall.Syscall6(getMessage, 4, uintptr(unsafe.Pointer(&msg)), 0, 0, 0, 0, 0)
		if int32(ret) == -1 {
			return err
		}
		if ret == 0 {
			return nil
		}

		ret, _, _ = syscall.Syscall(isDialogMessage, 2, wnd, uintptr(unsafe.Pointer(&msg)), 0)
		if ret == 0 {
			syscall.Syscall(translateMessage, 1, uintptr(unsafe.Pointer(&msg)), 0, 0)
			syscall.Syscall(dispatchMessage, 1, uintptr(unsafe.Pointer(&msg)), 0, 0)
		}
	}
}

func editBox(title, text string, opts options) (out string, ok bool, err error) {
	var wnd, textCtl, editCtl uintptr
	var okBtn, cancelBtn, extraBtn uintptr
	defWindowProc := defWindowProc.Addr()

	font := getFont()
	defer font.Delete()

	layout := func(dpi dpi) {
		hfont := font.ForDPI(dpi)
		sendMessage.Call(textCtl, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(editCtl, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(okBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(cancelBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
		setWindowPos.Call(wnd, 0, 0, 0, dpi.Scale(281), dpi.Scale(140), 0x6)                            // SWP_NOZORDER|SWP_NOMOVE
		setWindowPos.Call(textCtl, 0, dpi.Scale(12), dpi.Scale(10), dpi.Scale(241), dpi.Scale(16), 0x4) // SWP_NOZORDER
		setWindowPos.Call(editCtl, 0, dpi.Scale(12), dpi.Scale(30), dpi.Scale(241), dpi.Scale(24), 0x4) // SWP_NOZORDER
		if extraBtn == 0 {
			setWindowPos.Call(okBtn, 0, dpi.Scale(95), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4)      // SWP_NOZORDER
			setWindowPos.Call(cancelBtn, 0, dpi.Scale(178), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4) // SWP_NOZORDER
		} else {
			sendMessage.Call(extraBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
			setWindowPos.Call(okBtn, 0, dpi.Scale(12), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4)      // SWP_NOZORDER
			setWindowPos.Call(extraBtn, 0, dpi.Scale(95), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4)   // SWP_NOZORDER
			setWindowPos.Call(cancelBtn, 0, dpi.Scale(178), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4) // SWP_NOZORDER
		}
	}

	proc := func(wnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
		switch msg {
		case 0x0002: // WM_DESTROY
			postQuitMessage.Call(0)

		case 0x0010: // WM_CLOSE
			destroyWindow.Call(wnd)

		case 0x0111: // WM_COMMAND
			switch wparam {
			default:
				return 1
			case 1, 6: // IDOK, IDYES
				out = getWindowTextString(editCtl)
				ok = true
			case 2: // IDCANCEL
			case 7: // IDNO
				err = ErrExtraButton
			}
			destroyWindow.Call(wnd)

		case 0x02e0: // WM_DPICHANGED
			layout(dpi(uint32(wparam) >> 16))

		default:
			ret, _, _ := syscall.Syscall6(defWindowProc, 4, wnd, uintptr(msg), wparam, lparam, 0, 0)
			return ret
		}

		return 0
	}

	defer setup()()

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return "", false, err
	}

	cls, err := registerClass(instance, syscall.NewCallback(proc))
	if cls == 0 {
		return "", false, err
	}
	defer unregisterClass.Call(cls, instance)

	wnd, _, _ = createWindowEx.Call(0x10101, // WS_EX_CONTROLPARENT|WS_EX_WINDOWEDGE|WS_EX_DLGMODALFRAME
		cls, stringUintptr(title),
		0x84c80000, // WS_POPUPWINDOW|WS_CLIPSIBLINGS|WS_DLGFRAME
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		281, 140, 0, 0, instance)

	textCtl, _, _ = createWindowEx.Call(0,
		stringUintptr("STATIC"), stringUintptr(text),
		0x5002e080, // WS_CHILD|WS_VISIBLE|WS_GROUP|SS_WORDELLIPSIS|SS_EDITCONTROL|SS_NOPREFIX
		12, 10, 241, 16, wnd, 0, instance)

	var flags uintptr = 0x50030080 // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|ES_AUTOHSCROLL
	if opts.hideText {
		flags |= 0x20 // ES_PASSWORD
	}
	editCtl, _, _ = createWindowEx.Call(0x200, // WS_EX_CLIENTEDGE
		stringUintptr("EDIT"), stringUintptr(opts.entryText),
		flags,
		12, 30, 241, 24, wnd, 0, instance)

	okBtn, _, _ = createWindowEx.Call(0,
		stringUintptr("BUTTON"), stringUintptr(*opts.okLabel),
		0x50030001, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|BS_DEFPUSHBUTTON
		12, 65, 75, 24, wnd, 1 /* IDOK */, instance)
	cancelBtn, _, _ = createWindowEx.Call(0,
		stringUintptr("BUTTON"), stringUintptr(*opts.cancelLabel),
		0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
		12, 65, 75, 24, wnd, 2 /* IDCANCEL */, instance)
	if opts.extraButton != nil {
		extraBtn, _, _ = createWindowEx.Call(0,
			stringUintptr("BUTTON"), stringUintptr(*opts.extraButton),
			0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
			12, 65, 75, 24, wnd, 7 /* IDNO */, instance)
	}

	layout(getDPI(wnd))
	centerWindow(wnd)
	setFocus.Call(editCtl)
	showWindow.Call(wnd, 1 /* SW_SHOWNORMAL */, 0)
	sendMessage.Call(editCtl, 0xb1 /* EM_SETSEL */, 0, math.MaxUint64 /* -1 */)

	err = nil
	if err := messageLoop(wnd); err != nil {
		return "", false, err
	}
	return out, ok, err
}
