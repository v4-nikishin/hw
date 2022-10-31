package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	doStage := func(done, in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					out <- val
				}
			}
		}()
		return out
	}
	out := in
	for _, stage := range stages {
		out = doStage(done, stage(out))
	}
	return out
}
