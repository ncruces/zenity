package dialog

import (
	"path/filepath"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func OpenFile(title, defaultPath string, filters []FileFilter) (string, error) {
	comdlg32 := syscall.NewLazyDLL("comdlg32.dll")
	getOpenFileName := comdlg32.NewProc("GetOpenFileNameW")

	var args _OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))

	if title != "" {
		args.Title = syscall.StringToUTF16Ptr(title)
	}
	if defaultPath != "" {
		args.InitialDir = syscall.StringToUTF16Ptr(defaultPath)
	}
	args.Filter = &windowsFilters(filters)[0]

	res := [1024]uint16{}
	args.File = &res[0]
	args.MaxFile = uint32(len(res))

	_, _, _ = getOpenFileName.Call(uintptr(unsafe.Pointer(&args)))
	return syscall.UTF16ToString(res[:]), nil
}

func OpenFiles(title, defaultPath string, filters []FileFilter) ([]string, error) {
	comdlg32 := syscall.NewLazyDLL("comdlg32.dll")
	getOpenFileName := comdlg32.NewProc("GetOpenFileNameW")

	var args _OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))
	args.Flags = 0x00080200 // OFN_ALLOWMULTISELECT|OFN_EXPLORER

	if title != "" {
		args.Title = syscall.StringToUTF16Ptr(title)
	}
	if defaultPath != "" {
		args.InitialDir = syscall.StringToUTF16Ptr(defaultPath)
	}
	args.Filter = &windowsFilters(filters)[0]

	res := [65536]uint16{}
	args.File = &res[0]
	args.MaxFile = uint32(len(res))

	_, _, _ = getOpenFileName.Call(uintptr(unsafe.Pointer(&args)))

	var i int
	var nul bool
	var split []string
	for j, p := range res {
		if p == 0 {
			if nul {
				break
			}
			if i < j {
				split = append(split, string(utf16.Decode(res[i:j])))
			}
			i = j + 1
			nul = true
		} else {
			nul = false
		}
	}
	if len := len(split) - 1; len > 0 {
		base := split[0]
		for i := 0; i < len; i++ {
			split[i] = filepath.Join(base, string(split[i+1]))
		}
		split = split[:len]
	}
	return split, nil
}

func SaveFile(title, defaultPath string, confirmOverwrite bool, filters []FileFilter) (string, error) {
	comdlg32 := syscall.NewLazyDLL("comdlg32.dll")
	getSaveFileName := comdlg32.NewProc("GetSaveFileNameW")

	var args _OPENFILENAME
	args.StructSize = uint32(unsafe.Sizeof(args))

	if title != "" {
		args.Title = syscall.StringToUTF16Ptr(title)
	}
	if defaultPath != "" {
		args.InitialDir = syscall.StringToUTF16Ptr(defaultPath)
	}
	if confirmOverwrite {
		args.Flags |= 0x00000002 // OFN_OVERWRITEPROMPT
	}
	args.Filter = &windowsFilters(filters)[0]

	res := [1024]uint16{}
	args.File = &res[0]
	args.MaxFile = uint32(len(res))

	_, _, _ = getSaveFileName.Call(uintptr(unsafe.Pointer(&args)))
	return syscall.UTF16ToString(res[:]), nil
}

func PickFolder(title, defaultPath string) (string, error) {
	// TODO
	return "", nil
}

func windowsFilters(filters []FileFilter) []uint16 {
	var res []uint16
	for _, f := range filters {
		res = append(res, utf16.Encode([]rune(f.Name))...)
		res = append(res, 0)
		for _, e := range f.Exts {
			res = append(res, uint16('*'))
			res = append(res, utf16.Encode([]rune(e))...)
			res = append(res, uint16(';'))
		}
		res = append(res, 0)
	}
	if res != nil {
		res = append(res, 0)
	}
	return res
}

type _OPENFILENAME struct {
	StructSize      uint32
	Owner           uintptr
	Instance        uintptr
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
	CustData        uintptr
	FnHook          uintptr
	TemplateName    *uint16
	PvReserved      unsafe.Pointer
	DwReserved      uint32
	FlagsEx         uint32
}
