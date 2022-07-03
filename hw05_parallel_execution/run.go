package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 || m < 1 {
		return ErrErrorsLimitExceeded
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan struct{}, n)
	defer close(ch)

	var wg sync.WaitGroup
	var errorCount int64
	for _, task := range tasks {
		ch <- struct{}{}
		wg.Add(1)

		go func(task Task) {
			defer func() {
				wg.Done()
				<-ch
			}()

			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := task(); err != nil && atomic.AddInt64(&errorCount, 1) >= int64(m) {
				cancel()
			}
		}(task)
	}

	wg.Wait()

	if ctx.Err() != nil {
		return ErrErrorsLimitExceeded
	}

	return nil
}
