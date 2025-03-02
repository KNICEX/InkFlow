package ratelimit

import "context"

type Limiter interface {
	Limited(ctx context.Context) (bool, error)
}

type KeyLimiter interface {
	Limited(ctx context.Context, key string) (bool, error)
}

type Waiter interface {
	Wait(ctx context.Context) error
}

type KeyWaiter interface {
	Wait(ctx context.Context, key string) error
}
