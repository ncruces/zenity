//go:build windows

package win

const (
	WTS_CURRENT_SERVER_HANDLE = 0
	WTS_CURRENT_SESSION       = 0xffffffff
)

//sys WTSSendMessage(server Handle, sessionID uint32, title *uint16, titleLength int, message *uint16, messageLength int, style uint32, timeout int, response *uint32, wait bool) (err error) = wtsapi32.WTSSendMessageW
