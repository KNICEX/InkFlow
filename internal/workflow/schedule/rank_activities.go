package schedule

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink"
)

type RankActivities struct {
	rankingService ink.RankingService
}

func NewRankActivities(rankingService ink.RankingService) *RankActivities {
	return &RankActivities{
		rankingService: rankingService,
	}
}

func (r *RankActivities) RankInk(ctx context.Context, n int) error {
	return r.rankingService.TopNInk(ctx, n)
}

func (r *RankActivities) RankTag(ctx context.Context, n int) error {
	return r.rankingService.TopNTag(ctx, n)
}
