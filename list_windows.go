package zenity

import (
	"syscall"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

func list(text string, items []string, opts options) (string, error) {
	items, err := listDlg(text, items, false, opts)
	if len(items) == 1 {
		return items[0], err
	}
	return "", err
}

func listMultiple(text string, items []string, opts options) ([]string, error) {
	return listDlg(text, items, true, opts)
}

func listDlg(text string, items []string, multiple bool, opts options) ([]string, error) {
	if opts.title == nil {
		opts.title = stringPtr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = stringPtr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = stringPtr("Cancel")
	}

	dlg := &listDialog{
		items:    items,
		multiple: multiple,
	}
	return dlg.setup(text, opts)
}

type listDialog struct {
	items    []string
	multiple bool
	out      []string
	err      error

	wnd       win.HWND
	textCtl   win.HWND
	listCtl   win.HWND
	okBtn     win.HWND
	cancelBtn win.HWND
	extraBtn  win.HWND
	font      font
}

func (dlg *listDialog) setup(text string, opts options) ([]string, error) {
	defer setup()()
	dlg.font = getFont()
	defer dlg.font.delete()
	icon := getIcon(opts.windowIcon)
	defer icon.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}

	instance, err := win.GetModuleHandle(nil)
	if err != nil {
		return nil, err
	}

	cls, err := registerClass(instance, icon.handle, syscall.NewCallback(listProc))
	if err != nil {
		return nil, err
	}
	defer win.UnregisterClass(cls, instance)

	owner, _ := opts.attach.(win.HWND)
	dlg.wnd, _ = win.CreateWindowEx(_WS_EX_CONTROLPARENT|_WS_EX_WINDOWEDGE|_WS_EX_DLGMODALFRAME,
		cls, strptr(*opts.title),
		_WS_POPUPWINDOW|_WS_CLIPSIBLINGS|_WS_DLGFRAME,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		281, 281, owner, 0, instance, unsafe.Pointer(dlg))

	dlg.textCtl, _ = win.CreateWindowEx(0,
		strptr("STATIC"), strptr(text),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_SS_WORDELLIPSIS|_SS_EDITCONTROL|_SS_NOPREFIX,
		12, 10, 241, 16, dlg.wnd, 0, instance, nil)

	var flags uint32 = _WS_CHILD | _WS_VISIBLE | _WS_GROUP | _WS_TABSTOP | _WS_VSCROLL
	if dlg.multiple {
		flags |= _LBS_EXTENDEDSEL
	}
	dlg.listCtl, _ = win.CreateWindowEx(_WS_EX_CLIENTEDGE,
		strptr("LISTBOX"), strptr(opts.entryText),
		flags,
		12, 30, 241, 164, dlg.wnd, 0, instance, nil)

	dlg.okBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP|_BS_DEFPUSHBUTTON,
		12, 206, 75, 24, dlg.wnd, win.IDOK, instance, nil)
	dlg.cancelBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
		12, 206, 75, 24, dlg.wnd, win.IDCANCEL, instance, nil)
	if opts.extraButton != nil {
		dlg.extraBtn, _ = win.CreateWindowEx(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
			12, 206, 75, 24, dlg.wnd, win.IDNO, instance, nil)
	}

	for _, item := range dlg.items {
		win.SendMessagePointer(dlg.listCtl, win.LB_ADDSTRING, 0, unsafe.Pointer(strptr(item)))
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	win.SetFocus(dlg.listCtl)
	win.ShowWindow(dlg.wnd, _SW_NORMAL)

	if opts.ctx != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				win.SendMessage(dlg.wnd, win.WM_SYSCOMMAND, _SC_CLOSE, 0)
			case <-wait:
			}
		}()
	}

	if err := win.MessageLoop(win.HWND(dlg.wnd)); err != nil {
		return nil, err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}
	return dlg.out, dlg.err
}

func (dlg *listDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	win.SendMessage(dlg.textCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.listCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.okBtn, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.cancelBtn, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.extraBtn, win.WM_SETFONT, font, 1)
	win.SetWindowPos(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(281), _SWP_NOZORDER|_SWP_NOMOVE)
	win.SetWindowPos(dlg.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), _SWP_NOZORDER)
	win.SetWindowPos(dlg.listCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(164), _SWP_NOZORDER)
	if dlg.extraBtn == 0 {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	} else {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(12), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		win.SetWindowPos(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	}
}

func listProc(wnd win.HWND, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *listDialog
	switch msg {
	case win.WM_NCCREATE:
		saveBackRef(uintptr(wnd), *lparam)
		dlg = (*listDialog)(*lparam)
	case win.WM_NCDESTROY:
		deleteBackRef(uintptr(wnd))
	default:
		dlg = (*listDialog)(loadBackRef(uintptr(wnd)))
	}

	switch msg {
	case win.WM_DESTROY:
		win.PostQuitMessage(0)

	case win.WM_CLOSE:
		dlg.err = ErrCanceled
		win.DestroyWindow(wnd)

	case win.WM_COMMAND:
		switch wparam {
		default:
			return 1
		case win.IDOK, win.IDYES:
			if dlg.multiple {
				if len := win.SendMessage(dlg.listCtl, win.LB_GETSELCOUNT, 0, 0); int32(len) >= 0 {
					dlg.out = make([]string, len)
					if len > 0 {
						indices := make([]int32, len)
						win.SendMessage(dlg.listCtl, win.LB_GETSELITEMS, len, uintptr(unsafe.Pointer(&indices[0])))
						for i, idx := range indices {
							dlg.out[i] = dlg.items[idx]
						}
					}
				}
			} else {
				if idx := win.SendMessage(dlg.listCtl, win.LB_GETCURSEL, 0, 0); int32(idx) >= 0 {
					dlg.out = []string{dlg.items[idx]}
				} else {
					dlg.out = []string{}
				}
			}
		case win.IDCANCEL:
			dlg.err = ErrCanceled
		case win.IDNO:
			dlg.err = ErrExtraButton
		}
		win.DestroyWindow(wnd)

	case win.WM_DPICHANGED:
		dlg.layout(dpi(uint32(wparam) >> 16))

	default:
		return win.DefWindowProc(wnd, msg, wparam, unsafe.Pointer(lparam))
	}

	return 0
}
