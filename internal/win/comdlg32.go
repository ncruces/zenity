//go:build windows

package win

import (
	"fmt"

	"github.com/ncruces/zenity/internal/zenutil"
)

const (
	// ChooseColor flags
	CC_RGBINIT         = 0x00000001
	CC_FULLOPEN        = 0x00000002
	CC_PREVENTFULLOPEN = 0x00000004
)

// https://docs.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-choosecolorw-r1
type CHOOSECOLOR struct {
	StructSize   uint32
	Owner        HWND
	Instance     HWND
	RgbResult    uint32
	CustColors   *[16]uint32
	Flags        uint32
	CustData     Pointer
	FnHook       uintptr
	TemplateName *uint16
}

const (
	OFN_OVERWRITEPROMPT  = 0x00000002
	OFN_NOCHANGEDIR      = 0x00000008
	OFN_ALLOWMULTISELECT = 0x00000200
	OFN_PATHMUSTEXIST    = 0x00000800
	OFN_FILEMUSTEXIST    = 0x00001000
	OFN_CREATEPROMPT     = 0x00002000
	OFN_NOREADONLYRETURN = 0x00008000
	OFN_EXPLORER         = 0x00080000
	OFN_FORCESHOWHIDDEN  = 0x10000000
)

// https://docs.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-openfilenamew
type OPENFILENAME struct {
	StructSize      uint32
	Owner           HWND
	Instance        Handle
	Filter          *uint16
	CustomFilter    *uint16
	MaxCustomFilter uint32
	FilterIndex     uint32
	File            *uint16
	MaxFile         uint32
	FileTitle       *uint16
	MaxFileTitle    uint32
	InitialDir      *uint16
	Title           *uint16
	Flags           uint32
	FileOffset      uint16
	FileExtension   uint16
	DefExt          *uint16
	CustData        Pointer
	FnHook          uintptr
	TemplateName    *uint16
	pvReserved      uintptr
	dwReserved      uint32
	FlagsEx         uint32
}

func CommDlgError() error {
	if code := commDlgExtendedError(); code == 0 {
		return zenutil.ErrCanceled
	} else {
		return fmt.Errorf("common dialog error: %x", code)
	}
}

//sys ChooseColor(cc *CHOOSECOLOR) (ok bool) = comdlg32.ChooseColorW
//sys commDlgExtendedError() (code int) = comdlg32.CommDlgExtendedError
//sys GetOpenFileName(ofn *OPENFILENAME) (ok bool) = comdlg32.GetOpenFileNameW
//sys GetSaveFileName(ofn *OPENFILENAME) (ok bool) = comdlg32.GetSaveFileNameW
