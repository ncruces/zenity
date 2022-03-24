package zenity

import (
	"context"
	"sync"
	"syscall"
	"unsafe"
)

func progress(opts options) (ProgressDialog, error) {
	if opts.title == nil {
		opts.title = stringPtr("")
	}
	if opts.okLabel == nil {
		opts.okLabel = stringPtr("OK")
	}
	if opts.cancelLabel == nil {
		opts.cancelLabel = stringPtr("Cancel")
	}
	if opts.maxValue == 0 {
		opts.maxValue = 100
	}
	if opts.ctx == nil {
		opts.ctx = context.Background()
	} else if cerr := opts.ctx.Err(); cerr != nil {
		return nil, cerr
	}

	dlg := &progressDialog{
		done: make(chan struct{}),
		max:  opts.maxValue,
	}
	dlg.init.Add(1)

	go func() {
		dlg.err = dlg.setup(opts)
		close(dlg.done)
	}()

	dlg.init.Wait()
	return dlg, nil
}

type progressDialog struct {
	init sync.WaitGroup
	done chan struct{}
	max  int
	err  error

	wnd       uintptr
	textCtl   uintptr
	progCtl   uintptr
	okBtn     uintptr
	cancelBtn uintptr
	extraBtn  uintptr
	font      font
}

func (d *progressDialog) Text(text string) error {
	select {
	default:
		setWindowText.Call(d.textCtl, strptr(text))
		return nil
	case <-d.done:
		return d.err
	}
}

func (d *progressDialog) Value(value int) error {
	select {
	default:
		sendMessage.Call(d.progCtl, 0x402 /* PBM_SETPOS */, uintptr(value), 0)
		if value >= d.max {
			enableWindow.Call(d.okBtn, 1)
		}
		return nil
	case <-d.done:
		return d.err
	}
}

func (d *progressDialog) MaxValue() int {
	return d.max
}

func (d *progressDialog) Done() <-chan struct{} {
	return d.done
}

func (d *progressDialog) Complete() error {
	select {
	default:
		setWindowLong.Call(d.progCtl, intptr(-16) /* GWL_STYLE */, 0x50000001 /* WS_CHILD|WS_VISIBLE|PBS_SMOOTH */)
		sendMessage.Call(d.progCtl, 0x406 /* PBM_SETRANGE32 */, 0, 1)
		sendMessage.Call(d.progCtl, 0x402 /* PBM_SETPOS */, 1, 0)
		enableWindow.Call(d.okBtn, 1)
		enableWindow.Call(d.cancelBtn, 0)
		return nil
	case <-d.done:
		return d.err
	}
}

func (d *progressDialog) Close() error {
	sendMessage.Call(d.wnd, 0x0112 /* WM_SYSCOMMAND */, 0xf060 /* SC_CLOSE */, 0)
	<-d.done
	if d.err == ErrCanceled {
		return nil
	}
	return d.err
}

func (dlg *progressDialog) setup(opts options) error {
	done := false
	defer func() {
		if !done {
			dlg.init.Done()
		}
	}()

	defer setup()()
	dlg.font = getFont()
	defer dlg.font.delete()

	if opts.ctx != nil && opts.ctx.Err() != nil {
		return opts.ctx.Err()
	}

	instance, _, err := getModuleHandle.Call(0)
	if instance == 0 {
		return err
	}

	cls, err := registerClass(instance, syscall.NewCallback(progressProc))
	if cls == 0 {
		return err
	}
	defer unregisterClass.Call(cls, instance)

	dlg.wnd, _, _ = createWindowEx.Call(0x10101, // WS_EX_CONTROLPARENT|WS_EX_WINDOWEDGE|WS_EX_DLGMODALFRAME
		cls, strptr(*opts.title),
		0x84c80000, // WS_POPUPWINDOW|WS_CLIPSIBLINGS|WS_DLGFRAME
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		281, 133, 0, 0, instance, uintptr(unsafe.Pointer(dlg)))

	dlg.textCtl, _, _ = createWindowEx.Call(0,
		strptr("STATIC"), 0,
		0x5002e080, // WS_CHILD|WS_VISIBLE|WS_GROUP|SS_WORDELLIPSIS|SS_EDITCONTROL|SS_NOPREFIX
		12, 10, 241, 16, dlg.wnd, 0, instance, 0)

	var flags uintptr = 0x50000001 // WS_CHILD|WS_VISIBLE|PBS_SMOOTH
	if opts.maxValue < 0 {
		flags |= 0x8 // PBS_MARQUEE
	}
	dlg.progCtl, _, _ = createWindowEx.Call(0,
		strptr("msctls_progress32"), // PROGRESS_CLASS
		0, flags,
		12, 30, 241, 16, dlg.wnd, 0, instance, 0)

	dlg.okBtn, _, _ = createWindowEx.Call(0,
		strptr("BUTTON"), strptr(*opts.okLabel),
		0x58030001, // WS_CHILD|WS_VISIBLE|WS_DISABLED|WS_GROUP|WS_TABSTOP|BS_DEFPUSHBUTTON
		12, 58, 75, 24, dlg.wnd, 1 /* IDOK */, instance, 0)
	if !opts.noCancel {
		dlg.cancelBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.cancelLabel),
			0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
			12, 58, 75, 24, dlg.wnd, 2 /* IDCANCEL */, instance, 0)
	}
	if opts.extraButton != nil {
		dlg.extraBtn, _, _ = createWindowEx.Call(0,
			strptr("BUTTON"), strptr(*opts.extraButton),
			0x50010000, // WS_CHILD|WS_VISIBLE|WS_GROUP|WS_TABSTOP
			12, 58, 75, 24, dlg.wnd, 7 /* IDNO */, instance, 0)
	}

	dlg.layout(getDPI(dlg.wnd))
	centerWindow(dlg.wnd)
	showWindow.Call(dlg.wnd, 1 /* SW_SHOWNORMAL */, 0)
	if opts.maxValue < 0 {
		sendMessage.Call(dlg.progCtl, 0x40a /* PBM_SETMARQUEE */, 1, 0)
	} else {
		sendMessage.Call(dlg.progCtl, 0x406 /* PBM_SETRANGE32 */, 0, uintptr(opts.maxValue))
	}
	dlg.init.Done()
	done = true

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
		return err
	}
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return opts.ctx.Err()
	}
	return dlg.err
}

func (d *progressDialog) layout(dpi dpi) {
	font := d.font.forDPI(dpi)
	sendMessage.Call(d.textCtl, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(d.okBtn, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(d.cancelBtn, 0x0030 /* WM_SETFONT */, font, 1)
	sendMessage.Call(d.extraBtn, 0x0030 /* WM_SETFONT */, font, 1)
	setWindowPos.Call(d.wnd, 0, 0, 0, dpi.scale(281), dpi.scale(133), 0x6)                            // SWP_NOZORDER|SWP_NOMOVE
	setWindowPos.Call(d.textCtl, 0, dpi.scale(12), dpi.scale(10), dpi.scale(241), dpi.scale(16), 0x4) // SWP_NOZORDER
	setWindowPos.Call(d.progCtl, 0, dpi.scale(12), dpi.scale(30), dpi.scale(241), dpi.scale(16), 0x4) // SWP_NOZORDER
	if d.extraBtn == 0 {
		if d.cancelBtn == 0 {
			setWindowPos.Call(d.okBtn, 0, dpi.scale(178), dpi.scale(58), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
		} else {
			setWindowPos.Call(d.okBtn, 0, dpi.scale(95), dpi.scale(58), dpi.scale(75), dpi.scale(24), 0x4)      // SWP_NOZORDER
			setWindowPos.Call(d.cancelBtn, 0, dpi.scale(178), dpi.scale(58), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
		}
	} else {
		if d.cancelBtn == 0 {
			setWindowPos.Call(d.okBtn, 0, dpi.scale(95), dpi.scale(58), dpi.scale(75), dpi.scale(24), 0x4)     // SWP_NOZORDER
			setWindowPos.Call(d.extraBtn, 0, dpi.scale(178), dpi.scale(58), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
		} else {
			setWindowPos.Call(d.okBtn, 0, dpi.scale(12), dpi.scale(58), dpi.scale(75), dpi.scale(24), 0x4)      // SWP_NOZORDER
			setWindowPos.Call(d.extraBtn, 0, dpi.scale(95), dpi.scale(58), dpi.scale(75), dpi.scale(24), 0x4)   // SWP_NOZORDER
			setWindowPos.Call(d.cancelBtn, 0, dpi.scale(178), dpi.scale(58), dpi.scale(75), dpi.scale(24), 0x4) // SWP_NOZORDER
		}
	}
}

func progressProc(wnd uintptr, msg uint32, wparam uintptr, lparam *unsafe.Pointer) uintptr {
	var dlg *progressDialog
	switch msg {
	case 0x0081: // WM_NCCREATE
		saveBackRef(wnd, *lparam)
		dlg = (*progressDialog)(*lparam)
	case 0x0082: // WM_NCDESTROY
		deleteBackRef(wnd)
	default:
		dlg = (*progressDialog)(loadBackRef(wnd))
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
			//
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
