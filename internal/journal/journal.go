package journal

import (
	"fmt"
	"github.com/cbuschka/golf/internal/config"
	"github.com/kataras/golog"
)

type Journal interface {
	ListMessages(begin string, limit int, callback func(message *Message) (bool, error)) error
	WriteMessage(message *Message) error
	Close() error
}

func NewJournal(config *config.Config) (Journal, error) {
	golog.Debugf("Journal is %s.", config.Journal.Type)
	if config.Journal.Type == "simple" {
		return NewSimpleJournal(config)
	} else if config.Journal.Type == "pebble" {
		return NewPebbleJournal(config)
	} else {
		return nil, fmt.Errorf("Unknown journal type '%s'", config.Journal.Type)
	}
}
