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
		opts.title = ptr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = ptr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = ptr("Cancel")
	}

	dlg := &passwordDialog{}
	return dlg.setup(opts)
}

type passwordDialog struct {
	usr string
	pwd string
	err error

	wnd       win.HWND
	uTextCtl  win.HWND
	uEditCtl  win.HWND
	pTextCtl  win.HWND
	pEditCtl  win.HWND
	okBtn     win.HWND
	cancelBtn win.HWND
	extraBtn  win.HWND
	font      font
}

func (dlg *passwordDialog) setup(opts options) (string, string, error) {
	owner, _ := opts.attach.(win.HWND)
	defer setup(owner)()
	dlg.font = getFont()
	defer dlg.font.delete()
	icon, _ := getIcon(opts.windowIcon)
	defer icon.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", "", opts.ctx.Err()
	}

	instance, err := win.GetModuleHandle(nil)
	if err != nil {
		return "", "", err
	}

	cls, err := registerClass(instance, icon.handle, syscall.NewCallback(passwordProc))
	if err != nil {
		return "", "", err
	}
	defer win.UnregisterClass(cls, instance)

	dlg.wnd, _ = win.CreateWindowEx(_WS_EX_ZEN_DIALOG,
		cls, strptr(*opts.title), _WS_ZEN_DIALOG,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		281, 191, owner, 0, instance, unsafe.Pointer(dlg))

	dlg.uTextCtl, _ = win.CreateWindowEx(0,
		strptr("STATIC"), strptr("Username:"),
		_WS_ZEN_LABEL,
		12, 10, 241, 16, dlg.wnd, 0, instance, nil)

	dlg.uEditCtl, _ = win.CreateWindowEx(win.WS_EX_CLIENTEDGE,
		strptr("EDIT"), nil,
		_WS_ZEN_CONTROL|win.ES_AUTOHSCROLL,
		12, 30, 241, 24, dlg.wnd, 0, instance, nil)

	dlg.pTextCtl, _ = win.CreateWindowEx(0,
		strptr("STATIC"), strptr("Password:"), _WS_ZEN_LABEL,
		12, 60, 241, 16, dlg.wnd, 0, instance, nil)

	dlg.pEditCtl, _ = win.CreateWindowEx(win.WS_EX_CLIENTEDGE,
		strptr("EDIT"), nil,
		_WS_ZEN_CONTROL|win.ES_AUTOHSCROLL|win.ES_PASSWORD,
		12, 80, 241, 24, dlg.wnd, 0, instance, nil)

	dlg.okBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(quoteAccelerators(*opts.okLabel)),
		_WS_ZEN_BUTTON|win.BS_DEFPUSHBUTTON,
		12, 116, 75, 24, dlg.wnd, win.IDOK, instance, nil)
	dlg.cancelBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(quoteAccelerators(*opts.cancelLabel)),
		_WS_ZEN_BUTTON,
		12, 116, 75, 24, dlg.wnd, win.IDCANCEL, instance, nil)
	if opts.extraButton != nil {
		dlg.extraBtn, _ = win.CreateWindowEx(0,
			strptr("BUTTON"), strptr(quoteAccelerators(*opts.extraButton)),
			_WS_ZEN_BUTTON,
			12, 116, 75, 24, dlg.wnd, win.IDNO, instance, nil)
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	win.SetFocus(dlg.uEditCtl)
	win.ShowWindow(dlg.wnd, win.SW_NORMAL)
	win.SendMessage(dlg.uEditCtl, win.EM_SETSEL, 0, intptr(-1))

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
		return "", "", err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", "", opts.ctx.Err()
	}
	return dlg.usr, dlg.pwd, dlg.err
}

func (dlg *passwordDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	win.SendMessage(dlg.uTextCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.uEditCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.pTextCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.pEditCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.okBtn, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.cancelBtn, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.extraBtn, win.WM_SETFONT, font, 1)
	win.SetWindowPos(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(191), win.SWP_NOMOVE|win.SWP_NOZORDER)
	win.SetWindowPos(dlg.uTextCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), win.SWP_NOZORDER)
	win.SetWindowPos(dlg.uEditCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(24), win.SWP_NOZORDER)
	win.SetWindowPos(dlg.pTextCtl, 0, dpi.scale(12), dpi.scale(60), dpi.scale(241), dpi.scale(16), win.SWP_NOZORDER)
	win.SetWindowPos(dlg.pEditCtl, 0, dpi.scale(12), dpi.scale(80), dpi.scale(241), dpi.scale(24), win.SWP_NOZORDER)
	if dlg.extraBtn == 0 {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(95), dpi.scale(116), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(116), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
	} else {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(12), dpi.scale(116), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(116), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(116), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
	}
}

func passwordProc(wnd win.HWND, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *passwordDialog
	switch msg {
	case win.WM_NCCREATE:
		saveBackRef(uintptr(wnd), *lparam)
		dlg = (*passwordDialog)(*lparam)
	case win.WM_NCDESTROY:
		deleteBackRef(uintptr(wnd))
	default:
		dlg = (*passwordDialog)(loadBackRef(uintptr(wnd)))
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
			dlg.usr = win.GetWindowText(dlg.uEditCtl)
			dlg.pwd = win.GetWindowText(dlg.pEditCtl)
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
