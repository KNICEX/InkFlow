package cache

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, ids []int64) error
	Get(ctx context.Context, offset, limit int) ([]int64, error)
}

type RedisRankingCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
	l          logx.Logger
}

func NewRedisRankingCache(cmd redis.Cmdable, expiration time.Duration) RankingCache {
	return &RedisRankingCache{
		cmd:        cmd,
		expiration: expiration,
	}
}
func (r *RedisRankingCache) Set(ctx context.Context, ids []int64) error {
	zs := make([]redis.Z, 0, len(ids))
	n := len(ids)
	for i := len(ids) - 1; i >= 0; i-- {
		zs = append(zs, redis.Z{
			Score:  float64(n - i),
			Member: ids[i],
		})
	}

	pipeline := r.cmd.Pipeline()
	pipeline.Del(ctx, r.key())
	pipeline.ZAdd(ctx, r.key(), zs...)
	_, err := pipeline.Exec(ctx)
	return err
}

func (r *RedisRankingCache) Get(ctx context.Context, offset, limit int) ([]int64, error) {
	res, err := r.cmd.ZRange(ctx, r.key(), int64(offset), int64(limit)).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(res))
	for _, s := range res {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			r.l.Error("ink ranking cache parse id error", logx.Error(err), logx.String("id", s))
			return nil, err
		}

		ids = append(ids, id)
	}
	return ids, nil
}

func (r *RedisRankingCache) key() string {
	return "ranking:ink"
}
