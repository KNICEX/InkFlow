package gorse

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/pkg/gorsex"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/queuex"
	client "github.com/gorse-io/gorse-go"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"strconv"
	"sync"
)

type RecommendService struct {
	cli       *gorsex.Client
	l         logx.Logger
	followSvc relation.FollowService
	intrSvc   interactive.Service
}

func NewRecommendService(cli *gorsex.Client, followSvc relation.FollowService, intrSvc interactive.Service, l logx.Logger) *RecommendService {
	return &RecommendService{
		cli:       cli,
		l:         l,
		followSvc: followSvc,
		intrSvc:   intrSvc,
	}
}

func (svc *RecommendService) FindSimilarInk(ctx context.Context, inkId int64, offset, limit int) ([]int64, error) {
	// TODO 这个sdk没有实现offset
	scores, err := svc.cli.GetNeighbors(ctx, strconv.FormatInt(inkId, 10), limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(scores, func(item client.Score, index int) int64 {
		id, err := strconv.ParseInt(item.Id, 10, 64)
		svc.l.WithCtx(ctx).Error("gorse recommend ink parse id error", logx.String("id", item.Id), logx.Error(err))
		return id
	}), nil
}

func (svc *RecommendService) FindSimilarUser(ctx context.Context, userId int64, offset, limit int) ([]int64, error) {
	scores, err := svc.cli.GetNeighborsUsers(ctx, strconv.FormatInt(userId, 10), limit, offset)
	if err != nil {
		return nil, err
	}
	return lo.Map(scores, func(item client.Score, index int) int64 {
		id, err := strconv.ParseInt(item.Id, 10, 64)
		svc.l.WithCtx(ctx).Error("gorse recommend user parse id error", logx.String("id", item.Id), logx.Error(err))
		return id
	}), nil

}

func (svc *RecommendService) FindSimilarAuthor(ctx context.Context, authorId int64, offset, limit int) ([]int64, error) {
	return nil, nil
}

func (svc *RecommendService) FindPopular(ctx context.Context, offset, limit int) ([]int64, error) {
	return nil, nil
}

func (svc *RecommendService) FindRecommendInk(ctx context.Context, userId int64, offset, limit int) ([]int64, error) {
	ids, err := svc.cli.GetRecommend(ctx, strconv.FormatInt(userId, 10), "", limit, offset)
	if err != nil {
		return nil, err
	}
	return lo.Map(ids, func(item string, index int) int64 {
		id, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			svc.l.WithCtx(ctx).Error("gorse recommend ink parse id error", logx.String("id", item), logx.Error(err))
		}
		return id
	}), nil

}

func (svc *RecommendService) FindRecommendAuthor(ctx context.Context, userId int64, offset, limit int) ([]int64, error) {
	similarUids, err := svc.FindSimilarUser(ctx, userId, offset, limit)
	if err != nil {
		return nil, err
	}
	eg := errgroup.Group{}
	var following []int64
	eg.Go(func() error {
		var er error
		following, er = svc.followSvc.FollowingIds(ctx, userId, 0, 10000)
		return er
	})

	// 查找相近用户的关注列表
	similarFollowing := make(map[int64]relation.FollowStatistic)
	mu := sync.Mutex{}
	for _, uid := range similarUids {
		eg.Go(func() error {
			follow, er := svc.followSvc.FollowingList(ctx, uid, 0, 0, 100)
			if er != nil {
				return er
			}
			mu.Lock()
			defer mu.Unlock()
			for _, f := range follow {
				similarFollowing[f.Uid] = f
			}
			return nil
		})
	}
	if err = eg.Wait(); err != nil {
		return nil, err
	}

	// 过滤掉已经关注的
	for _, uid := range following {
		delete(similarFollowing, uid)
	}

	// 过滤掉自己
	delete(similarFollowing, userId)

	// 挑关注数量多的
	q := queuex.NewPriorityQueue(limit, func(src relation.FollowStatistic, dst relation.FollowStatistic) int {
		return int(src.Followers - dst.Followers)
	})
	for _, v := range similarFollowing {
		q.Enqueue(v)
	}
	mostPopular := q.All()
	return lo.Map(mostPopular, func(item relation.FollowStatistic, index int) int64 {
		return item.Uid
	}), nil
}
