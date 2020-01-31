package zenity

import (
	"context"
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

var (
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")
	wtsapi32 = syscall.NewLazyDLL("wtsapi32.dll")

	commDlgExtendedError = comdlg32.NewProc("CommDlgExtendedError")

	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")

	coInitializeEx   = ole32.NewProc("CoInitializeEx")
	coUninitialize   = ole32.NewProc("CoUninitialize")
	coCreateInstance = ole32.NewProc("CoCreateInstance")
	coTaskMemFree    = ole32.NewProc("CoTaskMemFree")

	sendMessage         = user32.NewProc("SendMessageW")
	getClassName        = user32.NewProc("GetClassNameW")
	setWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	callNextHookEx      = user32.NewProc("CallNextHookEx")
	enumChildWindows    = user32.NewProc("EnumChildWindows")
	getDlgCtrlID        = user32.NewProc("GetDlgCtrlID")
	setWindowText       = user32.NewProc("SetWindowTextW")
)

func commDlgError() error {
	s, _, _ := commDlgExtendedError.Call()
	if s == 0 {
		return nil
	} else {
		return fmt.Errorf("Common Dialog error: %x", s)
	}
}

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

	var mtx sync.Mutex
	var hook, wnd uintptr
	tid, _, _ := getCurrentThreadId.Call()
	hook, _, err = setWindowsHookEx.Call(12, // WH_CALLWNDPROCRET
		syscall.NewCallback(func(code int, wparam uintptr, lparam *_CWPRETSTRUCT) uintptr {
			if lparam.Message == 0x0110 { // WM_INITDIALOG
				name := [8]uint16{}
				getClassName.Call(lparam.Wnd, uintptr(unsafe.Pointer(&name)), uintptr(len(name)))
				if syscall.UTF16ToString(name[:]) == "#32770" { // The class for a dialog box
					var close bool

					mtx.Lock()
					if ctx != nil && ctx.Err() != nil {
						close = true
					} else {
						wnd = lparam.Wnd
					}
					mtx.Unlock()

					if close {
						sendMessage.Call(lparam.Wnd, 0x0112 /* WM_SYSCOMMAND */, 0xf060 /* SC_CLOSE */, 0)
					} else if initDialog != nil {
						initDialog(lparam.Wnd)
					}
				}
			}
			next, _, _ := callNextHookEx.Call(hook, uintptr(code), wparam, uintptr(unsafe.Pointer(lparam)))
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
			mtx.Lock()
			w := wnd
			mtx.Unlock()

			if w != 0 {
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

func hookDialogTitle(ctx context.Context, title string) (unhook context.CancelFunc, err error) {
	return hookDialog(ctx, func(wnd uintptr) {
		ptr := syscall.StringToUTF16Ptr(title)
		setWindowText.Call(wnd, uintptr(unsafe.Pointer(ptr)))
	})
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

type _IUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}
