//go:build windows

package win

const (
	BIF_RETURNONLYFSDIRS = 0x00000001

	BFFM_INITIALIZED  = 1
	BFFM_SETSELECTION = WM_USER + 103

	NIM_ADD    = 0
	NIM_DELETE = 2
)

// https://docs.microsoft.com/en-us/windows/win32/api/shlobj_core/ns-shlobj_core-browseinfow
type BROWSEINFO struct {
	Owner        HWND
	Root         Pointer
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
	Wnd             HWND
	ID              uint32
	Flags           uint32
	CallbackMessage uint32
	Icon            Handle
	Tip             [128]uint16 // NOTIFYICONDATAA_V1_SIZE
	State           uint32
	StateMask       uint32
	Info            [256]uint16
	Version         uint32
	InfoTitle       [64]uint16
	InfoFlags       uint32
	// GuidItem     [16]byte // NOTIFYICONDATAA_V2_SIZE
	// BalloonIcon  Handle   // NOTIFYICONDATAA_V3_SIZE
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

//sys SHBrowseForFolder(bi *BROWSEINFO) (ret Pointer) = shell32.SHBrowseForFolder
//sys SHCreateItemFromParsingName(path *uint16, bc *COMObject, iid uintptr, item **IShellItem) (res error) = shell32.SHCreateItemFromParsingName
//sys ShellNotifyIcon(message uint32, data *NOTIFYICONDATA) (ret int, err error) = shell32.Shell_NotifyIconW
//sys SHGetPathFromIDListEx(ptr Pointer, path *uint16, pathLen int, opts int) (ok bool) = shell32.SHGetPathFromIDListEx
