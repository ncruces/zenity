package zenity

import "testing"

const defaultPath = ""

func TestSelectFile(t *testing.T) {
	res, err := SelectFile(Filename(defaultPath), FileFilters{
		{"Go files", []string{"*.go"}},
		{"Web files", []string{"*.html", "*.js", "*.css"}},
		{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
	}.New())

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestSelectFileMutiple(t *testing.T) {
	res, err := SelectFileMutiple(Filename(defaultPath), FileFilters{
		{"Go files", []string{"*.go"}},
		{"Web files", []string{"*.html", "*.js", "*.css"}},
		{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
	}.New())

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestSelectFileSave(t *testing.T) {
	res, err := SelectFileSave(Filename(defaultPath), ConfirmOverwrite, FileFilters{
		{"Go files", []string{"*.go"}},
		{"Web files", []string{"*.html", "*.js", "*.css"}},
		{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
	}.New())

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestSelectDirectory(t *testing.T) {
	res, err := SelectFile(Directory, Filename(defaultPath))

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestSelectDirectoryMultiple(t *testing.T) {
	res, err := SelectFileMutiple(Directory, Filename(defaultPath))

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}
