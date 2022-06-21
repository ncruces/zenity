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
)

const (
	_WS_ZEN_DIALOG    = win.WS_POPUPWINDOW | win.WS_CLIPSIBLINGS | win.WS_DLGFRAME
	_WS_EX_ZEN_DIALOG = win.WS_EX_CONTROLPARENT | win.WS_EX_WINDOWEDGE | win.WS_EX_DLGMODALFRAME
	_WS_ZEN_LABEL     = win.WS_CHILD | win.WS_VISIBLE | win.WS_GROUP | win.SS_WORDELLIPSIS | win.SS_EDITCONTROL | win.SS_NOPREFIX
	_WS_ZEN_CONTROL   = win.WS_CHILD | win.WS_VISIBLE | win.WS_GROUP | win.WS_TABSTOP
	_WS_ZEN_BUTTON    = _WS_ZEN_CONTROL
)

const nullptr win.Pointer = 0

func intptr(i int64) uintptr  { return uintptr(i) }
func strptr(s string) *uint16 { return syscall.StringToUTF16Ptr(s) }

func hwnd(v reflect.Value) win.HWND {
	return win.HWND(uintptr(v.Uint()))
}

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
	hook  win.Handle
	done  chan struct{}
	icon  any
	title *string
	init  func(wnd win.HWND)
}

func newDialogHook(ctx context.Context, icon any, title *string, init func(wnd win.HWND)) (*dialogHook, error) {
	tid := win.GetCurrentThreadId()
	hk, err := win.SetWindowsHookEx(win.WH_CALLWNDPROCRET,
		syscall.NewCallback(dialogHookProc), 0, tid)
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
	if lparam.Message == win.WM_INITDIALOG {
		tid := win.GetCurrentThreadId()
		hook := (*dialogHook)(loadBackRef(uintptr(tid)))
		atomic.StoreUintptr(&hook.wnd, uintptr(lparam.Wnd))
		if hook.ctx != nil && hook.ctx.Err() != nil {
			win.SendMessage(lparam.Wnd, win.WM_SYSCOMMAND, win.SC_CLOSE, 0)
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
	return win.CallNextHookEx(0, code, wparam, unsafe.Pointer(lparam))
}

func (h *dialogHook) unhook() {
	deleteBackRef(uintptr(h.tid))
	if h.done != nil {
		close(h.done)
	}
	win.UnhookWindowsHookEx(h.hook)
}

func (h *dialogHook) wait() {
	select {
	case <-h.ctx.Done():
		if wnd := atomic.LoadUintptr(&h.wnd); wnd != 0 {
			win.SendMessage(win.HWND(wnd), win.WM_SYSCOMMAND, win.SC_CLOSE, 0)
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
	var metrics win.NONCLIENTMETRICS
	metrics.Size = uint32(unsafe.Sizeof(metrics))
	win.SystemParametersInfo(win.SPI_GETNONCLIENTMETRICS,
		unsafe.Sizeof(metrics), unsafe.Pointer(&metrics), 0)
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
	x := win.GetSystemMetrics(win.SM_CXSCREEN) - int(rect.Right-rect.Left)
	y := win.GetSystemMetrics(win.SM_CYSCREEN) - int(rect.Bottom-rect.Top)
	win.SetWindowPos(wnd, 0, x/2, y/2, 0, 0, win.SWP_NOZORDER|win.SWP_NOSIZE)
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
