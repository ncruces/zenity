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

func ExampleList() {
	zenity.List(
		"Select items from the list below:",
		[]string{"apples", "oranges", "bananas", "strawberries"},
		zenity.Title("Select items from the list"),
		zenity.DisallowEmpty(),
	)
	// Output:
}

func ExampleListItems() {
	zenity.ListItems(
		"Select items from the list below:",
		"apples", "oranges", "bananas", "strawberries")
	// Output:
}

func ExampleListMultiple() {
	zenity.ListMultiple(
		"Select items from the list below:",
		[]string{"apples", "oranges", "bananas", "strawberries"},
		zenity.Title("Select items from the list"),
		zenity.DefaultItems("apples", "bananas"),
	)
	// Output:
}

func ExampleListMultipleItems() {
	zenity.ListMultipleItems(
		"Select items from the list below:",
		"apples", "oranges", "bananas", "strawberries")
	// Output:
}

func TestList_timeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	_, err := zenity.List("", nil, zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestList_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.List("", nil, zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
