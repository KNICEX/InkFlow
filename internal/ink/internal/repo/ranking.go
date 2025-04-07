package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/cache"
)

type RankingRepo interface {
	ReplaceTopN(ctx context.Context, inks []domain.Ink) error
	FindTopN(ctx context.Context) ([]domain.Ink, error)
}
type rankingRepo struct {
	rankingCache cache.RankingCache
}

func NewRankingRepo(rankingCache cache.RankingCache) RankingRepo {
	return &rankingRepo{
		rankingCache: rankingCache,
	}
}
func (r *rankingRepo) ReplaceTopN(ctx context.Context, inks []domain.Ink) error {
	return r.rankingCache.Set(ctx, inks)
}

func (r *rankingRepo) FindTopN(ctx context.Context) ([]domain.Ink, error) {
	return r.rankingCache.Get(ctx)
}
