// +build !windows,!js

package zenutil

import (
	"strconv"
)

type progressDialog struct {
	err     error
	done    chan struct{}
	lines   chan string
	percent bool
	max     int
}

func (d *progressDialog) send(line string) error {
	select {
	case d.lines <- line:
		return nil
	case <-d.done:
		return d.err
	}
}

func (d *progressDialog) Close() error {
	close(d.lines)
	<-d.done
	return d.err
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
