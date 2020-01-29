package zenity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
)

func ExampleError() {
	zenity.Error("An error has occurred.",
		zenity.Title("Error"),
		zenity.Icon(zenity.ErrorIcon))
	// Output:
}

func ExampleInfo() {
	zenity.Info("All updates are complete.",
		zenity.Title("Information"),
		zenity.Icon(zenity.InfoIcon))
	// Output:
}

func ExampleWarning() {
	zenity.Warning("Are you sure you want to proceed?",
		zenity.Title("Warning"),
		zenity.Icon(zenity.WarningIcon))
	// Output:
}

func ExampleQuestion() {
	zenity.Question("Are you sure you want to proceed?",
		zenity.Title("Question"),
		zenity.Icon(zenity.QuestionIcon))
	// Output:
}

var msgFuncs = []func(string, ...zenity.Option) (bool, error){
	zenity.Error,
	zenity.Info,
	zenity.Warning,
	zenity.Question,
}

func TestMessageTimeout(t *testing.T) {
	for _, f := range msgFuncs {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)

		_, err := f("text", zenity.Context(ctx))
		if !os.IsTimeout(err) {
			t.Error("did not timeout:", err)
		}

		cancel()
	}
}

func TestMessageCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	for _, f := range msgFuncs {
		_, err := f("text", zenity.Context(ctx))
		if !errors.Is(err, context.Canceled) {
			t.Error("was not canceled:", err)
		}
	}
}
