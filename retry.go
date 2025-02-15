package retry

import (
	"sync"
	"sync/atomic"
	"time"
)

type RetryableFunc func() error

type RetryableFuncWithData[T any] func() (T, error)

type timerImpl struct {
	after atomic.Value
	mu    sync.Mutex
}

func (t *timerImpl) After(d time.Duration) <-chan time.Time {
	a := t.after.Load()
	if a != nil {
		return a.(<-chan time.Time)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	a = t.after.Load()
	if a == nil {
		a = time.After(d)
		t.after.Store(a)
	}
	return a.(<-chan time.Time)
}

func DoWithData[T any](fn RetryableFuncWithData[T], opts ...Option) (T, error) {
	var emptyT T
	conf := newDefalutRetryConf()
	for _, opt := range opts {
		opt(conf)
	}
	if err := conf.context.Err(); err != nil {
		return emptyT, err
	}
	var lastErr error
	if conf.attempts == 0 {
		for {
、
		}
	}
}
