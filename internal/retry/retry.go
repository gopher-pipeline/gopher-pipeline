package retry

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/gopher-pipeline/gopher-pipeline/internal/model"
)

type Config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	OnRetry      func(int, error)
}

func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts < 1 {
		return model.ErrRetryConfigAttempts
	}

	if cfg.InitialDelay <= 0 {
		return model.ErrRetryConfigInitialDelay
	}

	if cfg.MaxDelay <= 0 {
		return model.ErrRetryConfigMaxDelay
	}

	if cfg.Multiplier < 1 {
		return model.ErrRetryConfigMultiplier
	}

	var lastErr error = nil

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt == cfg.MaxAttempts {
			return lastErr
		}

		if cfg.OnRetry != nil {
			cfg.OnRetry(attempt, err)
		}

		delay := calculateDelay(attempt, cfg)
		delay = addJitter(delay)

		timer := time.NewTimer(delay)

		select {
		case <-timer.C:
			continue
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return ctx.Err()
		}
	}

	return lastErr

}

func calculateDelay(attempt int, cfg Config) time.Duration {
	multiToAttempt := math.Pow(cfg.Multiplier, float64(attempt-1))
	delay := time.Duration(float64(cfg.InitialDelay) * multiToAttempt)

	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}

	return delay
}

func addJitter(delay time.Duration) time.Duration {
	if delay < 0 {
		return 0
	}
	return time.Duration(rand.Int63n(int64(delay) + 1))
}
