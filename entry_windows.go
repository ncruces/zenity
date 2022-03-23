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

	defer setup()()
	wnd := &entryWnd{font: getFont()}
	defer wnd.font.delete()

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

	wnd.handle, _, _ = createWindowEx.Call(0x10101, // WS_EX_CONTROLPARENT|WS_EX_WINDOWEDGE|WS_EX_DLGMODALFRAME
		cls, strptr(*opts.title),
		0x84c80000, // WS_POPUPWINDOW|WS_CLIPSIBLINGS|WS_DLGFRAME
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		281, 141, 0, 0, instance, uintptr(unsafe.Pointer(wnd)))

	wnd.textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr(text),
		0x5002e080, // WS_CHILD|WS_VISIBLE|WS_GROUP|SS_WORDELLIPSIS|SS_EDITCONTROL|SS_NOPREFIX
		12, 10, 241, 16, wnd.handle, 0, instance, 0)

	var flags uintptr = 0x50030080 // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|ES_AUTOHSCROLL
	if opts.hideText {
		flags |= 0x20 // ES_PASSWORD
	}
	wnd.editCtl, _, _ = createWindowEx.Call(0x200, // WS_EX_CLIENTEDGE
		strptr("EDIT"), strptr(opts.entryText),
		flags,
		12, 30, 241, 24, wnd.handle, 0, instance, 0)

	wnd.okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		0x50030001, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP|BS_DEFPUSHBUTTON
		12, 66, 75, 24, wnd.handle, 1 /* IDOK */, instance, 0)
	wnd.cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
		12, 66, 75, 24, wnd.handle, 2 /* IDCANCEL */, instance, 0)
	if opts.extraButton != nil {
		wnd.extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
			12, 66, 75, 24, wnd.handle, 7 /* IDNO */, instance, 0)
	}

	wnd.layout(getDPI(wnd.handle))
	centerWindow(wnd.handle)
	setFocus.Call(wnd.editCtl)
	showWindow.Call(wnd.handle, 1 /* SW_SHOWNORMAL */, 0)
	sendMessage.Call(wnd.editCtl, 0xb1 /* EM_SETSEL */, 0, intptr(-1))

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
		return "", err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	return wnd.out, wnd.err
}

type entryWnd struct {
	handle    uintptr
	textCtl   uintptr
	editCtl   uintptr
	okBtn     uintptr
	cancelBtn uintptr
	extraBtn  uintptr
	out       string
	err       error
	font      font
}

func (wnd *entryWnd) layout(dpi dpi) {
	font := wnd.font.forDPI(dpi)
	sendMessage.Call(wnd.textCtl, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(wnd.editCtl, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(wnd.okBtn, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(wnd.cancelBtn, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(wnd.extraBtn, 0x0030 /* WM_SETFONT */, font, 1)
	setWindowPos.Call(wnd.handle, 0, 0, 0, dpi.scale(281), dpi.scale(141), 0x6)                         // SWP_NOZORDER|SWP_NOMOVE
	setWindowPos.Call(wnd.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), 0x4) // SWP_NOZORDER
	setWindowPos.Call(wnd.editCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(24), 0x4) // SWP_NOZORDER
	if wnd.extraBtn == 0 {
		setWindowPos.Call(wnd.okBtn, 0, dpi.scale(95), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4)      // SWP_NOZORDER
		setWindowPos.Call(wnd.cancelBtn, 0, dpi.scale(178), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
	} else {
		setWindowPos.Call(wnd.okBtn, 0, dpi.scale(12), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4)      // SWP_NOZORDER
		setWindowPos.Call(wnd.extraBtn, 0, dpi.scale(95), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4)   // SWP_NOZORDER
		setWindowPos.Call(wnd.cancelBtn, 0, dpi.scale(178), dpi.scale(66), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
	}
}

func entryProc(hwnd uintptr, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var wnd *entryWnd
	switch msg {
	case 0x0081: // WM_NCCREATE
		saveBackRef(hwnd, *lparam)
		wnd = (*entryWnd)(*lparam)
	case 0x0082: // WM_NCDESTROY
		deleteBackRef(hwnd)
	default:
		wnd = (*entryWnd)(loadBackRef(hwnd))
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
			wnd.out = getWindowString(wnd.editCtl)
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
