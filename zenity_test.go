package zenity_test

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func skip(err error) (error, bool) {
	if _, ok := err.(*exec.Error); ok {
		// zenity/osascript/etc were not found in path
		return err, true
	}
	if err != nil && os.Getenv("DISPLAY") == "" && !(runtime.GOOS == "windows" || runtime.GOOS == "darwin") {
		// no display
		return errors.New("no display"), true
	}
	return nil, false
}
