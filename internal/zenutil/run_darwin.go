package zenutil

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Run is internal.
func Run(ctx context.Context, script string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := scripts.ExecuteTemplate(&buf, script, data)
	if err != nil {
		return nil, err
	}

	if Command {
		// Try to use syscall.Exec, fallback to exec.Command.
		if path, err := exec.LookPath("osascript"); err != nil {
		} else if t, err := ioutil.TempFile("", ""); err != nil {
		} else if err := os.Remove(t.Name()); err != nil {
		} else if _, err := t.Write(buf.Bytes()); err != nil {
		} else if _, err := t.Seek(0, 0); err != nil {
		} else if err := syscall.Dup2(int(t.Fd()), syscall.Stdin); err != nil {
		} else if err := os.Stderr.Close(); err != nil {
		} else {
			syscall.Exec(path, []string{"osascript", "-l", "JavaScript"}, nil)
		}
	}

	if ctx != nil {
		cmd := exec.CommandContext(ctx, "osascript", "-l", "JavaScript")
		cmd.Stdin = &buf
		out, err := cmd.Output()
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return out, err
	}
	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = &buf
	return cmd.Output()
}

// RunProgress is internal.
func RunProgress(ctx context.Context, max int, data Progress) (dlg *progressDialog, err error) {
	var buf bytes.Buffer
	err = scripts.ExecuteTemplate(&buf, "progress", data)
	if err != nil {
		return nil, err
	}

	t, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if ctx != nil && ctx.Err() != nil {
				err = ctx.Err()
			}
			os.RemoveAll(t)
		}
	}()
	if ctx == nil {
		ctx = context.Background()
	}

	var cmd *exec.Cmd
	name := filepath.Join(t, "progress.app")

	cmd = exec.CommandContext(ctx, "osacompile", "-l", "JavaScript", "-o", name)
	cmd.Stdin = &buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	plist := filepath.Join(name, "Contents/Info.plist")

	cmd = exec.CommandContext(ctx, "defaults", "write", plist, "LSUIElement", "true")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.CommandContext(ctx, "defaults", "write", plist, "CFBundleName", "")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var executable string
	cmd = exec.CommandContext(ctx, "defaults", "read", plist, "CFBundleExecutable")
	if out, err := cmd.Output(); err != nil {
		return nil, err
	} else {
		out = bytes.TrimSuffix(out, []byte{'\n'})
		executable = filepath.Join(name, "Contents/MacOS", string(out))
	}

	cmd = exec.CommandContext(ctx, executable)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	dlg = &progressDialog{
		ctx:   ctx,
		cmd:   cmd,
		max:   max,
		lines: make(chan string),
		done:  make(chan struct{}),
	}
	go dlg.pipe(pipe)
	go func() {
		defer os.RemoveAll(t)
		dlg.wait(nil, nil)
	}()
	return dlg, nil
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

// DialogButtons is internal.
type DialogButtons struct {
	Buttons []string
	Default int
	Cancel  int
	Extra   int
}

// SetButtons is internal.
func (d *Dialog) SetButtons(btns DialogButtons) {
	d.Options.Buttons = btns.Buttons
	d.Options.Default = btns.Default
	d.Options.Cancel = btns.Cancel
	if btns.Extra > 0 {
		name := btns.Buttons[btns.Extra-1]
		d.Extra = &name
	}
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

// File is internal.
type File struct {
	Operation string
	Separator string
	Options   FileOptions
}

type FileOptions struct {
	Prompt     *string  `json:"withPrompt,omitempty"`
	Type       []string `json:"ofType,omitempty"`
	Name       string   `json:"defaultName,omitempty"`
	Location   string   `json:"defaultLocation,omitempty"`
	Multiple   bool     `json:"multipleSelectionsAllowed,omitempty"`
	Invisibles bool     `json:"invisibles,omitempty"`
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

// Progress is internal.
type Progress struct {
	Description *string
	Total       *int
}
