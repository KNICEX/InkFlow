package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/cache"
)

type RankingRepo interface {
	ReplaceTopN(ctx context.Context, ids []int64) error
	FindTop(ctx context.Context, offset, limit int) ([]int64, error)
}
type rankingRepo struct {
	rankingCache cache.RankingCache
}

func NewRankingRepo(rankingCache cache.RankingCache) RankingRepo {
	return &rankingRepo{
		rankingCache: rankingCache,
	}
}
func (r *rankingRepo) ReplaceTopN(ctx context.Context, ids []int64) error {
	return r.rankingCache.Set(ctx, ids)
}

func (r *rankingRepo) FindTop(ctx context.Context, offset, limit int) ([]int64, error) {
	return r.rankingCache.Get(ctx, offset, limit)
}
