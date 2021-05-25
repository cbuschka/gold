package journal

import (
	"encoding/binary"
	"encoding/json"
	"github.com/cbuschka/golf/internal/config"
	"github.com/dgraph-io/badger/v3"
)

type Journal struct {
	db *badger.DB
}

func NewJournal(config *config.Config) (*Journal, error) {
	opts := badger.DefaultOptions(config.DataDirPath)
	opts.InMemory = false
	opts.NumVersionsToKeep = 1
	opts.BaseTableSize = 1024 * 1024 * 1
	opts.BaseLevelSize = 1024 * 1024 * 1
	opts.ValueLogFileSize = 1024 * 1024 * 10
	opts.MemTableSize = 1024 * 1024 * 10
	opts.DetectConflicts = false
	opts.NumLevelZeroTables = 1
	opts.NumLevelZeroTablesStall = 2
	opts.IndexCacheSize = 1024 * 1024 * 1
	opts.BlockSize = 1024 * 4
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Journal{db: db}, nil
}

func (journal *Journal) ListMessages(begin int, limit int, callback func(message *Message) (bool, error)) error {

	return journal.db.View(func(tx *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := tx.NewIterator(opts)
		defer it.Close()
		if begin != -1 {
			it.Seek(itob(uint64(begin)))
		} else {
			it.Rewind()
		}
		count := 0
		for ; it.Valid(); it.Next() {
			if limit != -1 && count >= limit {
				break
			}
			item := it.Item()
			var message Message
			err := item.Value(func(v []byte) error {
				err := json.Unmarshal(v, &message)
				return err
			})

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
	})
}

func (journal *Journal) WriteMessage(message *Message) error {

	return journal.db.Update(func(tx *badger.Txn) error {

		seq, err := journal.db.GetSequence([]byte("idSeq"), 10)
		if err != nil {
			return err
		}
		defer seq.Release()

		id, err := seq.Next()
		if err != nil {
			return err
		}
		message.Id = id

		buf, err := json.Marshal(message)
		if err != nil {
			return err
		}

		err = tx.Set(itob(id), buf)
		if err != nil {
			return err
		}

		return nil
	})
}

func (journal *Journal) Close() error {
	return journal.db.Close()
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(b []byte) uint64 {
	v := binary.BigEndian.Uint64(b)
	return v
}
