package service

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/redis/go-redis/v9"
	"math"
	"strconv"
	"time"
)

const inkBiz = "ink"

type ScoreFunc func(likeCnt int64, updateTime time.Time) float64

type RankingService interface {
	TopN(ctx context.Context, n int) error
	FindTopN(ctx context.Context, offset int, limit int) ([]domain.Ink, error)
}
type BatchRankingService struct {
	inkSvc      inkService
	intrSvc     interactive.Service
	rankingRepo repo.RankingRepo
	cmd         redis.Cmdable
	scoreFunc   ScoreFunc
}

func NewBatchRankingService(inkSvc inkService, intr interactive.Service, cmd redis.Cmdable) RankingService {
	res := &BatchRankingService{
		inkSvc:  inkSvc,
		cmd:     cmd,
		intrSvc: intr,
	}
	res.scoreFunc = res.score
	return res
}
func (b *BatchRankingService) TopN(ctx context.Context, n int) error {
	ids, err := b.rankTopN(ctx, n)
	if err != nil {
		return err
	}
	return b.rankingRepo.ReplaceTopN(ctx, ids)
}

func (b *BatchRankingService) FindTopN(ctx context.Context, offset int, limit int) ([]domain.Ink, error) {
	inks, err := b.rankingRepo.FindTopN(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (b *BatchRankingService) rankTopN(ctx context.Context, n int) ([]int64, error) {
	var (
		lastId  int64
		zsetKey = b.rankingKey()
	)
	// 用于批量添加到 ZSET 的数据
	var zsetMembers []redis.Z

	for {
		inks, err := b.inkSvc.ListAllLive(ctx, lastId, n)
		if err != nil {
			return nil, fmt.Errorf("failed to list live inks: %w", err)
		}
		if len(inks) == 0 {
			break
		}
		if len(inks) == 0 || inks[len(inks)-1].UpdatedAt.Before(time.Now().Add(-time.Hour*24*7)) {
			break
		}
		for _, ink := range inks {
			intr, err := b.intrSvc.Get(ctx, inkBiz, ink.Id, 0)
			if err != nil {
				continue
			}

			score := b.score(intr.LikeCnt, ink.UpdatedAt)
			zsetMembers = append(zsetMembers, redis.Z{
				Score:  score,
				Member: ink.Id,
			})
		}

		// 批量添加到 ZSET
		if _, err := b.cmd.ZAdd(ctx, zsetKey, zsetMembers...).Result(); err != nil {
			return nil, fmt.Errorf("failed to add scores to ZSET: %w", err)
		}

		// 清理 ZSET，只保留前 N 个元素
		if _, err := b.cmd.ZRemRangeByRank(ctx, zsetKey, 0, int64(n-1)).Result(); err != nil {
			return nil, fmt.Errorf("failed to trim ZSET to top %d elements: %w", n, err)
		}
		zsetMembers = zsetMembers[:0]
		lastId = inks[len(inks)-1].Id
	}
	members, err := b.cmd.ZRangeWithScores(ctx, zsetKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all members from ZSET: %w", err)
	}
	var ids []int64
	for _, member := range members {
		inkId, err := strconv.ParseInt(member.Member.(string), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, inkId)
	}
	return ids, nil
}

func (b *BatchRankingService) score(likeCnt int64, updateTime time.Time) float64 {
	// 这个 factor 也可以做成一个参数
	const factor = 1.5
	return float64(likeCnt-1) /
		math.Pow(time.Since(updateTime).Hours()+2, factor)
}
func (b *BatchRankingService) rankingKey() string {
	return fmt.Sprintf("ranking:sort:%s", "ink")
}
