package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/cache"
)

type RankingRepo interface {
	ReplaceTopN(ctx context.Context, ids []int64) error
	FindTopN(ctx context.Context, offset, limit int) ([]domain.Ink, error)
}
type rankingRepo struct {
	rankingCache cache.RankingCache
	inkRepo      CachedLiveInkRepo
}

func NewRankingRepo(rankingCache cache.RankingCache, inkRepo CachedLiveInkRepo) RankingRepo {
	return &rankingRepo{
		rankingCache: rankingCache,
		inkRepo:      inkRepo,
	}
}
func (r *rankingRepo) ReplaceTopN(ctx context.Context, ids []int64) error {
	return r.rankingCache.Set(ctx, ids)
}

func (r *rankingRepo) FindTopN(ctx context.Context, offset, limit int) ([]domain.Ink, error) {

	ids, err := r.rankingCache.Get(ctx)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 || offset >= len(ids) {
		return nil, nil
	}
	if len(ids) >= offset+limit {
		ids = ids[offset : offset+limit]
	} else {
		ids = ids[offset:]
	}
	inksMap, err := r.inkRepo.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	inks := make([]domain.Ink, 0, len(ids))
	for _, id := range ids {
		if ink, ok := inksMap[id]; ok {
			inks = append(inks, ink)
		}
	}
	return inks, nil
}
