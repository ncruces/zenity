package zenity

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
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

func intptr(i int64) uintptr  { return uintptr(i) }
func strptr(s string) *uint16 { return syscall.StringToUTF16Ptr(s) }

func setup(wnd win.HWND) context.CancelFunc {
	if wnd == 0 {
		win.EnumWindows(syscall.NewCallback(setupEnumCallback), unsafe.Pointer(&wnd))
	}
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
	icc.ICC = win.ICC_STANDARD_CLASSES | win.ICC_DATE_CLASSES | win.ICC_PROGRESS_CLASS
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
	if ctx == nil && icon == nil && title == nil && init == nil {
		return func() {}, nil
	}
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
	icon  icon
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
	ico, _ := getIcon(icon)

	hook := dialogHook{
		ctx:   ctx,
		tid:   tid,
		hook:  hk,
		icon:  ico,
		title: title,
		init:  init,
	}
	if ctx != nil && ctx.Done() != nil {
		hook.done = make(chan struct{})
		go hook.wait()
	}

	saveBackRef(uintptr(tid), unsafe.Pointer(&hook))
	return &hook, nil
}

func dialogHookProc(code int32, wparam uintptr, lparam *win.CWPRETSTRUCT) uintptr {
	if lparam.Message == win.WM_INITDIALOG {
		tid := win.GetCurrentThreadId()
		hook := (*dialogHook)(loadBackRef(uintptr(tid)))
		atomic.StoreUintptr(&hook.wnd, uintptr(lparam.Wnd))
		if hook.ctx != nil && hook.ctx.Err() != nil {
			win.SendMessage(lparam.Wnd, win.WM_SYSCOMMAND, win.SC_CLOSE, 0)
		} else {
			if hook.icon.handle != 0 {
				win.SendMessage(lparam.Wnd, win.WM_SETICON, 0, uintptr(hook.icon.handle))
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
	h.icon.delete()
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

func getIcon(i any) (icon icon, err error) {
	var resource uintptr
	switch i {
	case ErrorIcon:
		resource = win.IDI_ERROR
	case QuestionIcon:
		resource = win.IDI_QUESTION
	case WarningIcon:
		resource = win.IDI_WARNING
	case InfoIcon:
		resource = win.IDI_INFORMATION
	}
	if resource != 0 {
		icon.handle, err = win.LoadIcon(0, resource)
		return icon, err
	}

	path, ok := i.(string)
	if !ok {
		return icon, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return icon, err
	}
	defer file.Close()

	var peek [8]byte
	_, err = file.ReadAt(peek[:], 0)
	if err != nil {
		return icon, err
	}

	if bytes.Equal(peek[:], []byte("\x89PNG\r\n\x1a\n")) {
		data, err := io.ReadAll(file)
		if err != nil {
			return icon, err
		}
		icon.handle, err = win.CreateIconFromResourceEx(
			data, true, 0x00030000, 0, 0,
			win.LR_DEFAULTSIZE)
		if err != nil {
			return icon, err
		}
	} else {
		instance, err := win.GetModuleHandle(nil)
		if err != nil {
			return icon, err
		}
		path, err = filepath.Abs(path)
		if err != nil {
			return icon, err
		}
		var i uint16
		icon.handle, err = win.ExtractAssociatedIcon(
			instance, strptr(path), &i)
		if err != nil {
			return icon, err
		}
	}
	icon.destroy = true
	return icon, nil
}

func (i *icon) delete() {
	if i.handle != 0 && i.destroy {
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

func registerClass(instance, icon win.Handle, proc uintptr) (*uint16, error) {
	name := "WC_" + strconv.FormatUint(uint64(proc), 16)

	var wcx win.WNDCLASSEX
	wcx.Size = uint32(unsafe.Sizeof(wcx))
	wcx.WndProc = proc
	wcx.Icon = icon
	wcx.Instance = instance
	wcx.Background = win.COLOR_WINDOW
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
