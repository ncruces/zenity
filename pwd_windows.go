package zenity

import (
	"syscall"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

func password(opts options) (string, string, error) {
	if !opts.username {
		opts.entryText = ""
		opts.hideText = true
		str, err := entry("Password:", opts)
		return "", str, err
	}

	if opts.title == nil {
		opts.title = stringPtr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = stringPtr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = stringPtr("Cancel")
	}

	dlg := &passwordDialog{}
	return dlg.setup(opts)
}

type passwordDialog struct {
	usr string
	pwd string
	err error

	wnd       uintptr
	uTextCtl  uintptr
	uEditCtl  uintptr
	pTextCtl  uintptr
	pEditCtl  uintptr
	okBtn     uintptr
	cancelBtn uintptr
	extraBtn  uintptr
	font      font
}

func (dlg *passwordDialog) setup(opts options) (string, string, error) {
	defer setup()()
	dlg.font = getFont()
	defer dlg.font.delete()
	icon := getIcon(opts.windowIcon)
	defer icon.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", "", opts.ctx.Err()
	}

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return "", "", err
	}

	cls, err := registerClass(instance, icon.handle, syscall.NewCallback(passwordProc))
	if cls == 0 {
		return "", "", err
	}
	defer unregisterClass.Call(cls, instance)

	owner, _ := opts.attach.(win.HWND)
	dlg.wnd, _, _ = createWindowEx.Call(_WS_EX_CONTROLPARENT|_WS_EX_WINDOWEDGE|_WS_EX_DLGMODALFRAME,
		cls, strptr(*opts.title),
		_WS_POPUPWINDOW|_WS_CLIPSIBLINGS|_WS_DLGFRAME,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		281, 191, uintptr(owner), 0, instance, uintptr(unsafe.Pointer(dlg)))

	dlg.uTextCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr("Username:"),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_SS_WORDELLIPSIS|_SS_EDITCONTROL|_SS_NOPREFIX,
		12, 10, 241, 16, dlg.wnd, 0, instance, 0)

	var flags uintptr = _WS_CHILD | _WS_VISIBLE | _WS_GROUP | _WS_TABSTOP | _ES_AUTOHSCROLL
	dlg.uEditCtl, _, _ = createWindowEx.Call(_WS_EX_CLIENTEDGE,
		strptr("EDIT"), 0,
		flags,
		12, 30, 241, 24, dlg.wnd, 0, instance, 0)

	dlg.pTextCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr("Password:"),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_SS_WORDELLIPSIS|_SS_EDITCONTROL|_SS_NOPREFIX,
		12, 60, 241, 16, dlg.wnd, 0, instance, 0)

	dlg.pEditCtl, _, _ = createWindowEx.Call(_WS_EX_CLIENTEDGE,
		strptr("EDIT"), 0,
		flags|_ES_PASSWORD,
		12, 80, 241, 24, dlg.wnd, 0, instance, 0)

	dlg.okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP|_BS_DEFPUSHBUTTON,
		12, 116, 75, 24, dlg.wnd, win.IDOK, instance, 0)
	dlg.cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
		12, 116, 75, 24, dlg.wnd, win.IDCANCEL, instance, 0)
	if opts.extraButton != nil {
		dlg.extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
			12, 116, 75, 24, dlg.wnd, win.IDNO, instance, 0)
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	setFocus.Call(dlg.uEditCtl)
	showWindow.Call(dlg.wnd, _SW_NORMAL, 0)
	sendMessage.Call(dlg.uEditCtl, win.EM_SETSEL, 0, intptr(-1))

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

	if err := messageLoop(dlg.wnd); err != nil {
		return "", "", err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", "", opts.ctx.Err()
	}
	return dlg.usr, dlg.pwd, dlg.err
}

func (dlg *passwordDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	sendMessage.Call(dlg.uTextCtl, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.uEditCtl, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.pTextCtl, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.pEditCtl, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.okBtn, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.cancelBtn, win.WM_SETFONT, font, 1)
	sendMessage.Call(dlg.extraBtn, win.WM_SETFONT, font, 1)
	setWindowPos.Call(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(191), _SWP_NOZORDER|_SWP_NOMOVE)
	setWindowPos.Call(dlg.uTextCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), _SWP_NOZORDER)
	setWindowPos.Call(dlg.uEditCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(24), _SWP_NOZORDER)
	setWindowPos.Call(dlg.pTextCtl, 0, dpi.scale(12), dpi.scale(60), dpi.scale(241), dpi.scale(16), _SWP_NOZORDER)
	setWindowPos.Call(dlg.pEditCtl, 0, dpi.scale(12), dpi.scale(80), dpi.scale(241), dpi.scale(24), _SWP_NOZORDER)
	if dlg.extraBtn == 0 {
		setWindowPos.Call(dlg.okBtn, 0, dpi.scale(95), dpi.scale(116), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(116), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	} else {
		setWindowPos.Call(dlg.okBtn, 0, dpi.scale(12), dpi.scale(116), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(116), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(116), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	}
}

func passwordProc(wnd uintptr, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *passwordDialog
	switch msg {
	case win.WM_NCCREATE:
		saveBackRef(wnd, *lparam)
		dlg = (*passwordDialog)(*lparam)
	case win.WM_NCDESTROY:
		deleteBackRef(wnd)
	default:
		dlg = (*passwordDialog)(loadBackRef(wnd))
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
			dlg.usr = getWindowString(dlg.uEditCtl)
			dlg.pwd = getWindowString(dlg.pEditCtl)
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
