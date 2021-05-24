package daemon

import (
        journalPkg "github.com/cbuschka/golf/internal/journal"
        worker "github.com/cbuschka/golf/internal/worker"
	"github.com/cbuschka/golf/internal/gelf_server"
	"github.com/cbuschka/golf/internal/command_server"
)

func startUdpServer(addr string, journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	workerPool.RunWork(func() error {
			return gelf_server.ServeUdp(addr, journal)
		})
}

func startTcpServer(addr string, journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	workerPool.RunWork(func() error {
			return gelf_server.ServeTcp(addr, journal, workerPool)
		})
}

func startUdsCommandServer(journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	workerPool.RunWork(func() error {
		return command_server.ServeUds("./tmp/golfd.sock", journal)
	})
}

func Run() error {
	journal, err := journalPkg.NewJournal()
	if err != nil {
		return err
	}
	defer journal.Close()

	workerPool := worker.NewWorkerPool()
	startUdpServer("127.0.0.1:12201", journal, workerPool)
	startTcpServer("127.0.0.1:12201", journal, workerPool)
	startUdsCommandServer(journal, workerPool)
	workerPool.Wait()

	return nil
}
