package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	results := make(Bi)
	stage := stages[0]

	go func() {
		defer close(results)
		for item := range stage(in) {
			select {
			case <-done:
				return
			default:
				results <- item
			}
		}
	}()

	return ExecutePipeline(results, done, stages[1:]...)
}
