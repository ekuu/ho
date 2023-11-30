package gopool

import (
	"context"
	"os"
)

var defaultPool Pool = New()

func SetDefault(p Pool) {
	defaultPool = p
}

func Go(f func(ctx context.Context)) {
	defaultPool.Go(f)
}

func GoErr(f func(ctx context.Context) error) {
	defaultPool.GoErr(f)
}

func Wait() {
	defaultPool.Wait()
}

func WaitSignal(signals ...os.Signal) {
	defaultPool.WaitSignal(signals...)
}

func SetWorkerCap(n int32) {
	defaultPool.SetWorkerCap(n)
}

func WorkerCount() int32 {
	return defaultPool.WorkerCount()
}

func TaskCount() int32 {
	return defaultPool.TaskCount()
}
