package dump

import (
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	worker "github.com/cbuschka/golf/internal/worker"
	"math"
	"time"
)

func StartPeriodicDump(journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	_ = schedule(func() {
		var firstTimestamp time.Time
		err := journal.ListMessages(-1, -1, func(message *journalPkg.Message) (bool, error) {
			firstTimestamp = message.ReceivedTimeUnix
			return false, nil
		})
		if err != nil {
			fmt.Printf("Error querying first timestamp: %v\n", err)
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

func timeFromFloat64(timeFloat float64) time.Time {
	sec, dec := math.Modf(timeFloat)
	return time.Unix(int64(sec), int64(dec*(1e9)))
}
