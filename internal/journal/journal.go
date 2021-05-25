package journal

import (
	"encoding/json"
	"github.com/cbuschka/golf/internal/config"
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbOpts "github.com/syndtr/goleveldb/leveldb/opt"
)

type Journal struct {
	db *leveldb.DB
}

func NewJournal(config *config.Config) (*Journal, error) {
	db, err := leveldb.OpenFile(config.DataDirPath, &leveldbOpts.Options{})
	if err != nil {
		return nil, err
	}

	return &Journal{db: db}, nil
}

func (journal *Journal) ListMessages(begin string, limit int, callback func(message *Message) (bool, error)) error {

	iter := journal.db.NewIterator(nil, nil)
	defer iter.Release()
	found := iter.First()
	if !found {
		return nil
	}

	if begin != "" {
		beginId, err := uuid.Parse(begin)
		if err != nil {
			return err
		}

		beginIdBytes, err := beginId.MarshalBinary()
		if err != nil {
			return err
		}

		found := iter.Seek(beginIdBytes)
		if !found {
			return nil
		}
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

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	message.Id = id.String()

	buf, err := json.Marshal(message)
	if err != nil {
		return err
	}

	idBytes, err := id.MarshalBinary()
	if err != nil {
		return err
	}

	err = journal.db.Put(idBytes, buf, nil)
	if err != nil {
		return err
	}

	return nil
}

func (journal *Journal) Close() error {
	return journal.db.Close()
}
