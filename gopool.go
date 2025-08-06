package gopool

import "sync"

// Pool is a lightweight goroutine pool that limits the number of concurrent workers.
type Pool struct {
	work chan func()
	sem  chan struct{}
	wg   sync.WaitGroup
}

// New creates a new Pool with a maximum number of concurrent workers specified by size.
func New(size int) *Pool {
	return &Pool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}

// Schedule schedules a task to be executed by the pool.
// If the number of workers is less than the limit, it will spawn a new worker.
// Otherwise, the task will be queued until a worker is available.
func (p *Pool) Schedule(task func()) {
	select {
	// If a worker is already running, send the task to the work queue.
	case p.work <- task:
		p.wg.Add(1)

	// If we haven't reached the worker limit, start a new worker goroutine.
	case p.sem <- struct{}{}:
		p.wg.Add(1)
		go p.worker(task)
	}
}

// Release stops accepting new tasks and closes the work channel.
// This will signal worker goroutines to exit once all tasks are processed.
func (p *Pool) Release() {
	close(p.work)
}

// Wait blocks until all scheduled tasks have completed.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// worker is a long-running goroutine that processes tasks from the work queue.
// It exits when the work channel is closed and all tasks are consumed.
func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()

	for {
		if task == nil {
			return
		}
		task()
		p.wg.Done()
		task = <-p.work
	}
}
