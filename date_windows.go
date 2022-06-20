package zenity

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

func calendar(text string, opts options) (time.Time, error) {
	if opts.title == nil {
		opts.title = stringPtr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = stringPtr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = stringPtr("Cancel")
	}

	dlg := &calendarDialog{}
	return dlg.setup(text, opts)
}

type calendarDialog struct {
	out time.Time
	err error

	wnd       win.HWND
	textCtl   win.HWND
	dateCtl   win.HWND
	okBtn     win.HWND
	cancelBtn win.HWND
	extraBtn  win.HWND
	font      font
}

func (dlg *calendarDialog) setup(text string, opts options) (time.Time, error) {
	defer setup()()
	dlg.font = getFont()
	defer dlg.font.delete()
	icon := getIcon(opts.windowIcon)
	defer icon.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return time.Time{}, opts.ctx.Err()
	}

	instance, err := win.GetModuleHandle(nil)
	if err != nil {
		return time.Time{}, err
	}

	cls, err := registerClass(instance, icon.handle, syscall.NewCallback(calendarProc))
	if err != nil {
		return time.Time{}, err
	}
	defer win.UnregisterClass(cls, instance)

	owner, _ := opts.attach.(win.HWND)
	dlg.wnd, _ = win.CreateWindowEx(_WS_EX_CONTROLPARENT|_WS_EX_WINDOWEDGE|_WS_EX_DLGMODALFRAME,
		cls, strptr(*opts.title),
		_WS_POPUPWINDOW|_WS_CLIPSIBLINGS|_WS_DLGFRAME,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		281, 281, owner, 0, instance, unsafe.Pointer(dlg))

	dlg.textCtl, _ = win.CreateWindowEx(0,
		strptr("STATIC"), strptr(text),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_SS_WORDELLIPSIS|_SS_EDITCONTROL|_SS_NOPREFIX,
		12, 10, 241, 16, dlg.wnd, 0, instance, nil)

	var flags uint32 = _WS_CHILD | _WS_VISIBLE | _WS_GROUP | _WS_TABSTOP | _MCS_NOTODAY
	dlg.dateCtl, _ = win.CreateWindowEx(0,
		strptr(_MONTHCAL_CLASS),
		nil, flags,
		12, 30, 241, 164, dlg.wnd, 0, instance, nil)

	dlg.okBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP|_BS_DEFPUSHBUTTON,
		12, 206, 75, 24, dlg.wnd, win.IDOK, instance, nil)
	dlg.cancelBtn, _ = win.CreateWindowEx(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
		12, 206, 75, 24, dlg.wnd, win.IDCANCEL, instance, nil)
	if opts.extraButton != nil {
		dlg.extraBtn, _ = win.CreateWindowEx(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
			12, 206, 75, 24, dlg.wnd, win.IDNO, instance, nil)
	}

	if opts.time != nil {
		var date _SYSTEMTIME
		year, month, day := opts.time.Date()
		date.year = uint16(year)
		date.month = uint16(month)
		date.day = uint16(day)
		win.SendMessagePointer(dlg.dateCtl, win.MCM_SETCURSEL, 0, unsafe.Pointer(&date))
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	win.SetFocus(dlg.dateCtl)
	win.ShowWindow(dlg.wnd, _SW_NORMAL)

	if opts.ctx != nil {
		wait := make(chan struct{})
		defer close(wait)
		go func() {
			select {
			case <-opts.ctx.Done():
				win.SendMessage(dlg.wnd, win.WM_SYSCOMMAND, _SC_CLOSE, 0)
			case <-wait:
			}
		}()
	}

	if err := win.MessageLoop(win.HWND(dlg.wnd)); err != nil {
		return time.Time{}, err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return time.Time{}, opts.ctx.Err()
	}
	return dlg.out, dlg.err
}

func (dlg *calendarDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	win.SendMessage(dlg.textCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.dateCtl, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.okBtn, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.cancelBtn, win.WM_SETFONT, font, 1)
	win.SendMessage(dlg.extraBtn, win.WM_SETFONT, font, 1)
	win.SetWindowPos(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(281), _SWP_NOZORDER|_SWP_NOMOVE)
	win.SetWindowPos(dlg.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), _SWP_NOZORDER)
	win.SetWindowPos(dlg.dateCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(164), _SWP_NOZORDER)
	if dlg.extraBtn == 0 {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	} else {
		win.SetWindowPos(dlg.okBtn, 0, dpi.scale(12), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		win.SetWindowPos(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		win.SetWindowPos(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	}
}

func calendarProc(wnd uintptr, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *calendarDialog
	switch msg {
	case win.WM_NCCREATE:
		saveBackRef(wnd, *lparam)
		dlg = (*calendarDialog)(*lparam)
	case win.WM_NCDESTROY:
		deleteBackRef(wnd)
	default:
		dlg = (*calendarDialog)(loadBackRef(wnd))
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
			var date _SYSTEMTIME
			win.SendMessagePointer(dlg.dateCtl, win.MCM_GETCURSEL, 0, unsafe.Pointer(&date))
			dlg.out = time.Date(int(date.year), time.Month(date.month), int(date.day), 0, 0, 0, 0, time.UTC)
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
