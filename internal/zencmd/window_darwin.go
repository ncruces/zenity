package zencmd

import (
	"strconv"

	"golang.org/x/sys/unix"
)

type any = interface{}

// ParseWindowId is internal.
func ParseWindowId(id string) any {
	if pid, err := strconv.ParseUint(id, 0, 64); err == nil {
		return int(pid)
	}
	return id
}

// GetParentWindowId is internal.
func GetParentWindowId(pid int) int {
	for {
		kinfo, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
		if err != nil {
			return 0
		}
		ppid := kinfo.Eproc.Ppid
		switch ppid {
		case 0:
			return 0
		case 1:
			return pid
		default:
			pid = int(ppid)
		}
	}
}
