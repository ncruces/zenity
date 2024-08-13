package zenity

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

func selectFile(opts options) (string, error) {
	name, _, shown, err := fileOpenDialog(opts, false)
	if shown || opts.ctx != nil && opts.ctx.Err() != nil {
		return name, err
	}
	if opts.directory {
		return browseForFolder(opts)
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
		args.Filter = initFilters(opts.fileFilters)
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
	_, list, shown, err := fileOpenDialog(opts, true)
	if shown || opts.ctx != nil && opts.ctx.Err() != nil {
		return list, err
	}
	if opts.directory {
		return nil, fmt.Errorf("%w: multiple directory", ErrUnsupported)
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
		args.Filter = initFilters(opts.fileFilters)
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
		return selectFile(opts)
	}
	name, shown, err := fileSaveDialog(opts)
	if shown || opts.ctx != nil && opts.ctx.Err() != nil {
		return name, err
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
		args.Filter = initFilters(opts.fileFilters)
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

func fileOpenDialog(opts options, multi bool) (string, []string, bool, error) {
	uninit, err := coInitialize()
	if err != nil {
		return "", nil, false, err
	}
	defer uninit()

	owner, _ := opts.attach.(win.HWND)
	defer setup(owner)()

	var dialog *win.IFileOpenDialog
	err = win.CoCreateInstance(
		win.CLSID_FileOpenDialog, nil, win.CLSCTX_ALL,
		win.IID_IFileOpenDialog, unsafe.Pointer(&dialog))
	if err != nil {
		return "", nil, false, err
	}
	defer dialog.Release()

	flgs, err := dialog.GetOptions()
	if err != nil {
		return "", nil, false, err
	}
	flgs |= win.FOS_NOCHANGEDIR | win.FOS_FILEMUSTEXIST | win.FOS_FORCEFILESYSTEM
	if multi {
		flgs |= win.FOS_ALLOWMULTISELECT
	}
	if opts.directory {
		flgs |= win.FOS_PICKFOLDERS
	}
	if opts.showHidden {
		flgs |= win.FOS_FORCESHOWHIDDEN
	}
	err = dialog.SetOptions(flgs)
	if err != nil {
		return "", nil, false, err
	}

	if opts.title != nil {
		dialog.SetTitle(strptr(*opts.title))
	}
	if opts.fileFilters != nil {
		dialog.SetFileTypes(initFileTypes(opts.fileFilters))
	}

	if opts.filename != "" {
		var item *win.IShellItem
		dir, name, _ := splitDirAndName(opts.filename)
		dialog.SetFileName(strptr(name))
		if ext := filepath.Ext(name); len(ext) > 1 {
			dialog.SetDefaultExtension(strptr(ext[1:]))
		}
		win.SHCreateItemFromParsingName(strptr(dir), nil, win.IID_IShellItem, &item)
		if item != nil {
			defer item.Release()
			dialog.SetFolder(item)
		}
	}

	unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
	if err != nil {
		return "", nil, false, err
	}
	defer unhook()

	if opts.ctx != nil && opts.ctx.Done() != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				dialog.Close(win.E_TIMEOUT)
			case <-wait:
			}
		}()
	}

	err = dialog.Show(owner)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", nil, true, opts.ctx.Err()
	}
	if err == win.E_CANCELED {
		return "", nil, true, ErrCanceled
	}
	if err != nil {
		return "", nil, true, err
	}

	if multi {
		items, err := dialog.GetResults()
		if err != nil {
			return "", nil, true, err
		}
		defer items.Release()

		count, err := items.GetCount()
		if err != nil {
			return "", nil, true, err
		}

		var lst []string
		for i := uint32(0); i < count && err == nil; i++ {
			str, err := shellItemPath(items.GetItemAt(i))
			if err != nil {
				return "", nil, true, err
			}
			lst = append(lst, str)
		}
		return "", lst, true, nil
	} else {
		str, err := shellItemPath(dialog.GetResult())
		if err != nil {
			return "", nil, true, err
		}
		return str, nil, true, nil
	}
}

func fileSaveDialog(opts options) (string, bool, error) {
	uninit, err := coInitialize()
	if err != nil {
		return "", false, err
	}
	defer uninit()

	owner, _ := opts.attach.(win.HWND)
	defer setup(owner)()

	var dialog *win.IFileSaveDialog
	err = win.CoCreateInstance(
		win.CLSID_FileSaveDialog, nil, win.CLSCTX_ALL,
		win.IID_IFileSaveDialog, unsafe.Pointer(&dialog))
	if err != nil {
		return "", false, err
	}
	defer dialog.Release()

	flgs, err := dialog.GetOptions()
	if err != nil {
		return "", false, err
	}
	flgs |= win.FOS_NOCHANGEDIR | win.FOS_PATHMUSTEXIST | win.FOS_NOREADONLYRETURN | win.FOS_FORCEFILESYSTEM
	if opts.confirmOverwrite {
		flgs |= win.FOS_OVERWRITEPROMPT
	}
	if opts.confirmCreate {
		flgs |= win.FOS_CREATEPROMPT
	}
	if opts.showHidden {
		flgs |= win.FOS_FORCESHOWHIDDEN
	}
	err = dialog.SetOptions(flgs)
	if err != nil {
		return "", false, err
	}

	if opts.title != nil {
		dialog.SetTitle(strptr(*opts.title))
	}
	if opts.fileFilters != nil {
		dialog.SetFileTypes(initFileTypes(opts.fileFilters))
	}

	if opts.filename != "" {
		var item *win.IShellItem
		dir, name, _ := splitDirAndName(opts.filename)
		dialog.SetFileName(strptr(name))
		if ext := filepath.Ext(name); len(ext) > 1 {
			dialog.SetDefaultExtension(strptr(ext[1:]))
		}
		win.SHCreateItemFromParsingName(strptr(dir), nil, win.IID_IShellItem, &item)
		if item != nil {
			defer item.Release()
			dialog.SetFolder(item)
		}
	}

	unhook, err := hookDialog(opts.ctx, opts.windowIcon, nil, nil)
	if err != nil {
		return "", false, err
	}
	defer unhook()

	if opts.ctx != nil && opts.ctx.Done() != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				dialog.Close(win.E_TIMEOUT)
			case <-wait:
			}
		}()
	}

	err = dialog.Show(owner)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", true, opts.ctx.Err()
	}
	if err == win.E_CANCELED {
		return "", true, ErrCanceled
	}
	if err != nil {
		return "", true, err
	}

	str, err := shellItemPath(dialog.GetResult())
	if err != nil {
		return "", true, err
	}
	return str, true, nil
}

func shellItemPath(item *win.IShellItem, err error) (string, error) {
	if err != nil {
		return "", err
	}
	defer item.Release()
	return item.GetDisplayName(win.SIGDN_FILESYSPATH)
}

func browseForFolder(opts options) (string, error) {
	uninit, err := coInitialize()
	if err != nil {
		return "", err
	}
	defer uninit()

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
		return "", err
	}
	defer unhook()

	ptr := win.SHBrowseForFolder(&args)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	if ptr == nil {
		return "", ErrCanceled
	}
	defer win.CoTaskMemFree(unsafe.Pointer(ptr))

	var res [32768]uint16
	win.SHGetPathFromIDListEx(ptr, &res[0], len(res), 0)

	str := syscall.UTF16ToString(res[:])
	return str, nil
}

func browseForFolderCallback(wnd win.HWND, msg uint32, lparam, data uintptr) uintptr {
	if msg == win.BFFM_INITIALIZED {
		win.SendMessage(wnd, win.BFFM_SETSELECTION, 1, data)
	}
	return 0
}

func coInitialize() (context.CancelFunc, error) {
	runtime.LockOSThread()
	// .NET uses MTA for all background threads, so do the same.
	// If someone needs STA because they're doing UI,
	// they should initialize COM themselves before.
	err := win.CoInitializeEx(0, win.COINIT_MULTITHREADED|win.COINIT_DISABLE_OLE1DDE)
	if err == win.S_FALSE {
		// COM was already initialized, we simply increased the ref count.
		// Make this a no-op by decreasing our ref count.
		win.CoUninitialize()
		return runtime.UnlockOSThread, nil
	}
	// Don't uninitialize COM; this is against the docs, but it's what .NET does.
	// Eventually all threads will have COM initialized.
	if err == nil || err == win.RPC_E_CHANGED_MODE {
		return runtime.UnlockOSThread, nil
	}
	runtime.UnlockOSThread()
	return nil, err
}

func initDirNameExt(filename string, name []uint16) (dir *uint16, ext *uint16) {
	d, n, _ := splitDirAndName(filename)
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

func initFilters(filters FileFilters) *uint16 {
	filters.simplify()
	filters.name()
	var res []uint16
	for _, f := range filters {
		if len(f.Patterns) == 0 {
			continue
		}
		res = append(res, syscall.StringToUTF16(f.Name)...)
		for _, p := range f.Patterns {
			res = append(res, syscall.StringToUTF16(p)...)
			res[len(res)-1] = ';'
		}
		res = append(res, 0)
	}
	if res != nil {
		res = append(res, 0)
		return &res[0]
	}
	return nil
}

func initFileTypes(filters FileFilters) (int, *win.COMDLG_FILTERSPEC) {
	filters.simplify()
	filters.name()
	var res []win.COMDLG_FILTERSPEC
	for _, f := range filters {
		if len(f.Patterns) == 0 {
			continue
		}
		var spec []uint16
		for i, p := range f.Patterns {
			spec = append(spec, syscall.StringToUTF16(p)...)
			if i != len(f.Patterns)-1 {
				spec[len(spec)-1] = ';'
			}
		}
		res = append(res, win.COMDLG_FILTERSPEC{
			Name: syscall.StringToUTF16Ptr(f.Name),
			Spec: &spec[0],
		})
	}
	return len(res), &res[0]
}
