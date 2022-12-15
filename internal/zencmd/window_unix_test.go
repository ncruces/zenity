//go:build !windows && !darwin

package zencmd

import (
	"os"
	"testing"
)

func TestParseWindowId(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		text string
		want int
	}{
		{name: "Zero", text: "0", want: 0},
		{name: "Dec", text: "10", want: 10},
		{name: "Hex", text: "0700", want: 0700},
		{name: "Oct", text: "0xFF", want: 0xff},
		{name: "Error", text: "a", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseWindowId(tt.text); got != tt.want {
				t.Errorf("ParseWindowId(%q) = %v; want %v", tt.text, got, tt.want)
			}
		})
	}
}

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
