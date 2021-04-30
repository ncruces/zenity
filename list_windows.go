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

func listDlg(text string, items []string, multiple bool, opts options) (out []string, err error) {
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
	font := getFont()
	defer font.Delete()
	defWindowProc := defWindowProc.Addr()

	var wnd, textCtl, listCtl uintptr
	var okBtn, cancelBtn, extraBtn uintptr

	layout := func(dpi dpi) {
		hfont := font.ForDPI(dpi)
		sendMessage.Call(textCtl, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(listCtl, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(okBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(cancelBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(extraBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
		setWindowPos.Call(wnd, 0, 0, 0, dpi.Scale(281), dpi.Scale(281), 0x6)                             // SWP_NOZORDER|SWP_NOMOVE
		setWindowPos.Call(textCtl, 0, dpi.Scale(12), dpi.Scale(10), dpi.Scale(241), dpi.Scale(16), 0x4)  // SWP_NOZORDER
		setWindowPos.Call(listCtl, 0, dpi.Scale(12), dpi.Scale(30), dpi.Scale(241), dpi.Scale(164), 0x4) // SWP_NOZORDER
		if extraBtn == 0 {
			setWindowPos.Call(okBtn, 0, dpi.Scale(95), dpi.Scale(206), dpi.Scale(75), dpi.Scale(24), 0x4)      // SWP_NOZORDER
			setWindowPos.Call(cancelBtn, 0, dpi.Scale(178), dpi.Scale(206), dpi.Scale(75), dpi.Scale(24), 0x4) // SWP_NOZORDER
		} else {
			setWindowPos.Call(okBtn, 0, dpi.Scale(12), dpi.Scale(206), dpi.Scale(75), dpi.Scale(24), 0x4)      // SWP_NOZORDER
			setWindowPos.Call(extraBtn, 0, dpi.Scale(95), dpi.Scale(206), dpi.Scale(75), dpi.Scale(24), 0x4)   // SWP_NOZORDER
			setWindowPos.Call(cancelBtn, 0, dpi.Scale(178), dpi.Scale(206), dpi.Scale(75), dpi.Scale(24), 0x4) // SWP_NOZORDER
		}
	}

	proc := func(wnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
		switch msg {
		case 0x0002: // WM_DESTROY
			postQuitMessage.Call(0)

		case 0x0010: // WM_CLOSE
			err = ErrCanceled
			destroyWindow.Call(wnd)

		case 0x0111: // WM_COMMAND
			switch wparam {
			default:
				return 1
			case 1, 6: // IDOK, IDYES
				if multiple {
					if len, _, _ := sendMessage.Call(listCtl, 0x190 /* LB_GETSELCOUNT */, 0, 0); int32(len) >= 0 {
						out = make([]string, len)
						if len > 0 {
							indices := make([]int32, len)
							sendMessage.Call(listCtl, 0x191 /* LB_GETSELITEMS */, len, uintptr(unsafe.Pointer(&indices[0])))
							for i, idx := range indices {
								out[i] = items[idx]
							}
						}
					}
				} else {
					if idx, _, _ := sendMessage.Call(listCtl, 0x188 /* LB_GETCURSEL */, 0, 0); int32(idx) >= 0 {
						out = []string{items[idx]}
					} else {
						out = []string{}
					}
				}
			case 2: // IDCANCEL
				err = ErrCanceled
			case 7: // IDNO
				err = ErrExtraButton
			}
			destroyWindow.Call(wnd)

		case 0x02e0: // WM_DPICHANGED
			layout(dpi(uint32(wparam) >> 16))

		default:
			res, _, _ := syscall.Syscall6(defWindowProc, 4, wnd, uintptr(msg), wparam, lparam, 0, 0)
			return res
		}

		return 0
	}

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return nil, err
	}

	cls, err := registerClass(instance, syscall.NewCallback(proc))
	if cls == 0 {
		return nil, err
	}
	defer unregisterClass.Call(cls, instance)

	wnd, _, _ = createWindowEx.Call(0x10101, // WS_EX_CONTROLPARENT|WS_EX_WINDOWEDGE|WS_EX_DLGMODALFRAME
		cls, strptr(*opts.title),
		0x84c80000, // WS_POPUPWINDOW|WS_CLIPSIBLINGS|WS_DLGFRAME
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		281, 281, 0, 0, instance, 0)

	textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr(text),
		0x5002e080, // WS_CHILD|WS_VISIBLE|WS_GROUP|SS_WORDELLIPSIS|SS_EDITCONTROL|SS_NOPREFIX
		12, 10, 241, 16, wnd, 0, instance, 0)

	var flags uintptr = 0x50320000 // WS_CHILD|WS_VISIBLE|WS_VSCROLL|WS_GROUP|WS_TABSTOP
	if multiple {
		flags |= 0x0800 // LBS_EXTENDEDSEL
	}
	listCtl, _, _ = createWindowEx.Call(0x200, // WS_EX_CLIENTEDGE
		strptr("LISTBOX"), strptr(opts.entryText),
		flags,
		12, 30, 241, 164, wnd, 0, instance, 0)

	okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		0x50030001, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|BS_DEFPUSHBUTTON
		12, 206, 75, 24, wnd, 1 /* IDOK */, instance, 0)
	cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
		12, 206, 75, 24, wnd, 2 /* IDCANCEL */, instance, 0)
	if opts.extraButton != nil {
		extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
			12, 206, 75, 24, wnd, 7 /* IDNO */, instance, 0)
	}

	for _, item := range items {
		sendMessage.Call(listCtl, 0x180 /* LB_ADDSTRING */, 0, strptr(item))
	}

	layout(getDPI(wnd))
	centerWindow(wnd)
	setFocus.Call(listCtl)
	showWindow.Call(wnd, 1 /* SW_SHOWNORMAL */, 0)

	if opts.ctx != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				sendMessage.Call(wnd, 0x0112 /* WM_SYSCOMMAND */, 0xf060 /* SC_CLOSE */, 0)
			case <-wait:
			}
		}()
	}

	// set default values
	out, err = nil, nil

	if err := messageLoop(wnd); err != nil {
		return nil, err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}
	return out, err
}
