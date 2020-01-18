package zenity

import (
	"image/color"
	"runtime"
	"sync"
	"unsafe"
)

var (
	chooseColor = comdlg32.NewProc("ChooseColorW")

	savedColors = [16]uint32{}
	colorsMutex sync.Mutex
)

func init() {
	for i := range savedColors {
		savedColors[i] = 0xffffff
	}
}

func SelectColor(options ...Option) (color.Color, error) {
	opts := optsParse(options)

	// load custom colors
	colorsMutex.Lock()
	customColors := savedColors
	colorsMutex.Unlock()

	var args _CHOOSECOLORW
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.CustColors = &customColors

	if opts.color != nil {
		args.Flags |= 0x1 // CC_RGBINIT
		r, g, b, _ := opts.color.RGBA()
		args.RgbResult = (r >> 8 << 0) | (g >> 8 << 8) | (b >> 8 << 16)
	}
	if opts.palette {
		args.Flags |= 4 // CC_PREVENTFULLOPEN
	} else {
		args.Flags |= 2 // CC_FULLOPEN
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	n, _, _ := chooseColor.Call(uintptr(unsafe.Pointer(&args)))

	// save custom colors back
	colorsMutex.Lock()
	savedColors = customColors
	colorsMutex.Unlock()

	if n == 0 {
		return nil, commDlgError()
	}

	r := uint8(args.RgbResult >> 0)
	g := uint8(args.RgbResult >> 8)
	b := uint8(args.RgbResult >> 16)
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

type _CHOOSECOLORW struct {
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
