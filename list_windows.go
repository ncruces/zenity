package zenity

import (
	"syscall"
	"unsafe"
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

	defer setup()()
	wnd := &listWnd{
		items:    items,
		multiple: multiple,
		font:     getFont(),
	}
	defer wnd.font.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return nil, err
	}

	cls, err := registerClass(instance, syscall.NewCallback(listProc))
	if cls == 0 {
		return nil, err
	}
	defer unregisterClass.Call(cls, instance)

	wnd.handle, _, _ = createWindowEx.Call(0x10101, // WS_EX_CONTROLPARENT|WS_EX_WINDOWEDGE|WS_EX_DLGMODALFRAME
		cls, strptr(*opts.title),
		0x84c80000, // WS_POPUPWINDOW|WS_CLIPSIBLINGS|WS_DLGFRAME
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		281, 281, 0, 0, instance, uintptr(unsafe.Pointer(wnd)))

	wnd.textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr(text),
		0x5002e080, // WS_CHILD|WS_VISIBLE|WS_GROUP|SS_WORDELLIPSIS|SS_EDITCONTROL|SS_NOPREFIX
		12, 10, 241, 16, wnd.handle, 0, instance, 0)

	var flags uintptr = 0x50320000 // WS_CHILD|WS_VISIBLE|WS_VSCROLL|WS_GROUP|WS_TABSTOP
	if multiple {
		flags |= 0x0800 // LBS_EXTENDEDSEL
	}
	wnd.listCtl, _, _ = createWindowEx.Call(0x200, // WS_EX_CLIENTEDGE
		strptr("LISTBOX"), strptr(opts.entryText),
		flags,
		12, 30, 241, 164, wnd.handle, 0, instance, 0)

	wnd.okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		0x50030001, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|BS_DEFPUSHBUTTON
		12, 206, 75, 24, wnd.handle, 1 /* IDOK */, instance, 0)
	wnd.cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
		12, 206, 75, 24, wnd.handle, 2 /* IDCANCEL */, instance, 0)
	if opts.extraButton != nil {
		wnd.extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
			12, 206, 75, 24, wnd.handle, 7 /* IDNO */, instance, 0)
	}

	for _, item := range items {
		sendMessage.Call(wnd.listCtl, 0x180 /* LB_ADDSTRING */, 0, strptr(item))
	}

	wnd.layout(getDPI(wnd.handle))
	centerWindow(wnd.handle)
	setFocus.Call(wnd.listCtl)
	showWindow.Call(wnd.handle, 1 /* SW_SHOWNORMAL */, 0)

	if opts.ctx != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				sendMessage.Call(wnd.handle, 0x0112 /* WM_SYSCOMMAND */, 0xf060 /* SC_CLOSE */, 0)
			case <-wait:
			}
		}()
	}

	if err := messageLoop(wnd.handle); err != nil {
		return nil, err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}
	return wnd.out, wnd.err
}

type listWnd struct {
	handle    uintptr
	textCtl   uintptr
	listCtl   uintptr
	okBtn     uintptr
	cancelBtn uintptr
	extraBtn  uintptr
	multiple  bool
	items     []string
	out       []string
	err       error
	font      font
}

func (wnd *listWnd) layout(dpi dpi) {
	font := wnd.font.forDPI(dpi)
	sendMessage.Call(wnd.textCtl, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(wnd.listCtl, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(wnd.okBtn, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(wnd.cancelBtn, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(wnd.extraBtn, 0x0030 /* WM_SETFONT */, font, 1)
	setWindowPos.Call(wnd.handle, 0, 0, 0, dpi.scale(281), dpi.scale(281), 0x6)                          // SWP_NOZORDER|SWP_NOMOVE
	setWindowPos.Call(wnd.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), 0x4)  // SWP_NOZORDER
	setWindowPos.Call(wnd.listCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(164), 0x4) // SWP_NOZORDER
	if wnd.extraBtn == 0 {
		setWindowPos.Call(wnd.okBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), 0x4)      // SWP_NOZORDER
		setWindowPos.Call(wnd.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
	} else {
		setWindowPos.Call(wnd.okBtn, 0, dpi.scale(12), dpi.scale(206), dpi.scale(75), dpi.scale(24), 0x4)      // SWP_NOZORDER
		setWindowPos.Call(wnd.extraBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), 0x4)   // SWP_NOZORDER
		setWindowPos.Call(wnd.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
	}
}

func listProc(hwnd uintptr, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var wnd *listWnd
	switch msg {
	case 0x0081: // WM_NCCREATE
		saveBackRef(hwnd, *lparam)
		wnd = (*listWnd)(*lparam)
	case 0x0082: // WM_NCDESTROY
		deleteBackRef(hwnd)
	default:
		wnd = (*listWnd)(loadBackRef(hwnd))
	}

	switch msg {
	case 0x0002: // WM_DESTROY
		postQuitMessage.Call(0)

	case 0x0010: // WM_CLOSE
		wnd.err = ErrCanceled
		destroyWindow.Call(hwnd)

	case 0x0111: // WM_COMMAND
		switch wparam {
		default:
			return 1
		case 1, 6: // IDOK, IDYES
			if wnd.multiple {
				if len, _, _ := sendMessage.Call(wnd.listCtl, 0x190 /* LB_GETSELCOUNT */, 0, 0); int32(len) >= 0 {
					wnd.out = make([]string, len)
					if len > 0 {
						indices := make([]int32, len)
						sendMessage.Call(wnd.listCtl, 0x191 /* LB_GETSELITEMS */, len, uintptr(unsafe.Pointer(&indices[0])))
						for i, idx := range indices {
							wnd.out[i] = wnd.items[idx]
						}
					}
				}
			} else {
				if idx, _, _ := sendMessage.Call(wnd.listCtl, 0x188 /* LB_GETCURSEL */, 0, 0); int32(idx) >= 0 {
					wnd.out = []string{wnd.items[idx]}
				} else {
					wnd.out = []string{}
				}
			}
		case 2: // IDCANCEL
			wnd.err = ErrCanceled
		case 7: // IDNO
			wnd.err = ErrExtraButton
		}
		destroyWindow.Call(hwnd)

	case 0x02e0: // WM_DPICHANGED
		wnd.layout(dpi(uint32(wparam) >> 16))

	default:
		res, _, _ := syscall.Syscall6(defWindowProc.Addr(), 4, hwnd, uintptr(msg), wparam, uintptr(unsafe.Pointer(lparam)), 0, 0)
		return res
	}

	return 0
}
