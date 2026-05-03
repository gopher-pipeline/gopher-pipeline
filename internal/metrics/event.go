package metrics

import "time"

type EventType string

const (
	EventJobStarted EventType = "job_started"
	EventJobSuccess EventType = "job_success"
	EventJobFailed  EventType = "job_failed"
	EventJobRetried EventType = "job_retried"
)

type Event struct {
	Type      EventType
	JobID     string
	WorkerID  int
	Timestamp time.Time
	Duration  time.Duration
	Err       error
}
