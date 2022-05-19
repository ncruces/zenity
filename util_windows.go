package zenity

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	comctl32 = windows.NewLazySystemDLL("comctl32.dll")
	comdlg32 = windows.NewLazySystemDLL("comdlg32.dll")
	gdi32    = windows.NewLazySystemDLL("gdi32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	ntdll    = windows.NewLazySystemDLL("ntdll.dll")
	ole32    = windows.NewLazySystemDLL("ole32.dll")
	shell32  = windows.NewLazySystemDLL("shell32.dll")
	user32   = windows.NewLazySystemDLL("user32.dll")
	wtsapi32 = windows.NewLazySystemDLL("wtsapi32.dll")

	commDlgExtendedError = comdlg32.NewProc("CommDlgExtendedError")
	initCommonControlsEx = comctl32.NewProc("InitCommonControlsEx")

	createFontIndirect = gdi32.NewProc("CreateFontIndirectW")
	deleteObject       = gdi32.NewProc("DeleteObject")
	getDeviceCaps      = gdi32.NewProc("GetDeviceCaps")

	activateActCtx     = kernel32.NewProc("ActivateActCtx")
	createActCtx       = kernel32.NewProc("CreateActCtxW")
	deactivateActCtx   = kernel32.NewProc("DeactivateActCtx")
	getConsoleWindow   = kernel32.NewProc("GetConsoleWindow")
	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	getModuleHandle    = kernel32.NewProc("GetModuleHandleW")
	getSystemDirectory = kernel32.NewProc("GetSystemDirectoryW")

	coCreateInstance = ole32.NewProc("CoCreateInstance")
	coInitializeEx   = ole32.NewProc("CoInitializeEx")
	coTaskMemFree    = ole32.NewProc("CoTaskMemFree")
	coUninitialize   = ole32.NewProc("CoUninitialize")

	callNextHookEx               = user32.NewProc("CallNextHookEx")
	createIconFromResource       = user32.NewProc("CreateIconFromResource")
	createWindowEx               = user32.NewProc("CreateWindowExW")
	defWindowProc                = user32.NewProc("DefWindowProcW")
	destroyIcon                  = user32.NewProc("DestroyIcon")
	destroyWindow                = user32.NewProc("DestroyWindow")
	dispatchMessage              = user32.NewProc("DispatchMessageW")
	enableWindow                 = user32.NewProc("EnableWindow")
	enumChildWindows             = user32.NewProc("EnumChildWindows")
	enumWindows                  = user32.NewProc("EnumWindows")
	getDpiForWindow              = user32.NewProc("GetDpiForWindow")
	getMessage                   = user32.NewProc("GetMessageW")
	getSystemMetrics             = user32.NewProc("GetSystemMetrics")
	getWindowDC                  = user32.NewProc("GetWindowDC")
	getWindowRect                = user32.NewProc("GetWindowRect")
	getWindowText                = user32.NewProc("GetWindowTextW")
	getWindowTextLength          = user32.NewProc("GetWindowTextLengthW")
	getWindowThreadProcessId     = user32.NewProc("GetWindowThreadProcessId")
	isDialogMessage              = user32.NewProc("IsDialogMessageW")
	loadIcon                     = user32.NewProc("LoadIconW")
	loadImage                    = user32.NewProc("LoadImageW")
	postQuitMessage              = user32.NewProc("PostQuitMessage")
	registerClassEx              = user32.NewProc("RegisterClassExW")
	releaseDC                    = user32.NewProc("ReleaseDC")
	sendMessage                  = user32.NewProc("SendMessageW")
	setFocus                     = user32.NewProc("SetFocus")
	setForegroundWindow          = user32.NewProc("SetForegroundWindow")
	setThreadDpiAwarenessContext = user32.NewProc("SetThreadDpiAwarenessContext")
	setWindowLong                = user32.NewProc("SetWindowLongW")
	setWindowPos                 = user32.NewProc("SetWindowPos")
	setWindowsHookEx             = user32.NewProc("SetWindowsHookExW")
	setWindowText                = user32.NewProc("SetWindowTextW")
	showWindow                   = user32.NewProc("ShowWindow")
	systemParametersInfo         = user32.NewProc("SystemParametersInfoW")
	translateMessage             = user32.NewProc("TranslateMessage")
	unhookWindowsHookEx          = user32.NewProc("UnhookWindowsHookEx")
	unregisterClass              = user32.NewProc("UnregisterClassW")
)

func intptr(i int64) uintptr {
	return uintptr(i)
}

func strptr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func setup() context.CancelFunc {
	var wnd uintptr
	enumWindows.Call(syscall.NewCallback(setupEnumCallback), uintptr(unsafe.Pointer(&wnd)))
	if wnd == 0 {
		wnd, _, _ = getConsoleWindow.Call()
	}
	if wnd != 0 {
		setForegroundWindow.Call(wnd)
	}

	runtime.LockOSThread()

	var restore uintptr
	cookie := enableVisualStyles()
	if setThreadDpiAwarenessContext.Find() == nil {
		// try:
		//   DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2
		//   DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE
		//   DPI_AWARENESS_CONTEXT_SYSTEM_AWARE
		for i := -4; i <= -2; i++ {
			restore, _, _ = setThreadDpiAwarenessContext.Call(uintptr(i))
			if restore != 0 {
				break
			}
		}
	}

	var icc _INITCOMMONCONTROLSEX
	icc.Size = uint32(unsafe.Sizeof(icc))
	icc.ICC = 0x00004020 // ICC_STANDARD_CLASSES|ICC_PROGRESS_CLASS
	initCommonControlsEx.Call(uintptr(unsafe.Pointer(&icc)))

	return func() {
		if restore != 0 {
			setThreadDpiAwarenessContext.Call(restore)
		}
		if cookie != 0 {
			deactivateActCtx.Call(0, cookie)
		}
		runtime.UnlockOSThread()
	}
}

func setupEnumCallback(wnd uintptr, lparam *uintptr) uintptr {
	var pid uintptr
	getWindowThreadProcessId.Call(wnd, uintptr(unsafe.Pointer(&pid)))
	if int(pid) == os.Getpid() {
		*lparam = wnd
		return 0 // stop enumeration
	}
	return 1 // continue enumeration
}

func commDlgError() error {
	s, _, _ := commDlgExtendedError.Call()
	if s == 0 {
		return ErrCanceled
	} else {
		return fmt.Errorf("Common Dialog error: %x", s)
	}
}

func hookDialog(ctx context.Context, initDialog func(wnd uintptr)) (unhook context.CancelFunc, err error) {
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}
	hook, err := newDialogHook(ctx, initDialog)
	if err != nil {
		return nil, err
	}
	return hook.unhook, nil
}

type dialogHook struct {
	ctx  context.Context
	tid  uintptr
	wnd  uintptr
	hook uintptr
	done chan struct{}
	init func(wnd uintptr)
}

func newDialogHook(ctx context.Context, initDialog func(wnd uintptr)) (*dialogHook, error) {
	tid, _, _ := getCurrentThreadId.Call()
	hk, _, err := setWindowsHookEx.Call(5, // WH_CBT
		syscall.NewCallback(dialogHookProc), 0, tid)
	if hk == 0 {
		return nil, err
	}

	hook := dialogHook{
		ctx:  ctx,
		tid:  tid,
		hook: hk,
		init: initDialog,
	}
	if ctx != nil {
		hook.done = make(chan struct{})
		go hook.wait()
	}

	saveBackRef(tid, unsafe.Pointer(&hook))
	return &hook, nil
}

func dialogHookProc(code int32, wparam, lparam uintptr) uintptr {
	if code == 5 { // HCBT_ACTIVATE
		tid, _, _ := getCurrentThreadId.Call()
		hook := (*dialogHook)(loadBackRef(tid))
		atomic.StoreUintptr(&hook.wnd, wparam)
		if hook.ctx != nil && hook.ctx.Err() != nil {
			sendMessage.Call(wparam, _WM_SYSCOMMAND, _SC_CLOSE, 0)
		} else if hook.init != nil {
			hook.init(wparam)
		}
	}
	next, _, _ := callNextHookEx.Call(
		0, uintptr(code), wparam, lparam)
	return next
}

func (h *dialogHook) unhook() {
	deleteBackRef(h.tid)
	if h.done != nil {
		close(h.done)
	}
	unhookWindowsHookEx.Call(h.hook)
}

func (h *dialogHook) wait() {
	select {
	case <-h.ctx.Done():
		if wnd := atomic.LoadUintptr(&h.wnd); wnd != 0 {
			sendMessage.Call(wnd, _WM_SYSCOMMAND, _SC_CLOSE, 0)
		}
	case <-h.done:
	}
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

var backRefs struct {
	sync.Mutex
	m map[uintptr]unsafe.Pointer
}

func saveBackRef(id uintptr, ptr unsafe.Pointer) {
	backRefs.Lock()
	defer backRefs.Unlock()
	if backRefs.m == nil {
		backRefs.m = map[uintptr]unsafe.Pointer{}
	} else if _, ok := backRefs.m[id]; ok {
		panic("saveBackRef")
	}
	backRefs.m[id] = ptr
}

func loadBackRef(id uintptr) unsafe.Pointer {
	backRefs.Lock()
	defer backRefs.Unlock()
	return backRefs.m[id]
}

func deleteBackRef(id uintptr) {
	backRefs.Lock()
	defer backRefs.Unlock()
	delete(backRefs.m, id)
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

func (d dpi) scale(dim uintptr) uintptr {
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

func (f *font) forDPI(dpi dpi) uintptr {
	if h := -int32(dpi.scale(12)); f.handle == 0 || f.logical.Height != h {
		f.delete()
		f.logical.Height = h
		f.handle, _, _ = createFontIndirect.Call(uintptr(unsafe.Pointer(&f.logical)))
	}
	return f.handle
}

func (f *font) delete() {
	if f.handle != 0 {
		deleteObject.Call(f.handle)
		f.handle = 0
	}
}

type icon struct {
	handle  uintptr
	destroy bool
}

func getIcon(i any) icon {
	var res icon
	var resource uintptr
	switch i {
	case ErrorIcon:
		resource = 32513 // IDI_ERROR
	case QuestionIcon:
		resource = 32514 // IDI_QUESTION
	case WarningIcon:
		resource = 32515 // IDI_WARNING
	case InfoIcon:
		resource = 32516 // IDI_INFORMATION
	}
	if resource != 0 {
		res.handle, _, _ = loadIcon.Call(0, resource)
		return res
	}

	path, ok := i.(string)
	if !ok {
		return res
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return res
	}

	switch {
	case bytes.HasPrefix(data, []byte("\x00\x00\x01\x00")):
		res.handle, _, _ = loadImage.Call(0,
			strptr(path),
			1, /*IMAGE_ICON*/
			0, 0,
			0x00008050 /*LR_LOADFROMFILE|LR_DEFAULTSIZE|LR_SHARED*/)
	case bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")):
		res.handle, _, _ = createIconFromResource.Call(
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(len(data)),
			1, 0x00030000)
		res.destroy = true
	}
	return res
}

func (i *icon) delete() {
	if i.handle != 0 {
		destroyIcon.Call(i.handle)
		i.handle = 0
	}
}

func centerWindow(wnd uintptr) {
	getMetric := func(i uintptr) int32 {
		n, _, _ := getSystemMetrics.Call(i)
		return int32(n)
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

	atom, _, err := registerClassEx.Call(uintptr(unsafe.Pointer(&wcx)))
	return atom, err
}

// https://docs.microsoft.com/en-us/windows/win32/winmsg/using-messages-and-message-queues
func messageLoop(wnd uintptr) error {
	getMessage := getMessage.Addr()
	isDialogMessage := isDialogMessage.Addr()
	translateMessage := translateMessage.Addr()
	dispatchMessage := dispatchMessage.Addr()

	for {
		var msg _MSG
		s, _, err := syscall.Syscall6(getMessage, 4, uintptr(unsafe.Pointer(&msg)), 0, 0, 0, 0, 0)
		if int32(s) == -1 {
			return err
		}
		if s == 0 {
			return nil
		}

		s, _, _ = syscall.Syscall(isDialogMessage, 2, wnd, uintptr(unsafe.Pointer(&msg)), 0)
		if s == 0 {
			syscall.Syscall(translateMessage, 1, uintptr(unsafe.Pointer(&msg)), 0, 0)
			syscall.Syscall(dispatchMessage, 1, uintptr(unsafe.Pointer(&msg)), 0, 0)
		}
	}
}

// https://stackoverflow.com/questions/4308503/how-to-enable-visual-styles-without-a-manifest
func enableVisualStyles() (cookie uintptr) {
	var dir [260]uint16
	n, _, _ := getSystemDirectory.Call(uintptr(unsafe.Pointer(&dir[0])), uintptr(len(dir)))
	if n == 0 || int(n) >= len(dir) {
		return
	}

	var ctx _ACTCTX
	ctx.Size = uint32(unsafe.Sizeof(ctx))
	ctx.Flags = 0x01c // ACTCTX_FLAG_RESOURCE_NAME_VALID|ACTCTX_FLAG_SET_PROCESS_DEFAULT|ACTCTX_FLAG_ASSEMBLY_DIRECTORY_VALID
	ctx.Source = syscall.StringToUTF16Ptr("shell32.dll")
	ctx.AssemblyDirectory = &dir[0]
	ctx.ResourceName = 124

	if h, _, _ := createActCtx.Call(uintptr(unsafe.Pointer(&ctx))); h != 0 {
		activateActCtx.Call(h, uintptr(unsafe.Pointer(&cookie)))
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-actctxw
type _ACTCTX struct {
	Size                  uint32
	Flags                 uint32
	Source                *uint16
	ProcessorArchitecture uint16
	LangId                uint16
	AssemblyDirectory     *uint16
	ResourceName          uintptr
	ApplicationName       *uint16
	Module                uintptr
}

// https://docs.microsoft.com/en-us/windows/win32/api/commctrl/ns-commctrl-initcommoncontrolsex
type _INITCOMMONCONTROLSEX struct {
	Size uint32
	ICC  uint32
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

// https://docs.microsoft.com/en-us/windows/win32/api/minwinbase/ns-minwinbase-systemtime
type _SYSTEMTIME struct {
	year         uint16
	month        uint16
	dayOfWeek    uint16
	day          uint16
	hour         uint16
	minute       uint16
	second       uint16
	milliseconds uint16
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

//go:uintptrescapes
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
