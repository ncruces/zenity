//go:build windows

package win

import "golang.org/x/sys/windows"

const (
	// CreateActCtx flags
	ACTCTX_FLAG_PROCESSOR_ARCHITECTURE_VALID = 0x001
	ACTCTX_FLAG_LANGID_VALID                 = 0x002
	ACTCTX_FLAG_ASSEMBLY_DIRECTORY_VALID     = 0x004
	ACTCTX_FLAG_RESOURCE_NAME_VALID          = 0x008
	ACTCTX_FLAG_SET_PROCESS_DEFAULT          = 0x010
	ACTCTX_FLAG_APPLICATION_NAME_VALID       = 0x020
	ACTCTX_FLAG_HMODULE_VALID                = 0x080

	// Control signals
	CTRL_C_EVENT        = windows.CTRL_C_EVENT
	CTRL_BREAK_EVENT    = windows.CTRL_BREAK_EVENT
	CTRL_CLOSE_EVENT    = windows.CTRL_CLOSE_EVENT
	CTRL_LOGOFF_EVENT   = windows.CTRL_LOGOFF_EVENT
	CTRL_SHUTDOWN_EVENT = windows.CTRL_SHUTDOWN_EVENT
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
func GenerateConsoleCtrlEvent(ctrlEvent uint32, processGroupID uint32) (err error) {
	return windows.GenerateConsoleCtrlEvent(ctrlEvent, processGroupID)
}

//sys ActivateActCtx(actCtx Handle, cookie *uintptr) (err error)
//sys CreateActCtx(actCtx *ACTCTX) (ret Handle, err error) [failretval==^Handle(0)] = CreateActCtxW
//sys DeactivateActCtx(flags uint32, cookie uintptr) (err error)
//sys GetConsoleWindow() (ret HWND)
//sys GetModuleHandle(moduleName *uint16) (ret Handle, err error) = GetModuleHandleW
//sys GlobalAlloc(flags uint32, bytes uintptr) (ret Handle, err error)
//sys GlobalFree(mem Handle) (err error) [failretval!=0]
//sys ReleaseActCtx(actCtx Handle)
