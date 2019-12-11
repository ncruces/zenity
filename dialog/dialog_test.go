package dialog

import "testing"

func TestOpenFile(t *testing.T) {
	ret, err := OpenFile("", "", []FileFilter{
		{"Go files", []string{".go"}},
		{"Web files", []string{".html", ".js", ".css"}},
		{"Image files", []string{".png", ".ico"}},
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", ret)
	}
}

func TestOpenFiles(t *testing.T) {
	ret, err := OpenFiles("", "", []FileFilter{
		{"Go files", []string{".go"}},
		{"Web files", []string{".html", ".js", ".css"}},
		{"Image files", []string{".png", ".ico"}},
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", ret)
	}
}

func TestSaveFile(t *testing.T) {
	ret, err := SaveFile("", "", true, []FileFilter{
		{"Go files", []string{".go"}},
		{"Web files", []string{".html", ".js", ".css"}},
		{"Image files", []string{".png", ".ico"}},
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", ret)
	}
}

func TestPickFolder(t *testing.T) {
	ret, err := PickFolder("", "")

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", ret)
	}
}
