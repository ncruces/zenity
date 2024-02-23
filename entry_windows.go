package zenity

import (
	"syscall"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

func entry(text string, opts options) (string, error) {
	if opts.title == nil {
		opts.title = ptr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = ptr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = ptr("Cancel")
	}

	dlg := &entryDialog{}
	return dlg.setup(text, opts)
}

type entryDialog struct {
	out string
	err error

	wnd       win.HWND
	textCtl   win.HWND
	editCtl   win.HWND
	okBtn     win.HWND
	cancelBtn win.HWND
	extraBtn  win.HWND
	font      font
}

func (dlg *entryDialog) setup(text string, opts options) (string, error) {
	owner, _ := opts.attach.(win.HWND)
	defer setup(owner)()
	dlg.font = getFont()
	defer dlg.font.delete()
	icon, _ := getIcon(opts.windowIcon)
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

	dlg.wnd, _ = win.CreateWindowEx(_WS_EX_ZEN_DIALOG,
		cls, strptr(*opts.title), _WS_ZEN_DIALOG,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		281, 141, owner, 0, instance, unsafe.Pointer(dlg))

	dlg.textCtl, _ = win.CreateWindowEx(0,
		strptr("STATIC"), strptr(text), _WS_ZEN_LABEL,
		12, 10, 241, 16, dlg.wnd, 0, instance, nil)

	var flags uint32 = _WS_ZEN_CONTROL | win.ES_AUTOHSCROLL
	if opts.hideText {
		flags |= win.ES_PASSWORD
	}
	dlg.editCtl, _ = win.CreateWindowEx(win.WS_EX_CLIENTEDGE,
		strptr("EDIT"), strptr(opts.entryText),
		flags,
		12, 30, 241, 24, dlg.wnd, 0, instance, nil)

	dlg.okBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(quoteAccelerators(*opts.okLabel)),
		_WS_ZEN_BUTTON|win.BS_DEFPUSHBUTTON,
		12, 66, 75, 24, dlg.wnd, win.IDOK, instance, nil)
	dlg.cancelBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(quoteAccelerators(*opts.cancelLabel)),
		_WS_ZEN_BUTTON,
		12, 66, 75, 24, dlg.wnd, win.IDCANCEL, instance, nil)
	if opts.extraButton != nil {
		dlg.extraBtn, _ = win.CreateWindowEx(0,
			strptr("BUTTON"), strptr(quoteAccelerators(*opts.extraButton)),
			_WS_ZEN_BUTTON,
			12, 66, 75, 24, dlg.wnd, win.IDNO, instance, nil)
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	win.SetFocus(dlg.editCtl)
	win.ShowWindow(dlg.wnd, win.SW_NORMAL)
	win.SendMessage(dlg.editCtl, win.EM_SETSEL, 0, intptr(-1))

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
		return "", err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return "", opts.ctx.Err()
	}
	return dlg.out, dlg.err
}

func (dlg *entryDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	win.SendMessage(dlg.textCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.editCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.okBtn, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.cancelBtn, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.extraBtn, win.WM_SETFONT, font, 1)
	win.SetWindowPos(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(141), win.SWP_NOMOVE|win.SWP_NOZORDER)
	win.SetWindowPos(dlg.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), win.SWP_NOZORDER)
	win.SetWindowPos(dlg.editCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(24), win.SWP_NOZORDER)
	if dlg.extraBtn == 0 {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(95), dpi.scale(66), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(66), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
	} else {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(12), dpi.scale(66), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(66), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(66), dpi.scale(75), dpi.scale(24), win.SWP_NOZORDER)
	}
}

func entryProc(wnd win.HWND, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *entryDialog
	switch msg {
	case win.WM_NCCREATE:
		saveBackRef(uintptr(wnd), *lparam)
		dlg = (*entryDialog)(*lparam)
	case win.WM_NCDESTROY:
		deleteBackRef(uintptr(wnd))
	default:
		dlg = (*entryDialog)(loadBackRef(uintptr(wnd)))
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
			dlg.out = win.GetWindowText(dlg.editCtl)
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
