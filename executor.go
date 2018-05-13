package crawler

import (
	"context"
	"sync"
)

func ExecutorWithContext(ctx context.Context) Executor {
	childCtx, cancel := context.WithCancel(ctx)
	return &executorImpl{
		ctx:        childCtx,
		cancel:     cancel,
		firstError: nil,
	}
}

type Executor interface {
	Go(f func(ctx context.Context) error, count int)
	Cancel()
	Wait() error
}

type executorImpl struct {
	ctx       context.Context
	cancel    context.CancelFunc
	mtxCancel sync.Mutex

	wg        sync.WaitGroup
	onceError sync.Once

	firstError error
}

func (e *executorImpl) Go(f func(ctx context.Context) error, count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		e.wg.Add(1)
		ctx, _ := context.WithCancel(e.ctx)
		go func(ctx context.Context) {
			defer e.wg.Done()
			if err := f(ctx); err != nil {
				e.onceError.Do(func() {
					e.firstError = err
				})
			}
		}(ctx)
	}
}

func (e *executorImpl) Cancel() {
	e.mtxCancel.Lock()
	defer e.mtxCancel.Unlock()
	e.cancel()
}

func (e *executorImpl) Wait() error {
	e.Cancel()
	e.wg.Wait()
	return e.firstError
}
