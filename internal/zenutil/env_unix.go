//go:build !windows && !darwin

package zenutil

import "strconv"

// These are internal.
var (
	Separator = "\x1e"
	LineBreak = "\n"
)

func ParseWindowId(id string) int {
	wid, _ := strconv.ParseUint(id, 0, 64)
	return int(wid)
}
