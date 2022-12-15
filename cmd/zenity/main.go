//go:build windows || darwin || dev

package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/ncruces/go-strftime"
	"github.com/ncruces/zenity"
	"github.com/ncruces/zenity/internal/zencmd"
	"github.com/ncruces/zenity/internal/zenutil"
)

const unspecified = "\x00"

var (
	// Application Options
	errorDlg          bool
	infoDlg           bool
	warningDlg        bool
	questionDlg       bool
	entryDlg          bool
	listDlg           bool
	calendarDlg       bool
	passwordDlg       bool
	fileSelectionDlg  bool
	colorSelectionDlg bool
	progressDlg       bool
	notification      bool

	// General options
	title         string
	width         uint
	height        uint
	okLabel       string
	cancelLabel   string
	extraButton   string
	text          string
	icon          string
	windowIcon    string
	attach        string
	modal         bool
	multiple      bool
	defaultCancel bool

	// Message options
	noWrap    bool
	noMarkup  bool
	ellipsize bool

	// Entry options
	entryText string
	hideText  bool

	// Password options
	username bool

	// List options
	columns       int
	checklist     bool
	radiolist     bool
	disallowEmpty bool

	// Calendar options
	year  uint
	month uint
	day   uint

	// File selection options
	save             bool
	directory        bool
	confirmOverwrite bool
	confirmCreate    bool
	showHidden       bool
	filename         string
	fileFilters      zenity.FileFilters

	// Color selection options
	defaultColor string
	showPalette  bool

	// Progress options
	percentage    float64
	pulsate       bool
	autoClose     bool
	autoKill      bool
	noCancel      bool
	timeRemaining bool

	// Notify options
	listen bool

	// Windows specific options
	unixeol bool
	cygpath bool
	wslpath bool

	// Command options
	version bool
)

func main() {
	args := parseFlags()
	opts := loadFlags()
	validateFlags()

	zenutil.Command = true
	zenutil.DateUTS35 = func() (string, error) { return strftime.UTS35(zenutil.DateFormat) }
	zenutil.DateParse = func(s string) (time.Time, error) { return strftime.Parse(zenutil.DateFormat, s) }
	if unixeol {
		zenutil.LineBreak = "\n"
	}
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if zenutil.Timeout > 0 {
		c, cancel := context.WithTimeout(ctx, time.Duration(zenutil.Timeout)*time.Second)
		defer cancel()
		ctx = c
	}
	opts = append(opts, zenity.Context(ctx))

	switch {
	case errorDlg:
		errResult(zenity.Error(text, opts...))
	case infoDlg:
		errResult(zenity.Info(text, opts...))
	case warningDlg:
		errResult(zenity.Warning(text, opts...))
	case questionDlg:
		errResult(zenity.Question(text, opts...))

	case entryDlg:
		strResult(zenity.Entry(text, opts...))

	case listDlg:
		if columns > 1 {
			var n int
			for i := 1; i < len(args); {
				args[n] = args[i]
				i += columns
				n += 1
			}
			args = args[:n:n]
		}
		if multiple {
			lstResult(zenity.ListMultiple(text, args, opts...))
		} else {
			strResult(zenity.List(text, args, opts...))
		}

	case calendarDlg:
		calResult(zenity.Calendar(text, opts...))

	case passwordDlg:
		pwdResult(zenity.Password(opts...))

	case fileSelectionDlg:
		switch {
		default:
			strResult(egestPath(zenity.SelectFile(opts...)))
		case save:
			strResult(egestPath(zenity.SelectFileSave(opts...)))
		case multiple:
			lstResult(egestPaths(zenity.SelectFileMultiple(opts...)))
		}

	case colorSelectionDlg:
		colResult(zenity.SelectColor(opts...))

	case progressDlg:
		errResult(progress(opts...))

	case notification:
		errResult(notify(opts...))

	default:
		panic("unreachable")
	}
}

func parseFlags() []string {
	fset := flag.NewFlagSet("zenity", flag.ContinueOnError)

	// Application Options
	fset.BoolVar(&errorDlg, "error", false, "Display error dialog")
	fset.BoolVar(&infoDlg, "info", false, "Display info dialog")
	fset.BoolVar(&warningDlg, "warning", false, "Display warning dialog")
	fset.BoolVar(&questionDlg, "question", false, "Display question dialog")
	fset.BoolVar(&entryDlg, "entry", false, "Display text entry dialog")
	fset.BoolVar(&listDlg, "list", false, "Display list dialog")
	fset.BoolVar(&calendarDlg, "calendar", false, "Display calendar dialog")
	fset.BoolVar(&passwordDlg, "password", false, "Display password dialog")
	fset.BoolVar(&fileSelectionDlg, "file-selection", false, "Display file selection dialog")
	fset.BoolVar(&colorSelectionDlg, "color-selection", false, "Display color selection dialog")
	fset.BoolVar(&progressDlg, "progress", false, "Display progress indication dialog")
	fset.BoolVar(&notification, "notification", false, "Display notification")

	// General options
	fset.StringVar(&title, "title", "", "Set the dialog `title`")
	fset.UintVar(&width, "width", 0, "Set the `width` (Unix only)")
	fset.UintVar(&height, "height", 0, "Set the `height` (Unix only)")
	fset.StringVar(&okLabel, "ok-label", "", "Set the `label` of the OK button")
	fset.StringVar(&cancelLabel, "cancel-label", "", "Set the `label` of the Cancel button")
	fset.Func("extra-button", "Add an extra `button`", setExtraButton)
	fset.StringVar(&text, "text", "", "Set the dialog `text`")
	fset.StringVar(&windowIcon, "window-icon", "", "Set the window `icon` (error, info, question, warning)")
	fset.StringVar(&attach, "attach", "", "Set the parent `window` to attach to")
	fset.BoolVar(&modal, "modal", runtime.GOOS == "darwin", "Set the modal hint")
	fset.BoolVar(&multiple, "multiple", false, "Allow multiple items to be selected")
	fset.BoolVar(&defaultCancel, "default-cancel", false, "Give Cancel button focus by default")

	// Message options
	fset.StringVar(&icon, "icon-name", "", "Set the dialog `icon` (dialog-error, dialog-information, dialog-question, dialog-warning)")
	fset.BoolVar(&noWrap, "no-wrap", false, "Do not enable text wrapping (Unix only)")
	fset.BoolVar(&noMarkup, "no-markup", false, "Do not enable Pango markup")
	fset.BoolVar(&ellipsize, "ellipsize", false, "Enable ellipsizing in the dialog text (Unix only)")

	// Entry options
	fset.StringVar(&entryText, "entry-text", "", "Set the entry `text`")
	fset.BoolVar(&hideText, "hide-text", false, "Hide the entry text")

	// Password options
	fset.BoolVar(&username, "username", false, "Display the username option")

	// List options
	fset.Func("column", "Set the column `header`", addColumn)
	fset.Bool("hide-header", true, "Hide the column headers")
	fset.BoolVar(&checklist, "checklist", false, "Use check boxes for the first column (Unix only)")
	fset.BoolVar(&radiolist, "radiolist", false, "Use radio buttons for the first column (Unix only)")
	fset.BoolVar(&disallowEmpty, "disallow-empty", false, "Disallow empty selection (Windows and macOS only)")

	// Calendar options
	fset.UintVar(&year, "year", 0, "Set the calendar `year`")
	fset.UintVar(&month, "month", 0, "Set the calendar `month`")
	fset.UintVar(&day, "day", 0, "Set the calendar `day`")
	fset.StringVar(&zenutil.DateFormat, "date-format", "%m/%d/%Y", "Set the `format` for the returned date")

	// File selection options
	fset.BoolVar(&save, "save", false, "Activate save mode")
	fset.BoolVar(&directory, "directory", false, "Activate directory-only selection")
	fset.BoolVar(&confirmOverwrite, "confirm-overwrite", false, "Confirm file selection if filename already exists")
	fset.BoolVar(&confirmCreate, "confirm-create", false, "Confirm file selection if filename does not yet exist (Windows only)")
	fset.BoolVar(&showHidden, "show-hidden", false, "Show hidden files (Windows and macOS only)")
	fset.StringVar(&filename, "filename", "", "Set the `filename`")
	fset.Func("file-filter", "Set a filename `filter` (NAME | PATTERN1 PATTERN2 ...)", addFileFilter)

	// Color selection options
	fset.StringVar(&defaultColor, "color", "", "Set the `color`")
	fset.BoolVar(&showPalette, "show-palette", false, "Show the palette")

	// Progress options
	fset.Float64Var(&percentage, "percentage", 0, "Set initial `percentage`")
	fset.BoolVar(&pulsate, "pulsate", false, "Pulsate progress bar")
	fset.BoolVar(&autoClose, "auto-close", false, "Dismiss the dialog when 100% has been reached")
	fset.BoolVar(&autoKill, "auto-kill", false, "Kill parent process if Cancel button is pressed (macOS and Unix only)")
	fset.BoolVar(&noCancel, "no-cancel", false, "Hide Cancel button (Windows and Unix only)")
	fset.BoolVar(&timeRemaining, "time-remaining", false, "Estimate when progress will reach 100% (Unix only)")

	// Notify options
	fset.BoolVar(&listen, "listen", false, "Listen for commands on stdin")

	// Windows specific options
	if runtime.GOOS == "windows" {
		fset.BoolVar(&unixeol, "unixeol", false, "Use Unix line endings (Windows only)")
		fset.BoolVar(&cygpath, "cygpath", false, "Use cygpath for path translation (Windows only)")
		fset.BoolVar(&wslpath, "wslpath", false, "Use wslpath for path translation (Windows only)")
	}

	// Command options
	fset.BoolVar(&version, "version", false, "Show version of program")
	fset.IntVar(&zenutil.Timeout, "timeout", 0, "Set dialog `timeout` in seconds")
	fset.StringVar(&zenutil.Separator, "separator", "|", "Set output `separator` character")

	// Detect unspecified values
	title = unspecified
	okLabel = unspecified
	cancelLabel = unspecified
	extraButton = unspecified
	text = unspecified
	icon = unspecified
	windowIcon = unspecified

	fset.Usage = func() {}
	err := fset.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println("usage: zenity [options...]")
		fset.PrintDefaults()
		os.Exit(0)
	}
	if err != nil {
		os.Exit(-1)
	}
	return fset.Args()
}

func validateFlags() {
	if version {
		fmt.Printf("zenity %s %s/%s\n", getVersion(), runtime.GOOS, runtime.GOARCH)
		fmt.Println("https://github.com/ncruces/zenity")
		os.Exit(0)
	}

	var n int
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
	if entryDlg {
		n++
	}
	if listDlg {
		n++
	}
	if calendarDlg {
		n++
	}
	if passwordDlg {
		n++
	}
	if fileSelectionDlg {
		n++
	}
	if colorSelectionDlg {
		n++
	}
	if progressDlg {
		n++
	}
	if notification {
		n++
	}
	if n == 0 {
		os.Stderr.WriteString("no dialog type specified; try 'zenity --help'\n")
		os.Exit(-1)
	}
	if n >= 2 {
		os.Stderr.WriteString("two or more dialogs types specified\n")
		os.Exit(-1)
	}

	if checklist {
		multiple = true
	}
	if radiolist {
		multiple = false
	}
	if checklist && radiolist {
		os.Stderr.WriteString("two or more list dialog types specified\n")
		os.Exit(-1)
	}
	if !checklist && !radiolist && columns > 1 || columns > 2 {
		os.Stderr.WriteString("multiple columns not supported\n")
		os.Exit(-1)
	}
}

func loadFlags() []zenity.Option {
	var opts []zenity.Option

	// Defaults

	setDefault := func(s *string, val string) {
		if *s == unspecified {
			*s = val
		}
	}
	switch {
	case errorDlg:
		setDefault(&title, "Error")
		setDefault(&text, "An error has occurred.")
		setDefault(&okLabel, "OK")
	case infoDlg:
		setDefault(&title, "Information")
		setDefault(&text, "All updates are complete.")
		setDefault(&okLabel, "OK")
	case warningDlg:
		setDefault(&title, "Warning")
		setDefault(&text, "Are you sure you want to proceed?")
		setDefault(&okLabel, "OK")
	case questionDlg:
		setDefault(&title, "Question")
		setDefault(&text, "Are you sure you want to proceed?")
		setDefault(&okLabel, "Yes")
		setDefault(&cancelLabel, "No")
	case entryDlg:
		setDefault(&title, "Add a new entry")
		setDefault(&text, "Enter new text:")
		setDefault(&okLabel, "OK")
		setDefault(&cancelLabel, "Cancel")
	case listDlg:
		setDefault(&title, "Select items from the list")
		setDefault(&text, "Select items from the list below:")
		setDefault(&okLabel, "OK")
		setDefault(&cancelLabel, "Cancel")
	case calendarDlg:
		setDefault(&title, "Calendar selection")
		setDefault(&text, "Select a date from below:")
		setDefault(&okLabel, "OK")
		setDefault(&cancelLabel, "Cancel")
	case passwordDlg:
		setDefault(&title, "Type your password")
		setDefault(&okLabel, "OK")
		setDefault(&cancelLabel, "Cancel")
	case progressDlg:
		setDefault(&title, "Progress")
		setDefault(&text, "Running…")
		setDefault(&okLabel, "OK")
		setDefault(&cancelLabel, "Cancel")
	}

	// General options

	if title != unspecified {
		opts = append(opts, zenity.Title(title))
	}
	opts = append(opts, zenity.Width(width))
	opts = append(opts, zenity.Height(height))
	if okLabel != unspecified {
		opts = append(opts, zenity.OKLabel(okLabel))
	}
	if cancelLabel != unspecified {
		opts = append(opts, zenity.CancelLabel(cancelLabel))
	}
	if extraButton != unspecified {
		opts = append(opts, zenity.ExtraButton(extraButton))
	}
	if defaultCancel {
		opts = append(opts, zenity.DefaultCancel())
	}

	if notification {
		icon = windowIcon
	}

	switch icon {
	case "error", "dialog-error":
		opts = append(opts, zenity.ErrorIcon)
	case "info", "dialog-information":
		opts = append(opts, zenity.InfoIcon)
	case "question", "dialog-question":
		opts = append(opts, zenity.QuestionIcon)
	case "important", "warning", "dialog-warning":
		opts = append(opts, zenity.WarningIcon)
	case "dialog-password":
		opts = append(opts, zenity.PasswordIcon)
	case "":
		opts = append(opts, zenity.NoIcon)
	case unspecified:
		//
	default:
		opts = append(opts, zenity.Icon(ingestPath(icon)))
	}

	switch windowIcon {
	case "error", "dialog-error":
		opts = append(opts, zenity.WindowIcon(zenity.ErrorIcon))
	case "info", "dialog-information":
		opts = append(opts, zenity.WindowIcon(zenity.InfoIcon))
	case "question", "dialog-question":
		opts = append(opts, zenity.WindowIcon(zenity.QuestionIcon))
	case "important", "warning", "dialog-warning":
		opts = append(opts, zenity.WindowIcon(zenity.WarningIcon))
	case "dialog-password":
		opts = append(opts, zenity.WindowIcon(zenity.PasswordIcon))
	case "", unspecified:
		//
	default:
		opts = append(opts, zenity.WindowIcon(ingestPath(windowIcon)))
	}

	if attach != "" {
		opts = append(opts, zenity.Attach(zencmd.ParseWindowId(attach)))
	} else if modal {
		if id := zencmd.GetParentWindowId(os.Getppid()); id != 0 {
			opts = append(opts, zenity.Attach(id))
		}
	}
	if modal {
		opts = append(opts, zenity.Modal())
	}

	// Message options

	if noWrap {
		opts = append(opts, zenity.NoWrap())
	}
	if ellipsize {
		opts = append(opts, zenity.Ellipsize())
	}
	if !noMarkup {
		switch {
		case errorDlg, infoDlg, warningDlg, questionDlg:
			text = zencmd.StripMarkup(text)
		}
	}

	// Entry options

	opts = append(opts, zenity.EntryText(entryText))
	if hideText {
		opts = append(opts, zenity.HideText())
	}

	// Password options

	if username {
		opts = append(opts, zenity.Username())
	}

	// List options

	if checklist {
		opts = append(opts, zenity.CheckList())
	}
	if radiolist {
		opts = append(opts, zenity.RadioList())
	}
	if disallowEmpty {
		opts = append(opts, zenity.DisallowEmpty())
	}

	// Calendar options

	y, m, d := time.Now().Date()
	if month != 0 {
		m = time.Month(month)
	}
	if day != 0 {
		d = int(day)
	}
	if year != 0 {
		y = int(year)
	}
	opts = append(opts, zenity.DefaultDate(y, m, d))

	// File selection options

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
	if filename != "" {
		opts = append(opts, zenity.Filename(ingestPath(filename)))
	}
	opts = append(opts, fileFilters)

	// Color selection options

	if defaultColor != "" {
		opts = append(opts, zenity.Color(zenutil.ParseColor(defaultColor)))
	}
	if showPalette {
		opts = append(opts, zenity.ShowPalette())
	}

	// Progress options

	if pulsate {
		opts = append(opts, zenity.Pulsate())
	}
	if noCancel {
		opts = append(opts, zenity.NoCancel())
	}
	if timeRemaining {
		opts = append(opts, zenity.TimeRemaining())
	}

	return opts
}

func errResult(err error) {
	if os.IsTimeout(err) {
		os.Exit(5)
	}
	if err == zenity.ErrCanceled || err == context.Canceled {
		os.Exit(1)
	}
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
}

func strResult(s string, err error) {
	errResult(err)
	os.Stdout.WriteString(s)
	os.Stdout.WriteString(zenutil.LineBreak)
}

func lstResult(l []string, err error) {
	errResult(err)
	os.Stdout.WriteString(strings.Join(l, zenutil.Separator))
	if len(l) > 0 {
		os.Stdout.WriteString(zenutil.LineBreak)
	}
}

func calResult(d time.Time, err error) {
	errResult(err)
	os.Stdout.WriteString(strftime.Format(zenutil.DateFormat, d))
	os.Stdout.WriteString(zenutil.LineBreak)
}

func colResult(c color.Color, err error) {
	errResult(err)
	os.Stdout.WriteString(zenutil.UnparseColor(c))
	os.Stdout.WriteString(zenutil.LineBreak)
}

func pwdResult(u, p string, err error) {
	errResult(err)
	if username {
		os.Stdout.WriteString(u)
		os.Stdout.WriteString(zenutil.Separator)
	}
	os.Stdout.WriteString(p)
	os.Stdout.WriteString(zenutil.LineBreak)
}

func ingestPath(path string) string {
	if runtime.GOOS == "windows" && path != "" {
		var args []string
		switch {
		case wslpath:
			args = []string{"wsl", "wslpath", "-w"}
		case cygpath:
			args = []string{"cygpath", "-C", "UTF8", "-w"}
		}
		if args != nil {
			args = append(args, path)
			out, err := exec.Command(args[0], args[1:]...).Output()
			if err == nil {
				path = string(bytes.TrimSuffix(out, []byte{'\n'}))
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
			if err == nil {
				path = string(bytes.TrimSuffix(out, []byte{'\n'}))
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

func setExtraButton(s string) error {
	if extraButton != unspecified {
		return errors.New("multiple extra buttons not supported")
	}
	extraButton = s
	return nil
}

func addColumn(s string) error {
	columns++
	return nil
}

func addFileFilter(s string) error {
	var filter zenity.FileFilter

	if head, tail, cut := strings.Cut(s, "|"); cut {
		filter.Name = head
		s = tail
	}

	filter.Patterns = strings.Split(strings.Trim(s, " "), " ")
	fileFilters = append(fileFilters, filter)
	return nil
}

func getVersion() string {
	if tag != "" {
		return tag
	}

	rev := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, kv := range info.Settings {
			if kv.Key == "vcs.modified" && kv.Value == "true" {
				return "custom"
			}
			if kv.Key == "vcs.revision" {
				rev = kv.Value
			}
		}
	}
	return rev
}
