package zenity

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
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

	getModuleHandle    = kernel32.NewProc("GetModuleHandleW")
	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	getConsoleWindow   = kernel32.NewProc("GetConsoleWindow")

	coInitializeEx   = ole32.NewProc("CoInitializeEx")
	coUninitialize   = ole32.NewProc("CoUninitialize")
	coCreateInstance = ole32.NewProc("CoCreateInstance")
	coTaskMemFree    = ole32.NewProc("CoTaskMemFree")

	getMessage                   = user32.NewProc("GetMessageW")
	sendMessage                  = user32.NewProc("SendMessageW")
	getClassName                 = user32.NewProc("GetClassNameW")
	setWindowsHookEx             = user32.NewProc("SetWindowsHookExW")
	unhookWindowsHookEx          = user32.NewProc("UnhookWindowsHookEx")
	callNextHookEx               = user32.NewProc("CallNextHookEx")
	enumWindows                  = user32.NewProc("EnumWindows")
	enumChildWindows             = user32.NewProc("EnumChildWindows")
	setWindowText                = user32.NewProc("SetWindowTextW")
	getWindowText                = user32.NewProc("GetWindowTextW")
	getWindowTextLength          = user32.NewProc("GetWindowTextLengthW")
	setForegroundWindow          = user32.NewProc("SetForegroundWindow")
	getWindowThreadProcessId     = user32.NewProc("GetWindowThreadProcessId")
	setThreadDpiAwarenessContext = user32.NewProc("SetThreadDpiAwarenessContext")
)

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
			runtime.UnlockOSThread()
		}
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

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-cwpretstruct
type _CWPRETSTRUCT struct {
	Result  uintptr
	LParam  uintptr
	WParam  uintptr
	Message uint32
	Wnd     uintptr
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

func strptr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func intptr(i int64) uintptr {
	return uintptr(i)
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
	self := uintptr(unsafe.Pointer(o))
	nargs := uintptr(len(a))
	switch nargs {
	case 0:
		return syscall.Syscall(trap, nargs+1, self, 0, 0)
	case 1:
		return syscall.Syscall(trap, nargs+1, self, a[0], 0)
	case 2:
		return syscall.Syscall(trap, nargs+1, self, a[0], a[1])
	default:
		panic("COM call with too many arguments.")
	}
}
