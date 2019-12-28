package dialog

import "testing"

const defaultPath = ""

func TestOpenFile(t *testing.T) {
	res, err := OpenFile("", defaultPath, []FileFilter{
		{"Go files", []string{".go"}},
		{"Web files", []string{".html", ".js", ".css"}},
		{"Image files", []string{".png", ".gif", ".ico", ".jpg", ".webp"}},
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestOpenFiles(t *testing.T) {
	res, err := OpenFiles("", defaultPath, []FileFilter{
		{"Go files", []string{".go"}},
		{"Web files", []string{".html", ".js", ".css"}},
		{"Image files", []string{".png", ".gif", ".ico", ".jpg", ".webp"}},
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestSaveFile(t *testing.T) {
	res, err := SaveFile("", defaultPath, true, []FileFilter{
		{"Go files", []string{".go"}},
		{"Web files", []string{".html", ".js", ".css"}},
		{"Image files", []string{".png", ".gif", ".ico", ".jpg", ".webp"}},
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestPickFolder(t *testing.T) {
	res, err := PickFolder("", defaultPath)

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}
