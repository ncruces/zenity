package zenity

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_splitDirAndName(t *testing.T) {
	tempDir := os.TempDir()
	tests := []struct {
		path     string
		wantDir  string
		wantName string
	}{
		// filepath.Split test cases
		{"a/b", "a/", "b"},
		{"a/b/", "a/b/", ""},
		{"a/", "a/", ""},
		{"a", "", "a"},
		{"/", "/", ""},
		// we split differently if we know it's a directory
		{tempDir, tempDir, ""},
		{filepath.Clean(tempDir), filepath.Clean(tempDir), ""},
		{filepath.Join(tempDir, "a"), filepath.Clean(tempDir) + string(filepath.Separator), "a"},
	}

	for i, tt := range tests {
		gotDir, gotName := splitDirAndName(tt.path)
		if gotDir != tt.wantDir {
			t.Errorf("splitDirAndName[%d].dir = %q; want %q", i, gotDir, tt.wantDir)
		}
		if gotName != tt.wantName {
			t.Errorf("splitDirAndName[%d].name = %q; want %q", i, gotName, tt.wantName)
		}
	}
}
