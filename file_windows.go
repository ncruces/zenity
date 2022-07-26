package zenity

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

func selectFile(opts options) (string, error) {
	if opts.directory {
		res, _, err := pickFolders(opts, false)
		return res, err
	}

	var args win.OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Owner, _ = opts.attach.(win.HWND)
	args.Flags = win.OFN_NOCHANGEDIR | win.OFN_FILEMUSTEXIST | win.OFN_EXPLORER

	if opts.title != nil {
		args.Title = strptr(*opts.title)
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

	defer setup(args.Owner)()
	unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
	if err != nil {
		return "", err
	}
	defer unhook()

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
	args.Owner, _ = opts.attach.(win.HWND)
	args.Flags = win.OFN_NOCHANGEDIR | win.OFN_ALLOWMULTISELECT | win.OFN_FILEMUSTEXIST | win.OFN_EXPLORER

	if opts.title != nil {
		args.Title = strptr(*opts.title)
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

	defer setup(args.Owner)()
	unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
	if err != nil {
		return nil, err
	}
	defer unhook()

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
	args.Owner, _ = opts.attach.(win.HWND)
	args.Flags = win.OFN_NOCHANGEDIR | win.OFN_PATHMUSTEXIST | win.OFN_NOREADONLYRETURN | win.OFN_EXPLORER

	if opts.title != nil {
		args.Title = strptr(*opts.title)
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

	defer setup(args.Owner)()
	unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
	if err != nil {
		return "", err
	}
	defer unhook()

	ok := win.GetSaveFileName(&args)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	if !ok {
		return "", win.CommDlgError()
	}
	return syscall.UTF16ToString(res[:]), nil
}

func pickFolders(opts options, multi bool) (string, []string, error) {
	owner, _ := opts.attach.(win.HWND)
	defer setup(owner)()

	err := win.CoInitializeEx(0, win.COINIT_APARTMENTTHREADED|win.COINIT_DISABLE_OLE1DDE)
	if err != win.RPC_E_CHANGED_MODE {
		if err != nil {
			return "", nil, err
		}
		defer win.CoUninitialize()
	}

	var dialog *win.IFileOpenDialog
	err = win.CoCreateInstance(
		win.CLSID_FileOpenDialog, nil, win.CLSCTX_ALL,
		win.IID_IFileOpenDialog, unsafe.Pointer(&dialog))
	if err != nil {
		if multi {
			return "", nil, fmt.Errorf("%w: multiple directory", ErrUnsupported)
		}
		return browseForFolder(opts)
	}
	defer dialog.Release()

	flgs, err := dialog.GetOptions()
	if err != nil {
		return "", nil, err
	}
	flgs |= win.FOS_NOCHANGEDIR | win.FOS_PICKFOLDERS | win.FOS_FORCEFILESYSTEM
	if multi {
		flgs |= win.FOS_ALLOWMULTISELECT
	}
	if opts.showHidden {
		flgs |= win.FOS_FORCESHOWHIDDEN
	}
	err = dialog.SetOptions(flgs)
	if err != nil {
		return "", nil, err
	}

	if opts.title != nil {
		dialog.SetTitle(strptr(*opts.title))
	}

	if opts.filename != "" {
		var item *win.IShellItem
		win.SHCreateItemFromParsingName(strptr(opts.filename), nil, win.IID_IShellItem, &item)
		if item != nil {
			defer item.Release()
			dialog.SetFolder(item)
		}
	}

	unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
	if err != nil {
		return "", nil, err
	}
	defer unhook()

	err = dialog.Show(owner)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", nil, opts.ctx.Err()
	}
	if err == win.E_CANCELED {
		return "", nil, ErrCanceled
	}
	if err != nil {
		return "", nil, err
	}

	if multi {
		items, err := dialog.GetResults()
		if err != nil {
			return "", nil, err
		}
		defer items.Release()

		count, err := items.GetCount()
		if err != nil {
			return "", nil, err
		}

		var lst []string
		for i := uint32(0); i < count && err == nil; i++ {
			str, err := shellItemPath(items.GetItemAt(i))
			if err != nil {
				return "", nil, err
			}
			lst = append(lst, str)
		}
		return "", lst, nil
	} else {
		str, err := shellItemPath(dialog.GetResult())
		if err != nil {
			return "", nil, err
		}
		return str, nil, nil
	}
}

func shellItemPath(item *win.IShellItem, err error) (string, error) {
	if err != nil {
		return "", err
	}
	defer item.Release()
	return item.GetDisplayName(win.SIGDN_FILESYSPATH)
}

func browseForFolder(opts options) (string, []string, error) {
	var args win.BROWSEINFO
	args.Owner, _ = opts.attach.(win.HWND)
	args.Flags = win.BIF_RETURNONLYFSDIRS

	if opts.title != nil {
		args.Title = strptr(*opts.title)
	}
	if opts.filename != "" {
		args.LParam = strptr(opts.filename)
		args.CallbackFunc = syscall.NewCallback(browseForFolderCallback)
	}

	unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
	if err != nil {
		return "", nil, err
	}
	defer unhook()

	ptr := win.SHBrowseForFolder(&args)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", nil, opts.ctx.Err()
	}
	if ptr == nil {
		return "", nil, ErrCanceled
	}
	defer win.CoTaskMemFree(unsafe.Pointer(ptr))

	var res [32768]uint16
	win.SHGetPathFromIDListEx(ptr, &res[0], len(res), 0)

	str := syscall.UTF16ToString(res[:])
	return str, []string{str}, nil
}

func browseForFolderCallback(wnd win.HWND, msg uint32, lparam, data uintptr) uintptr {
	if msg == win.BFFM_INITIALIZED {
		win.SendMessage(wnd, win.BFFM_SETSELECTION, 1, data)
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
		dir = strptr(d)
	}
	if len(e) > 1 {
		ext = strptr(e[1:])
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
