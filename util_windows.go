package zenity

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync/atomic"
	"syscall"
	"unsafe"
)

var (
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")
	wtsapi32 = syscall.NewLazyDLL("wtsapi32.dll")

	commDlgExtendedError = comdlg32.NewProc("CommDlgExtendedError")

	deleteObject       = gdi32.NewProc("DeleteObject")
	getDeviceCaps      = gdi32.NewProc("GetDeviceCaps")
	createFontIndirect = gdi32.NewProc("CreateFontIndirectW")

	getModuleHandle    = kernel32.NewProc("GetModuleHandleW")
	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	getConsoleWindow   = kernel32.NewProc("GetConsoleWindow")

	coInitializeEx   = ole32.NewProc("CoInitializeEx")
	coUninitialize   = ole32.NewProc("CoUninitialize")
	coCreateInstance = ole32.NewProc("CoCreateInstance")
	coTaskMemFree    = ole32.NewProc("CoTaskMemFree")

	getMessage                   = user32.NewProc("GetMessageW")
	sendMessage                  = user32.NewProc("SendMessageW")
	postQuitMessage              = user32.NewProc("PostQuitMessage")
	isDialogMessage              = user32.NewProc("IsDialogMessageW")
	dispatchMessage              = user32.NewProc("DispatchMessageW")
	translateMessage             = user32.NewProc("TranslateMessage")
	getClassName                 = user32.NewProc("GetClassNameW")
	unhookWindowsHookEx          = user32.NewProc("UnhookWindowsHookEx")
	setWindowsHookEx             = user32.NewProc("SetWindowsHookExW")
	callNextHookEx               = user32.NewProc("CallNextHookEx")
	enumWindows                  = user32.NewProc("EnumWindows")
	enumChildWindows             = user32.NewProc("EnumChildWindows")
	setWindowText                = user32.NewProc("SetWindowTextW")
	getWindowText                = user32.NewProc("GetWindowTextW")
	getWindowTextLength          = user32.NewProc("GetWindowTextLengthW")
	setForegroundWindow          = user32.NewProc("SetForegroundWindow")
	getWindowThreadProcessId     = user32.NewProc("GetWindowThreadProcessId")
	setThreadDpiAwarenessContext = user32.NewProc("SetThreadDpiAwarenessContext")
	getDpiForWindow              = user32.NewProc("GetDpiForWindow")
	releaseDC                    = user32.NewProc("ReleaseDC")
	getWindowDC                  = user32.NewProc("GetWindowDC")
	systemParametersInfo         = user32.NewProc("SystemParametersInfoW")
	setWindowPos                 = user32.NewProc("SetWindowPos")
	getWindowRect                = user32.NewProc("GetWindowRect")
	getSystemMetrics             = user32.NewProc("GetSystemMetrics")
	unregisterClass              = user32.NewProc("UnregisterClassW")
	registerClassEx              = user32.NewProc("RegisterClassExW")
	destroyWindow                = user32.NewProc("DestroyWindow")
	createWindowEx               = user32.NewProc("CreateWindowExW")
	showWindow                   = user32.NewProc("ShowWindow")
	setFocus                     = user32.NewProc("SetFocus")
	defWindowProc                = user32.NewProc("DefWindowProcW")
)

func intptr(i int64) uintptr {
	return uintptr(i)
}

func strptr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func setup() context.CancelFunc {
	var hwnd uintptr
	enumWindows.Call(syscall.NewCallback(func(wnd, lparam uintptr) uintptr {
		var pid uintptr
		getWindowThreadProcessId.Call(wnd, uintptr(unsafe.Pointer(&pid)))
		if int(pid) == os.Getpid() {
			hwnd = wnd
			return 0
		}
		return 1
	}), 0)
	if hwnd == 0 {
		hwnd, _, _ = getConsoleWindow.Call()
	}
	if hwnd != 0 {
		setForegroundWindow.Call(hwnd)
	}

	var old uintptr
	runtime.LockOSThread()
	if setThreadDpiAwarenessContext.Find() == nil {
		// try:
		//   DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2
		//   DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE
		//   DPI_AWARENESS_CONTEXT_SYSTEM_AWARE
		for i := -4; i <= -2; i++ {
			restore, _, _ := setThreadDpiAwarenessContext.Call(uintptr(i))
			if restore != 0 {
				break
			}
		}
	}

	return func() {
		if old != 0 {
			setThreadDpiAwarenessContext.Call(old)
		}
		runtime.UnlockOSThread()
	}
}

func commDlgError() error {
	s, _, _ := commDlgExtendedError.Call()
	if s == 0 {
		return nil
	} else {
		return fmt.Errorf("Common Dialog error: %x", s)
	}
}

func hookDialog(ctx context.Context, initDialog func(wnd uintptr)) (unhook context.CancelFunc, err error) {
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var hook, wnd uintptr
	callNextHookEx := callNextHookEx.Addr()
	tid, _, _ := getCurrentThreadId.Call()
	hook, _, err = setWindowsHookEx.Call(12, // WH_CALLWNDPROCRET
		syscall.NewCallback(func(code int32, wparam uintptr, lparam *_CWPRETSTRUCT) uintptr {
			if lparam.Message == 0x0110 { // WM_INITDIALOG
				name := [8]uint16{}
				getClassName.Call(lparam.Wnd, uintptr(unsafe.Pointer(&name)), uintptr(len(name)))
				if syscall.UTF16ToString(name[:]) == "#32770" { // The class for a dialog box
					var close bool

					if ctx != nil && ctx.Err() != nil {
						close = true
					} else {
						atomic.StoreUintptr(&wnd, lparam.Wnd)
					}

					if close {
						sendMessage.Call(lparam.Wnd, 0x0112 /* WM_SYSCOMMAND */, 0xf060 /* SC_CLOSE */, 0)
					} else if initDialog != nil {
						initDialog(lparam.Wnd)
					}
				}
			}
			next, _, _ := syscall.Syscall6(callNextHookEx, 4,
				hook, uintptr(code), wparam, uintptr(unsafe.Pointer(lparam)),
				0, 0)
			return next
		}), 0, tid)

	if hook == 0 {
		return nil, err
	}
	if ctx == nil {
		return func() { unhookWindowsHookEx.Call(hook) }, nil
	}

	wait := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			if w := atomic.LoadUintptr(&wnd); w != 0 {
				sendMessage.Call(w, 0x0112 /* WM_SYSCOMMAND */, 0xf060 /* SC_CLOSE */, 0)
			}
		case <-wait:
		}
	}()
	return func() {
		unhookWindowsHookEx.Call(hook)
		close(wait)
	}, nil
}

func hookDialogTitle(ctx context.Context, title *string) (unhook context.CancelFunc, err error) {
	var init func(wnd uintptr)
	if title != nil {
		init = func(wnd uintptr) {
			setWindowText.Call(wnd, strptr(*title))
		}
	}
	return hookDialog(ctx, init)
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
	handle  uintptr
	logical _LOGFONT
}

func getFont() font {
	var metrics _NONCLIENTMETRICS
	metrics.Size = uint32(unsafe.Sizeof(metrics))
	systemParametersInfo.Call(0x29, // SPI_GETNONCLIENTMETRICS
		unsafe.Sizeof(metrics), uintptr(unsafe.Pointer(&metrics)), 0)
	return font{logical: metrics.MessageFont}
}

func (f *font) ForDPI(dpi dpi) uintptr {
	if h := -int32(dpi.Scale(12)); f.handle == 0 || f.logical.Height != h {
		f.Delete()
		f.logical.Height = h
		f.handle, _, _ = createFontIndirect.Call(uintptr(unsafe.Pointer(&f.logical)))
	}
	return f.handle
}

func (f *font) Delete() {
	if f.handle != 0 {
		deleteObject.Call(f.handle)
		f.handle = 0
	}
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

func getWindowString(wnd uintptr) string {
	len, _, _ := getWindowTextLength.Call(wnd)
	buf := make([]uint16, len+1)
	getWindowText.Call(wnd, uintptr(unsafe.Pointer(&buf[0])), len+1)
	return syscall.UTF16ToString(buf)
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

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-cwpretstruct
type _CWPRETSTRUCT struct {
	Result  uintptr
	LParam  uintptr
	WParam  uintptr
	Message uint32
	Wnd     uintptr
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

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type _MSG struct {
	Owner   syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      _POINT
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

// https://github.com/wine-mirror/wine/blob/master/include/unknwn.idl

type _IUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

func uuid(s string) uintptr {
	return (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
}

type _COMObject struct{}

func (o *_COMObject) Call(trap uintptr, a ...uintptr) (r1, r2 uintptr, lastErr error) {
	switch nargs := uintptr(len(a)); nargs {
	case 0:
		return syscall.Syscall(trap, nargs+1, uintptr(unsafe.Pointer(o)), 0, 0)
	case 1:
		return syscall.Syscall(trap, nargs+1, uintptr(unsafe.Pointer(o)), a[0], 0)
	case 2:
		return syscall.Syscall(trap, nargs+1, uintptr(unsafe.Pointer(o)), a[0], a[1])
	default:
		panic("COM call with too many arguments.")
	}
}
