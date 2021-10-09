package xerrors

const (
	InvalidChatID  Error = "Chat ID is invalid"
	ChatNotFound   Error = "Chat not found"
	InvalidMessageID  Error = "Message ID is invalid"
	MessageNotFound   Error = "Message not found"
	MissingParam Error = "Missing param"
)

const (
	CodeInvalidData = iota
)
