package zenity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
	"go.uber.org/goleak"
)

func ExampleError() {
	zenity.Error("An error has occurred.",
		zenity.Title("Error"),
		zenity.ErrorIcon)
	// Output:
}

func ExampleInfo() {
	zenity.Info("All updates are complete.",
		zenity.Title("Information"),
		zenity.InfoIcon)
	// Output:
}

func ExampleWarning() {
	zenity.Warning("Are you sure you want to proceed?",
		zenity.Title("Warning"),
		zenity.WarningIcon)
	// Output:
}

func ExampleQuestion() {
	zenity.Question("Are you sure you want to proceed?",
		zenity.Title("Question"),
		zenity.QuestionIcon)
	// Output:
}

var msgFuncs = []func(string, ...zenity.Option) error{
	zenity.Error,
	zenity.Info,
	zenity.Warning,
	zenity.Question,
}

func TestMessage_timeout(t *testing.T) {
	for _, f := range msgFuncs {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)

		err := f("text", zenity.Context(ctx))
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if !os.IsTimeout(err) {
			t.Error("did not timeout:", err)
		}

		cancel()
		goleak.VerifyNone(t)
	}
}

func TestMessage_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	for _, f := range msgFuncs {
		err := f("text", zenity.Context(ctx))
		if skip, err := skip(err); skip {
			t.Skip("skipping:", err)
		}
		if !errors.Is(err, context.Canceled) {
			t.Error("was not canceled:", err)
		}
	}
}
