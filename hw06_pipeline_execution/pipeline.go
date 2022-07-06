package hw06pipelineexecution

import (
	"context"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	results := make(map[int]interface{})
	ctx, cancel := context.WithCancel(context.Background())

	var i int
	var mx sync.Mutex
	var wg sync.WaitGroup
	for datum := range in {
		wg.Add(1)
		ch := make(Bi, 1)
		ch <- datum
		close(ch)

		go func(item Out, i int) {
			defer wg.Done()
			for _, stage := range stages {
				select {
				case <-ctx.Done():
					return
				default:
					item = stage(item)
				}
			}

			mx.Lock()
			defer mx.Unlock()
			results[i] = <-item
		}(ch, i)

		i++
	}

	wait := make(chan struct{})
	go func() {
		wg.Wait()
		close(wait)
	}()

	out := make(Bi)
	go func() {
		defer close(out)

		select {
		case <-wait:
		case <-done:
			cancel()
			return
		}

		for k := 0; ; k++ {
			v, ok := results[k]
			if !ok {
				break
			}

			out <- v
		}
	}()

	return out
}
