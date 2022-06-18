package zenity

import (
	"fmt"
	"path/filepath"
	"reflect"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

const (
	_FOS_NOCHANGEDIR      = 0x00000008
	_FOS_PICKFOLDERS      = 0x00000020
	_FOS_FORCEFILESYSTEM  = 0x00000040
	_FOS_ALLOWMULTISELECT = 0x00000200
	_FOS_FORCESHOWHIDDEN  = 0x10000000
)

func selectFile(opts options) (string, error) {
	if opts.directory {
		res, _, err := pickFolders(opts, false)
		return res, err
	}

	var args win.OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Owner, _ = opts.attach.(uintptr)
	args.Flags = win.OFN_NOCHANGEDIR | win.OFN_FILEMUSTEXIST | win.OFN_EXPLORER

	if opts.title != nil {
		args.Title = syscall.StringToUTF16Ptr(*opts.title)
	}
	if opts.showHidden {
		args.Flags |= win.OFN_FORCESHOWHIDDEN
	}
	if opts.fileFilters != nil {
		args.Filter = &initFilters(opts.fileFilters)[0]
	}

	var res [32768]uint16
	args.File = &res[0]
	args.MaxFile = uint32(len(res))
	args.InitialDir, args.DefExt = initDirNameExt(opts.filename, res[:])

	defer setup()()

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
		if err != nil {
			return "", err
		}
		defer unhook()
	}

	ok := win.GetOpenFileName(&args)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	if !ok {
		return "", win.CommDlgError()
	}
	return syscall.UTF16ToString(res[:]), nil
}

func selectFileMultiple(opts options) ([]string, error) {
	if opts.directory {
		_, res, err := pickFolders(opts, true)
		return res, err
	}

	var args win.OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Owner, _ = opts.attach.(uintptr)
	args.Flags = win.OFN_NOCHANGEDIR | win.OFN_ALLOWMULTISELECT | win.OFN_FILEMUSTEXIST | win.OFN_EXPLORER

	if opts.title != nil {
		args.Title = syscall.StringToUTF16Ptr(*opts.title)
	}
	if opts.showHidden {
		args.Flags |= win.OFN_FORCESHOWHIDDEN
	}
	if opts.fileFilters != nil {
		args.Filter = &initFilters(opts.fileFilters)[0]
	}

	var res [32768 + 1024*256]uint16
	args.File = &res[0]
	args.MaxFile = uint32(len(res))
	args.InitialDir, args.DefExt = initDirNameExt(opts.filename, res[:])

	defer setup()()

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
		if err != nil {
			return nil, err
		}
		defer unhook()
	}

	ok := win.GetOpenFileName(&args)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}
	if !ok {
		return nil, win.CommDlgError()
	}

	var i int
	var nul bool
	var split []string
	for j, p := range res {
		if p == 0 {
			if nul {
				break
			}
			if i < j {
				split = append(split, string(utf16.Decode(res[i:j])))
			}
			i = j + 1
			nul = true
		} else {
			nul = false
		}
	}
	if len := len(split) - 1; len > 0 {
		base := split[0]
		for i := 0; i < len; i++ {
			split[i] = filepath.Join(base, string(split[i+1]))
		}
		split = split[:len]
	}
	return split, nil
}

func selectFileSave(opts options) (string, error) {
	if opts.directory {
		res, _, err := pickFolders(opts, false)
		return res, err
	}

	var args win.OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Owner, _ = opts.attach.(uintptr)
	args.Flags = win.OFN_NOCHANGEDIR | win.OFN_PATHMUSTEXIST | win.OFN_NOREADONLYRETURN | win.OFN_EXPLORER

	if opts.title != nil {
		args.Title = syscall.StringToUTF16Ptr(*opts.title)
	}
	if opts.confirmOverwrite {
		args.Flags |= win.OFN_OVERWRITEPROMPT
	}
	if opts.confirmCreate {
		args.Flags |= win.OFN_CREATEPROMPT
	}
	if opts.showHidden {
		args.Flags |= win.OFN_FORCESHOWHIDDEN
	}
	if opts.fileFilters != nil {
		args.Filter = &initFilters(opts.fileFilters)[0]
	}

	var res [32768]uint16
	args.File = &res[0]
	args.MaxFile = uint32(len(res))
	args.InitialDir, args.DefExt = initDirNameExt(opts.filename, res[:])

	defer setup()()

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
		if err != nil {
			return "", err
		}
		defer unhook()
	}

	ok := win.GetSaveFileName(&args)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	if !ok {
		return "", win.CommDlgError()
	}
	return syscall.UTF16ToString(res[:]), nil
}

func pickFolders(opts options, multi bool) (str string, lst []string, err error) {
	defer setup()()

	err = win.CoInitializeEx(0, 0x6) // COINIT_APARTMENTTHREADED|COINIT_DISABLE_OLE1DDE
	if err != win.RPC_E_CHANGED_MODE {
		if err != nil {
			return "", nil, err
		}
		defer win.CoUninitialize()
	}

	var dialog *_IFileOpenDialog
	err = win.CoCreateInstance(
		_CLSID_FileOpenDialog, nil, 0x17, // CLSCTX_ALL
		_IID_IFileOpenDialog, unsafe.Pointer(&dialog))
	if err != nil {
		if multi {
			return "", nil, fmt.Errorf("%w: multiple directory", ErrUnsupported)
		}
		return browseForFolder(opts)
	}
	defer dialog.Call(dialog.Release)

	var flgs int
	hr, _, _ := dialog.Call(dialog.GetOptions, uintptr(unsafe.Pointer(&flgs)))
	if int32(hr) < 0 {
		return "", nil, syscall.Errno(hr)
	}
	flgs |= _FOS_NOCHANGEDIR | _FOS_PICKFOLDERS | _FOS_FORCEFILESYSTEM
	if multi {
		flgs |= _FOS_ALLOWMULTISELECT
	}
	if opts.showHidden {
		flgs |= _FOS_FORCESHOWHIDDEN
	}
	hr, _, _ = dialog.Call(dialog.SetOptions, uintptr(flgs))
	if int32(hr) < 0 {
		return "", nil, syscall.Errno(hr)
	}

	if opts.title != nil {
		ptr := syscall.StringToUTF16Ptr(*opts.title)
		dialog.Call(dialog.SetTitle, uintptr(unsafe.Pointer(ptr)))
	}

	if opts.filename != "" {
		var item *win.IShellItem
		ptr := syscall.StringToUTF16Ptr(opts.filename)
		win.SHCreateItemFromParsingName(ptr, nil, _IID_IShellItem, &item)

		if int32(hr) >= 0 && item != nil {
			dialog.Call(dialog.SetFolder, uintptr(unsafe.Pointer(item)))
			item.Call(item.Release)
		}
	}

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
		if err != nil {
			return "", nil, err
		}
		defer unhook()
	}

	owner, _ := opts.attach.(uintptr)
	hr, _, _ = dialog.Call(dialog.Show, owner)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", nil, opts.ctx.Err()
	}
	if hr == 0x800704c7 { // ERROR_CANCELLED
		return "", nil, ErrCanceled
	}
	if int32(hr) < 0 {
		return "", nil, syscall.Errno(hr)
	}

	shellItemPath := func(obj *win.COMObject, trap uintptr, a ...uintptr) error {
		var item *win.IShellItem
		hr, _, _ := obj.Call(trap, append(a, uintptr(unsafe.Pointer(&item)))...)
		if int32(hr) < 0 {
			return syscall.Errno(hr)
		}
		defer item.Call(item.Release)

		var ptr uintptr
		hr, _, _ = item.Call(item.GetDisplayName,
			0x80058000, // SIGDN_FILESYSPATH
			uintptr(unsafe.Pointer(&ptr)))
		if int32(hr) < 0 {
			return syscall.Errno(hr)
		}
		defer win.CoTaskMemFree(ptr)

		var res []uint16
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&res))
		hdr.Data, hdr.Len, hdr.Cap = uintptr(ptr), 32768, 32768
		str = syscall.UTF16ToString(res)
		lst = append(lst, str)
		return nil
	}

	if multi {
		var items *_IShellItemArray
		hr, _, _ = dialog.Call(dialog.GetResults, uintptr(unsafe.Pointer(&items)))
		if int32(hr) < 0 {
			return "", nil, syscall.Errno(hr)
		}
		defer items.Call(items.Release)

		var count uint32
		hr, _, _ = items.Call(items.GetCount, uintptr(unsafe.Pointer(&count)))
		if int32(hr) < 0 {
			return "", nil, syscall.Errno(hr)
		}
		for i := uintptr(0); i < uintptr(count) && err == nil; i++ {
			err = shellItemPath(&items.COMObject, items.GetItemAt, i)
		}
	} else {
		err = shellItemPath(&dialog.COMObject, dialog.GetResult)
	}
	return
}

func browseForFolder(opts options) (string, []string, error) {
	var args win.BROWSEINFO
	args.Owner, _ = opts.attach.(uintptr)
	args.Flags = 0x1 // BIF_RETURNONLYFSDIRS

	if opts.title != nil {
		args.Title = syscall.StringToUTF16Ptr(*opts.title)
	}
	if opts.filename != "" {
		args.LParam = syscall.StringToUTF16Ptr(opts.filename)
		args.CallbackFunc = syscall.NewCallback(browseForFolderCallback)
	}

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
		if err != nil {
			return "", nil, err
		}
		defer unhook()
	}

	ptr := win.SHBrowseForFolder(&args)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", nil, opts.ctx.Err()
	}
	if ptr == 0 {
		return "", nil, ErrCanceled
	}
	defer win.CoTaskMemFree(ptr)

	var res [32768]uint16
	win.SHGetPathFromIDListEx(ptr, &res[0], len(res), 0)

	str := syscall.UTF16ToString(res[:])
	return str, []string{str}, nil
}

func browseForFolderCallback(wnd uintptr, msg uint32, lparam, data uintptr) uintptr {
	if msg == 1 { // BFFM_INITIALIZED
		sendMessage.Call(wnd, 1024+103 /* BFFM_SETSELECTIONW */, 1 /* TRUE */, data)
	}
	return 0
}

func initDirNameExt(filename string, name []uint16) (dir *uint16, ext *uint16) {
	d, n := splitDirAndName(filename)
	e := filepath.Ext(n)
	if n != "" {
		copy(name, syscall.StringToUTF16(n))
	}
	if d != "" {
		dir = syscall.StringToUTF16Ptr(d)
	}
	if len(e) > 1 {
		ext = syscall.StringToUTF16Ptr(e[1:])
	}
	return
}

func initFilters(filters FileFilters) []uint16 {
	filters.simplify()
	filters.name()
	var res []uint16
	for _, f := range filters {
		res = append(res, utf16.Encode([]rune(f.Name))...)
		res = append(res, 0)
		for _, p := range f.Patterns {
			res = append(res, utf16.Encode([]rune(p))...)
			res = append(res, uint16(';'))
		}
		res = append(res, 0)
	}
	if res != nil {
		res = append(res, 0)
	}
	return res
}

// https://github.com/wine-mirror/wine/blob/master/include/shobjidl.idl

var (
	_IID_IShellItem       = uuid("\x1e\x6d\x82\x43\x18\xe7\xee\x42\xbc\x55\xa1\xe2\x61\xc3\x7b\xfe")
	_IID_IFileOpenDialog  = uuid("\x88\x72\x7c\xd5\xad\xd4\x68\x47\xbe\x02\x9d\x96\x95\x32\xd9\x60")
	_CLSID_FileOpenDialog = uuid("\x9c\x5a\x1c\xdc\x8a\xe8\xde\x4d\xa5\xa1\x60\xf8\x2a\x20\xae\xf7")
)

type _IFileOpenDialog struct {
	win.COMObject
	*_IFileOpenDialogVtbl
}

type _IShellItemArray struct {
	win.COMObject
	*_IShellItemArrayVtbl
}

type _IFileOpenDialogVtbl struct {
	_IFileDialogVtbl
	GetResults       uintptr
	GetSelectedItems uintptr
}

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

type _IModalWindowVtbl struct {
	_IUnknownVtbl
	Show uintptr
}

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
