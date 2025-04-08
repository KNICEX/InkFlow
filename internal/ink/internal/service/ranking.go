package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/pkg/queuex"
	"github.com/samber/lo"
	"math"
	"time"
)

const inkBiz = "ink"

type ScoreFunc func(likeCnt, favoriteCnt int64, createdAt time.Time) float64

type RankingService interface {
	TopNInk(ctx context.Context, n int) error
	FindTopNInk(ctx context.Context, offset int, limit int) ([]domain.Ink, error)

	TopNTag(ctx context.Context, n int) error
	FindTopNTag(ctx context.Context, offset int, limit int) ([]domain.TagStats, error)
}
type BatchRankingService struct {
	inkRepo     repo.LiveInkRepo
	rankingRepo repo.RankingRepo
	intrSvc     interactive.Service
	scoreFunc   ScoreFunc
}

func NewBatchRankingService(inkRepo repo.LiveInkRepo, rankRepo repo.RankingRepo, intrSvc interactive.Service) RankingService {
	return &BatchRankingService{
		inkRepo:     inkRepo,
		rankingRepo: rankRepo,
		intrSvc:     intrSvc,
		scoreFunc: func(likeCnt, favoriteCnt int64, createdAt time.Time) float64 {
			// 这个 factor 也可以做成一个参数
			const factor = 1.5
			return float64(likeCnt-1) + float64(favoriteCnt-1)*1.2/
				math.Pow(time.Since(createdAt).Hours()+2, factor)
		},
	}
}
func (b *BatchRankingService) TopNInk(ctx context.Context, n int) error {
	ids, err := b.rankTopN(ctx, n, time.Now().Add(-time.Hour*24*7))
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	return b.rankingRepo.ReplaceTopNInks(ctx, ids)
}

func (b *BatchRankingService) FindTopNInk(ctx context.Context, offset int, limit int) ([]domain.Ink, error) {
	ids, err := b.rankingRepo.FindTopInk(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	inkMap, err := b.inkRepo.FindByIds(ctx, ids, domain.InkStatusPublished)
	inks := make([]domain.Ink, 0, len(ids))
	for _, id := range ids {
		if ink, ok := inkMap[id]; ok {
			inks = append(inks, ink)
		}
	}
	return inks, nil
}

func (b *BatchRankingService) rankTopN(ctx context.Context, n int, startTime time.Time) ([]int64, error) {
	var (
		maxId     int64 = 0
		batchSize       = 100
	)
	zq := queuex.NewZQueue[float64, domain.Ink](n)

	for {
		inks, err := b.inkRepo.FindAll(ctx, maxId, batchSize, domain.InkStatusPublished)
		if err != nil {
			return nil, err
		}

		if len(inks) == 0 {
			break
		}

		intrs, err := b.intrSvc.GetMulti(ctx, inkBiz, lo.Map(inks, func(item domain.Ink, index int) int64 {
			return item.Id
		}), 0)
		if err != nil {
			return nil, err
		}

		for _, ink := range inks {
			if ink.CreatedAt.Before(startTime) {
				break
			}

			score := b.scoreFunc(intrs[ink.Id].LikeCnt, intrs[ink.Id].FavoriteCnt, ink.CreatedAt)
			if score > 0 {
				zq.Enqueue(score, ink)
			}
		}
		maxId = inks[len(inks)-1].Id
	}

	return lo.Map(zq.AllValues(), func(item domain.Ink, index int) int64 {
		return item.Id
	}), nil
}

func (b *BatchRankingService) TopNTag(ctx context.Context, n int) error {
	tags, err := b.rankTopNTag(ctx, n, time.Now().Add(-time.Hour*24))
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return nil
	}
	return b.rankingRepo.ReplaceTopNTags(ctx, tags)
}

func (b *BatchRankingService) rankTopNTag(ctx context.Context, n int, startTime time.Time) ([]domain.TagStats, error) {
	var (
		maxId     int64 = 0
		batchSize       = 100
	)

	tagMap := make(map[string]int64)

	for {
		inks, err := b.inkRepo.FindAll(ctx, maxId, batchSize)
		if err != nil {
			return nil, err
		}

		if len(inks) == 0 {
			break
		}

		for _, ink := range inks {
			if ink.CreatedAt.Before(startTime) {
				break
			}
			for _, tag := range ink.Tags {
				tagMap[tag]++
			}
		}
		maxId = inks[len(inks)-1].Id
	}

	zq := queuex.NewZQueue[int64, string](n)
	for tag, cnt := range tagMap {
		zq.Enqueue(cnt, tag)
	}
	return lo.Map(zq.All(), func(item queuex.ZQueueItem[int64, string], index int) domain.TagStats {
		return domain.TagStats{
			Name:      item.Value,
			Reference: item.Score,
		}
	}), nil
}

func (b *BatchRankingService) FindTopNTag(ctx context.Context, offset int, limit int) ([]domain.TagStats, error) {
	return b.rankingRepo.FindTopTag(ctx, offset, limit)
}
