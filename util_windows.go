package zenity

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")

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
	n, _, _ := commDlgExtendedError.Call()
	if n == 0 {
		return nil
	} else {
		return fmt.Errorf("Common Dialog error: %x", n)
	}
}

type _CWPRETSTRUCT struct {
	Result  uintptr
	LParam  uintptr
	WParam  uintptr
	Message uint32
	Wnd     uintptr
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
