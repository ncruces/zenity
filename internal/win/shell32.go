//go:build windows

package win

import (
	"reflect"
	"syscall"
	"unsafe"
)

const (
	BIF_RETURNONLYFSDIRS = 0x00000001

	BFFM_INITIALIZED  = 1
	BFFM_SETSELECTION = WM_USER + 103

	// ShellNotifyIcon messages
	NIM_ADD        = 0
	NIM_MODIFY     = 1
	NIM_DELETE     = 2
	NIM_SETFOCUS   = 3
	NIM_SETVERSION = 4

	// NOTIFYICONDATA flags
	NIF_MESSAGE  = 0x01
	NIF_ICON     = 0x02
	NIF_TIP      = 0x04
	NIF_STATE    = 0x08
	NIF_INFO     = 0x10
	NIF_GUID     = 0x20
	NIF_REALTIME = 0x40
	NIF_SHOWTIP  = 0x80

	// NOTIFYICONDATA info flags
	NIIF_NONE               = 0x00
	NIIF_INFO               = 0x01
	NIIF_WARNING            = 0x02
	NIIF_ERROR              = 0x03
	NIIF_USER               = 0x04
	NIIF_NOSOUND            = 0x10
	NIIF_LARGE_ICON         = 0x20
	NIIF_RESPECT_QUIET_TIME = 0x80
	NIIF_ICON_MASK          = 0x0F

	// NOTIFYICONDATA state
	NIS_HIDDEN     = 0x1
	NIS_SHAREDICON = 0x2

	// IFileOpenDialog options
	FOS_OVERWRITEPROMPT          = 0x2
	FOS_STRICTFILETYPES          = 0x4
	FOS_NOCHANGEDIR              = 0x8
	FOS_PICKFOLDERS              = 0x20
	FOS_FORCEFILESYSTEM          = 0x40
	FOS_ALLNONSTORAGEITEMS       = 0x80
	FOS_NOVALIDATE               = 0x100
	FOS_ALLOWMULTISELECT         = 0x200
	FOS_PATHMUSTEXIST            = 0x800
	FOS_FILEMUSTEXIST            = 0x1000
	FOS_CREATEPROMPT             = 0x2000
	FOS_SHAREAWARE               = 0x4000
	FOS_NOREADONLYRETURN         = 0x8000
	FOS_NOTESTFILECREATE         = 0x10000
	FOS_HIDEMRUPLACES            = 0x20000
	FOS_HIDEPINNEDPLACES         = 0x40000
	FOS_NODEREFERENCELINKS       = 0x100000
	FOS_OKBUTTONNEEDSINTERACTION = 0x200000
	FOS_DONTADDTORECENT          = 0x2000000
	FOS_FORCESHOWHIDDEN          = 0x10000000
	FOS_DEFAULTNOMINIMODE        = 0x20000000
	FOS_FORCEPREVIEWPANEON       = 0x40000000
	FOS_SUPPORTSTREAMABLEITEMS   = 0x80000000

	// IShellItem.GetDisplayName forms
	SIGDN_NORMALDISPLAY               = 0x00000000
	SIGDN_PARENTRELATIVEPARSING       = ^(^0x18001 + 0x80000000)
	SIGDN_DESKTOPABSOLUTEPARSING      = ^(^0x28000 + 0x80000000)
	SIGDN_PARENTRELATIVEEDITING       = ^(^0x31001 + 0x80000000)
	SIGDN_DESKTOPABSOLUTEEDITING      = ^(^0x4c000 + 0x80000000)
	SIGDN_FILESYSPATH                 = ^(^0x58000 + 0x80000000)
	SIGDN_URL                         = ^(^0x68000 + 0x80000000)
	SIGDN_PARENTRELATIVEFORADDRESSBAR = ^(^0x7c001 + 0x80000000)
	SIGDN_PARENTRELATIVE              = ^(^0x80001 + 0x80000000)
	SIGDN_PARENTRELATIVEFORUI         = ^(^0x94001 + 0x80000000)
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

type IDLIST struct{}

// https://github.com/wine-mirror/wine/blob/master/include/shobjidl.idl

var (
	IID_IShellItem       = uuid("\x1e\x6d\x82\x43\x18\xe7\xee\x42\xbc\x55\xa1\xe2\x61\xc3\x7b\xfe")
	IID_IFileOpenDialog  = uuid("\x88\x72\x7c\xd5\xad\xd4\x68\x47\xbe\x02\x9d\x96\x95\x32\xd9\x60")
	CLSID_FileOpenDialog = uuid("\x9c\x5a\x1c\xdc\x8a\xe8\xde\x4d\xa5\xa1\x60\xf8\x2a\x20\xae\xf7")
)

type IFileOpenDialog struct{ IFileDialog }
type iFileOpenDialogVtbl struct {
	iFileDialogVtbl
	GetResults       uintptr
	GetSelectedItems uintptr
}

func (u *IFileOpenDialog) GetResults() (res *IShellItemArray, err error) {
	vtbl := *(**iFileOpenDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.GetResults, uintptr(unsafe.Pointer(&res)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

type IFileDialog struct{ IModalWindow }
type iFileDialogVtbl struct {
	iModalWindowVtbl
	SetFileTypes        uintptr
	SetFileTypeIndex    uintptr
	GetFileTypeIndex    uintptr
	Advise              uintptr
	Unadvise            uintptr
	SetOptions          uintptr
	GetOptions          uintptr
	SetDefaultFolder    uintptr
	SetFolder           uintptr
	GetFolder           uintptr
	GetCurrentSelection uintptr
	SetFileName         uintptr
	GetFileName         uintptr
	SetTitle            uintptr
	SetOkButtonLabel    uintptr
	SetFileNameLabel    uintptr
	GetResult           uintptr
	AddPlace            uintptr
	SetDefaultExtension uintptr
	Close               uintptr
	SetClientGuid       uintptr
	ClearClientData     uintptr
	SetFilter           uintptr
}

func (u *IFileDialog) SetOptions(fos int) (err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.SetOptions, uintptr(fos))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (u *IFileDialog) GetOptions() (fos int, err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.GetOptions, uintptr(unsafe.Pointer(&fos)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (u *IFileDialog) SetFolder(item *IShellItem) (err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.SetFolder, uintptr(unsafe.Pointer(item)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (u *IFileDialog) SetTitle(title *uint16) (err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.SetTitle, uintptr(unsafe.Pointer(title)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (u *IFileDialog) GetResult() (item *IShellItem, err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.GetResult, uintptr(unsafe.Pointer(&item)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

type IModalWindow struct{ IUnknown }
type iModalWindowVtbl struct {
	iUnknownVtbl
	Show uintptr
}

func (u *IModalWindow) Show(wnd HWND) (err error) {
	vtbl := *(**iModalWindowVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.Show, uintptr(wnd))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

type IShellItem struct{ IUnknown }
type iShellItemVtbl struct {
	iUnknownVtbl
	BindToHandler  uintptr
	GetParent      uintptr
	GetDisplayName uintptr
	GetAttributes  uintptr
	Compare        uintptr
}

func (u *IShellItem) GetDisplayName(name int) (res string, err error) {
	var ptr uintptr
	vtbl := *(**iShellItemVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.GetDisplayName, uintptr(name), uintptr(unsafe.Pointer(&ptr)))
	if hr != 0 {
		err = syscall.Errno(hr)
	} else {
		var buf []uint16
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
		hdr.Data, hdr.Len, hdr.Cap = uintptr(ptr), 32768, 32768
		res = syscall.UTF16ToString(buf)
	}
	return
}

type IShellItemArray struct{ IUnknown }
type iShellItemArrayVtbl struct {
	iUnknownVtbl
	BindToHandler              uintptr
	GetPropertyStore           uintptr
	GetPropertyDescriptionList uintptr
	GetAttributes              uintptr
	GetCount                   uintptr
	GetItemAt                  uintptr
	EnumItems                  uintptr
}

func (u *IShellItemArray) GetCount() (numItems uint32, err error) {
	vtbl := *(**iShellItemArrayVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.GetCount, uintptr(unsafe.Pointer(&numItems)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (u *IShellItemArray) GetItemAt(index uint32) (item *IShellItem, err error) {
	vtbl := *(**iShellItemArrayVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.GetItemAt, uintptr(index), uintptr(unsafe.Pointer(&item)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

//sys ExtractAssociatedIcon(instance Handle, path *uint16, icon *uint16) (ret Handle, err error) = shell32.ExtractAssociatedIconW
//sys SHBrowseForFolder(bi *BROWSEINFO) (ret *IDLIST) = shell32.SHBrowseForFolder
//sys SHCreateItemFromParsingName(path *uint16, bc *IBindCtx, iid uintptr, item **IShellItem) (res error) = shell32.SHCreateItemFromParsingName
//sys ShellNotifyIcon(message uint32, data *NOTIFYICONDATA) (ok bool) = shell32.Shell_NotifyIconW
//sys SHGetPathFromIDListEx(ptr *IDLIST, path *uint16, pathLen int, opts int) (ok bool) = shell32.SHGetPathFromIDListEx
