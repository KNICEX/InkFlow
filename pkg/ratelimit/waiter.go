package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrClosed = errors.New("limiter is closed")

type SliceWindowWaiter struct {
	ch       chan time.Time
	interval time.Duration
	rate     int
	closed   bool
	once     sync.Once
}

func NewSliceWindowWaiter(interval time.Duration, rate int) (waiter Waiter, closeFunc func()) {
	w := &SliceWindowWaiter{
		ch:       make(chan time.Time, rate),
		interval: interval,
		rate:     rate,
	}

	go func() {
		for t := range w.ch {
			diff := time.Now().Sub(t)
			if diff > interval {
				continue
			}
			time.Sleep(interval - diff)
		}
	}()

	return w, w.Close
}
func (s *SliceWindowWaiter) Wait(ctx context.Context) error {
	if s.closed {
		// 关于closed后是
		return ErrClosed
	}
	select {
	case s.ch <- time.Now():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *SliceWindowWaiter) Close() {
	s.once.Do(func() {
		s.closed = true
		close(s.ch)
	})
}
