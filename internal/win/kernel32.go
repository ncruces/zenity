//go:build windows

package win

import "golang.org/x/sys/windows"

const (
	ACTCTX_FLAG_PROCESSOR_ARCHITECTURE_VALID = 0x001
	ACTCTX_FLAG_LANGID_VALID                 = 0x002
	ACTCTX_FLAG_ASSEMBLY_DIRECTORY_VALID     = 0x004
	ACTCTX_FLAG_RESOURCE_NAME_VALID          = 0x008
	ACTCTX_FLAG_SET_PROCESS_DEFAULT          = 0x010
	ACTCTX_FLAG_APPLICATION_NAME_VALID       = 0x020
	ACTCTX_FLAG_HMODULE_VALID                = 0x080
)

// https://docs.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-actctxw
type ACTCTX struct {
	Size                  uint32
	Flags                 uint32
	Source                *uint16
	ProcessorArchitecture uint16
	LangId                uint16
	AssemblyDirectory     *uint16
	ResourceName          uintptr
	ApplicationName       *uint16
	Module                Handle
}

func GetCurrentThreadId() (id uint32)     { return windows.GetCurrentThreadId() }
func GetSystemDirectory() (string, error) { return windows.GetSystemDirectory() }

//sys ActivateActCtx(actCtx Handle, cookie *uintptr) (err error) = kernel32.ActivateActCtx
//sys CreateActCtx(actCtx *ACTCTX) (ret Handle, err error) = kernel32.CreateActCtxW
//sys DeactivateActCtx(flags uint32, cookie uintptr) (err error) = kernel32.DeactivateActCtx
//sys GetConsoleWindow() (ret HWND) = kernel32.GetConsoleWindow
//sys GetModuleHandle(moduleName *uint16) (ret Handle, err error) = kernel32.GetModuleHandleW
//sys ReleaseActCtx(actCtx Handle) = kernel32.ReleaseActCtx
