package zenity

import (
	"errors"
	"os/exec"
	"reflect"
	"runtime"
	"testing"

	"github.com/ncruces/zenity/internal/zenutil"
)

func Test_quoteAccelerators(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		text string
		want string
	}{
		{name: "None", text: "abc", want: "abc"},
		{name: "One", text: "&abc", want: "&&abc"},
		{name: "Two", text: "&a&bc", want: "&&a&&bc"},
		{name: "Three", text: "ab&&c", want: "ab&&&&c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quoteAccelerators(tt.text); got != tt.want {
				t.Errorf("quoteAccelerators(%q) = %q; want %q", tt.text, got, tt.want)
			}
		})
	}
}

func Test_quoteMnemonics(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		text string
		want string
	}{
		{name: "None", text: "abc", want: "abc"},
		{name: "One", text: "_abc", want: "__abc"},
		{name: "Two", text: "_a_bc", want: "__a__bc"},
		{name: "Three", text: "ab__c", want: "ab____c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quoteMnemonics(tt.text); got != tt.want {
				t.Errorf("quoteMnemonics(%q) = %q; want %q", tt.text, got, tt.want)
			}
		})
	}
}

func Test_quoteMarkup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		text string
		want string
	}{
		{name: "None", text: `abc`, want: "abc"},
		{name: "LT", text: `<`, want: "&lt;"},
		{name: "Amp", text: `&`, want: "&amp;"},
		{name: "Quot", text: `"`, want: "&#34;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quoteMarkup(tt.text); got != tt.want {
				t.Errorf("quoteMarkup(%q) = %q; want %q", tt.text, got, tt.want)
			}
		})
	}
}

func Test_appendGeneral(t *testing.T) {
	t.Parallel()
	got := appendGeneral(nil, options{
		title:   ptr("Title"),
		attach:  12345,
		modal:   true,
		display: ":1",
		class:   "Class",
		name:    "Name",
	})
	want := []string{
		"--title", "Title",
		"--attach", "12345",
		"--modal",
		"--display", ":1",
		"--class", "Class",
		"--name", "Name",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("appendTitle() = %v; want %v", got, want)
	}
}

func Test_appendButtons(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		opts options
		want []string
	}{
		{name: "OK", opts: options{okLabel: ptr("OK")}, want: []string{"--ok-label", "OK"}},
		{name: "Cancel", opts: options{cancelLabel: ptr("Cancel")}, want: []string{"--cancel-label", "Cancel"}},
		{name: "Extra", opts: options{extraButton: ptr("Extra")}, want: []string{"--extra-button", "Extra"}},
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
	t.Parallel()
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
	t.Parallel()
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
	cancel := exit1Cmd().Run()
	t.Parallel()

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
	sentinel := errors.New("sentinel")
	cancel := exit1Cmd().Run()
	zenutil.Separator = "|"
	t.Parallel()

	if out, err := lstResult(options{}, []byte(""), nil); !reflect.DeepEqual(out, []string{}) || err != nil {
		t.Errorf(`lstResult("", nil) = %v, %v`, out, err)
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

func Test_pwdResult(t *testing.T) {
	username := options{username: true}
	sentinel := errors.New("sentinel")
	cancel := exit1Cmd().Run()
	t.Parallel()

	if usr, pwd, err := pwdResult("|", options{}, []byte(""), nil); usr != "" || pwd != "" || err != nil {
		t.Errorf(`pwdResult("", nil) = %v, %q, %q`, usr, pwd, err)
	}
	if usr, pwd, err := pwdResult("|", options{}, []byte("out"), nil); usr != "" || pwd != "out" || err != nil {
		t.Errorf(`pwdResult("out", nil) = %v, %q, %q`, usr, pwd, err)
	}
	if usr, pwd, err := pwdResult("|", username, []byte("one|two"), nil); usr != "one" || pwd != "two" || err != nil {
		t.Errorf(`pwdResult("one|two", nil) = %v, %q, %q`, usr, pwd, err)
	}
	if usr, pwd, err := pwdResult("|", options{}, []byte("one|two"), nil); usr != "" || pwd != "one|two" || err != nil {
		t.Errorf(`pwdResult("one|two", nil) = %v, %q, %q`, usr, pwd, err)
	}
	if usr, pwd, err := pwdResult("|", options{}, []byte("out"), sentinel); usr != "" || pwd != "" || err != sentinel {
		t.Errorf(`pwdResult("out", error) = %v, %q, %q`, usr, pwd, err)
	}
	if usr, pwd, err := pwdResult("|", options{}, []byte("out"), cancel); usr != "" || pwd != "" || err != ErrCanceled {
		t.Errorf(`pwdResult("out", cancel) = %v, %q, %q`, usr, pwd, err)
	}
}

func exit1Cmd() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/k", "exit", "1")
	}
	return exec.Command("false")
}
