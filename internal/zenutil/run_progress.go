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

func (m *progressDialog) send(line string) error {
	select {
	case m.lines <- line:
		return nil
	case <-m.done:
		return m.err
	}
}

func (m *progressDialog) Close() error {
	close(m.lines)
	<-m.done
	return m.err
}

func (m *progressDialog) Text(text string) error {
	return m.send("#" + text)
}

func (m *progressDialog) Value(value int) error {
	if m.percent {
		return m.send(strconv.FormatFloat(100*float64(value)/float64(m.max), 'f', -1, 64))
	} else {
		return m.send(strconv.Itoa(value))
	}
}

func (m *progressDialog) MaxValue() int {
	return m.max
}

func (m *progressDialog) Done() <-chan struct{} {
	return m.done
}
