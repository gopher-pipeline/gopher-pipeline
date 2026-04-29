package pipeline

import (
	"context"
	"sync"

	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
	"github.com/gopher-pipeline/gopher-pipeline/internal/proccesor"
)

type Pipeline struct {
	jobsCh     chan model.Job
	resultsCh  chan model.Result
	errCh      chan error
	numWorkers int
	wg         sync.WaitGroup
}

func NewPipeline(
	jobsCh chan model.Job,
	resultCh chan model.Result,
	errCh chan error, numWorkers int,
) *Pipeline {

	return &Pipeline{
		jobsCh:     jobsCh,
		resultsCh:  resultCh,
		errCh:      errCh,
		numWorkers: numWorkers,
	}
}

func (p *Pipeline) Run(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)

		go func(workerID int) {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case jobs, ok := <-p.jobsCh:
					if !ok {
						return
					}

					res, err := proccesor.Transform(jobs)
					if err != nil {
						p.errCh <- err
					} else {
						p.resultsCh <- res
					}
				}
			}
		}(i)
	}
}

func (p *Pipeline) Stop() <-chan model.Result {
	close(p.jobsCh)

	go func() {
		p.wg.Wait()
		close(p.errCh)
		close(p.resultsCh)
	}()

	return p.resultsCh
}
