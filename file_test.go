package zenity_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ncruces/zenity"
	"go.uber.org/goleak"
)

const defaultPath = ``
const defaultName = ``

func ExampleSelectFile() {
	zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		})
}

func ExampleSelectFileMultiple() {
	zenity.SelectFileMultiple(
		zenity.Filename(defaultPath),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		})
}

func ExampleSelectFileSave() {
	zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename(defaultName),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		})
}

func ExampleSelectFile_directory() {
	zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.Directory())
}

func ExampleSelectFileMultiple_directory() {
	zenity.SelectFileMultiple(
		zenity.Filename(defaultPath),
		zenity.Directory())
}

var fileFuncs = []struct {
	name string
	fn   func(...zenity.Option) (string, error)
}{
	{"Open", zenity.SelectFile},
	{"Save", zenity.SelectFileSave},
	{"Directory", func(o ...zenity.Option) (string, error) {
		return zenity.SelectFile(append(o, zenity.Directory())...)
	}},
	{"Multiple", func(o ...zenity.Option) (string, error) {
		_, err := zenity.SelectFileMultiple(append(o, zenity.Directory())...)
		return "", err
	}},
	{"MultipleDirectory", func(o ...zenity.Option) (string, error) {
		_, err := zenity.SelectFileMultiple(o...)
		return "", err
	}},
}

func TestSelectFile_timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	for _, tt := range fileFuncs {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second/5)
			defer cancel()

			_, err := tt.fn(zenity.Context(ctx))
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			if !os.IsTimeout(err) {
				t.Error("did not timeout:", err)
			}
		})
	}
}

func TestSelectFile_cancel(t *testing.T) {
	for _, tt := range fileFuncs {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := tt.fn(zenity.Context(ctx))
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			if !errors.Is(err, context.Canceled) {
				t.Error("was not canceled:", err)
			}
		})
	}
}

func TestSelectFile_script(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Run("Cancel", func(t *testing.T) {
		zenity.Info(fmt.Sprintf("In the file selection dialog, cancel."))
		str, err := zenity.SelectFile()
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if str != "" || err != zenity.ErrCanceled {
			t.Errorf("SelectFile() = %q, %v; want %q, %v", str, err, "", zenity.ErrCanceled)
		}
	})
	t.Run("File", func(t *testing.T) {
		zenity.Info(fmt.Sprintf("In the file selection dialog, pick any file."))
		str, err := zenity.SelectFile()
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if str == "" || err != nil {
			t.Errorf("SelectFile() = %q, %v; want [path], nil", str, err)
		}
		if _, serr := os.Stat(str); serr != nil {
			t.Errorf("SelectFile() = %q, %v; %v", str, err, serr)
		}
	})
	t.Run("Directory", func(t *testing.T) {
		zenity.Info(fmt.Sprintf("In the file selection dialog, pick any directory."))
		str, err := zenity.SelectFile(zenity.Directory())
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if str == "" || err != nil {
			t.Errorf("SelectFile() = %q, %v; want [path], nil", str, err)
		}
		if s, serr := os.Stat(str); serr != nil {
			t.Errorf("SelectFile() = %q, %v; %v", str, err, serr)
		} else if !s.IsDir() {
			t.Errorf("SelectFile() = %q, %v; not a directory", str, err)
		}
	})
}

func TestSelectFileMultiple_script(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Run("Cancel", func(t *testing.T) {
		zenity.Info(fmt.Sprintf("In the file selection dialog, cancel."))
		lst, err := zenity.SelectFileMultiple()
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if lst != nil || err != zenity.ErrCanceled {
			t.Errorf("SelectFileMultiple() = %v, %v; want nil, %v", lst, err, zenity.ErrCanceled)
		}
	})
	t.Run("Files", func(t *testing.T) {
		zenity.Info(fmt.Sprintf("In the file selection dialog, pick two files."))
		lst, err := zenity.SelectFileMultiple()
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if lst == nil || err != nil {
			t.Errorf("SelectFileMultiple() = %v, %v; want [path, path], nil", lst, err)
		}
		for _, str := range lst {
			if _, serr := os.Stat(str); serr != nil {
				t.Errorf("SelectFileMultiple() = %q, %v; %v", lst, err, serr)
			}
		}
	})
	t.Run("Directories", func(t *testing.T) {
		zenity.Info(fmt.Sprintf("In the file selection dialog, pick two directories."))
		lst, err := zenity.SelectFileMultiple(zenity.Directory())
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if errors.Is(err, zenity.ErrUnsupported) {
			t.Skip("was not unsupported:", err)
		}
		if lst == nil || err != nil {
			t.Errorf("SelectFileMultiple() = %v, %v; want [path, path], nil", lst, err)
		}
		for _, str := range lst {
			if s, serr := os.Stat(str); serr != nil {
				t.Errorf("SelectFileMultiple() = %q, %v; %v", str, err, serr)
			} else if !s.IsDir() {
				t.Errorf("SelectFileMultiple() = %q, %v; not a directory", str, err)
			}
		}
	})
}

func TestSelectFileSave_script(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Run("Cancel", func(t *testing.T) {
		zenity.Info(fmt.Sprintf("In the file save dialog, cancel."))
		str, err := zenity.SelectFileSave()
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if str != "" || err != zenity.ErrCanceled {
			t.Errorf("SelectFileSave() = %q, %v; want %q, %v", str, err, "", zenity.ErrCanceled)
		}
	})
	t.Run("Name", func(t *testing.T) {
		zenity.Info(fmt.Sprintf("In the file save dialog, press OK."))
		str, err := zenity.SelectFileSave(
			zenity.ConfirmOverwrite(),
			zenity.Filename("Χρτο.go"),
			zenity.FileFilter{"Go files", []string{"*.go"}},
		)
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if _, name := filepath.Split(str); name != "Χρτο.go" || err != nil {
			t.Errorf("SelectFileSave() = %q, %v; want %q, nil", str, err, "Χρτο.go")
		}
	})
}
