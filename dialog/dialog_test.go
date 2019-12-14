package dialog

import "testing"

func TestOpenFile(t *testing.T) {
	res, err := OpenFile("", "", []FileFilter{
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
	res, err := OpenFiles("", "", []FileFilter{
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
	res, err := SaveFile("", "", true, []FileFilter{
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
	res, err := PickFolder("", "")

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}
