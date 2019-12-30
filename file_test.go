package zenity

import "testing"

const defaultPath = ""

func TestSelectFile(t *testing.T) {
	res, err := SelectFile("", defaultPath, []FileFilter{
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

func TestSelectFileMutiple(t *testing.T) {
	res, err := SelectFileMutiple("", defaultPath, []FileFilter{
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

func TestSelectFileSave(t *testing.T) {
	res, err := SelectFileSave("", defaultPath, true, []FileFilter{
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

func TestSelectDirectory(t *testing.T) {
	res, err := SelectDirectory("", defaultPath)

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}
