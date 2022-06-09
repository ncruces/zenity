package zenutil

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// ParseWindowId is internal.
func ParseWindowId(id string) any {
	if pid, err := strconv.ParseUint(id, 0, 64); err == nil {
		return int(pid)
	}
	return id
}

// GetParentWindowId is internal.
func GetParentWindowId() int {
	buf, err := sysctlKernProcAll()
	if err != nil {
		return 0
	}

	type kinfo_proc struct {
		_    [40]byte
		Pid  int32
		_    [516]byte
		PPid int32
		_    [84]byte
	}
	const size = int(unsafe.Sizeof(kinfo_proc{}))

	ppids := map[int]int{}
	for i := 0; i+size < len(buf); i += size {
		p := (*kinfo_proc)(unsafe.Pointer(&buf[i]))
		ppids[int(p.Pid)] = int(p.PPid)
	}

	pid := os.Getppid()
	for {
		ppid := ppids[pid]
		switch ppid {
		case 0:
			return 0
		case 1:
			return pid
		default:
			pid = ppid
		}
	}
}

func sysctlKernProcAll() ([]byte, error) {
	const (
		CTRL_KERN     = 1
		KERN_PROC     = 14
		KERN_PROC_ALL = 0
	)

	mib := [4]int32{CTRL_KERN, KERN_PROC, KERN_PROC_ALL, 0}
	size := uintptr(0)

	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(len(mib)),
		0,
		uintptr(unsafe.Pointer(&size)),
		0,
		0)

	if errno != 0 {
		return nil, errno
	}

	bs := make([]byte, size)
	_, _, errno = syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(len(mib)),
		uintptr(unsafe.Pointer(&bs[0])),
		uintptr(unsafe.Pointer(&size)),
		0,
		0)

	if errno != 0 {
		return nil, errno
	}
	return bs[0:size], nil
}
