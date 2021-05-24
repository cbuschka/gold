package query

import (
	"fmt"
	jsonPkg "encoding/json"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	)

func ListMessages() error {
	journal, err := journalPkg.NewJournal()
	if err!= nil {
		return err
	}
	defer journal.Close()
	return journal.ListMessages(func(id uint64, message *gelf.Message) (bool, error) {
		json, err := jsonPkg.Marshal(message)
		if err != nil {
			return false, err
		}
		fmt.Printf("id=%d, message=%s\n", id, json)
		return true, nil
	})
}
