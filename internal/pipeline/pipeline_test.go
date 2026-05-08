package pipeline

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
	"github.com/gopher-pipeline/gopher-pipeline/internal/ratelimit"
	"github.com/gopher-pipeline/gopher-pipeline/internal/retry"
)

func TestPipeline_Run(t *testing.T) {
	jobsCh := make(chan model.Job, 10)
	resultCh := make(chan model.Result, 10)
	errCh := make(chan error, 10)
	limiter := ratelimit.Limiter{}
	config := retry.Config{
		MaxAttempts:  5,
		InitialDelay: time.Duration(5 * time.Second),
		MaxDelay:     time.Duration(10 * time.Second),
		Multiplier:   2.0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pipeline := NewPipeline(jobsCh, resultCh, errCh, 3, &limiter, config)

	go pipeline.Run(ctx)
	go func() {
		for err := range errCh {
			t.Log("Job error:", err)
		}
	}()

	testJobs := []model.Job{
		{
			ID:       uuid.New(),
			Filename: "testfilename",
			Value:    rand.Int(),
			Status:   "pending",
		},
		{
			ID:       uuid.New(),
			Filename: "testfilename",
			Value:    rand.Int(),
			Status:   "pending",
		},
		{
			ID:       uuid.New(),
			Filename: "testfilename",
			Value:    rand.Int(),
			Status:   "pending",
		},
		{
			ID:       uuid.New(),
			Filename: "testfilename",
			Value:    rand.Int(),
			Status:   "pending",
		},
		{
			ID:       uuid.New(),
			Filename: "testfilename",
			Value:    rand.Int(),
			Status:   "pending",
		},
	}
	for _, job := range testJobs {
		jobsCh <- job
	}

	close(jobsCh)
	results := make([]model.Result, 0)

	for res := range resultCh {
		results = append(results, res)
	}

	if len(results) != len(testJobs) {
		t.Errorf("expected %d results, got %d results", len(testJobs), len(results))
	}
}
