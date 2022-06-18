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

	wnd       uintptr
	textCtl   uintptr
	dateCtl   uintptr
	okBtn     uintptr
	cancelBtn uintptr
	extraBtn  uintptr
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

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return time.Time{}, err
	}

	cls, err := registerClass(instance, icon.handle, syscall.NewCallback(calendarProc))
	if cls == 0 {
		return time.Time{}, err
	}
	defer unregisterClass.Call(cls, instance)

	owner, _ := opts.attach.(win.HWND)
	dlg.wnd, _, _ = createWindowEx.Call(_WS_EX_CONTROLPARENT|_WS_EX_WINDOWEDGE|_WS_EX_DLGMODALFRAME,
		cls, strptr(*opts.title),
		_WS_POPUPWINDOW|_WS_CLIPSIBLINGS|_WS_DLGFRAME,
		_CW_USEDEFAULT, _CW_USEDEFAULT,
		281, 281, uintptr(owner), 0, instance, uintptr(unsafe.Pointer(dlg)))

	dlg.textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), strptr(text),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_SS_WORDELLIPSIS|_SS_EDITCONTROL|_SS_NOPREFIX,
		12, 10, 241, 16, dlg.wnd, 0, instance, 0)

	var flags uintptr = _WS_CHILD | _WS_VISIBLE | _WS_GROUP | _WS_TABSTOP | _MCS_NOTODAY
	dlg.dateCtl, _, _ = createWindowEx.Call(0,
		strptr(_MONTHCAL_CLASS),
		0, flags,
		12, 30, 241, 164, dlg.wnd, 0, instance, 0)

	dlg.okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP|_BS_DEFPUSHBUTTON,
		12, 206, 75, 24, dlg.wnd, win.IDOK, instance, 0)
	dlg.cancelBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.cancelLabel),
		_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
		12, 206, 75, 24, dlg.wnd, win.IDCANCEL, instance, 0)
	if opts.extraButton != nil {
		dlg.extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			_WS_CHILD|_WS_VISIBLE|_WS_GROUP|_WS_TABSTOP,
			12, 206, 75, 24, dlg.wnd, win.IDNO, instance, 0)
	}

	if opts.time != nil {
		var date _SYSTEMTIME
		year, month, day := opts.time.Date()
		date.year = uint16(year)
		date.month = uint16(month)
		date.day = uint16(day)
		sendMessage.Call(dlg.dateCtl, _MCM_SETCURSEL, 0, uintptr(unsafe.Pointer(&date)))
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	setFocus.Call(dlg.dateCtl)
	showWindow.Call(dlg.wnd, _SW_NORMAL, 0)

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
		return time.Time{}, err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return time.Time{}, opts.ctx.Err()
	}
	return dlg.out, dlg.err
}

func (dlg *calendarDialog) layout(dpi dpi) {
	font := dlg.font.forDPI(dpi)
	sendMessage.Call(dlg.textCtl, _WM_SETFONT, font, 1)
	sendMessage.Call(dlg.dateCtl, _WM_SETFONT, font, 1)
	sendMessage.Call(dlg.okBtn, _WM_SETFONT, font, 1)
	sendMessage.Call(dlg.cancelBtn, _WM_SETFONT, font, 1)
	sendMessage.Call(dlg.extraBtn, _WM_SETFONT, font, 1)
	setWindowPos.Call(dlg.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(281), _SWP_NOZORDER|_SWP_NOMOVE)
	setWindowPos.Call(dlg.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), _SWP_NOZORDER)
	setWindowPos.Call(dlg.dateCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(164), _SWP_NOZORDER)
	if dlg.extraBtn == 0 {
		setWindowPos.Call(dlg.okBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	} else {
		setWindowPos.Call(dlg.okBtn, 0, dpi.scale(12), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.extraBtn, 0, dpi.scale(95), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
		setWindowPos.Call(dlg.cancelBtn, 0, dpi.scale(178), dpi.scale(206), dpi.scale(75), dpi.scale(24), _SWP_NOZORDER)
	}
}

func calendarProc(wnd uintptr, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *calendarDialog
	switch msg {
	case _WM_NCCREATE:
		saveBackRef(wnd, *lparam)
		dlg = (*calendarDialog)(*lparam)
	case _WM_NCDESTROY:
		deleteBackRef(wnd)
	default:
		dlg = (*calendarDialog)(loadBackRef(wnd))
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
		case win.IDOK, win.IDYES:
			var date _SYSTEMTIME
			sendMessage.Call(dlg.dateCtl, _MCM_GETCURSEL, 0, uintptr(unsafe.Pointer(&date)))
			dlg.out = time.Date(int(date.year), time.Month(date.month), int(date.day), 0, 0, 0, 0, time.UTC)
		case win.IDCANCEL:
			dlg.err = ErrCanceled
		case win.IDNO:
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
