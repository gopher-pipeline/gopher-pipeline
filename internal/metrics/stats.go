package metrics

import (
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
