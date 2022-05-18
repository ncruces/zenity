package zenity_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
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
}

func ExampleListItems() {
	zenity.ListItems(
		"Select items from the list below:",
		"apples", "oranges", "bananas", "strawberries")
}

func ExampleListMultiple() {
	zenity.ListMultiple(
		"Select items from the list below:",
		[]string{"apples", "oranges", "bananas", "strawberries"},
		zenity.Title("Select items from the list"),
		zenity.DefaultItems("apples", "bananas"),
	)
}

func ExampleListMultipleItems() {
	zenity.ListMultipleItems(
		"Select items from the list below:",
		"apples", "oranges", "bananas", "strawberries")
}

func TestList_timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/5)
	defer cancel()

	_, err := zenity.List("", []string{""}, zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestList_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.List("", []string{""}, zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}

func TestList_script(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	items := []string{"apples", "oranges", "bananas", "strawberries"}
	tests := []struct {
		name string
		call string
		want string
		err  error
	}{
		{name: "Cancel", call: "cancel", err: zenity.ErrCanceled},
		{name: "Apples", call: "select apples", want: "apples"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := zenity.ListItems(fmt.Sprintf("Please, %s.", tt.call), items...)
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			if got != tt.want || err != tt.err {
				t.Errorf("List() = %q, %v; want %q, %v", got, err, tt.want, tt.err)
			}
		})
	}
}

func TestListMultiple_script(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	items := []string{"apples", "oranges", "bananas", "strawberries"}
	tests := []struct {
		name string
		call string
		want []string
		err  error
	}{
		{name: "Cancel", call: "cancel", err: zenity.ErrCanceled},
		{name: "Nothing", call: "select nothing", want: []string{}},
		{name: "Apples", call: "select apples", want: []string{"apples"}},
		{name: "Apples & Oranges", call: "select apples and oranges",
			want: []string{"apples", "oranges"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := zenity.ListMultipleItems(fmt.Sprintf("Please, %s.", tt.call), items...)
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			if !reflect.DeepEqual(got, tt.want) || err != tt.err {
				t.Errorf("ListMultiple() = %q, %v; want %v, %v", got, err, tt.want, tt.err)
			}
		})
	}
}
