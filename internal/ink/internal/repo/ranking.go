package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/cache"
)

type RankingRepo interface {
	ReplaceTopNInks(ctx context.Context, ids []int64) error
	FindTopInk(ctx context.Context, offset, limit int) ([]int64, error)

	ReplaceTopNTags(ctx context.Context, tags []domain.TagStats) error
	FindTopTag(ctx context.Context, offset, limit int) ([]domain.TagStats, error)
}
type rankingRepo struct {
	rankingCache cache.RankingCache
}

func NewRankingRepo(rankingCache cache.RankingCache) RankingRepo {
	return &rankingRepo{
		rankingCache: rankingCache,
	}
}
func (r *rankingRepo) ReplaceTopNInks(ctx context.Context, ids []int64) error {
	return r.rankingCache.SetInkIds(ctx, ids)
}

func (r *rankingRepo) FindTopInk(ctx context.Context, offset, limit int) ([]int64, error) {
	return r.rankingCache.GetInkIds(ctx, offset, limit)
}

func (r *rankingRepo) ReplaceTopNTags(ctx context.Context, tags []domain.TagStats) error {
	return r.rankingCache.SetTags(ctx, tags)
}

func (r *rankingRepo) FindTopTag(ctx context.Context, offset, limit int) ([]domain.TagStats, error) {
	return r.rankingCache.GetTags(ctx, offset, limit)
}
