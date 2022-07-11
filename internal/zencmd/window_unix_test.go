//go:build !windows && !darwin

package zencmd

import (
	"os"
	"testing"
)

func Test_getPidToPpidMap(t *testing.T) {
	got, err := getPidToPpidMap()
	if err != nil {
		t.Fatalf("getPidToPpidMap() error = %v", err)
	}
	if ppid := got[os.Getpid()]; ppid != os.Getppid() {
		t.Errorf("getPidToPpidMap()[%d] = %d; want %d", os.Getpid(), ppid, os.Getppid())
	}
}

func Test_getPidToWindowMap(t *testing.T) {
	got, err := getPidToWindowMap()
	if err != nil {
		if os.Getenv("DISPLAY") == "" {
			t.Skip("skipping:", err)
		}
		t.Fatalf("getPidToWindowMap() error = %v", err)
	}
	if len(got) == 0 {
		t.Errorf("getPidToWindowMap() %v", got)
	}
}
