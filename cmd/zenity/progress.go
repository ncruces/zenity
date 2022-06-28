//go:build windows || darwin || dev

package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ncruces/zenity"
	"github.com/ncruces/zenity/internal/zencmd"
)

func progress(opts ...zenity.Option) (err error) {
	dlg, err := zenity.Progress(opts...)
	if err != nil {
		return err
	}

	if autoKill {
		defer func() {
			if err == zenity.ErrCanceled {
				zencmd.KillParent()
			}
		}()
	}

	if err := dlg.Text(text); err != nil {
		return err
	}
	if err := dlg.Value(int(math.Round(percentage))); err != nil {
		return err
	}

	lines := make(chan string)

	go func() {
		defer close(lines)
		for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
			lines <- scanner.Text()
		}
	}()

	for {
		select {
		case line, ok := <-lines:
			if !ok {
				break
			}
			if len(line) > 1 && line[0] == '#' {
				if err := dlg.Text(strings.TrimSpace(line[1:])); err != nil {
					return err
				}
			} else if v, err := strconv.ParseFloat(line, 64); err == nil {
				if err := dlg.Value(int(math.Round(v))); err != nil {
					return err
				}
				if v >= 100 && autoClose {
					return dlg.Close()
				}
			}
			continue
		case <-dlg.Done():
		}
		break
	}

	if err := dlg.Complete(); err != nil {
		return err
	}
	if !autoClose {
		<-dlg.Done()
	}
	return dlg.Close()
}
