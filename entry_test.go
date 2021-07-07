package zenity_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
	"go.uber.org/goleak"
)

func ExampleEntry() {
	zenity.Entry("Enter new text:",
		zenity.Title("Add a new entry"))
	// Output:
}

func TestEntry_timeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	_, err := zenity.Entry("", zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestEntry_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.Entry("", zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}

func TestEntry_script(t *testing.T) {
	tests := []struct {
		name string
		call string
		opts []zenity.Option
		want string
		err  error
	}{
		{name: "Cancel", call: "cancel", want: "", err: zenity.ErrCanceled},
		{name: "123", call: "enter 123", want: "123", err: nil},
		{name: "abc", call: "enter abc", want: "abc", err: nil},
		{name: "Password", call: "press OK", want: "xpto", err: nil,
			opts: []zenity.Option{zenity.HideText(), zenity.EntryText("xpto")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := zenity.Entry(fmt.Sprintf("Please, %s.", tt.call), tt.opts...)
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			if text != tt.want || err != tt.err {
				t.Errorf("Entry() = %q, %v; want %q, %v", text, err, tt.want, tt.err)
			}
		})
	}
}
