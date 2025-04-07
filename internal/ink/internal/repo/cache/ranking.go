package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, ids []int64) error
	Get(ctx context.Context) ([]int64, error)
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
func (r *RedisRankingCache) Set(ctx context.Context, ids []int64) error {

	val, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	// 过期时间要设置得比定时计算的间隔长
	return r.cmd.Set(ctx, r.key(), val,
		r.expiration).Err()
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]int64, error) {
	val, err := r.cmd.Get(ctx, r.key()).Bytes()
	if err != nil {
		return nil, err
	}
	var ids []int64
	err = json.Unmarshal(val, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *RedisRankingCache) key() string {
	return "ranking:ink"
}
