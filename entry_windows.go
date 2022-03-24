package zenity

import (
	"syscall"
	"unsafe"
)

func entry(text string, opts options) (string, error) {
	if opts.title == nil {
		opts.title = stringPtr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = stringPtr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = stringPtr("Cancel")
	}

	dlg := &entryDialog{}
	return dlg.setup(text, opts)
}

type entryDialog struct {
	out string
	err error

	wnd       uintptr
	textCtl   uintptr
	editCtl   uintptr
	okBtn     uintptr
	cancelBtn uintptr
	extraBtn  uintptr
	font      font
}

func (dlg *entryDialog) setup(text string, opts options) (string, error) {
	defer setup()()
	dlg.font = getFont()
	defer dlg.font.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return "", err
	}

	cls, err := registerClass(instance, syscall.NewCallback(entryProc))
	if cls == 0 {
		return "", err
	}
	defer unregisterClass.Call(cls, instance)

	dlg.wnd, _, _ = createWindowEx.Call(0x10101, // WS_EX_CONTROLPARENT|WS_EX_WINDOWEDGE|WS_EX_DLGMODALFRAME
		cls, strptr(*opts.title),
		0x84c80000, // WS_POPUPWINDOW|WS_CLIPSIBLINGS|WS_DLGFRAME
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		281, 141, 0, 0, instance, uintptr(unsafe.Pointer(dlg)))

	dlg.textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr(text),
		0x5002e080, // WS_CHILD|WS_VISIBLE|WS_GROUP|SS_WORDELLIPSIS|SS_EDITCONTROL|SS_NOPREFIX
		12, 10, 241, 16, dlg.wnd, 0, instance, 0)

	var flags uintptr = 0x50030080 // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|ES_AUTOHSCROLL
	if opts.hideText {
		flags |= 0x20 // ES_PASSWORD
	}
	dlg.editCtl, _, _ = createWindowEx.Call(0x200, // WS_EX_CLIENTEDGE
		strptr("EDIT"), strptr(opts.entryText),
		flags,
		12, 30, 241, 24, dlg.wnd, 0, instance, 0)

	dlg.okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		0x50030001, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|BS_DEFPUSHBUTTON
		12, 66, 75, 24, dlg.wnd, 1 /* IDOK */, instance, 0)
	dlg.cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
		12, 66, 75, 24, dlg.wnd, 2 /* IDCANCEL */, instance, 0)
	if opts.extraButton != nil {
		dlg.extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
			12, 66, 75, 24, dlg.wnd, 7 /* IDNO */, instance, 0)
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	setFocus.Call(dlg.editCtl)
	showWindow.Call(dlg.wnd, 1 /* SW_SHOWNORMAL */, 0)
	sendMessage.Call(dlg.editCtl, 0xb1 /* EM_SETSEL */, 0, intptr(-1))

	if opts.ctx != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				sendMessage.Call(dlg.wnd, 0x0112 /* WM_SYSCOMMAND */, 0xf060 /* SC_CLOSE */, 0)
			case <-wait:
			}
		}()
	}

	if err := messageLoop(dlg.wnd); err != nil {
		return "", err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	return dlg.out, dlg.err
}

func (dlg *entryDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	sendMessage.Call(dlg.textCtl, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(dlg.editCtl, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(dlg.okBtn, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(dlg.cancelBtn, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(dlg.extraBtn, 0x0030 /* WM_SETFONT */, font, 1)
	setWindowPos.Call(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(141), 0x6)                            // SWP_NOZORDER|SWP_NOMOVE
	setWindowPos.Call(dlg.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), 0x4) // SWP_NOZORDER
	setWindowPos.Call(dlg.editCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(24), 0x4) // SWP_NOZORDER
	if dlg.extraBtn == 0 {
		setWindowPos.Call(dlg.okBtn, 0, dpi.scale(95), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4)      // SWP_NOZORDER
		setWindowPos.Call(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
	} else {
		setWindowPos.Call(dlg.okBtn, 0, dpi.scale(12), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4)      // SWP_NOZORDER
		setWindowPos.Call(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4)   // SWP_NOZORDER
		setWindowPos.Call(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
	}
}

func entryProc(wnd uintptr, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *entryDialog
	switch msg {
	case 0x0081: // WM_NCCREATE
		saveBackRef(wnd, *lparam)
		dlg = (*entryDialog)(*lparam)
	case 0x0082: // WM_NCDESTROY
		deleteBackRef(wnd)
	default:
		dlg = (*entryDialog)(loadBackRef(wnd))
	}

	switch msg {
	case 0x0002: // WM_DESTROY
		postQuitMessage.Call(0)

	case 0x0010: // WM_CLOSE
		dlg.err = ErrCanceled
		destroyWindow.Call(wnd)

	case 0x0111: // WM_COMMAND
		switch wparam {
		default:
			return 1
		case 1, 6: // IDOK, IDYES
			dlg.out = getWindowString(dlg.editCtl)
		case 2: // IDCANCEL
			dlg.err = ErrCanceled
		case 7: // IDNO
			dlg.err = ErrExtraButton
		}
		destroyWindow.Call(wnd)

	case 0x02e0: // WM_DPICHANGED
		dlg.layout(dpi(uint32(wparam) >> 16))

	default:
		res, _, _ := defWindowProc.Call(wnd, uintptr(msg), wparam, uintptr(unsafe.Pointer(lparam)))
		return res
	}

	return 0
}
