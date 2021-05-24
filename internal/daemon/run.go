package daemon

import (
        journalPkg "github.com/cbuschka/golf/internal/journal"
	"sync"
	"github.com/cbuschka/golf/internal/gelf_server"
	"github.com/cbuschka/golf/internal/command_server"
)

func runWork(f func(), waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	go (func() {
		defer waitGroup.Done()
		f()
	})()
}

func startUdpServer(addr string, journal *journalPkg.Journal, waitGroup *sync.WaitGroup) {
	runWork(func() {
			gelf_server.ServeUdp(addr, journal)
		}, waitGroup)
}

func startUdsCommandServer(journal *journalPkg.Journal, waitGroup *sync.WaitGroup) {
	runWork(func() {
		command_server.ServeUds("./tmp/golfd.sock", journal)
	}, waitGroup)
}

func Run() error {
	journal, err := journalPkg.NewJournal()
	if err != nil {
		return err
	}
	defer journal.Close()

	var waitGroup sync.WaitGroup
	startUdpServer("127.0.0.1:12201", journal, &waitGroup)
	startUdsCommandServer(journal, &waitGroup)
	waitGroup.Wait()

	return nil
}
