package cache

import (
	"context"
	"encoding/json"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, inks []domain.Ink) error
	Get(ctx context.Context) ([]domain.Ink, error)
}

type RedisRankingCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisRankingCache(cmd redis.Cmdable, expiration time.Duration) RankingCache {
	return &RedisRankingCache{
		cmd:        cmd,
		expiration: expiration,
	}
}
func (r *RedisRankingCache) Set(ctx context.Context, inks []domain.Ink) error {
	for i := 0; i < len(inks); i++ {
		inks[i].ContentHtml = inks[i].Abstract()
	}
	val, err := json.Marshal(inks)
	if err != nil {
		return err
	}
	// 过期时间要设置得比定时计算的间隔长
	return r.cmd.Set(ctx, r.key(), val,
		r.expiration).Err()
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Ink, error) {
	val, err := r.cmd.Get(ctx, r.key()).Bytes()
	if err != nil {
		return nil, err
	}
	var inks []domain.Ink
	err = json.Unmarshal(val, &inks)
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (r *RedisRankingCache) key() string {
	return "ranking:ink"
}
