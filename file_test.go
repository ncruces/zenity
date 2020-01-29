package zenity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
)

const defaultPath = ""
const defaultName = ""

func ExampleSelectFile() {
	zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		})
	// Output:
}

func ExampleSelectFileMutiple() {
	zenity.SelectFileMutiple(
		zenity.Filename(defaultPath),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		})
	// Output:
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
	// Output:
}

func ExampleSelectFile_directory() {
	zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.Directory())
	// Output:
}

func ExampleSelectFileMutiple_directory() {
	zenity.SelectFileMutiple(
		zenity.Filename(defaultPath),
		zenity.Directory())
	// Output:
}

var fileFuncs = []func(...zenity.Option) (string, error){
	zenity.SelectFile,
	zenity.SelectFileSave,
	func(o ...zenity.Option) (string, error) {
		return zenity.SelectFile(append(o, zenity.Directory())...)
	},
	func(o ...zenity.Option) (string, error) {
		_, err := zenity.SelectFileMutiple(append(o, zenity.Directory())...)
		return "", err
	},
	func(o ...zenity.Option) (string, error) {
		_, err := zenity.SelectFileMutiple(o...)
		return "", err
	},
}

func TestFileTimeout(t *testing.T) {
	for _, f := range fileFuncs {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)

		_, err := f(zenity.Context(ctx))
		if !os.IsTimeout(err) {
			t.Error("did not timeout:", err)
		}

		cancel()
	}
}

func TestFileCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	for _, f := range fileFuncs {
		_, err := f(zenity.Context(ctx))
		if !errors.Is(err, context.Canceled) {
			t.Error("was not canceled:", err)
		}
	}
}
