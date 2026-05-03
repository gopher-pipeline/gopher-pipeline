package ratelimit

import (
	"context"
	"time"
)

type Limiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	done   chan struct{}
}

func NewLimiter(rps int) *Limiter {
	if rps <= 0 {
		rps = 5
	}

	tokens := make(chan struct{}, rps)
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	done := make(chan struct{}, 1)

	tokenBucket := &Limiter{
		tokens: tokens,
		ticker: ticker,
		done:   done,
	}

	go func() {
		for {
			select {
			case <-tokenBucket.ticker.C:
				select {
				case tokenBucket.tokens <- struct{}{}:
				default:
					// Bucket full
				}
			case <-tokenBucket.done:
				return

			}
		}
	}()

	return tokenBucket
}

func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.tokens:
		return nil
	}
}

func (l *Limiter) Stop() {
	l.ticker.Stop()
	close(l.done)
}
