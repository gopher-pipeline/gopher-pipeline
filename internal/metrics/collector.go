package metrics

import (
	"context"
	"sort"
	"sync"
	"time"
)

type Collector struct {
	eventCh   chan Event
	done      chan struct{}
	once      sync.Once
	wg        sync.WaitGroup
	stats     Stats
	durations []time.Duration
	mu        sync.RWMutex
}

func NewCollector(bufSize int) *Collector {
	return &Collector{
		eventCh: make(chan Event, bufSize),
		done:    make(chan struct{}),
	}
}

func (c *Collector) Send(e Event) {
	select {
	case <-c.done:
		return
	case c.eventCh <- e:
		return
	default:
		return
	}
}

func (c *Collector) Wait() {
	c.wg.Wait()
}

func (c *Collector) Run(ctx context.Context) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.done:
				return
			case e := <-c.eventCh:
				c.mu.Lock()
				switch e.Type {
				case EventJobStarted:
					c.stats.Total++
				case EventJobSuccess:
					c.stats.Succeeded++
					c.durations = append(c.durations, e.Duration)
				case EventJobFailed:
					c.stats.Failed++
					c.durations = append(c.durations, e.Duration)
				case EventJobRetried:
					c.stats.Retried++
				}
				c.mu.Unlock()
			}
		}
	}()

}

func (c *Collector) Stats() Stats {
	c.mu.RLock()
	defer c.mu.Unlock()

	statsCopy := c.stats

	durationsCopy := make([]time.Duration, len(c.durations))
	copy(durationsCopy, c.durations)

	if len(durationsCopy) == 0 {
		return statsCopy
	}

	sort.Slice(durationsCopy, func(i, j int) bool {
		return durationsCopy[i] < durationsCopy[j]
	})

	sum := time.Duration(0)
	for _, d := range durationsCopy {
		sum += d
	}

	statsCopy.AvgDuration = sum / time.Duration(len(durationsCopy))

	idx := int(float64(len(durationsCopy)) * 0.95)
	if idx >= len(durationsCopy) {
		idx = len(durationsCopy) - 1
	}
	statsCopy.P95Duration = durationsCopy[idx]

	return statsCopy
}

func (c *Collector) Stop() {
	c.once.Do(func() {
		close(c.done)
		c.wg.Wait()
	})
}
