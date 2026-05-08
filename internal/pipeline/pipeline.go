package pipeline

import (
	"context"
	"sync"

	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
	"github.com/gopher-pipeline/gopher-pipeline/internal/proccesor"
	"github.com/gopher-pipeline/gopher-pipeline/internal/ratelimit"
	"github.com/gopher-pipeline/gopher-pipeline/internal/retry"
)

type Pipeline struct {
	jobsCh      <-chan model.Job
	resultsCh   chan<- model.Result
	errCh       chan<- error
	numWorkers  int
	wg          sync.WaitGroup
	limiter     *ratelimit.Limiter
	retryConfig retry.Config
}

func NewPipeline(
	jobsCh chan model.Job,
	resultCh chan model.Result,
	errCh chan error,
	numWorkers int,
	limiter *ratelimit.Limiter,
	config retry.Config,
) *Pipeline {

	return &Pipeline{
		jobsCh:      jobsCh,
		resultsCh:   resultCh,
		errCh:       errCh,
		numWorkers:  numWorkers,
		limiter:     limiter,
		retryConfig: config,
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

					err := p.limiter.Wait(ctx)
					if err != nil {
						select {
						case p.errCh <- err:
						case <-ctx.Done():
						}
						return
					}

					var result model.Result
					err = retry.Do(ctx, p.retryConfig, func() error {
						tempResult, err := proccesor.Transform(jobs)
						if err != nil {
							return err
						}
						result = tempResult
						return nil
					})

					if err != nil {
						select {
						case p.errCh <- err:
						case <-ctx.Done():
						}
						continue
					}

					select {
					case p.resultsCh <- result:
					case <-ctx.Done():
						return
					}
				}
			}
		}(i)
	}

	p.wg.Wait()
	close(p.resultsCh)
	close(p.errCh)
}
