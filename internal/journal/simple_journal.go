package journal

import (
	"encoding/json"
	"github.com/cbuschka/golf/internal/config"
	"github.com/google/uuid"
)

type SimpleJournal struct {
	segmentManager *SegmentManager
}

func NewSimpleJournal(config *config.Config) (*SimpleJournal, error) {
	segmentManager := SegmentManager{basedir: config.Journal.DataDirPath}
	if err := segmentManager.collectFiles(); err != nil {
		return nil, err
	}

	return &SimpleJournal{segmentManager: &segmentManager}, nil
}

func (journal *SimpleJournal) ListMessages(begin string, limit int, callback func(message *Message) (bool, error)) error {

	return nil
}

func (journal *SimpleJournal) WriteMessage(message *Message) error {

	id := uuid.New()

	message.Id = id.String()

	buf, err := json.Marshal(message)
	if err != nil {
		return err
	}

	segment, err := journal.segmentManager.getLatestSegment()
	if err != nil {
		return err
	}

	if err := segment.Append(buf); err != nil {
		return err
	}

	return nil
}

func (journal *SimpleJournal) Close() error {
	return journal.segmentManager.Close()
}
