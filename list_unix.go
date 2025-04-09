//go:build !windows && !darwin

package zenity

import "github.com/ncruces/zenity/internal/zenutil"

func list(text string, items []string, opts options) (string, error) {
	args := []string{"--list", "--hide-header", "--text", text}
	args = appendGeneral(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendWindowIcon(args, opts)
	if opts.listKind == radioListKind {
		args = append(args, "--radiolist", "--column=", "--column=")
		for _, i := range items {
			args = append(args, "", i)
		}
	} else {
		args = append(args, "--column=")
		args = append(args, items...)
	}
	if opts.midSearch {
		args = append(args, "--mid-search")
	}

	out, err := zenutil.Run(opts.ctx, args)
	return strResult(opts, out, err)
}

func isSelected(defaults []string, value string) string {
	for _, d := range defaults {
		if d == value {
			return "TRUE"
		}
	}
	return "FALSE"
}

func listMultiple(text string, items []string, opts options) ([]string, error) {
	args := []string{"--list", "--hide-header", "--text", text, "--multiple", "--separator", zenutil.Separator}
	args = appendGeneral(args, opts)
	args = appendButtons(args, opts)
	args = appendWidthHeight(args, opts)
	args = appendWindowIcon(args, opts)

	// Having multiple items selected by default is only supported for checklists.
	// In case user provides non-empty list of default items, checklist will be enforced to avoid confusion.
	if opts.listKind == checkListKind || len(opts.defaultItems) > 0 {
		args = append(args, "--checklist", "--column=", "--column=")
		for _, i := range items {
			args = append(args, isSelected(opts.defaultItems, i), i)
		}
	} else {
		args = append(args, "--column=")
		args = append(args, items...)
	}

	out, err := zenutil.Run(opts.ctx, args)
	return lstResult(opts, out, err)
}
