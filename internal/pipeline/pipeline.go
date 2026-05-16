package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/gopher-pipeline/gopher-pipeline/internal/metrics"
	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
	"github.com/gopher-pipeline/gopher-pipeline/internal/proccesor"
)

type Pipeline struct {
	jobsCh     <-chan model.Job
	resultsCh  chan<- model.Result
	errCh      chan<- error
	numWorkers int
	wg         sync.WaitGroup
	collector  *metrics.Collector
	done       chan struct{}
}

func NewPipeline(
	jobsCh chan model.Job,
	resultCh chan model.Result,
	errCh chan error,
	numWorkers int,
	collector *metrics.Collector,
) *Pipeline {

	return &Pipeline{
		jobsCh:     jobsCh,
		resultsCh:  resultCh,
		errCh:      errCh,
		numWorkers: numWorkers,
		collector:  collector,
	}
}

func (p *Pipeline) Done() <-chan struct{} {
	return p.done
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
					start := time.Now()
					if p.collector != nil {
						p.collector.Send(metrics.Event{
							Type:      metrics.EventJobStarted,
							JobID:     jobs.ID,
							WorkerID:  workerID,
							Timestamp: start,
						})
					}

					res, err := proccesor.Transform(jobs)
					if err != nil {
						if p.collector != nil {
							p.collector.Send(metrics.Event{
								Type:      metrics.EventJobFailed,
								JobID:     jobs.ID,
								WorkerID:  workerID,
								Timestamp: time.Now(),
								Duration:  time.Since(start),
								Err:       err,
							})
						}
						p.errCh <- err
						continue
					} else {
						if p.collector != nil {
							p.collector.Send(metrics.Event{
								Type:      metrics.EventJobSuccess,
								JobID:     jobs.ID,
								WorkerID:  workerID,
								Timestamp: time.Now(),
								Duration:  time.Since(start),
							})
						}
						p.resultsCh <- res
					}

				}
			}
		}(i)
	}

	p.wg.Wait()
	close(p.resultsCh)
	close(p.errCh)
}
