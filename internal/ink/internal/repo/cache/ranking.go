package cache

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type RankingCache interface {
	SetInkIds(ctx context.Context, ids []int64) error
	GetInkIds(ctx context.Context, offset, limit int) ([]int64, error)

	SetTags(ctx context.Context, tags []domain.TagStats) error
	GetTags(ctx context.Context, offset, limit int) ([]domain.TagStats, error)
}

type RedisRankingCache struct {
	cmd redis.Cmdable
	l   logx.Logger
}

func NewRedisRankingCache(cmd redis.Cmdable, l logx.Logger) RankingCache {
	return &RedisRankingCache{
		cmd: cmd,
		l:   l,
	}
}
func (r *RedisRankingCache) SetInkIds(ctx context.Context, ids []int64) error {
	zs := make([]redis.Z, 0, len(ids))
	n := len(ids)
	for i := len(ids) - 1; i >= 0; i-- {
		zs = append(zs, redis.Z{
			Score:  float64(n - i),
			Member: ids[i],
		})
	}

	pipeline := r.cmd.Pipeline()
	pipeline.Del(ctx, r.inkKey())
	pipeline.ZAdd(ctx, r.inkKey(), zs...)
	_, err := pipeline.Exec(ctx)
	return err
}

func (r *RedisRankingCache) GetInkIds(ctx context.Context, offset, limit int) ([]int64, error) {
	res, err := r.cmd.ZRevRange(ctx, r.inkKey(), int64(offset), int64(offset+limit-1)).Result()
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

func (r *RedisRankingCache) SetTags(ctx context.Context, tags []domain.TagStats) error {
	zs := make([]redis.Z, 0, len(tags))
	for _, tag := range tags {
		zs = append(zs, redis.Z{
			Score:  float64(tag.Reference),
			Member: tag.Name,
		})
	}

	pipeline := r.cmd.Pipeline()
	pipeline.Del(ctx, r.tagKey())
	pipeline.ZAdd(ctx, r.tagKey(), zs...)
	_, err := pipeline.Exec(ctx)
	return err
}

func (r *RedisRankingCache) GetTags(ctx context.Context, offset, limit int) ([]domain.TagStats, error) {
	res, err := r.cmd.ZRevRangeWithScores(ctx, r.tagKey(), int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, err
	}
	tags := make([]domain.TagStats, 0, len(res))
	for _, s := range res {
		tag := domain.TagStats{
			Name:      s.Member.(string),
			Reference: int64(s.Score),
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *RedisRankingCache) inkKey() string {
	return "ranking:ink"
}

func (r *RedisRankingCache) tagKey() string {
	return "ranking:tag"
}
