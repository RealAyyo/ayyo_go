package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	current := in

	for _, stage := range stages {
		current = stage(process(current, done))
	}

	return current
}

func process(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)

		for {
			select {
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			case <-done:
				return
			}
		}
	}()

	return out
}
