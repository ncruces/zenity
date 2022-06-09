//go:build !windows && !darwin

package zenutil

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

// ParseWindowId is internal.
func ParseWindowId(id string) int {
	wid, _ := strconv.ParseUint(id, 0, 64)
	return int(wid)
}

// GetParentWindowId is internal.
func GetParentWindowId() int {
	buf, err := exec.Command("ps", "-xo", "pid=,ppid=").CombinedOutput()
	if err != nil {
		return 0
	}

	ppids := map[int]int{}
	reader := bytes.NewReader(buf)
	for {
		var pid, ppid int
		_, err := fmt.Fscan(reader, &pid, &ppid)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0
		}
		ppids[pid] = ppid
	}

	// Find the relevant pid and window id.
	return 0
}
