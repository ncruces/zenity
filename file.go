package zenity

import (
	"os"
	"path/filepath"
)

func splitDirAndName(path string) (dir, name string) {
	path = filepath.Clean(path)
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return path, ""
	}
	return filepath.Split(path)
}
