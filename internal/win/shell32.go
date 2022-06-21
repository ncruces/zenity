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

// https://github.com/wine-mirror/wine/blob/master/include/shobjidl.idl

var (
	IID_IShellItem       = uuid("\x1e\x6d\x82\x43\x18\xe7\xee\x42\xbc\x55\xa1\xe2\x61\xc3\x7b\xfe")
	IID_IFileOpenDialog  = uuid("\x88\x72\x7c\xd5\xad\xd4\x68\x47\xbe\x02\x9d\x96\x95\x32\xd9\x60")
	CLSID_FileOpenDialog = uuid("\x9c\x5a\x1c\xdc\x8a\xe8\xde\x4d\xa5\xa1\x60\xf8\x2a\x20\xae\xf7")
)

type IFileOpenDialog struct {
	_IFileOpenDialogBase
	_ *_IFileOpenDialogVtbl
}

type _IFileOpenDialogBase struct{ _IFileDialogBase }
type _IFileOpenDialogVtbl struct {
	_IFileDialogVtbl
	GetResults       uintptr
	GetSelectedItems uintptr
}

func (c *_IFileOpenDialogBase) GetResults() (res *IShellItemArray, err error) {
	vtbl := (*(**_IFileOpenDialogVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.GetResults, uintptr(unsafe.Pointer(&res)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

type _IFileDialogBase struct{ _IModalWindowBase }
type _IFileDialogVtbl struct {
	_IModalWindowVtbl
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

func (c *_IFileDialogBase) SetOptions(fos int) (err error) {
	vtbl := (*(**_IFileDialogVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.SetOptions, uintptr(fos))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (c *_IFileDialogBase) GetOptions() (fos int, err error) {
	vtbl := (*(**_IFileDialogVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.GetOptions, uintptr(unsafe.Pointer(&fos)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (c *_IFileDialogBase) SetFolder(item *IShellItem) (err error) {
	vtbl := (*(**_IFileDialogVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.SetFolder, uintptr(unsafe.Pointer(item)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (c *_IFileDialogBase) SetTitle(title *uint16) (err error) {
	vtbl := (*(**_IFileDialogVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.SetTitle, uintptr(unsafe.Pointer(title)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (c *_IFileDialogBase) GetResult() (item *IShellItem, err error) {
	vtbl := (*(**_IFileDialogVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.GetResult, uintptr(unsafe.Pointer(&item)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

type _IModalWindowBase struct{ _IUnknownBase }
type _IModalWindowVtbl struct {
	_IUnknownVtbl
	Show uintptr
}

func (c *_IModalWindowBase) Show(wnd HWND) (err error) {
	vtbl := (*(**_IModalWindowVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.Show, uintptr(wnd))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

type IShellItem struct {
	_IShellItemBase
	_ *_IShellItemVtbl
}

type _IShellItemBase struct{ _IUnknownBase }
type _IShellItemVtbl struct {
	_IUnknownVtbl
	BindToHandler  uintptr
	GetParent      uintptr
	GetDisplayName uintptr
	GetAttributes  uintptr
	Compare        uintptr
}

func (c *_IShellItemBase) GetDisplayName(name int) (res string, err error) {
	var ptr uintptr
	vtbl := (*(**_IShellItemVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.GetDisplayName, uintptr(name), uintptr(unsafe.Pointer(&ptr)))
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

type IShellItemArray struct {
	_IShellItemArrayBase
	_ *_IShellItemArrayVtbl
}

type _IShellItemArrayBase struct{ _IUnknownBase }
type _IShellItemArrayVtbl struct {
	_IUnknownVtbl
	BindToHandler              uintptr
	GetPropertyStore           uintptr
	GetPropertyDescriptionList uintptr
	GetAttributes              uintptr
	GetCount                   uintptr
	GetItemAt                  uintptr
	EnumItems                  uintptr
}

func (c *_IShellItemArrayBase) GetCount() (numItems uint32, err error) {
	vtbl := (*(**_IShellItemArrayVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.GetCount, uintptr(unsafe.Pointer(&numItems)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

func (c *_IShellItemArrayBase) GetItemAt(index uint32) (item *IShellItem, err error) {
	vtbl := (*(**_IShellItemArrayVtbl)(unsafe.Pointer(c)))
	hr, _, _ := c.Call(vtbl.GetItemAt, uintptr(index), uintptr(unsafe.Pointer(&item)))
	if hr != 0 {
		err = syscall.Errno(hr)
	}
	return
}

//sys SHBrowseForFolder(bi *BROWSEINFO) (ret Pointer) = shell32.SHBrowseForFolder
//sys SHCreateItemFromParsingName(path *uint16, bc *COMObject, iid uintptr, item **IShellItem) (res error) = shell32.SHCreateItemFromParsingName
//sys ShellNotifyIcon(message uint32, data *NOTIFYICONDATA) (ret int, err error) = shell32.Shell_NotifyIconW
//sys SHGetPathFromIDListEx(ptr Pointer, path *uint16, pathLen int, opts int) (ok bool) = shell32.SHGetPathFromIDListEx
