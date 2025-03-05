package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/code/internal/repo/cache"
	"time"
)

var (
	ErrCodeSendTooMany = cache.ErrCodeSendTooMany
	ErrCodeVerifyLimit = cache.ErrCodeVerifyLimit
)

type CodeRepo interface {
	Store(ctx context.Context, biz, recipient, code string, effectiveTime time.Duration, resendInterval time.Duration, maxRetry int) error
	Verify(ctx context.Context, biz, recipient, inputCode string) (bool, error)
}

var _ CodeRepo = (*CachedCodeRepo)(nil)

type CachedCodeRepo struct {
	cache cache.CodeCache
}

func NewCodeRepo(cache cache.CodeCache) CodeRepo {
	return &CachedCodeRepo{cache: cache}
}

func (r *CachedCodeRepo) Store(ctx context.Context, biz, recipient, code string, effectiveTime time.Duration, resendInterval time.Duration, maxRetry int) error {
	return r.cache.Set(ctx, biz, recipient, code, effectiveTime, resendInterval, maxRetry)
}
func (r *CachedCodeRepo) Verify(ctx context.Context, biz, recipient, inputCode string) (bool, error) {
	return r.cache.Verify(ctx, biz, recipient, inputCode)
}
