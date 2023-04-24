package zenity

import (
	"syscall"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

func forms(text string, opts options) ([]string, error) {
	if opts.title == nil {
		opts.title = ptr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = ptr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = ptr("Cancel")
	}

	dlg := &formsDialog{}
	return dlg.setup(text, opts)
}

type uiCtrl struct {
	hwnd   win.HWND
	x      int
	y      int
	width  int
	height int
}

type formsDialog struct {
	out []string
	err error

	height       int
	wnd          win.HWND
	textCtl      win.HWND
	okBtn        win.HWND
	cancelBtn    win.HWND
	extraBtn     win.HWND
	pwdLabelList []win.HWND
	pwdEditList  []win.HWND
	dnyUIList    []uiCtrl
	font         font
}

func (dlg *formsDialog) setup(text string, opts options) ([]string, error) {
	owner, _ := opts.attach.(win.HWND)
	defer setup(owner)()
	dlg.font = getFont()
	defer dlg.font.delete()
	icon, _ := getIcon(opts.windowIcon)
	defer icon.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}

	instance, err := win.GetModuleHandle(nil)
	if err != nil {
		return nil, err
	}

	cls, err := registerClass(instance, icon.handle, syscall.NewCallback(formsProc))
	if err != nil {
		return nil, err
	}
	defer win.UnregisterClass(cls, instance)

	dlg.wnd, _ = win.CreateWindowEx(_WS_EX_ZEN_DIALOG,
		cls, strptr(*opts.title), _WS_ZEN_DIALOG,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		281, 281, owner, 0, instance, unsafe.Pointer(dlg))

	// coordinates
	x := 12
	y := 10
	w := 241
	h := 16

	dlg.textCtl, _ = win.CreateWindowEx(0,
		strptr("STATIC"), strptr(text), _WS_ZEN_LABEL,
		x, y, w, h, dlg.wnd, 0, instance, nil)

	// fields
	for _, field := range opts.fields {
		switch field.kind {
		case FormFieldEntry:
			// label
			y += h + 4
			h = 16
			ctl, _ := win.CreateWindowEx(0,
				strptr("STATIC"), strptr(field.name), _WS_ZEN_LABEL,
				x, y, w, h, dlg.wnd, 0, instance, nil)
			dlg.dnyUIList = append(dlg.dnyUIList, uiCtrl{hwnd: ctl, x: x, y: y, width: w, height: h})

			// Edit
			y += h + 4
			h = 24
			ctl, _ = win.CreateWindowEx(win.WS_EX_CLIENTEDGE,
				strptr("EDIT"), nil,
				_WS_ZEN_CONTROL|win.ES_AUTOHSCROLL,
				x, y, w, h, dlg.wnd, 0, instance, nil)
			dlg.dnyUIList = append(dlg.dnyUIList, uiCtrl{hwnd: ctl, x: x, y: y, width: w, height: h})
			dlg.pwdEditList = append(dlg.pwdEditList, ctl)
		case FormFieldPassword:
			// label
			y += h + 4
			h = 16
			ctl, _ := win.CreateWindowEx(0,
				strptr("STATIC"), strptr(field.name), _WS_ZEN_LABEL,
				x, y, w, h, dlg.wnd, 0, instance, nil)
			dlg.pwdLabelList = append(dlg.pwdLabelList, ctl)
			dlg.dnyUIList = append(dlg.dnyUIList, uiCtrl{hwnd: ctl, x: x, y: y, width: w, height: h})

			// edit
			y += h + 4
			h = 24
			ctl, _ = win.CreateWindowEx(win.WS_EX_CLIENTEDGE,
				strptr("EDIT"), nil,
				_WS_ZEN_CONTROL|win.ES_AUTOHSCROLL|win.ES_PASSWORD,
				x, y, w, h, dlg.wnd, 0, instance, nil)
			dlg.pwdEditList = append(dlg.pwdEditList, ctl)
			dlg.dnyUIList = append(dlg.dnyUIList, uiCtrl{hwnd: ctl, x: x, y: y, width: w, height: h})
		case FormFieldCalendar:
		case FormFieldComboBox:
		case FormFieldList:
		}
	}

	x = 95
	y += h + 12
	h = 24
	w = 75
	ctl, _ := win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(quoteAccelerators(*opts.okLabel)),
		_WS_ZEN_BUTTON|win.BS_DEFPUSHBUTTON,
		x, y, w, h, dlg.wnd, win.IDOK, instance, nil)
	dlg.dnyUIList = append(dlg.dnyUIList, uiCtrl{hwnd: ctl, x: x, y: y, width: w, height: h})
	dlg.okBtn = ctl

	x += w + 8
	h = 24
	w = 75
	ctl, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(quoteAccelerators(*opts.cancelLabel)),
		_WS_ZEN_BUTTON,
		x, y, w, h, dlg.wnd, win.IDCANCEL, instance, nil)
	dlg.dnyUIList = append(dlg.dnyUIList, uiCtrl{hwnd: ctl, x: x, y: y, width: w, height: h})
	dlg.cancelBtn = ctl

	if opts.extraButton != nil {
		x = 12
		// y += h + 4
		h = 24
		ctl, _ = win.CreateWindowEx(0,
			strptr("BUTTON"), strptr(quoteAccelerators(*opts.extraButton)),
			_WS_ZEN_BUTTON,
			x, y, w, h, dlg.wnd, win.IDNO, instance, nil)
		dlg.dnyUIList = append(dlg.dnyUIList, uiCtrl{hwnd: ctl, x: x, y: y, width: w, height: h})
		dlg.extraBtn = ctl
	}

	dlg.height = y + h + 50

	dlg.update()
	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	if len(dlg.dnyUIList) > 0 {
		win.SetFocus(dlg.dnyUIList[0].hwnd)
	}
	win.ShowWindow(dlg.wnd, win.SW_NORMAL)

	if opts.ctx != nil {
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
		return nil, err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}

	return dlg.out, dlg.err
}

func (dlg *formsDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	win.SendMessage(dlg.textCtl, win.WM_SETFONT, font, 1)
	win.SetWindowPos(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(dlg.height), win.SWP_NOMOVE|win.SWP_NOZORDER)
	win.SetWindowPos(dlg.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), win.SWP_NOZORDER)

	// dynamic ui list
	for _, ctl := range dlg.dnyUIList {
		win.SendMessage(ctl.hwnd, win.WM_SETFONT, font, 1)
		win.SetWindowPos(ctl.hwnd, 0, dpi.scale(ctl.x), dpi.scale(ctl.y), dpi.scale(ctl.width), dpi.scale(ctl.height), win.SWP_NOZORDER)
	}
}

func (dlg *formsDialog) update() {
	/*if dlg.disallowEmpty {
		var enable bool
		if dlg.multiple {
			len := win.SendMessage(dlg.listCtl, win.LB_GETSELCOUNT, 0, 0)
			enable = int32(len) > 0
		} else {
			idx := win.SendMessage(dlg.listCtl, win.LB_GETCURSEL, 0, 0)
			enable = int32(idx) >= 0
		}
		win.EnableWindow(dlg.okBtn, enable)
	}*/
}

func formsProc(wnd win.HWND, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *formsDialog
	switch msg {
	case win.WM_NCCREATE:
		saveBackRef(uintptr(wnd), *lparam)
		dlg = (*formsDialog)(*lparam)
	case win.WM_NCDESTROY:
		deleteBackRef(uintptr(wnd))
	default:
		dlg = (*formsDialog)(loadBackRef(uintptr(wnd)))
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
			dlg.update()
			return 1
		case win.IDOK, win.IDYES:
			for _, ctl := range dlg.pwdEditList {
				dlg.out = append(dlg.out, win.GetWindowText(ctl))
			}
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
