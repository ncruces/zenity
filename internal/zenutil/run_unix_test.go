// +build !windows,!darwin

package zenutil

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestRun(t *testing.T) {
	_, err := Run(nil, []string{"--version"})
	if err, skip := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun_context(t *testing.T) {
	_, err := Run(context.TODO(), []string{"--version"})
	if err, skip := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunProgress(t *testing.T) {
	_, err := RunProgress(nil, 100, nil, []string{"--version"})
	if err, skip := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func skip(err error) (error, bool) {
	if _, ok := err.(*exec.Error); ok {
		// zenity/osascript/etc were not found in path
		return err, true
	}
	if err != nil && os.Getenv("DISPLAY") == "" {
		// no display
		return errors.New("no display"), true
	}
	return nil, false
}
