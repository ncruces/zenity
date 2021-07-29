// +build !windows,!darwin

package zenutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestRun(t *testing.T) {
	_, err := Run(nil, []string{"--version"})
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun_context(t *testing.T) {
	_, err := Run(context.TODO(), []string{"--version"})
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunProgress(t *testing.T) {
	_, err := RunProgress(nil, 100, nil, []string{"--version"})
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func skip(err error) (bool, error) {
	if _, ok := err.(*exec.Error); ok {
		// zenity was not found in path
		return true, err
	}
	if err != nil && os.Getenv("DISPLAY") == "" && os.Getenv("WSL_DISTRO_NAME") == "" {
		// no display, not WSL
		return true, fmt.Errorf("no display: %w", err)
	}
	return false, err
}
