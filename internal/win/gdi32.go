//go:build windows

package win

const (
	LOGPIXELSX = 88
	LOGPIXELSY = 90
)

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-logfontw
type LOGFONT struct {
	Height         int32
	Width          int32
	Escapement     int32
	Orientation    int32
	Weight         int32
	Italic         byte
	Underline      byte
	StrikeOut      byte
	CharSet        byte
	OutPrecision   byte
	ClipPrecision  byte
	Quality        byte
	PitchAndFamily byte
	FaceName       [32]uint16
}

//sys CreateFontIndirect(lf *LOGFONT) (ret Handle) = gdi32.CreateFontIndirectW
//sys DeleteObject(o Handle) (ok bool) = gdi32.DeleteObject
//sys GetDeviceCaps(dc Handle, index int) (ret int) = gdi32.GetDeviceCaps
