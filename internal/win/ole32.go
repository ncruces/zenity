//go:build windows

package win

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	RPC_E_CHANGED_MODE syscall.Errno = 0x80010106
)

func CoInitializeEx(reserved uintptr, coInit uint32) error {
	return windows.CoInitializeEx(reserved, coInit)
}

func CoUninitialize() { windows.CoUninitialize() }

// https://github.com/wine-mirror/wine/blob/master/include/unknwn.idl

type IUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

type COMObject struct{}

//go:uintptrescapes
func (o *COMObject) Call(trap uintptr, a ...uintptr) (r1, r2 uintptr, lastErr error) {
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

//sys CoTaskMemFree(address uintptr) = ole32.CoTaskMemFree
//sys CoCreateInstance(clsid uintptr, unkOuter unsafe.Pointer, clsContext int32, iid uintptr, address unsafe.Pointer) (res error) = ole32.CoCreateInstance
