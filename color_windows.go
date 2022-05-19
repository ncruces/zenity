package zenity

import (
	"image/color"
	"sync"
	"unsafe"
)

var (
	chooseColor = comdlg32.NewProc("ChooseColorW")

	savedColors [16]uint32
	colorsMutex sync.Mutex
)

func init() {
	for i := range savedColors {
		savedColors[i] = 0xffffff
	}
}

const (
	_CC_RGBINIT         = 0x00000001
	_CC_FULLOPEN        = 0x00000002
	_CC_PREVENTFULLOPEN = 0x00000004
)

func selectColor(opts options) (color.Color, error) {
	// load custom colors
	colorsMutex.Lock()
	customColors := savedColors
	colorsMutex.Unlock()

	var args _CHOOSECOLOR
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.CustColors = &customColors

	if opts.color != nil {
		args.Flags |= _CC_RGBINIT
		n := color.NRGBAModel.Convert(opts.color).(color.NRGBA)
		args.RgbResult = uint32(n.R) | uint32(n.G)<<8 | uint32(n.B)<<16
	}
	if opts.showPalette {
		args.Flags |= _CC_PREVENTFULLOPEN
	} else {
		args.Flags |= _CC_FULLOPEN
	}

	defer setup()()

	if opts.ctx != nil || opts.title != nil {
		unhook, err := hookDialog(opts.ctx, opts.windowIcon, opts.title, nil)
		if err != nil {
			return nil, err
		}
		defer unhook()
	}

	s, _, _ := chooseColor.Call(uintptr(unsafe.Pointer(&args)))
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}
	if s == 0 {
		return nil, commDlgError()
	}

	// save custom colors back
	colorsMutex.Lock()
	savedColors = customColors
	colorsMutex.Unlock()

	r := uint8(args.RgbResult >> 0)
	g := uint8(args.RgbResult >> 8)
	b := uint8(args.RgbResult >> 16)
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

// https://docs.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-choosecolorw-r1
type _CHOOSECOLOR struct {
	StructSize   uint32
	Owner        uintptr
	Instance     uintptr
	RgbResult    uint32
	CustColors   *[16]uint32
	Flags        uint32
	CustData     uintptr
	FnHook       uintptr
	TemplateName *uint16
}
