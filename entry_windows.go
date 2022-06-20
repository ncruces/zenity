package zenity

import (
	"syscall"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
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

	instance, err := win.GetModuleHandle(nil)
	if err != nil {
		return "", err
	}

	cls, err := registerClass(instance, icon.handle, syscall.NewCallback(entryProc))
	if err != nil {
		return "", err
	}
	defer win.UnregisterClass(cls, instance)

	owner, _ := opts.attach.(win.HWND)
	dlg.wnd, _, _ = createWindowEx.Call(_WS_EX_CONTROLPARENT|_WS_EX_WINDOWEDGE|_WS_EX_DLGMODALFRAME,
		uintptr(cls), strptr(*opts.title),
		_WS_POPUPWINDOW|_WS_CLIPSIBLINGS|_WS_DLGFRAME,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		281, 141, uintptr(owner), 0, uintptr(instance), uintptr(unsafe.Pointer(dlg)))

	dlg.textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr(text),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_SS_WORDELLIPSIS|_SS_EDITCONTROL|_SS_NOPREFIX,
		12, 10, 241, 16, dlg.wnd, 0, uintptr(instance), 0)

	var flags uintptr = _WS_CHILD | _WS_VISIBLE | _WS_GROUP | _WS_TABSTOP | _ES_AUTOHSCROLL
	if opts.hideText {
		flags |= _ES_PASSWORD
	}
	dlg.editCtl, _, _ = createWindowEx.Call(_WS_EX_CLIENTEDGE,
		strptr("EDIT"), strptr(opts.entryText),
		flags,
		12, 30, 241, 24, dlg.wnd, 0, uintptr(instance), 0)

	dlg.okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP|_BS_DEFPUSHBUTTON,
		12, 66, 75, 24, dlg.wnd, win.IDOK, uintptr(instance), 0)
	dlg.cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
		12, 66, 75, 24, dlg.wnd, win.IDCANCEL, uintptr(instance), 0)
	if opts.extraButton != nil {
		dlg.extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
			12, 66, 75, 24, dlg.wnd, win.IDNO, uintptr(instance), 0)
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	setFocus.Call(dlg.editCtl)
	showWindow.Call(dlg.wnd, _SW_NORMAL, 0)
	sendMessage.Call(dlg.editCtl, win.EM_SETSEL, 0, intptr(-1))

	if opts.ctx != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				sendMessage.Call(dlg.wnd, win.WM_SYSCOMMAND, _SC_CLOSE, 0)
			case <-wait:
			}
		}()
	}

	if err := win.MessageLoop(win.HWND(dlg.wnd)); err != nil {
		return "", err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	return dlg.out, dlg.err
}

func (dlg *entryDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	sendMessage.Call(dlg.textCtl, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.editCtl, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.okBtn, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.cancelBtn, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.extraBtn, win.WM_SETFONT, font, 1)
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
	case win.WM_NCCREATE:
		saveBackRef(wnd, *lparam)
		dlg = (*entryDialog)(*lparam)
	case win.WM_NCDESTROY:
		deleteBackRef(wnd)
	default:
		dlg = (*entryDialog)(loadBackRef(wnd))
	}

	switch msg {
	case win.WM_DESTROY:
		postQuitMessage.Call(0)

	case win.WM_CLOSE:
		dlg.err = ErrCanceled
		destroyWindow.Call(wnd)

	case win.WM_COMMAND:
		switch wparam {
		default:
			return 1
		case win.IDOK, win.IDYES:
			dlg.out = getWindowString(dlg.editCtl)
		case win.IDCANCEL:
			dlg.err = ErrCanceled
		case win.IDNO:
			dlg.err = ErrExtraButton
		}
		destroyWindow.Call(wnd)

	case win.WM_DPICHANGED:
		dlg.layout(dpi(uint32(wparam) >> 16))

	default:
		res, _, _ := defWindowProc.Call(wnd, uintptr(msg), wparam, uintptr(unsafe.Pointer(lparam)))
		return res
	}

	return 0
}
