package daemon

import (
        journalPkg "github.com/cbuschka/golf/internal/journal"
	"sync"
)

func runWork(f func(), waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	go (func() {
		defer waitGroup.Done()
		f()
	})()
}

func Run() error {
	journal, err := journalPkg.NewJournal()
	if err != nil {
		return err
	}
	defer journal.Close()

	var waitGroup sync.WaitGroup
	runWork(func() {
			serveUdp("127.0.0.1:12201", journal)
		}, &waitGroup);
	waitGroup.Wait()

	return nil
}
