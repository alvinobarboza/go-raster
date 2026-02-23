package main

import "sync"

type WorkerPool struct {
	wg    sync.WaitGroup
	tasks chan func()
}

func NewWorkerPool(threads int) *WorkerPool {
	if threads <= 0 {
		threads = 1
	}

	pool := &WorkerPool{
		tasks: make(chan func(), threads),
	}

	pool.wg.Add(threads)
	for range threads {
		go func() {
			defer pool.wg.Done()
			for task := range pool.tasks {
				task()
			}
		}()
	}

	return pool
}

func (p *WorkerPool) Submit(task func()) {
	p.tasks <- task
}

func (p *WorkerPool) Close() {
	close(p.tasks)
	p.wg.Wait()
}
