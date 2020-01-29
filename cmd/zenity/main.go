package main

import (
	"context"
	"flag"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ncruces/zenity"
	"github.com/ncruces/zenity/internal/zenutil"
)

//go:generate go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo -platform-specific -manifest=win.manifest

var (
	// Application Options
	notification      bool
	errorDlg          bool
	infoDlg           bool
	warningDlg        bool
	questionDlg       bool
	fileSelectionDlg  bool
	colorSelectionDlg bool

	// General options
	title string

	// Message options
	text          string
	icon          string
	okLabel       string
	cancelLabel   string
	extraButton   string
	noWrap        bool
	ellipsize     bool
	defaultCancel bool

	// File selection options
	save             bool
	multiple         bool
	directory        bool
	confirmOverwrite bool
	confirmCreate    bool
	showHidden       bool
	filename         string
	fileFilters      FileFilters

	// Color selection options
	defaultColor string
	showPalette  bool

	// Windows specific options
	cygpath bool
	wslpath bool
)

func init() {
	prevUsage := flag.Usage
	flag.Usage = func() {
		prevUsage()
		os.Exit(-1)
	}
}

func main() {
	setupFlags()
	flag.Parse()
	validateFlags()
	opts := loadFlags()
	zenutil.Command = true
	if zenutil.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(zenutil.Timeout)*time.Second)
		opts = append(opts, zenity.Context(ctx))
		_ = cancel
	}

	switch {
	case notification:
		errResult(zenity.Notify(text, opts...))

	case errorDlg:
		msgResult(zenity.Error(text, opts...))
	case infoDlg:
		msgResult(zenity.Info(text, opts...))
	case warningDlg:
		msgResult(zenity.Warning(text, opts...))
	case questionDlg:
		msgResult(zenity.Question(text, opts...))

	case fileSelectionDlg:
		switch {
		default:
			strResult(egestPath(zenity.SelectFile(opts...)))
		case save:
			strResult(egestPath(zenity.SelectFileSave(opts...)))
		case multiple:
			listResult(egestPaths(zenity.SelectFileMutiple(opts...)))
		}

	case colorSelectionDlg:
		colorResult(zenity.SelectColor(opts...))
	}

	flag.Usage()
}

func setupFlags() {
	// Application Options
	flag.BoolVar(&notification, "notification", false, "Display notification")
	flag.BoolVar(&errorDlg, "error", false, "Display error dialog")
	flag.BoolVar(&infoDlg, "info", false, "Display info dialog")
	flag.BoolVar(&warningDlg, "warning", false, "Display warning dialog")
	flag.BoolVar(&questionDlg, "question", false, "Display question dialog")
	flag.BoolVar(&fileSelectionDlg, "file-selection", false, "Display file selection dialog")
	flag.BoolVar(&colorSelectionDlg, "color-selection", false, "Display color selection dialog")

	// General options
	flag.StringVar(&title, "title", "", "Set the dialog title")
	flag.StringVar(&icon, "window-icon", "", "Set the window icon (error, info, question, warning)")

	// Message options
	flag.StringVar(&text, "text", "", "Set the dialog text")
	flag.StringVar(&icon, "icon-name", "", "Set the dialog icon (error, info, question, warning)")
	flag.StringVar(&okLabel, "ok-label", "", "Set the label of the OK button")
	flag.StringVar(&cancelLabel, "cancel-label", "", "Set the label of the Cancel button")
	flag.StringVar(&extraButton, "extra-button", "", "Add an extra button")
	flag.BoolVar(&noWrap, "no-wrap", false, "Do not enable text wrapping")
	flag.BoolVar(&ellipsize, "ellipsize", false, "Enable ellipsizing in the dialog text")
	flag.BoolVar(&defaultCancel, "default-cancel", false, "Give Cancel button focus by default")

	// File selection options
	flag.BoolVar(&save, "save", false, "Activate save mode")
	flag.BoolVar(&multiple, "multiple", false, "Allow multiple files to be selected")
	flag.BoolVar(&directory, "directory", false, "Activate directory-only selection")
	flag.BoolVar(&confirmOverwrite, "confirm-overwrite", false, "Confirm file selection if filename already exists")
	flag.BoolVar(&confirmCreate, "confirm-create", false, "Confirm file selection if filename does not yet exist (Windows only)")
	flag.BoolVar(&showHidden, "show-hidden", false, "Show hidden files (Windows and macOS only)")
	flag.StringVar(&filename, "filename", "", "Set the filename")
	flag.Var(&fileFilters, "file-filter", "Set a filename filter (NAME | PATTERN1 PATTERN2 ...)")

	// Color selection options
	flag.StringVar(&defaultColor, "color", "", "Set the color")
	flag.BoolVar(&showPalette, "show-palette", false, "Show the palette")

	// Windows specific options
	if runtime.GOOS == "windows" {
		flag.BoolVar(&cygpath, "cygpath", false, "Use cygpath for path translation (Windows only)")
		flag.BoolVar(&wslpath, "wslpath", false, "Use wslpath for path translation (Windows only)")
	}

	// Internal options
	flag.IntVar(&zenutil.Timeout, "timeout", 0, "Set dialog timeout in seconds")
	flag.StringVar(&zenutil.Separator, "separator", "|", "Set output separator character")
}

func validateFlags() {
	var n int
	if notification {
		n++
	}
	if errorDlg {
		n++
	}
	if infoDlg {
		n++
	}
	if warningDlg {
		n++
	}
	if questionDlg {
		n++
	}
	if fileSelectionDlg {
		n++
	}
	if colorSelectionDlg {
		n++
	}
	if n != 1 {
		flag.Usage()
	}
}

func loadFlags() []zenity.Option {
	var opts []zenity.Option

	// General options

	opts = append(opts, zenity.Title(title))

	// Message options

	var ico zenity.DialogIcon
	switch icon {
	case "error", "dialog-error":
		ico = zenity.ErrorIcon
	case "info", "dialog-information":
		ico = zenity.InfoIcon
	case "question", "dialog-question":
		ico = zenity.QuestionIcon
	case "important", "warning", "dialog-warning":
		ico = zenity.WarningIcon
	}

	opts = append(opts, zenity.Icon(ico))
	opts = append(opts, zenity.OKLabel(okLabel))
	opts = append(opts, zenity.CancelLabel(cancelLabel))
	opts = append(opts, zenity.ExtraButton(extraButton))
	if noWrap {
		opts = append(opts, zenity.NoWrap())
	}
	if ellipsize {
		opts = append(opts, zenity.Ellipsize())
	}
	if defaultCancel {
		opts = append(opts, zenity.DefaultCancel())
	}

	// File selection options

	opts = append(opts, fileFilters)
	if filename != "" {
		opts = append(opts, zenity.Filename(ingestPath(filename)))
	}
	if directory {
		opts = append(opts, zenity.Directory())
	}
	if confirmOverwrite {
		opts = append(opts, zenity.ConfirmOverwrite())
	}
	if confirmCreate {
		opts = append(opts, zenity.ConfirmCreate())
	}
	if showHidden {
		opts = append(opts, zenity.ShowHidden())
	}

	// Color selection options

	if defaultColor != "" {
		opts = append(opts, zenity.Color(zenutil.ParseColor(defaultColor)))
	}
	if showPalette {
		opts = append(opts, zenity.ShowPalette())
	}

	return opts
}

func errResult(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString(zenutil.LineBreak)
		os.Exit(-1)
	}
	os.Exit(0)
}

func msgResult(ok bool, err error) {
	if err == zenity.ErrExtraButton {
		os.Stdout.WriteString(extraButton)
		os.Stdout.WriteString(zenutil.LineBreak)
		os.Exit(1)
	}
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString(zenutil.LineBreak)
		os.Exit(-1)
	}
	if ok {
		os.Exit(0)
	}
	os.Exit(1)
}

func strResult(s string, err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString(zenutil.LineBreak)
		os.Exit(-1)
	}
	if s == "" {
		os.Exit(1)
	}
	os.Stdout.WriteString(s)
	os.Stdout.WriteString(zenutil.LineBreak)
	os.Exit(0)
}

func listResult(l []string, err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString(zenutil.LineBreak)
		os.Exit(-1)
	}
	os.Stdout.WriteString(strings.Join(l, zenutil.Separator))
	os.Stdout.WriteString(zenutil.LineBreak)
	if l == nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func colorResult(c color.Color, err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString(zenutil.LineBreak)
		os.Exit(-1)
	}
	if c == nil {
		os.Exit(1)
	}
	os.Stdout.WriteString(zenutil.UnparseColor(c))
	os.Stdout.WriteString(zenutil.LineBreak)
	os.Exit(0)
}

func ingestPath(path string) string {
	if runtime.GOOS == "windows" && path != "" {
		var args []string
		switch {
		case wslpath:
			args = []string{"wsl", "wslpath", "-m"}
		case cygpath:
			args = []string{"cygpath", "-C", "UTF8", "-m"}
		}
		if args != nil {
			args = append(args, path)
			out, err := exec.Command(args[0], args[1:]...).Output()
			if len(out) > 0 && err == nil {
				path = string(out[:len(out)-1])
			}
		}
	}
	return path
}

func egestPath(path string, err error) (string, error) {
	if runtime.GOOS == "windows" && path != "" && err == nil {
		var args []string
		switch {
		case wslpath:
			args = []string{"wsl", "wslpath", "-u"}
		case cygpath:
			args = []string{"cygpath", "-C", "UTF8", "-u"}
		}
		if args != nil {
			var out []byte
			args = append(args, filepath.ToSlash(path))
			out, err = exec.Command(args[0], args[1:]...).Output()
			if len(out) > 0 && err == nil {
				path = string(out[:len(out)-1])
			}
		}
	}
	return path, err
}

func egestPaths(paths []string, err error) ([]string, error) {
	if runtime.GOOS == "windows" && err == nil && (wslpath || cygpath) {
		paths = append(paths[:0:0], paths...)
		for i, p := range paths {
			paths[i], err = egestPath(p, nil)
			if err != nil {
				break
			}
		}
	}
	return paths, err
}

// FileFilters is internal.
type FileFilters struct {
	zenity.FileFilters
}

// String is internal.
func (f *FileFilters) String() string {
	return "zenity.FileFilters"
}

// Set is internal.
func (f *FileFilters) Set(s string) error {
	var filter zenity.FileFilter

	if split := strings.SplitN(s, "|", 2); len(split) > 1 {
		filter.Name = split[0]
		s = split[1]
	}

	filter.Patterns = strings.Split(s, " ")
	f.FileFilters = append(f.FileFilters, filter)

	return nil
}
