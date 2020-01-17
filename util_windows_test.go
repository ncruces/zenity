package zenity

func init() {
	user32.NewProc("SetProcessDPIAware").Call()
}
