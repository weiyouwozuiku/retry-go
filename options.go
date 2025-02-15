package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

type OnRetryFunc func(attempt uint, err error)
type RetryIfFunc func(error) bool
type DelayTypeFunc func(n uint, err error, config *Config) time.Duration
type Timer interface {
	After(time.Duration) <-chan time.Time
}
type Config struct {
	attempts         uint
	attemptsForError map[error]uint
	delay            time.Duration
	maxDelay         time.Duration
	onRetry          OnRetryFunc
	retryIf          RetryIfFunc
	delayType        DelayTypeFunc
	lastErrorOnly    bool
	context          context.Context
	timer            Timer
	maxBackOffN      uint
}
type Option func(*Config)

func newDefalutRetryConf() *Config {
	return &Config{
		attempts:         uint(5),
		attemptsForError: make(map[error]uint),
		delay:            time.Second,
		maxDelay:         time.Second * 30,
		onRetry:          func(n uint, err error) {},
		retryIf:          IsRecoverable,
		delayType:        BackOffDelay,
		context:          context.Background(),
		timer:            &timerImpl{},
	}
}

type unrecoverableErr struct {
	error
}

func IsRecoverable(err error) bool {
	return !errors.Is(err, unrecoverableErr{})
}

func Unrecoverable(err error) error {
	return unrecoverableErr{err}
}

// BackOffDelay calculates the delay for a given retry attempt.
// It uses a binary back-off strategy, where the delay is doubled for each
// attempt until the maximum delay is reached.
func BackOffDelay(n uint, _ error, config *Config) time.Duration {
	// 1 << 63 would overflow signed int64 (time.Duration), thus 62.
	const max uint = 62
	if config.maxBackOffN == 0 {
		if config.delay <= 0 {
			config.delay = 1
		}
		// Calculate the maximum backoff number.
		// The formula is: 2^maxBackOffN <= delay, so maxBackOffN = log2(delay)
		config.maxBackOffN = max - uint(math.Floor(math.Log2(float64(config.delay))))
	}
	if n > config.maxBackOffN {
		n = config.maxBackOffN
	}
	// Calculate the delay using the binary back-off strategy.
	// The formula is: delay = config.delay << n
	return config.delay << n
}
