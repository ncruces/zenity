package zenity

import "syscall"

var (
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")

	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")

	coInitializeEx   = ole32.NewProc("CoInitializeEx")
	coUninitialize   = ole32.NewProc("CoUninitialize")
	coCreateInstance = ole32.NewProc("CoCreateInstance")
	coTaskMemFree    = ole32.NewProc("CoTaskMemFree")

	getClassName        = user32.NewProc("GetClassNameA")
	setWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	callNextHookEx      = user32.NewProc("CallNextHookEx")
	enumChildWindows    = user32.NewProc("EnumChildWindows")
	getDlgCtrlID        = user32.NewProc("GetDlgCtrlID")
	setWindowText       = user32.NewProc("SetWindowTextW")
)

type _CWPRETSTRUCT struct {
	Result  uintptr
	LParam  uintptr
	WParam  uintptr
	Message uint32
	HWnd    uintptr
}
