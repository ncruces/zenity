package zenity

import (
	"syscall"
)

func entry(text string, opts options) (out string, ok bool, err error) {
	var title string
	if opts.title != nil {
		title = *opts.title
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

	var wnd, textCtl, editCtl uintptr
	var okBtn, cancelBtn, extraBtn uintptr

	layout := func(dpi dpi) {
		hfont := font.ForDPI(dpi)
		sendMessage.Call(textCtl, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(editCtl, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(okBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
		sendMessage.Call(cancelBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
		setWindowPos.Call(wnd, 0, 0, 0, dpi.Scale(281), dpi.Scale(140), 0x6)                            // SWP_NOZORDER|SWP_NOMOVE
		setWindowPos.Call(textCtl, 0, dpi.Scale(12), dpi.Scale(10), dpi.Scale(241), dpi.Scale(16), 0x4) // SWP_NOZORDER
		setWindowPos.Call(editCtl, 0, dpi.Scale(12), dpi.Scale(30), dpi.Scale(241), dpi.Scale(24), 0x4) // SWP_NOZORDER
		if extraBtn == 0 {
			setWindowPos.Call(okBtn, 0, dpi.Scale(95), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4)      // SWP_NOZORDER
			setWindowPos.Call(cancelBtn, 0, dpi.Scale(178), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4) // SWP_NOZORDER
		} else {
			sendMessage.Call(extraBtn, 0x0030 /* WM_SETFONT */, hfont, 1)
			setWindowPos.Call(okBtn, 0, dpi.Scale(12), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4)      // SWP_NOZORDER
			setWindowPos.Call(extraBtn, 0, dpi.Scale(95), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4)   // SWP_NOZORDER
			setWindowPos.Call(cancelBtn, 0, dpi.Scale(178), dpi.Scale(65), dpi.Scale(75), dpi.Scale(24), 0x4) // SWP_NOZORDER
		}
	}

	proc := func(wnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
		switch msg {
		case 0x0002: // WM_DESTROY
			postQuitMessage.Call(0)

		case 0x0010: // WM_CLOSE
			destroyWindow.Call(wnd)

		case 0x0111: // WM_COMMAND
			switch wparam {
			default:
				return 1
			case 1, 6: // IDOK, IDYES
				out = getWindowString(editCtl)
				ok = true
			case 2: // IDCANCEL
			case 7: // IDNO
				err = ErrExtraButton
			}
			destroyWindow.Call(wnd)

		case 0x02e0: // WM_DPICHANGED
			layout(dpi(uint32(wparam) >> 16))

		default:
			ret, _, _ := syscall.Syscall6(defWindowProc, 4, wnd, uintptr(msg), wparam, lparam, 0, 0)
			return ret
		}

		return 0
	}

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", false, opts.ctx.Err()
	}

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return "", false, err
	}

	cls, err := registerClass(instance, syscall.NewCallback(proc))
	if cls == 0 {
		return "", false, err
	}
	defer unregisterClass.Call(cls, instance)

	wnd, _, _ = createWindowEx.Call(0x10101, // WS_EX_CONTROLPARENT|WS_EX_WINDOWEDGE|WS_EX_DLGMODALFRAME
		cls, strptr(title),
		0x84c80000, // WS_POPUPWINDOW|WS_CLIPSIBLINGS|WS_DLGFRAME
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		281, 140, 0, 0, instance)

	textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr(text),
		0x5002e080, // WS_CHILD|WS_VISIBLE|WS_GROUP|SS_WORDELLIPSIS|SS_EDITCONTROL|SS_NOPREFIX
		12, 10, 241, 16, wnd, 0, instance)

	var flags uintptr = 0x50030080 // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|ES_AUTOHSCROLL
	if opts.hideText {
		flags |= 0x20 // ES_PASSWORD
	}
	editCtl, _, _ = createWindowEx.Call(0x200, // WS_EX_CLIENTEDGE
		strptr("EDIT"), strptr(opts.entryText),
		flags,
		12, 30, 241, 24, wnd, 0, instance)

	okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		0x50030001, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|BS_DEFPUSHBUTTON
		12, 65, 75, 24, wnd, 1 /* IDOK */, instance)
	cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
		12, 65, 75, 24, wnd, 2 /* IDCANCEL */, instance)
	if opts.extraButton != nil {
		extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
			12, 65, 75, 24, wnd, 7 /* IDNO */, instance)
	}

	layout(getDPI(wnd))
	centerWindow(wnd)
	setFocus.Call(editCtl)
	showWindow.Call(wnd, 1 /* SW_SHOWNORMAL */, 0)
	sendMessage.Call(editCtl, 0xb1 /* EM_SETSEL */, 0, intptr(-1))

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
	out, ok, err = "", false, nil

	if err := messageLoop(wnd); err != nil {
		return "", false, err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", false, opts.ctx.Err()
	}
	return out, ok, err
}

func password(opts options) (string, string, bool, error) {
	opts.hideText = true
	str, ok, err := entry("Password:", opts)
	return "", str, ok, err
}
