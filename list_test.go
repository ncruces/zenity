package zenity_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
)

func ExampleList() {
	indices, err := zenity.List("Choose from the list:",
		[]string{"Yes", "No", "Skip", "Skip all"},
		zenity.Title("Chooser"),
		zenity.Icon(zenity.InfoIcon),
		zenity.MultipleSelection())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: '%s'\n", err.Error())
	} else {
		fmt.Fprintf(os.Stderr, "Indices: '%v'\n", indices)
	}
	// Output:
}
func TestListTimeout(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	_, err := zenity.List("text", []string{"Yes", "No", "Skip", "Skip all"}, zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
	cancel()
}

func TestListCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.List("text", []string{"Yes", "No", "Skip", "Skip all"}, zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}

}
