package zenity

import (
	"path/filepath"
	"reflect"
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	getOpenFileName             = comdlg32.NewProc("GetOpenFileNameW")
	getSaveFileName             = comdlg32.NewProc("GetSaveFileNameW")
	shBrowseForFolder           = shell32.NewProc("SHBrowseForFolderW")
	shGetPathFromIDListEx       = shell32.NewProc("SHGetPathFromIDListEx")
	shCreateItemFromParsingName = shell32.NewProc("SHCreateItemFromParsingName")
)

func selectFile(opts options) (string, error) {
	if opts.directory {
		res, _, err := pickFolders(opts, false)
		return res, err
	}

	var args _OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Flags = 0x81008 // OFN_NOCHANGEDIR|OFN_FILEMUSTEXIST|OFN_EXPLORER

	if opts.title != nil {
		args.Title = syscall.StringToUTF16Ptr(*opts.title)
	}
	if opts.showHidden {
		args.Flags |= 0x10000000 // OFN_FORCESHOWHIDDEN
	}
	if opts.fileFilters != nil {
		args.Filter = &initFilters(opts.fileFilters)[0]
	}

	res := [32768]uint16{}
	args.File = &res[0]
	args.MaxFile = uint32(len(res))
	args.InitialDir, args.DefExt = initDirNameExt(opts.filename, res[:])

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, nil)
		if err != nil {
			return "", err
		}
		defer unhook()
	}

	activate()
	s, _, _ := getOpenFileName.Call(uintptr(unsafe.Pointer(&args)))
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	if s == 0 {
		return "", commDlgError()
	}
	return syscall.UTF16ToString(res[:]), nil
}

func selectFileMutiple(opts options) ([]string, error) {
	if opts.directory {
		_, res, err := pickFolders(opts, true)
		return res, err
	}

	var args _OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Flags = 0x81208 // OFN_NOCHANGEDIR|OFN_ALLOWMULTISELECT|OFN_FILEMUSTEXIST|OFN_EXPLORER

	if opts.title != nil {
		args.Title = syscall.StringToUTF16Ptr(*opts.title)
	}
	if opts.showHidden {
		args.Flags |= 0x10000000 // OFN_FORCESHOWHIDDEN
	}
	if opts.fileFilters != nil {
		args.Filter = &initFilters(opts.fileFilters)[0]
	}

	res := [32768 + 1024*256]uint16{}
	args.File = &res[0]
	args.MaxFile = uint32(len(res))
	args.InitialDir, args.DefExt = initDirNameExt(opts.filename, res[:])

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, nil)
		if err != nil {
			return nil, err
		}
		defer unhook()
	}

	activate()
	s, _, _ := getOpenFileName.Call(uintptr(unsafe.Pointer(&args)))
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}
	if s == 0 {
		return nil, commDlgError()
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

	var args _OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Flags = 0x88808 // OFN_NOCHANGEDIR|OFN_PATHMUSTEXIST|OFN_NOREADONLYRETURN|OFN_EXPLORER

	if opts.title != nil {
		args.Title = syscall.StringToUTF16Ptr(*opts.title)
	}
	if opts.confirmOverwrite {
		args.Flags |= 0x2 // OFN_OVERWRITEPROMPT
	}
	if opts.confirmCreate {
		args.Flags |= 0x2000 // OFN_CREATEPROMPT
	}
	if opts.showHidden {
		args.Flags |= 0x10000000 // OFN_FORCESHOWHIDDEN
	}
	if opts.fileFilters != nil {
		args.Filter = &initFilters(opts.fileFilters)[0]
	}

	res := [32768]uint16{}
	args.File = &res[0]
	args.MaxFile = uint32(len(res))
	args.InitialDir, args.DefExt = initDirNameExt(opts.filename, res[:])

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, nil)
		if err != nil {
			return "", err
		}
		defer unhook()
	}

	activate()
	s, _, _ := getSaveFileName.Call(uintptr(unsafe.Pointer(&args)))
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	if s == 0 {
		return "", commDlgError()
	}
	return syscall.UTF16ToString(res[:]), nil
}

func pickFolders(opts options, multi bool) (str string, lst []string, err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	hr, _, _ := coInitializeEx.Call(0, 0x6) // COINIT_APARTMENTTHREADED|COINIT_DISABLE_OLE1DDE
	if hr != 0x80010106 {                   // RPC_E_CHANGED_MODE
		if int32(hr) < 0 {
			return "", nil, syscall.Errno(hr)
		}
		defer coUninitialize.Call()
	}

	var dialog *_IFileOpenDialog
	hr, _, _ = coCreateInstance.Call(
		_CLSID_FileOpenDialog, 0, 0x17, // CLSCTX_ALL
		_IID_IFileOpenDialog, uintptr(unsafe.Pointer(&dialog)))
	if int32(hr) < 0 {
		return browseForFolder(opts)
	}
	defer dialog.Call(dialog.vtbl.Release)

	var flgs int
	hr, _, _ = dialog.Call(dialog.vtbl.GetOptions, uintptr(unsafe.Pointer(&flgs)))
	if int32(hr) < 0 {
		return "", nil, syscall.Errno(hr)
	}
	if multi {
		flgs |= 0x200 // FOS_ALLOWMULTISELECT
	}
	if opts.showHidden {
		flgs |= 0x10000000 // FOS_FORCESHOWHIDDEN
	}
	hr, _, _ = dialog.Call(dialog.vtbl.SetOptions, uintptr(flgs|0x68)) // FOS_NOCHANGEDIR|FOS_PICKFOLDERS|FOS_FORCEFILESYSTEM
	if int32(hr) < 0 {
		return "", nil, syscall.Errno(hr)
	}

	if opts.title != nil {
		ptr := syscall.StringToUTF16Ptr(*opts.title)
		dialog.Call(dialog.vtbl.SetTitle, uintptr(unsafe.Pointer(ptr)))
	}

	if opts.filename != "" {
		var item *_IShellItem
		ptr := syscall.StringToUTF16Ptr(opts.filename)
		hr, _, _ = shCreateItemFromParsingName.Call(
			uintptr(unsafe.Pointer(ptr)), 0,
			_IID_IShellItem,
			uintptr(unsafe.Pointer(&item)))

		if int32(hr) >= 0 && item != nil {
			dialog.Call(dialog.vtbl.SetFolder, uintptr(unsafe.Pointer(item)))
			item.Call(item.vtbl.Release)
		}
	}

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, nil)
		if err != nil {
			return "", nil, err
		}
		defer unhook()
	}

	activate()
	hr, _, _ = dialog.Call(dialog.vtbl.Show, 0)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", nil, opts.ctx.Err()
	}
	if hr == 0x800704c7 { // ERROR_CANCELLED
		return "", nil, nil
	}
	if int32(hr) < 0 {
		return "", nil, syscall.Errno(hr)
	}

	shellItemPath := func(obj *_COMObject, trap uintptr, a ...uintptr) error {
		var item *_IShellItem
		hr, _, _ := obj.Call(trap, append(a, uintptr(unsafe.Pointer(&item)))...)
		if int32(hr) < 0 {
			return syscall.Errno(hr)
		}
		defer item.Call(item.vtbl.Release)

		var ptr uintptr
		hr, _, _ = item.Call(item.vtbl.GetDisplayName,
			0x80058000, // SIGDN_FILESYSPATH
			uintptr(unsafe.Pointer(&ptr)))
		if int32(hr) < 0 {
			return syscall.Errno(hr)
		}
		defer coTaskMemFree.Call(ptr)

		var res []uint16
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&res))
		hdr.Data, hdr.Len, hdr.Cap = ptr, 32768, 32768
		str = syscall.UTF16ToString(res)
		lst = append(lst, str)
		return nil
	}

	if multi {
		var items *_IShellItemArray
		hr, _, _ = dialog.Call(dialog.vtbl.GetResults, uintptr(unsafe.Pointer(&items)))
		if int32(hr) < 0 {
			return "", nil, syscall.Errno(hr)
		}
		defer items.Call(items.vtbl.Release)

		var count uint32
		hr, _, _ = items.Call(items.vtbl.GetCount, uintptr(unsafe.Pointer(&count)))
		if int32(hr) < 0 {
			return "", nil, syscall.Errno(hr)
		}
		for i := uintptr(0); i < uintptr(count) && err == nil; i++ {
			err = shellItemPath(&items._COMObject, items.vtbl.GetItemAt, i)
		}
	} else {
		err = shellItemPath(&dialog._COMObject, dialog.vtbl.GetResult)
	}
	return
}

func browseForFolder(opts options) (string, []string, error) {
	var args _BROWSEINFO
	args.Flags = 0x1 // BIF_RETURNONLYFSDIRS

	if opts.title != nil {
		args.Title = syscall.StringToUTF16Ptr(*opts.title)
	}
	if opts.filename != "" {
		ptr := syscall.StringToUTF16Ptr(opts.filename)
		args.LParam = uintptr(unsafe.Pointer(ptr))
		args.CallbackFunc = syscall.NewCallback(func(wnd uintptr, msg uint32, lparam, data uintptr) uintptr {
			if msg == 1 { // BFFM_INITIALIZED
				sendMessage.Call(wnd, 1024+103 /* BFFM_SETSELECTIONW */, 1 /* TRUE */, data)
			}
			return 0
		})
	}

	if opts.ctx != nil {
		unhook, err := hookDialog(opts.ctx, nil)
		if err != nil {
			return "", nil, err
		}
		defer unhook()
	}

	activate()
	ptr, _, _ := shBrowseForFolder.Call(uintptr(unsafe.Pointer(&args)))
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", nil, opts.ctx.Err()
	}
	if ptr == 0 {
		return "", nil, nil
	}
	defer coTaskMemFree.Call(ptr)

	res := [32768]uint16{}
	shGetPathFromIDListEx.Call(ptr, uintptr(unsafe.Pointer(&res[0])), uintptr(len(res)), 0)

	str := syscall.UTF16ToString(res[:])
	return str, []string{str}, nil
}

func initDirNameExt(filename string, name []uint16) (dir *uint16, ext *uint16) {
	d, n := splitDirAndName(filename)
	e := filepath.Ext(n)
	if n != "" {
		copy(name, syscall.StringToUTF16(filename))
	}
	if d != "" {
		dir = syscall.StringToUTF16Ptr(d)
	}
	if len(e) > 1 {
		ext = syscall.StringToUTF16Ptr(e[1:])
	}
	return
}

func initFilters(filters []FileFilter) []uint16 {
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

type _OPENFILENAME struct {
	StructSize      uint32
	Owner           uintptr
	Instance        uintptr
	Filter          *uint16
	CustomFilter    *uint16
	MaxCustomFilter uint32
	FilterIndex     uint32
	File            *uint16
	MaxFile         uint32
	FileTitle       *uint16
	MaxFileTitle    uint32
	InitialDir      *uint16
	Title           *uint16
	Flags           uint32
	FileOffset      uint16
	FileExtension   uint16
	DefExt          *uint16
	CustData        uintptr
	FnHook          uintptr
	TemplateName    *uint16
	PvReserved      uintptr
	DwReserved      uint32
	FlagsEx         uint32
}

type _BROWSEINFO struct {
	Owner        uintptr
	Root         uintptr
	DisplayName  *uint16
	Title        *uint16
	Flags        uint32
	CallbackFunc uintptr
	LParam       uintptr
	Image        int32
}

func uuid(s string) uintptr {
	return (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
}

var (
	_IID_IShellItem       = uuid("\x1e\x6d\x82\x43\x18\xe7\xee\x42\xbc\x55\xa1\xe2\x61\xc3\x7b\xfe")
	_IID_IFileOpenDialog  = uuid("\x88\x72\x7c\xd5\xad\xd4\x68\x47\xbe\x02\x9d\x96\x95\x32\xd9\x60")
	_CLSID_FileOpenDialog = uuid("\x9c\x5a\x1c\xdc\x8a\xe8\xde\x4d\xa5\xa1\x60\xf8\x2a\x20\xae\xf7")
)

type _IFileOpenDialog struct {
	_COMObject
	vtbl *_IFileOpenDialogVtbl
}

type _IShellItem struct {
	_COMObject
	vtbl *_IShellItemVtbl
}

type _IShellItemArray struct {
	_COMObject
	vtbl *_IShellItemArrayVtbl
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

type _IShellItemVtbl struct {
	_IUnknownVtbl
	BindToHandler  uintptr
	GetParent      uintptr
	GetDisplayName uintptr
	GetAttributes  uintptr
	Compare        uintptr
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
