package zenity_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ncruces/zenity"
	"go.uber.org/goleak"
)

func ExampleProgress() {
	dlg, err := zenity.Progress(
		zenity.Title("Update System Logs"))
	if err != nil {
		return
	}
	defer dlg.Close()

	dlg.Text("Scanning mail logs...")
	dlg.Value(0)
	time.Sleep(time.Second)

	dlg.Value(25)
	time.Sleep(time.Second)

	dlg.Text("Updating mail logs...")
	dlg.Value(50)
	time.Sleep(time.Second)

	dlg.Text("Resetting cron jobs...")
	dlg.Value(75)
	time.Sleep(time.Second)

	dlg.Text("Rebooting system...")
	dlg.Value(100)
	time.Sleep(time.Second)

	dlg.Complete()
	time.Sleep(time.Second)
}

func ExampleProgress_pulsate() {
	dlg, err := zenity.Progress(
		zenity.Title("Update System Logs"),
		zenity.Pulsate())
	if err != nil {
		return
	}
	defer dlg.Close()

	dlg.Text("Scanning mail logs...")
	time.Sleep(time.Second)

	dlg.Text("Updating mail logs...")
	time.Sleep(time.Second)

	dlg.Text("Resetting cron jobs...")
	time.Sleep(time.Second)

	dlg.Text("Rebooting system...")
	time.Sleep(time.Second)

	dlg.Complete()
	time.Sleep(time.Second)
}

func TestProgress_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.Progress(zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}

func TestProgress_cancelAfter(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())

	dlg, err := zenity.Progress(zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if err != nil {
		t.Fatal(err)
	}

	go cancel()
	<-dlg.Done()
	err = dlg.Close()
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}

func TestProgress_examples(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ExampleProgress()
	ExampleProgress_pulsate()
}
