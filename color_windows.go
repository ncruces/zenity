package zenity

import (
	"image/color"
	"sync"
	"unsafe"

	"github.com/ncruces/zenity/internal/win"
)

var (
	savedColors [16]uint32
	colorsMutex sync.Mutex
)

func init() {
	for i := range savedColors {
		savedColors[i] = 0xffffff
	}
}

func selectColor(opts options) (color.Color, error) {
	// load custom colors
	colorsMutex.Lock()
	customColors := savedColors
	colorsMutex.Unlock()

	var args win.CHOOSECOLOR
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Owner, _ = opts.attach.(win.HWND)
	args.CustColors = &customColors

	if opts.color != nil {
		args.Flags |= win.CC_RGBINIT
		n := color.NRGBAModel.Convert(opts.color).(color.NRGBA)
		args.RgbResult = uint32(n.R) | uint32(n.G)<<8 | uint32(n.B)<<16
	}
	if opts.showPalette {
		args.Flags |= win.CC_PREVENTFULLOPEN
	} else {
		args.Flags |= win.CC_FULLOPEN
	}

	defer setup(args.Owner)()
	unhook, err := hookDialog(opts.ctx, opts.windowIcon, opts.title, nil)
	if err != nil {
		return nil, err
	}
	defer unhook()

	ok := win.ChooseColor(&args)
	if opts.ctx != nil && opts.ctx.Err() != nil {
		return nil, opts.ctx.Err()
	}
	if !ok {
		return nil, win.CommDlgError()
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
