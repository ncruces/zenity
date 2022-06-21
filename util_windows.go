package zenity

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")

	callNextHookEx       = user32.NewProc("CallNextHookEx")
	setWindowsHookEx     = user32.NewProc("SetWindowsHookExW")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
	unhookWindowsHookEx  = user32.NewProc("UnhookWindowsHookEx")
)

func intptr(i int64) uintptr  { return uintptr(i) }
func strptr(s string) *uint16 { return syscall.StringToUTF16Ptr(s) }
func hwnd(i uint64) win.HWND  { return win.HWND(uintptr(i)) }

func setup() context.CancelFunc {
	var wnd win.HWND
	win.EnumWindows(syscall.NewCallback(setupEnumCallback), unsafe.Pointer(&wnd))
	if wnd == 0 {
		wnd = win.GetConsoleWindow()
	}
	if wnd != 0 {
		win.SetForegroundWindow(wnd)
	}

	runtime.LockOSThread()

	var restore uintptr
	cookie := enableVisualStyles()
	for dpi := win.DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2; dpi <= win.DPI_AWARENESS_CONTEXT_SYSTEM_AWARE; dpi++ {
		var err error
		restore, err = win.SetThreadDpiAwarenessContext(dpi)
		if restore != 0 || err != nil {
			break
		}
	}

	var icc win.INITCOMMONCONTROLSEX
	icc.Size = uint32(unsafe.Sizeof(icc))
	icc.ICC = 0x00004020 // ICC_STANDARD_CLASSES|ICC_PROGRESS_CLASS
	win.InitCommonControlsEx(&icc)

	return func() {
		if restore != 0 {
			win.SetThreadDpiAwarenessContext(restore)
		}
		if cookie != 0 {
			win.DeactivateActCtx(0, cookie)
		}
		runtime.UnlockOSThread()
	}
}

func setupEnumCallback(wnd win.HWND, lparam *win.HWND) uintptr {
	var pid uint32
	win.GetWindowThreadProcessId(wnd, &pid)
	if int(pid) == os.Getpid() {
		*lparam = wnd
		return 0 // stop enumeration
	}
	return 1 // continue enumeration
}

func hookDialog(ctx context.Context, icon any, title *string, init func(wnd win.HWND)) (unhook context.CancelFunc, err error) {
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}
	hook, err := newDialogHook(ctx, icon, title, init)
	if err != nil {
		return nil, err
	}
	return hook.unhook, nil
}

type dialogHook struct {
	ctx   context.Context
	tid   uint32
	wnd   uintptr
	hook  uintptr
	done  chan struct{}
	icon  any
	title *string
	init  func(wnd win.HWND)
}

func newDialogHook(ctx context.Context, icon any, title *string, init func(wnd win.HWND)) (*dialogHook, error) {
	tid := win.GetCurrentThreadId()
	hk, _, err := setWindowsHookEx.Call(12, // WH_CALLWNDPROCRET
		syscall.NewCallback(dialogHookProc), 0, uintptr(tid))
	if hk == 0 {
		return nil, err
	}

	hook := dialogHook{
		ctx:   ctx,
		tid:   tid,
		hook:  hk,
		icon:  icon,
		title: title,
		init:  init,
	}
	if ctx != nil {
		hook.done = make(chan struct{})
		go hook.wait()
	}

	saveBackRef(uintptr(tid), unsafe.Pointer(&hook))
	return &hook, nil
}

func dialogHookProc(code int32, wparam uintptr, lparam *_CWPRETSTRUCT) uintptr {
	if lparam.Message == 0x0110 { // WM_INITDIALOG
		tid := win.GetCurrentThreadId()
		hook := (*dialogHook)(loadBackRef(uintptr(tid)))
		atomic.StoreUintptr(&hook.wnd, uintptr(lparam.Wnd))
		if hook.ctx != nil && hook.ctx.Err() != nil {
			win.SendMessage(lparam.Wnd, win.WM_SYSCOMMAND, _SC_CLOSE, 0)
		} else {
			if hook.icon != nil {
				icon := getIcon(hook.icon)
				if icon.handle != 0 {
					defer icon.delete()
					win.SendMessage(lparam.Wnd, win.WM_SETICON, 0, uintptr(icon.handle))
				}
			}
			if hook.title != nil {
				win.SetWindowText(lparam.Wnd, strptr(*hook.title))
			}
			if hook.init != nil {
				hook.init(lparam.Wnd)
			}
		}
	}
	next, _, _ := callNextHookEx.Call(
		0, uintptr(code), wparam, uintptr(unsafe.Pointer(lparam)))
	return next
}

func (h *dialogHook) unhook() {
	deleteBackRef(uintptr(h.tid))
	if h.done != nil {
		close(h.done)
	}
	unhookWindowsHookEx.Call(h.hook)
}

func (h *dialogHook) wait() {
	select {
	case <-h.ctx.Done():
		if wnd := atomic.LoadUintptr(&h.wnd); wnd != 0 {
			win.SendMessage(win.HWND(wnd), win.WM_SYSCOMMAND, _SC_CLOSE, 0)
		}
	case <-h.done:
	}
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

type dpi int

func getDPI(wnd win.HWND) dpi {
	res, _ := win.GetDpiForWindow(wnd)
	if res != 0 {
		return dpi(res)
	}

	if dc := win.GetWindowDC(wnd); dc != 0 {
		res = win.GetDeviceCaps(dc, win.LOGPIXELSY)
		win.ReleaseDC(wnd, dc)
	}

	if res == 0 {
		return win.USER_DEFAULT_SCREEN_DPI
	}
	return dpi(res)
}

func (d dpi) scale(dim int) int {
	if d == 0 {
		return dim
	}
	return dim * int(d) / win.USER_DEFAULT_SCREEN_DPI
}

type font struct {
	handle  win.Handle
	logical win.LOGFONT
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
		f.handle = win.CreateFontIndirect(&f.logical)
	}
	return uintptr(f.handle)
}

func (f *font) delete() {
	if f.handle != 0 {
		win.DeleteObject(f.handle)
		f.handle = 0
	}
}

type icon struct {
	handle  win.Handle
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
		res.handle, _ = win.LoadIcon(0, resource)
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
		res.handle, _ = win.LoadImage(0,
			strptr(path),
			win.IMAGE_ICON, 0, 0,
			win.LR_LOADFROMFILE|win.LR_DEFAULTSIZE|win.LR_SHARED)
	case bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")):
		res.handle, err = win.CreateIconFromResource(
			data, true, 0x00030000)
		res.destroy = true
	}
	return res
}

func (i *icon) delete() {
	if i.handle != 0 {
		win.DestroyIcon(i.handle)
		i.handle = 0
	}
}

func centerWindow(wnd win.HWND) {
	var rect win.RECT
	win.GetWindowRect(wnd, &rect)
	x := (win.GetSystemMetrics(0 /* SM_CXSCREEN */) - int(rect.Right-rect.Left)) / 2
	y := (win.GetSystemMetrics(1 /* SM_CYSCREEN */) - int(rect.Bottom-rect.Top)) / 2
	win.SetWindowPos(wnd, 0, x, y, 0, 0, 0x5) // SWP_NOZORDER|SWP_NOSIZE
}

func getWindowString(wnd win.HWND) string {
	len, _ := win.GetWindowTextLength(wnd)
	buf := make([]uint16, len+1)
	win.GetWindowText(wnd, &buf[0], len+1)
	return syscall.UTF16ToString(buf)
}

func registerClass(instance, icon win.Handle, proc uintptr) (*uint16, error) {
	name := "WC_" + strconv.FormatUint(uint64(proc), 16)

	var wcx win.WNDCLASSEX
	wcx.Size = uint32(unsafe.Sizeof(wcx))
	wcx.WndProc = proc
	wcx.Icon = icon
	wcx.Instance = instance
	wcx.Background = 5 // COLOR_WINDOW
	wcx.ClassName = strptr(name)

	if err := win.RegisterClassEx(&wcx); err != nil {
		return nil, err
	}
	return wcx.ClassName, nil
}

// https://stackoverflow.com/questions/4308503/how-to-enable-visual-styles-without-a-manifest
func enableVisualStyles() (cookie uintptr) {
	dir, err := win.GetSystemDirectory()
	if err != nil {
		return
	}

	var ctx win.ACTCTX
	ctx.Size = uint32(unsafe.Sizeof(ctx))
	ctx.Flags = win.ACTCTX_FLAG_RESOURCE_NAME_VALID | win.ACTCTX_FLAG_SET_PROCESS_DEFAULT | win.ACTCTX_FLAG_ASSEMBLY_DIRECTORY_VALID
	ctx.Source = strptr("shell32.dll")
	ctx.AssemblyDirectory = strptr(dir)
	ctx.ResourceName = 124

	if hnd, err := win.CreateActCtx(&ctx); err == nil {
		win.ActivateActCtx(hnd, &cookie)
		win.ReleaseActCtx(hnd)
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-cwpretstruct
type _CWPRETSTRUCT struct {
	Result  uintptr
	LParam  uintptr
	WParam  uintptr
	Message uint32
	Wnd     win.HWND
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-nonclientmetricsw
type _NONCLIENTMETRICS struct {
	Size            uint32
	BorderWidth     int32
	ScrollWidth     int32
	ScrollHeight    int32
	CaptionWidth    int32
	CaptionHeight   int32
	CaptionFont     win.LOGFONT
	SmCaptionWidth  int32
	SmCaptionHeight int32
	SmCaptionFont   win.LOGFONT
	MenuWidth       int32
	MenuHeight      int32
	MenuFont        win.LOGFONT
	StatusFont      win.LOGFONT
	MessageFont     win.LOGFONT
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
