//go:build windows

package win

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
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
)

// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/ne-shobjidl_core-_fileopendialogoptions
type _FILEOPENDIALOGOPTIONS uint32

const (
	FOS_OVERWRITEPROMPT          _FILEOPENDIALOGOPTIONS = 0x2
	FOS_STRICTFILETYPES          _FILEOPENDIALOGOPTIONS = 0x4
	FOS_NOCHANGEDIR              _FILEOPENDIALOGOPTIONS = 0x8
	FOS_PICKFOLDERS              _FILEOPENDIALOGOPTIONS = 0x20
	FOS_FORCEFILESYSTEM          _FILEOPENDIALOGOPTIONS = 0x40
	FOS_ALLNONSTORAGEITEMS       _FILEOPENDIALOGOPTIONS = 0x80
	FOS_NOVALIDATE               _FILEOPENDIALOGOPTIONS = 0x100
	FOS_ALLOWMULTISELECT         _FILEOPENDIALOGOPTIONS = 0x200
	FOS_PATHMUSTEXIST            _FILEOPENDIALOGOPTIONS = 0x800
	FOS_FILEMUSTEXIST            _FILEOPENDIALOGOPTIONS = 0x1000
	FOS_CREATEPROMPT             _FILEOPENDIALOGOPTIONS = 0x2000
	FOS_SHAREAWARE               _FILEOPENDIALOGOPTIONS = 0x4000
	FOS_NOREADONLYRETURN         _FILEOPENDIALOGOPTIONS = 0x8000
	FOS_NOTESTFILECREATE         _FILEOPENDIALOGOPTIONS = 0x10000
	FOS_HIDEMRUPLACES            _FILEOPENDIALOGOPTIONS = 0x20000
	FOS_HIDEPINNEDPLACES         _FILEOPENDIALOGOPTIONS = 0x40000
	FOS_NODEREFERENCELINKS       _FILEOPENDIALOGOPTIONS = 0x100000
	FOS_OKBUTTONNEEDSINTERACTION _FILEOPENDIALOGOPTIONS = 0x200000
	FOS_DONTADDTORECENT          _FILEOPENDIALOGOPTIONS = 0x2000000
	FOS_FORCESHOWHIDDEN          _FILEOPENDIALOGOPTIONS = 0x10000000
	FOS_DEFAULTNOMINIMODE        _FILEOPENDIALOGOPTIONS = 0x20000000
	FOS_FORCEPREVIEWPANEON       _FILEOPENDIALOGOPTIONS = 0x40000000
	FOS_SUPPORTSTREAMABLEITEMS   _FILEOPENDIALOGOPTIONS = 0x80000000
)

// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/ne-shobjidl_core-sigdn
type SIGDN int32

const (
	SIGDN_NORMALDISPLAY               SIGDN = 0x00000000
	SIGDN_PARENTRELATIVEPARSING       SIGDN = ^(^0x18001 + 0x80000000)
	SIGDN_DESKTOPABSOLUTEPARSING      SIGDN = ^(^0x28000 + 0x80000000)
	SIGDN_PARENTRELATIVEEDITING       SIGDN = ^(^0x31001 + 0x80000000)
	SIGDN_DESKTOPABSOLUTEEDITING      SIGDN = ^(^0x4c000 + 0x80000000)
	SIGDN_FILESYSPATH                 SIGDN = ^(^0x58000 + 0x80000000)
	SIGDN_URL                         SIGDN = ^(^0x68000 + 0x80000000)
	SIGDN_PARENTRELATIVEFORADDRESSBAR SIGDN = ^(^0x7c001 + 0x80000000)
	SIGDN_PARENTRELATIVE              SIGDN = ^(^0x80001 + 0x80000000)
	SIGDN_PARENTRELATIVEFORUI         SIGDN = ^(^0x94001 + 0x80000000)
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

// https://docs.microsoft.com/en-us/windows/win32/api/shtypes/ns-shtypes-comdlg_filterspec
type COMDLG_FILTERSPEC struct {
	Name *uint16
	Spec *uint16
}

// https://docs.microsoft.com/en-us/windows/win32/api/shtypes/ns-shtypes-itemidlist
type ITEMIDLIST struct{}

// https://github.com/wine-mirror/wine/blob/master/include/shobjidl.idl

var (
	IID_IShellItem       = guid("\x1e\x6d\x82\x43\x18\xe7\xee\x42\xbc\x55\xa1\xe2\x61\xc3\x7b\xfe")
	IID_IFileOpenDialog  = guid("\x88\x72\x7c\xd5\xad\xd4\x68\x47\xbe\x02\x9d\x96\x95\x32\xd9\x60")
	IID_IFileSaveDialog  = guid("\x23\xcd\xbc\x84\xde\x5f\xdb\x4c\xae\xa4\xaf\x64\xb8\x3d\x78\xab")
	CLSID_FileOpenDialog = guid("\x9c\x5a\x1c\xdc\x8a\xe8\xde\x4d\xa5\xa1\x60\xf8\x2a\x20\xae\xf7")
	CLSID_FileSaveDialog = guid("\xf3\xe2\xb4\xc0\x21\xba\x73\x47\x8d\xba\x33\x5e\xc9\x46\xeb\x8b")
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

type IFileSaveDialog struct{ IFileDialog }
type IFileSaveDialogVtbl struct {
	iFileDialogVtbl
	SetSaveAsItem          uintptr
	SetProperties          uintptr
	SetCollectedProperties uintptr
	GetProperties          uintptr
	ApplyProperties        uintptr
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

func (u *IFileDialog) SetFileTypes(cFileTypes int, rgFilterSpec *COMDLG_FILTERSPEC) (err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.SetFileTypes, uintptr(cFileTypes), uintptr(unsafe.Pointer(rgFilterSpec)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (u *IFileDialog) SetOptions(fos _FILEOPENDIALOGOPTIONS) (err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.SetOptions, uintptr(fos))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (u *IFileDialog) GetOptions() (fos _FILEOPENDIALOGOPTIONS, err error) {
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

func (u *IFileDialog) SetFileName(name *uint16) (err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.SetFileName, uintptr(unsafe.Pointer(name)))
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

func (u *IFileDialog) SetDefaultExtension(extension *uint16) (err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.SetDefaultExtension, uintptr(unsafe.Pointer(extension)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (u *IFileDialog) Close(res syscall.Errno) (err error) {
	vtbl := *(**iFileDialogVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.Close, uintptr(res))
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

func (u *IShellItem) GetDisplayName(name SIGDN) (res string, err error) {
	var ptr *uint16
	vtbl := *(**iShellItemVtbl)(unsafe.Pointer(u))
	hr, _, _ := u.call(vtbl.GetDisplayName, uintptr(name), uintptr(unsafe.Pointer(&ptr)))
	if hr != 0 {
		err = syscall.Errno(hr)
	} else {
		res = windows.UTF16PtrToString(ptr)
		CoTaskMemFree(unsafe.Pointer(ptr))
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
//sys SHBrowseForFolder(bi *BROWSEINFO) (ret *ITEMIDLIST) = shell32.SHBrowseForFolder
//sys SHCreateItemFromParsingName(path *uint16, bc *IBindCtx, iid *GUID, item **IShellItem) (res error) = shell32.SHCreateItemFromParsingName
//sys ShellNotifyIcon(message uint32, data *NOTIFYICONDATA) (ok bool) = shell32.Shell_NotifyIconW
//sys SHGetPathFromIDListEx(ptr *ITEMIDLIST, path *uint16, pathLen int, opts int) (ok bool) = shell32.SHGetPathFromIDListEx
