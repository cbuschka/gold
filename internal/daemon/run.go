package daemon

import (
	"github.com/cbuschka/golf/internal/command_server"
	configPkg "github.com/cbuschka/golf/internal/config"
	"github.com/cbuschka/golf/internal/dump"
	"github.com/cbuschka/golf/internal/gelf_server"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	worker "github.com/cbuschka/golf/internal/worker"
)

func startGelfUdpListener(addr string, journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	workerPool.RunWork(func() error {
		return gelf_server.ServeUdp(addr, journal)
	})
}

func startGelfTcpListener(addr string, journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	workerPool.RunWork(func() error {
		return gelf_server.ServeTcp(addr, journal, workerPool)
	})
}

func startUdsCommandServer(socketPath string, journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	workerPool.RunWork(func() error {
		return command_server.ServeUds(socketPath, journal)
	})
}

func startGelfHttpListener(bindAddr string, journal *journalPkg.Journal, workerPool *worker.WorkerPool) {
	workerPool.RunWork(func() error {
		return gelf_server.ServeHttp(bindAddr, journal)
	})
}

func Run() error {
	config := configPkg.GetConfig()

	journal, err := journalPkg.NewJournal(config)
	if err != nil {
		return err
	}
	defer journal.Close()

	workerPool := worker.NewWorkerPool()
	for _, bindAddr := range config.GelfUdpListeners {
		startGelfUdpListener(bindAddr, journal, workerPool)
	}
	for _, bindAddr := range config.GelfTcpListeners {
		startGelfTcpListener(bindAddr, journal, workerPool)
	}
	for _, bindAddr := range config.GelfHttpListeners {
		startGelfHttpListener(bindAddr, journal, workerPool)
	}
	startUdsCommandServer(config.CommandDomainSocketPath, journal, workerPool)
	dump.StartPeriodicDump(journal, workerPool)
	workerPool.Wait()

	return nil
}
