package worker

import (
	"github.com/kataras/golog"
	"sync"
)

type WorkerPool struct {
	waitGroup *sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	var waitGroup sync.WaitGroup
	return &WorkerPool{waitGroup: &waitGroup}
}

func (workerPool *WorkerPool) Wait() {
	workerPool.waitGroup.Wait()
}

func (workerPool *WorkerPool) RunWork(f func() error) {
	workerPool.waitGroup.Add(1)
	go (func() {
		defer workerPool.waitGroup.Done()
		err := f()
		if err != nil {
			golog.Errorf("Executing work failed: %v", err)
		}
	})()
}
