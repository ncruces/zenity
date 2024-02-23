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
		opts.title = ptr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = ptr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = ptr("Cancel")
	}

	dlg := &listDialog{
		items:         items,
		multiple:      multiple,
		disallowEmpty: opts.disallowEmpty,
	}
	return dlg.setup(text, opts)
}

type listDialog struct {
	items         []string
	multiple      bool
	disallowEmpty bool
	out           []string
	err           error

	wnd       win.HWND
	textCtl   win.HWND
	listCtl   win.HWND
	okBtn     win.HWND
	cancelBtn win.HWND
	extraBtn  win.HWND
	font      font
}

func (dlg *listDialog) setup(text string, opts options) ([]string, error) {
	owner, _ := opts.attach.(win.HWND)
	defer setup(owner)()
	dlg.font = getFont()
	defer dlg.font.delete()
	icon, _ := getIcon(opts.windowIcon)
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

	dlg.wnd, _ = win.CreateWindowEx(_WS_EX_ZEN_DIALOG,
		cls, strptr(*opts.title), _WS_ZEN_DIALOG,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		281, 281, owner, 0, instance, unsafe.Pointer(dlg))

	dlg.textCtl, _ = win.CreateWindowEx(0,
		strptr("STATIC"), strptr(text), _WS_ZEN_LABEL,
		12, 10, 241, 16, dlg.wnd, 0, instance, nil)

	var flags uint32 = _WS_ZEN_CONTROL | win.WS_VSCROLL | win.LBS_NOTIFY
	if dlg.multiple {
		flags |= win.LBS_EXTENDEDSEL
	}
	dlg.listCtl, _ = win.CreateWindowEx(win.WS_EX_CLIENTEDGE,
		strptr("LISTBOX"), nil, flags,
		12, 30, 241, 164, dlg.wnd, 0, instance, nil)

	dlg.okBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(quoteAccelerators(*opts.okLabel)),
		_WS_ZEN_BUTTON|win.BS_DEFPUSHBUTTON,
		12, 206, 75, 24, dlg.wnd, win.IDOK, instance, nil)
	dlg.cancelBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(quoteAccelerators(*opts.cancelLabel)),
		_WS_ZEN_BUTTON,
		12, 206, 75, 24, dlg.wnd, win.IDCANCEL, instance, nil)
	if opts.extraButton != nil {
		dlg.extraBtn, _ = win.CreateWindowEx(0,
			strptr("BUTTON"), strptr(quoteAccelerators(*opts.extraButton)),
			_WS_ZEN_BUTTON,
			12, 206, 75, 24, dlg.wnd, win.IDNO, instance, nil)
	}

	for i, item := range dlg.items {
		win.SendMessagePointer(dlg.listCtl, win.LB_ADDSTRING, 0, unsafe.Pointer(strptr(item)))
		for _, def := range opts.defaultItems {
			if def == item {
				if dlg.multiple {
					win.SendMessage(dlg.listCtl, win.LB_SETSEL, 1, uintptr(i))
				} else {
					win.SendMessage(dlg.listCtl, win.LB_SETCURSEL, uintptr(i), 0)
				}
			}
		}
	}

	dlg.update()
	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	win.SetFocus(dlg.listCtl)
	win.ShowWindow(dlg.wnd, win.SW_NORMAL)

	if opts.ctx != nil && opts.ctx.Done() != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				win.SendMessage(dlg.wnd, win.WM_SYSCOMMAND, win.SC_CLOSE, 0)
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
	win.SetWindowPos(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(281), win.SWP_NOMOVE|win.SWP_NOZORDER)
	win.SetWindowPos(dlg.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), win.SWP_NOZORDER)
	win.SetWindowPos(dlg.listCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(164), win.SWP_NOZORDER)
	if dlg.extraBtn == 0 {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
	} else {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(12), dpi.scale(206), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
	}
}

func (dlg *listDialog) update() {
	if dlg.disallowEmpty {
		var enable bool
		if dlg.multiple {
			len := win.SendMessage(dlg.listCtl, win.LB_GETSELCOUNT, 0, 0)
			enable = int32(len) > 0
		} else {
			idx := win.SendMessage(dlg.listCtl, win.LB_GETCURSEL, 0, 0)
			enable = int32(idx) >= 0
		}
		win.EnableWindow(dlg.okBtn, enable)
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
			dlg.update()
			return 1
		case win.IDOK, win.IDYES:
			if dlg.multiple {
				if len := win.SendMessage(dlg.listCtl, win.LB_GETSELCOUNT, 0, 0); int32(len) >= 0 {
					dlg.out = make([]string, len)
					if len > 0 {
						indices := make([]int32, len)
						win.SendMessagePointer(dlg.listCtl, win.LB_GETSELITEMS, len, unsafe.Pointer(&indices[0]))
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
