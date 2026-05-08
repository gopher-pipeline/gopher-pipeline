package retry

import "errors"

var (
	ErrRetryConfigAttempts     = errors.New("retry config attempts must be more or equal 1")
	ErrRetryConfigInitialDelay = errors.New("retry config init delay must me more or equal 1")
	ErrRetryConfigMaxDelay     = errors.New("retry config max delay must be more or equal 1")
	ErrRetryConfigMultiplier   = errors.New("retry config multiplier must be more or equal 1")
)
