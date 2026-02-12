package gopool

import (
	"errors"
	"log"
	"sync"
)

var ErrPoolOverload = errors.New("Pool Overload")

// Pool is a lightweight goroutine pool that limits the number of concurrent workers.
type Pool struct {
	work            chan func()
	sem             chan struct{} // pool
	wg              sync.WaitGroup
	WithNonblocking bool // When true, an ErrPoolOverload error will be returned when the pool is full
}

// New creates a new Pool with a maximum number of concurrent workers specified by size.
func New(size int) *Pool {
	return &Pool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}

// Schedule ignore blocking errors
func (p *Pool) Schedule(task func()) {
	p.Submit(task)
}

// Submit schedules a task to be executed by the pool.
// If the number of workers is less than the limit, it will spawn a new worker.
// Otherwise, the task will be queued until a worker is available.
func (p *Pool) Submit(task func()) error {
	p.wg.Add(1)
	select {
	// 1.If we haven't reached the worker limit, start a new worker goroutine.
	case p.sem <- struct{}{}:
		go p.worker(task)

	// 2.If a worker is already running, send the task to the work queue.
	case p.work <- task:

	default:
		if p.WithNonblocking {
			p.wg.Done()
			return ErrPoolOverload
		}
		p.work <- task
	}
	return nil
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
		// if p.work closed
		if task == nil {
			return
		}

		// keep safe run
		p.handle(task)

		// waiting add task
		task = <-p.work
	}
}

func (p *Pool) handle(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover:", err)
		}
		p.wg.Done()
	}()
	fn()
}
