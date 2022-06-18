//go:build windows

package win

const (
	NIM_ADD    = 0
	NIM_DELETE = 2
)

// https://docs.microsoft.com/en-us/windows/win32/api/shlobj_core/ns-shlobj_core-browseinfow
type BROWSEINFO struct {
	Owner        uintptr
	Root         uintptr
	DisplayName  *uint16
	Title        *uint16
	Flags        uint32
	CallbackFunc uintptr
	LParam       *uint16
	Image        int32
}

// https://docs.microsoft.com/en-us/windows/win32/api/shellapi/ns-shellapi-notifyicondataw
type NOTIFYICONDATA struct {
	StructSize      uint32
	Wnd             uintptr
	ID              uint32
	Flags           uint32
	CallbackMessage uint32
	Icon            uintptr
	Tip             [128]uint16 // NOTIFYICONDATAA_V1_SIZE
	State           uint32
	StateMask       uint32
	Info            [256]uint16
	Version         uint32
	InfoTitle       [64]uint16
	InfoFlags       uint32
	// GuidItem     [16]byte       // NOTIFYICONDATAA_V2_SIZE
	// BalloonIcon  syscall.Handle // NOTIFYICONDATAA_V3_SIZE
}

type IShellItem struct {
	COMObject
	*_IShellItemVtbl
}

type _IShellItemVtbl struct {
	IUnknownVtbl
	BindToHandler  uintptr
	GetParent      uintptr
	GetDisplayName uintptr
	GetAttributes  uintptr
	Compare        uintptr
}

//sys SHBrowseForFolder(bi *BROWSEINFO) (ptr uintptr) = shell32.SHBrowseForFolder
//sys SHCreateItemFromParsingName(path *uint16, bc unsafe.Pointer, iid uintptr, item **IShellItem) (err error) = shell32.SHCreateItemFromParsingName
//sys SHGetPathFromIDListEx(ptr uintptr, path *uint16, pathLen int, opts int) (err error) = shell32.SHGetPathFromIDListEx
//sys ShellNotifyIcon(message uint32, data *NOTIFYICONDATA) (ret int, err error) = shell32.Shell_NotifyIconW
