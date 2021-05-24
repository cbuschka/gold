package journal

import (
	"encoding/binary"
	bolt "go.etcd.io/bbolt"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
	"encoding/json"
)

const bucketName = "MessagesV1"

type Journal struct {
	db *bolt.DB
}

func NewJournal() (*Journal, error) {
	db, err := bolt.Open("./tmp/my.db", 0640, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Journal{db: db}, nil
}

func (journal* Journal) ListMessages(callback func(uint64, *gelf.Message) (bool, error)) error {

	return journal.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			id := btoi(key)
			var message gelf.Message
			err := json.Unmarshal(value, &message)
			if err != nil {
				return err
			}
			goon, err := callback(id, &message)
			if err != nil {
				return err
			}
			if !goon {
				return nil
			}
		}
		return nil
	})
}

func (journal *Journal) WriteMessage(message *gelf.Message) error {

	return journal.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		id, _ := bucket.NextSequence()

		buf, err := json.Marshal(message)
		if err != nil {
			return err
	        }

		err = bucket.Put(itob(id), buf)
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
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) uint64 {
	v := binary.BigEndian.Uint64(b)
	return v
}
