package dump

import (
	"fmt"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
        journalPkg "github.com/cbuschka/golf/internal/journal"
	"math"
	"time"
        worker "github.com/cbuschka/golf/internal/worker"
	)

func StartPeriodicDump(journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	_ = schedule(func() {
		var firstTimestamp time.Time
		err := journal.ListMessages(-1, 1, func(id uint64, message *gelf.Message) (bool, error) {
			firstTimestamp = fromUnixTimeFloat(message.TimeUnix)
			return false, nil
		})
		if err != nil {
			fmt.Printf("Error querying first timestamp: %v\n", err)
		}
	}, 1*time.Minute);
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

func fromUnixTimeFloat(timeFloat float64) time.Time {
	sec, dec := math.Modf(timeFloat)
	return time.Unix(int64(sec), int64(dec*(1e9)))
}
