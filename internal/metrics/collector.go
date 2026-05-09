package metrics

import (
	"context"
	"time"
)

type Stats struct {
	Total       int
	Succeeded   int
	Failed      int
	Retried     int
	AvgDuration time.Duration
	P95Duration time.Duration
}

type Collector struct {
	eventCh chan Event
	done    chan struct{}
}

func NewCollector(bufSize int) *Collector

func (c *Collector) Send(e Event) {
	select {
	// case <- c.done():
	// 	return
	case c.eventCh <- e:
		return
	default:
		return
	}
}

func (c *Collector) Wait() {} // TODO: implement the func

func (c *Collector) Run(ctx context.Context) {} // TODO: implement the func

func (c *Collector) Stats() {} // TODO: implement the func

func (c *Collector) Stop() {
	close(c.done) //wrong solution TODO: correct the func
}
