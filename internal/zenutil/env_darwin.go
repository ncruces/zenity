package zenutil

import "strconv"

// These are internal.
var (
	Separator = "\x00"
	LineBreak = "\n"
)

func ParseWindowId(id string) any {
	if pid, err := strconv.ParseUint(id, 0, 64); err == nil {
		return int(pid)
	}
	return id
}
