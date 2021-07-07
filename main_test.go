package zenity_test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func skip(err error) (bool, error) {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		if _, ok := err.(*exec.Error); ok {
			// zenity was not found in path
			return true, err
		}
		if err != nil && os.Getenv("DISPLAY") == "" {
			// no display
			return true, fmt.Errorf("no display: %w", err)
		}
	}
	return false, err
}
