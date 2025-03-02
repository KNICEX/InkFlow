package ratelimit

import (
	"context"
	"github.com/ecodeclub/ekit/syncx"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type WindowSliceLimiter struct {
	lim *rate.Limiter
}

func NewWindowSliceLimiter(interval time.Duration, rateNum int) Limiter {
	l := rate.NewLimiter(rate.Every(interval), rateNum)
	return &WindowSliceLimiter{
		lim: l,
	}
}

func (w WindowSliceLimiter) Limited(ctx context.Context) (bool, error) {
	return !w.lim.Allow(), nil
}

type SliceWindowKeyLimiter struct {
	interval time.Duration
	rateNum  int
	limMap   map[string]Limiter
	expMap   syncx.Map[string, time.Time]
	mu       sync.RWMutex
}

func NewSliceWindowKeyLimiter(interval time.Duration, rateNum int) KeyLimiter {
	l := &SliceWindowKeyLimiter{
		interval: interval,
		rateNum:  rateNum,
		limMap:   make(map[string]Limiter),
	}
	return l
}

func (k *SliceWindowKeyLimiter) Limited(ctx context.Context, key string) (bool, error) {
	k.mu.RLock()
	lim, ok := k.limMap[key]
	k.mu.RUnlock()
	if !ok {
		k.mu.Lock()
		defer k.mu.Unlock()
		// double check
		if lim, ok = k.limMap[key]; !ok {
			lim = NewWindowSliceLimiter(k.interval, k.rateNum)
			k.limMap[key] = lim
		}
	}
	// 重置过期时间
	k.expMap.Store(key, time.Now().Add(k.interval))
	defer k.expKey(key)
	return lim.Limited(ctx)
}

// expKey 删除过期的key
func (k *SliceWindowKeyLimiter) expKey(key string) {
	// 给2倍时间, 避免可能误差
	time.AfterFunc(k.interval*2, func() {
		expireAt, ok := k.expMap.Load(key)
		if !ok {
			return
		}
		if time.Now().After(expireAt) {
			k.expMap.Delete(key)
			k.mu.Lock()
			defer k.mu.Unlock()
			delete(k.limMap, key)
		}
	})
}
