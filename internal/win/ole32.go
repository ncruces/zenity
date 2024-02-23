//go:build windows

package win

import (
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	COINIT_MULTITHREADED     = windows.COINIT_MULTITHREADED
	COINIT_APARTMENTTHREADED = windows.COINIT_APARTMENTTHREADED
	COINIT_DISABLE_OLE1DDE   = windows.COINIT_DISABLE_OLE1DDE
	COINIT_SPEED_OVER_MEMORY = windows.COINIT_SPEED_OVER_MEMORY

	CLSCTX_INPROC_SERVER  = windows.CLSCTX_INPROC_SERVER
	CLSCTX_INPROC_HANDLER = windows.CLSCTX_INPROC_HANDLER
	CLSCTX_LOCAL_SERVER   = windows.CLSCTX_LOCAL_SERVER
	CLSCTX_REMOTE_SERVER  = windows.CLSCTX_REMOTE_SERVER
	CLSCTX_ALL            = windows.CLSCTX_INPROC_SERVER | windows.CLSCTX_INPROC_HANDLER | windows.CLSCTX_LOCAL_SERVER | windows.CLSCTX_REMOTE_SERVER

	E_CANCELED         = windows.ERROR_CANCELLED | windows.FACILITY_WIN32<<16 | 0x80000000
	E_TIMEOUT          = windows.ERROR_TIMEOUT | windows.FACILITY_WIN32<<16 | 0x80000000
	RPC_E_CHANGED_MODE = syscall.Errno(windows.RPC_E_CHANGED_MODE)
	S_FALSE            = syscall.Errno(windows.S_FALSE)
)

func CoInitializeEx(reserved uintptr, coInit uint32) error {
	return windows.CoInitializeEx(reserved, coInit)
}

func CoUninitialize() { windows.CoUninitialize() }

func CoTaskMemFree(address unsafe.Pointer) { windows.CoTaskMemFree(address) }

// https://github.com/wine-mirror/wine/blob/master/include/unknwn.idl

type IUnknown struct{}
type iUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

func (u *IUnknown) Release() {
	vtbl := *(**iUnknownVtbl)(unsafe.Pointer(u))
	u.call(vtbl.Release)
}

//go:uintptrescapes
func (u *IUnknown) call(trap uintptr, a ...uintptr) (r1, r2 uintptr, lastErr error) {
	return syscall.SyscallN(trap, append([]uintptr{uintptr(unsafe.Pointer(u))}, a...)...)
}

// https://github.com/wine-mirror/wine/blob/master/include/objidl.idl

type IBindCtx struct{ IUnknown }

//sys CoCreateInstance(clsid *GUID, unkOuter *IUnknown, clsContext int32, iid *GUID, address unsafe.Pointer) (res error) = ole32.CoCreateInstance

func guid(s string) *GUID {
	// TODO: use unsafe.StringData after 1.20
	return (*GUID)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data))
}
