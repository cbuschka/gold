package journal

type Journal interface {
	ListMessages(begin string, limit int, callback func(message *Message) (bool, error)) error
	WriteMessage(message *Message) error
	Close() error
}
