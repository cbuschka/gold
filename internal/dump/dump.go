package dump

import (
	journalPkg "github.com/cbuschka/golf/internal/journal"
	"github.com/cbuschka/golf/internal/worker"
	"github.com/kataras/golog"
	"time"
)

func StartPeriodicDump(journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	_ = schedule(func() {
		var firstTimestamp time.Time
		err := journal.ListMessages("", -1, func(message *journalPkg.Message) (bool, error) {
			firstTimestamp = message.ReceivedTimeUnix
			return false, nil
		})
		if err != nil {
			golog.Errorf("Error querying first timestamp: %v", err)
		}
	}, 1*time.Minute)
}

func schedule(work func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			work()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
