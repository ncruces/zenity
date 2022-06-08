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
	icon := getIcon(opts.windowIcon)
	defer icon.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return "", err
	}

	cls, err := registerClass(instance, icon.handle, syscall.NewCallback(entryProc))
	if cls == 0 {
		return "", err
	}
	defer unregisterClass.Call(cls, instance)

	owner, _ := opts.attach.(uintptr)
	dlg.wnd, _, _ = createWindowEx.Call(_WS_EX_CONTROLPARENT|_WS_EX_WINDOWEDGE|_WS_EX_DLGMODALFRAME,
		cls, strptr(*opts.title),
		_WS_POPUPWINDOW|_WS_CLIPSIBLINGS|_WS_DLGFRAME,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		281, 141, owner, 0, instance, uintptr(unsafe.Pointer(dlg)))

	dlg.textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr(text),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_SS_WORDELLIPSIS|_SS_EDITCONTROL|_SS_NOPREFIX,
		12, 10, 241, 16, dlg.wnd, 0, instance, 0)

	var flags uintptr = _WS_CHILD | _WS_VISIBLE | _WS_GROUP | _WS_TABSTOP | _ES_AUTOHSCROLL
	if opts.hideText {
		flags |= _ES_PASSWORD
	}
	dlg.editCtl, _, _ = createWindowEx.Call(_WS_EX_CLIENTEDGE,
		strptr("EDIT"), strptr(opts.entryText),
		flags,
		12, 30, 241, 24, dlg.wnd, 0, instance, 0)

	dlg.okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP|_BS_DEFPUSHBUTTON,
		12, 66, 75, 24, dlg.wnd, _IDOK, instance, 0)
	dlg.cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
		12, 66, 75, 24, dlg.wnd, _IDCANCEL, instance, 0)
	if opts.extraButton != nil {
		dlg.extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
			12, 66, 75, 24, dlg.wnd, _IDNO, instance, 0)
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	setFocus.Call(dlg.editCtl)
	showWindow.Call(dlg.wnd, _SW_NORMAL, 0)
	sendMessage.Call(dlg.editCtl, _EM_SETSEL, 0, intptr(-1))

	if opts.ctx != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				sendMessage.Call(dlg.wnd, _WM_SYSCOMMAND, _SC_CLOSE, 0)
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
	sendMessage.Call(dlg.textCtl, _WM_SETFONT, font, 1)
	sendMessage.Call(dlg.editCtl, _WM_SETFONT, font, 1)
	sendMessage.Call(dlg.okBtn, _WM_SETFONT, font, 1)
	sendMessage.Call(dlg.cancelBtn, _WM_SETFONT, font, 1)
	sendMessage.Call(dlg.extraBtn, _WM_SETFONT, font, 1)
	setWindowPos.Call(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(141), _SWP_NOZORDER|_SWP_NOMOVE)
	setWindowPos.Call(dlg.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), _SWP_NOZORDER)
	setWindowPos.Call(dlg.editCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(24), _SWP_NOZORDER)
	if dlg.extraBtn == 0 {
		setWindowPos.Call(dlg.okBtn, 0, dpi.scale(95), dpi.scale(66), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(66), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	} else {
		setWindowPos.Call(dlg.okBtn, 0, dpi.scale(12), dpi.scale(66), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(66), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(66), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	}
}

func entryProc(wnd uintptr, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *entryDialog
	switch msg {
	case _WM_NCCREATE:
		saveBackRef(wnd, *lparam)
		dlg = (*entryDialog)(*lparam)
	case _WM_NCDESTROY:
		deleteBackRef(wnd)
	default:
		dlg = (*entryDialog)(loadBackRef(wnd))
	}

	switch msg {
	case _WM_DESTROY:
		postQuitMessage.Call(0)

	case _WM_CLOSE:
		dlg.err = ErrCanceled
		destroyWindow.Call(wnd)

	case _WM_COMMAND:
		switch wparam {
		default:
			return 1
		case _IDOK, _IDYES:
			dlg.out = getWindowString(dlg.editCtl)
		case _IDCANCEL:
			dlg.err = ErrCanceled
		case _IDNO:
			dlg.err = ErrExtraButton
		}
		destroyWindow.Call(wnd)

	case _WM_DPICHANGED:
		dlg.layout(dpi(uint32(wparam) >> 16))

	default:
		res, _, _ := defWindowProc.Call(wnd, uintptr(msg), wparam, uintptr(unsafe.Pointer(lparam)))
		return res
	}

	return 0
}
