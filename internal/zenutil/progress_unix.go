//go:build !windows

package zenutil

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

type progressDialog struct {
	ctx     context.Context
	cmd     *exec.Cmd
	max     int
	percent bool
	closed  int32
	lines   chan string
	done    chan struct{}
	err     error
}

func (d *progressDialog) send(line string) error {
	select {
	case d.lines <- line:
		return nil
	case <-d.done:
		return d.err
	}
}

func (d *progressDialog) Text(text string) error {
	return d.send("#" + text)
}

func (d *progressDialog) Value(value int) error {
	if d.percent {
		return d.send(strconv.FormatFloat(100*float64(value)/float64(d.max), 'f', -1, 64))
	} else {
		return d.send(strconv.Itoa(value))
	}
}

func (d *progressDialog) MaxValue() int {
	return d.max
}

func (d *progressDialog) Done() <-chan struct{} {
	return d.done
}

func (d *progressDialog) Complete() error {
	err := d.Value(d.max)
	close(d.lines)
	return err
}

func (d *progressDialog) Close() error {
	atomic.StoreInt32(&d.closed, 1)
	d.cmd.Process.Signal(os.Interrupt)
	<-d.done
	return d.err
}

func (d *progressDialog) wait(extra *string, out *bytes.Buffer) {
	err := d.cmd.Wait()
	if cerr := d.ctx.Err(); cerr != nil {
		err = cerr
	}
	if eerr, ok := err.(*exec.ExitError); ok {
		switch {
		case eerr.ExitCode() == -1 && atomic.LoadInt32(&d.closed) != 0:
			err = nil
		case eerr.ExitCode() == 1:
			if extra != nil && *extra+"\n" == string(out.Bytes()) {
				err = ErrExtraButton
			} else {
				err = ErrCanceled
			}
		}
	}
	d.err = err
	close(d.done)
}

func (d *progressDialog) pipe(w io.WriteCloser) {
	defer w.Close()
	var timeout = time.Second
	if runtime.GOOS == "darwin" {
		timeout = 40 * time.Millisecond
	}
	for {
		var line string
		select {
		case s, ok := <-d.lines:
			if !ok {
				return
			}
			line = s
		case <-d.ctx.Done():
			return
		case <-d.done:
			return
		case <-time.After(timeout):
			// line = ""
		}
		if _, err := w.Write([]byte(line + "\n")); err != nil {
			return
		}
	}
}
