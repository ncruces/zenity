//go:build !windows && !darwin

package zencmd

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

// ParseWindowId is internal.
func ParseWindowId(id string) int {
	wid, _ := strconv.ParseUint(id, 0, 64)
	return int(wid & math.MaxInt)
}

// GetParentWindowId is internal.
func GetParentWindowId(pid int) int {
	winids, err := getPidToWindowMap()
	if err != nil {
		return 0
	}

	ppids, err := getPidToPpidMap()
	if err != nil {
		return 0
	}

	for {
		if winid, ok := winids[pid]; ok {
			id, _ := strconv.Atoi(winid)
			return id
		}
		if ppid, ok := ppids[pid]; ok {
			pid = ppid
		} else {
			return 0
		}
	}
}

func getPidToPpidMap() (map[int]int, error) {
	out, err := exec.Command("ps", "-xo", "pid=,ppid=").Output()
	if err != nil {
		return nil, err
	}

	ppids := map[int]int{}
	reader := bytes.NewReader(out)
	for {
		var pid, ppid int
		_, err := fmt.Fscan(reader, &pid, &ppid)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ppids[pid] = ppid
	}
	return ppids, nil
}

func getPidToWindowMap() (map[int]string, error) {
	ids, err := getWindowIDs()
	if err != nil {
		return nil, err
	}

	var pid int
	winids := map[int]string{}
	for _, id := range ids {
		pid, err = getWindowPid(id)
		if err != nil {
			continue
		}
		winids[pid] = id
	}
	if err != nil && len(winids) == 0 {
		return nil, err
	}
	return winids, nil
}

func getWindowIDs() ([]string, error) {
	out, err := exec.Command("xprop", "-root", "0i", "\t$0+", "_NET_CLIENT_LIST").Output()
	if err != nil {
		return nil, err
	}

	if _, out, cut := bytes.Cut(out, []byte("\t")); cut {
		return strings.Split(string(out), ", "), nil
	}
	return nil, fmt.Errorf("xprop: unexpected output: %q", out)
}

func getWindowPid(id string) (int, error) {
	out, err := exec.Command("xprop", "-id", id, "0i", "\t$0", "_NET_WM_PID").Output()
	if err != nil {
		return 0, err
	}

	if _, out, cut := bytes.Cut(out, []byte("\t")); cut {
		return strconv.Atoi(string(out))
	}
	return 0, fmt.Errorf("xprop: unexpected output: %q", out)
}
