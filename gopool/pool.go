package gopool

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/ekuu/ho/internal/linkedlist"
)

type Pool interface {
	Go(func(ctx context.Context))
	GoErr(func(ctx context.Context) error)
	Wait()
	WaitSignal(signals ...os.Signal)
	SetWorkerCap(n int32)
	WorkerCount() int32
	TaskCount() int32
}

// pool
// A worker pool implementation
//
//go:generate gogen option -n pool -s ctx,workerCap,recoverFunc -p _ --with-init
type pool struct {
	// sync.RWMutex for thread-safe operations
	mu sync.RWMutex
	// linked list of CtxErrFunc objects
	list linkedlist.Single[func(ctx context.Context) error]
	//capacity of workers. -1 (default): unlimited
	workerCap *int32
	// current count of workers
	workerCount int32
	// panic handle function
	recoverFunc func(ctx context.Context, r any)
	// waiting for all goroutine to finish
	wg sync.WaitGroup
	// context for the pool
	ctx context.Context
	// cancel function for the context
	cancel context.CancelFunc
}

func (p *pool) init() {
	p.list = linkedlist.NewSingle[func(ctx context.Context) error]()
	if p.ctx == nil {
		p.ctx = context.Background()
	}
	p.ctx, p.cancel = context.WithCancel(p.ctx)
	if p.recoverFunc == nil {
		p.recoverFunc = func(ctx context.Context, r any) {
			slog.Error("gopool panic", "recover", r)
		}
	}
	if p.workerCap == nil {
		var workerCap int32 = -1
		p.workerCap = &workerCap
	}
}

func (p *pool) Go(f func(ctx context.Context)) {
	p.GoErr(func(ctx context.Context) error {
		f(ctx)
		return nil
	})
}

func (p *pool) GoErr(f func(ctx context.Context) error) {
	p.wg.Add(1)
	p.mu.Lock()
	p.list.Append(func(ctx context.Context) error {
		defer p.wg.Done()
		return f(ctx)
	})
	p.mu.Unlock()
	workerCap := atomic.LoadInt32(p.workerCap)
	if p.WorkerCount() == 0 || workerCap == -1 || p.WorkerCount() < workerCap {
		p.runWorker()
	}
}

func (p *pool) Wait() {
	p.wg.Wait()
}

func (p *pool) WaitSignal(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	}
	signal.Notify(ch, signals...)
	<-ch
	p.cancel()
	p.wg.Wait()
}

func (p *pool) runWorker() {
	p.incrWorkerCount()
	go func() {
		defer p.decrWorkerCount()
		for {
			p.mu.Lock()
			fn, ok := p.list.Shift()
			p.mu.Unlock()
			if !ok {
				return
			}
			p.runTask(fn)
		}
	}()
}

func (p *pool) runTask(f func(ctx context.Context) error) {
	defer func() {
		if r := recover(); r != nil {
			p.recoverFunc(p.ctx, r)
		}
	}()
	if err := f(p.ctx); err != nil {
		slog.ErrorContext(p.ctx, "gopool run task", "error", err)
	}
}

// SetWorkerCap set worker max
func (p *pool) SetWorkerCap(n int32) {
	atomic.StoreInt32(p.workerCap, n)
	if taskCount := p.TaskCount(); taskCount < n {
		n = taskCount
	}
	newWorkerCount := n - p.WorkerCount()
	if newWorkerCount < 1 {
		return
	}
	for n = 0; n < newWorkerCount; n++ {
		p.runWorker()
	}
}

// WorkerCount current count of workers
func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

// taskCount current count of tasks
func (p *pool) taskCount() int32 {
	return int32(p.list.Len())
}

// TaskCount current count of tasks
func (p *pool) TaskCount() int32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.taskCount()
}

func (p *pool) incrWorkerCount() int32 {
	return atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) decrWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}
