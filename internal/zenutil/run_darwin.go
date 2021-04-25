package zenutil

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Run is internal.
func Run(ctx context.Context, script string, data interface{}) ([]byte, error) {
	var buf strings.Builder

	err := scripts.ExecuteTemplate(&buf, script, data)
	if err != nil {
		return nil, err
	}

	script = buf.String()
	if Command {
		// Try to use syscall.Exec, fallback to exec.Command.
		if path, err := exec.LookPath("osascript"); err != nil {
		} else if t, err := ioutil.TempFile("", ""); err != nil {
		} else if err := os.Remove(t.Name()); err != nil {
		} else if _, err := t.WriteString(script); err != nil {
		} else if _, err := t.Seek(0, 0); err != nil {
		} else if err := syscall.Dup2(int(t.Fd()), syscall.Stdin); err != nil {
		} else if err := os.Stderr.Close(); err != nil {
		} else {
			syscall.Exec(path, []string{"osascript", "-l", "JavaScript"}, nil)
		}
	}

	if ctx != nil {
		cmd := exec.CommandContext(ctx, "osascript", "-l", "JavaScript")
		cmd.Stdin = strings.NewReader(script)
		out, err := cmd.Output()
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return out, err
	}
	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = strings.NewReader(script)
	return cmd.Output()
}

func RunProgress(ctx context.Context, env []string) (m *progressDialog, err error) {
	t, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(t)
		}
	}()

	var cmd *exec.Cmd
	name := filepath.Join(t, "progress.app")

	cmd = exec.Command("osacompile", "-l", "JavaScript", "-o", name)
	cmd.Stdin = strings.NewReader(progress)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	plist := filepath.Join(name, "Contents/Info.plist")

	cmd = exec.Command("defaults", "write", plist, "LSUIElement", "true")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("defaults", "write", plist, "CFBundleName", "")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var executable string
	cmd = exec.Command("defaults", "read", plist, "CFBundleExecutable")
	if out, err := cmd.Output(); err != nil {
		return nil, err
	} else {
		out = bytes.TrimSuffix(out, []byte{'\n'})
		executable = filepath.Join(name, "Contents/MacOS", string(out))
	}

	cmd = exec.Command(executable)
	cmd.Env = env
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	m = &progressDialog{
		done:  make(chan struct{}),
		lines: make(chan string),
	}
	go func() {
		defer func() {
			pipe.Close()
			err := cmd.Wait()
			if cerr := ctx.Err(); cerr != nil {
				err = cerr
			}
			m.err = err
			close(m.done)
			os.RemoveAll(t)
		}()
		for {
			var line string
			select {
			case s, ok := <-m.lines:
				if !ok {
					return
				}
				line = s
			case <-ctx.Done():
				return
			case <-time.After(40 * time.Millisecond):
			}
			if _, err := pipe.Write([]byte(line + "\n")); err != nil {
				return
			}
		}
	}()
	return
}

type progressDialog struct {
	err   error
	done  chan struct{}
	lines chan string
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
	return m.send(strconv.Itoa(value))
}

// File is internal.
type File struct {
	Operation string
	Separator string
	Options   FileOptions
}

// FileOptions is internal.
type FileOptions struct {
	Prompt     *string  `json:"withPrompt,omitempty"`
	Type       []string `json:"ofType,omitempty"`
	Name       string   `json:"defaultName,omitempty"`
	Location   string   `json:"defaultLocation,omitempty"`
	Multiple   bool     `json:"multipleSelectionsAllowed,omitempty"`
	Invisibles bool     `json:"invisibles,omitempty"`
}

// Dialog is internal.
type Dialog struct {
	Operation string
	Text      string
	Extra     *string
	Options   DialogOptions
}

// DialogOptions is internal.
type DialogOptions struct {
	Message string   `json:"message,omitempty"`
	As      string   `json:"as,omitempty"`
	Answer  *string  `json:"defaultAnswer,omitempty"`
	Hidden  bool     `json:"hiddenAnswer,omitempty"`
	Title   *string  `json:"withTitle,omitempty"`
	Icon    string   `json:"withIcon,omitempty"`
	Buttons []string `json:"buttons,omitempty"`
	Cancel  int      `json:"cancelButton,omitempty"`
	Default int      `json:"defaultButton,omitempty"`
	Timeout int      `json:"givingUpAfter,omitempty"`
}

// List is internal.
type List struct {
	Items     []string
	Separator string
	Options   ListOptions
}

// ListOptions is internal.
type ListOptions struct {
	Title    *string  `json:"withTitle,omitempty"`
	Prompt   *string  `json:"withPrompt,omitempty"`
	OK       *string  `json:"okButtonName,omitempty"`
	Cancel   *string  `json:"cancelButtonName,omitempty"`
	Default  []string `json:"defaultItems,omitempty"`
	Multiple bool     `json:"multipleSelectionsAllowed,omitempty"`
	Empty    bool     `json:"emptySelectionAllowed,omitempty"`
}

// Notify is internal.
type Notify struct {
	Text    string
	Options NotifyOptions
}

// NotifyOptions is internal.
type NotifyOptions struct {
	Title    *string `json:"withTitle,omitempty"`
	Subtitle string  `json:"subtitle,omitempty"`
}

type Buttons struct {
	Buttons []string
	Default int
	Cancel  int
	Extra   int
}

func (d *Dialog) SetButtons(btns Buttons) {
	d.Options.Buttons = btns.Buttons
	d.Options.Default = btns.Default
	d.Options.Cancel = btns.Cancel
	if btns.Extra > 0 {
		name := btns.Buttons[btns.Extra-1]
		d.Extra = &name
	}
}
