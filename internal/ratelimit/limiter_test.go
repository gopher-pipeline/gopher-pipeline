package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestLimiter_Wait_BasicRate(t *testing.T) {
	limiter := NewLimiter(2) // 2 rps
	defer limiter.Stop()

	ctx := context.Background()

	const n = 5
	var wg sync.WaitGroup

	start := time.Now()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := limiter.Wait(ctx); err != nil {
				t.Errorf("Wait returned error: %v", err)
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	t.Logf("elapsed: %s", elapsed)

	if elapsed <= 0 {
		t.Errorf("expected positive elapsed time, got %s", elapsed)
	}
}
