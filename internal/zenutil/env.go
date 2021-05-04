package zenutil

// These are internal.
const (
	ErrCanceled    = stringErr("dialog canceled")
	ErrExtraButton = stringErr("extra button pressed")
	ErrUnsupported = stringErr("unsupported option")
)

type stringErr string

func (e stringErr) Error() string { return string(e) }
