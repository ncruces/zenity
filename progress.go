package zenity

type ProgressMonitor interface {
	Message(string) error
	Progress(int) error
	Close() error
}
