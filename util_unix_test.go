//go:build !windows

package zenity

import (
	"errors"
	"os/exec"
	"reflect"
	"testing"

	"github.com/ncruces/zenity/internal/zenutil"
)

func Test_appendGeneral(t *testing.T) {
	got := appendGeneral(nil, options{
		title:  stringPtr("Title"),
		attach: 12345,
		modal:  true,
	})
	want := []string{"--title", "Title", "--attach", "12345", "--modal"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("appendTitle() = %v; want %v", got, want)
	}
}

func Test_appendButtons(t *testing.T) {
	tests := []struct {
		name string
		opts options
		want []string
	}{
		{name: "OK", opts: options{okLabel: stringPtr("OK")}, want: []string{"--ok-label", "OK"}},
		{name: "Cancel", opts: options{cancelLabel: stringPtr("Cancel")}, want: []string{"--cancel-label", "Cancel"}},
		{name: "Extra", opts: options{extraButton: stringPtr("Extra")}, want: []string{"--extra-button", "Extra"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendButtons(nil, tt.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendButtons() = %v; want %v", got, tt.want)
			}
		})
	}
}

func Test_appendWidthHeight(t *testing.T) {
	tests := []struct {
		name string
		opts options
		want []string
	}{
		{name: "Width", opts: options{width: 100}, want: []string{"--width", "100"}},
		{name: "Height", opts: options{height: 100}, want: []string{"--height", "100"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendWidthHeight(nil, tt.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendWidthHeight() = %v; want %v", got, tt.want)
			}
		})
	}
}

func Test_appendWindowIcon(t *testing.T) {
	tests := []struct {
		name string
		opts options
		want []string
	}{
		{name: "NoIcon", opts: options{windowIcon: NoIcon}, want: nil},
		{name: "Info", opts: options{windowIcon: InfoIcon}, want: []string{"--window-icon=info"}},
		{name: "Error", opts: options{windowIcon: ErrorIcon}, want: []string{"--window-icon=error"}},
		{name: "Warning", opts: options{windowIcon: WarningIcon}, want: []string{"--window-icon=warning"}},
		{name: "Question", opts: options{windowIcon: QuestionIcon}, want: []string{"--window-icon=question"}},
		{name: "Password", opts: options{windowIcon: PasswordIcon}, want: nil},
		{name: "File", opts: options{windowIcon: "png"}, want: []string{"--window-icon", "png"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendWindowIcon(nil, tt.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendWindowIcon() = %v; want %v", got, tt.want)
			}
		})
	}
}

func Test_strResult(t *testing.T) {
	sentinel := errors.New("sentinel")
	cancel := exec.Command("false").Run()

	if out, err := strResult(options{}, []byte("out"), nil); out != "out" || err != nil {
		t.Errorf(`strResult("out", nil) = %q, %v`, out, err)
	}
	if out, err := strResult(options{}, []byte("out"), sentinel); out != "" || err != sentinel {
		t.Errorf(`strResult("out", error) = %q, %v`, out, err)
	}
	if out, err := strResult(options{}, []byte("out"), cancel); out != "" || err != ErrCanceled {
		t.Errorf(`strResult("out", cancel) = %q, %v`, out, err)
	}
}

func Test_lstResult(t *testing.T) {
	zenutil.Separator = "|"
	sentinel := errors.New("sentinel")
	cancel := exec.Command("false").Run()

	if out, err := lstResult(options{}, []byte(""), nil); !reflect.DeepEqual(out, []string{}) || err != nil {
		t.Errorf(`lstResult("out", nil) = %v, %v`, out, err)
	}
	if out, err := lstResult(options{}, []byte("out"), nil); !reflect.DeepEqual(out, []string{"out"}) || err != nil {
		t.Errorf(`lstResult("out", nil) = %v, %v`, out, err)
	}
	if out, err := lstResult(options{}, []byte("one|two"), nil); !reflect.DeepEqual(out, []string{"one", "two"}) || err != nil {
		t.Errorf(`lstResult("one|two", nil) = %v, %v`, out, err)
	}
	if out, err := lstResult(options{}, []byte("out"), sentinel); out != nil || err != sentinel {
		t.Errorf(`lstResult("out", error) = %v, %v`, out, err)
	}
	if out, err := lstResult(options{}, []byte("out"), cancel); out != nil || err != ErrCanceled {
		t.Errorf(`lstResult("out", cancel) = %v, %v`, out, err)
	}
}
