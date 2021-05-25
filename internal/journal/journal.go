package journal

import (
	"encoding/json"
	"github.com/cbuschka/golf/internal/config"
	"github.com/cockroachdb/pebble"
	"github.com/google/uuid"
)

type Journal struct {
	db *pebble.DB
}

func NewJournal(config *config.Config) (*Journal, error) {
	db, err := pebble.Open(config.DataDirPath, &pebble.Options{})
	if err != nil {
		return nil, err
	}

	return &Journal{db: db}, nil
}

func (journal *Journal) ListMessages(begin string, limit int, callback func(message *Message) (bool, error)) error {

	iter := journal.db.NewIter(nil)
	defer iter.Close()
	if begin != "" {
		beginId, err := uuid.Parse(begin)
		if err != nil {
			return err
		}

		beginIdBytes, err := beginId.MarshalBinary()
		if err != nil {
			return err
		}

		iter.SetBounds(beginIdBytes, nil)
	}

	found := iter.First()
	if !found {
		return nil
	}

	count := 0
	for ; iter.Valid(); iter.Next() {
		if limit != -1 && count >= limit {
			break
		}

		v := iter.Value()
		var message Message
		err := json.Unmarshal(v, &message)
		if err != nil {
			return nil
		}

		goon, err := callback(&message)
		if err != nil {
			return err
		}
		if !goon {
			return nil
		}
		count = count + 1
	}

	return nil
}

func (journal *Journal) WriteMessage(message *Message) error {

	id := uuid.New()

	message.Id = id.String()

	buf, err := json.Marshal(message)
	if err != nil {
		return err
	}

	idBytes, err := id.MarshalBinary()
	if err != nil {
		return err
	}

	err = journal.db.Set(idBytes, buf, pebble.NoSync)
	if err != nil {
		return err
	}

	return nil
}

func (journal *Journal) Close() error {
	return journal.db.Close()
}
