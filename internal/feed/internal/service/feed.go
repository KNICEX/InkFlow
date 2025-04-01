package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/action"
	"github.com/KNICEX/InkFlow/internal/feed/internal/domain"
	"github.com/KNICEX/InkFlow/internal/feed/internal/repo"
	"github.com/KNICEX/InkFlow/internal/relation"
	"golang.org/x/sync/errgroup"
	"sort"
	"time"
)

type FeedService interface {
	CreateFeed(ctx context.Context, feed domain.Feed) error
	FollowFeedInkList(ctx context.Context, uid int64, maxId, timestamp int64, limit int) ([]domain.Feed, error)
}

type feedService struct {
	repo      repo.FeedRepo
	followSvc relation.FollowService
	actionSvc action.Service
}

func NewFeedService(repo repo.FeedRepo, followSvc relation.FollowService, actionSvc action.Service) FeedService {
	return &feedService{
		repo:      repo,
		followSvc: followSvc,
		actionSvc: actionSvc,
	}
}

func (f *feedService) CreateFeed(ctx context.Context, feed domain.Feed) error {
	if err := f.repo.CreatePullFeed(ctx, feed); err != nil {
		return err
	}

	// TODO 现在采用全推模型，后续考虑混合
	followerIds, err := f.followSvc.FollowerIds(ctx, feed.UserId, 0, 10000)
	if err != nil {
		return err
	}
	// 30 天内活跃用户
	activeUsers, err := f.actionSvc.FindActiveUser(ctx, followerIds, time.Now().Add(-time.Hour*24*30))
	if err != nil {
		return err
	}
	pushFeeds := make([]domain.Feed, 0, len(activeUsers))
	for _, user := range activeUsers {
		feed.UserId = user.Id
		pushFeeds = append(pushFeeds, feed)
	}
	return f.repo.CreatePushFeed(ctx, pushFeeds)
}

func (f *feedService) FollowFeedInkList(ctx context.Context, uid, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	var pushFeeds []domain.Feed
	var pullFeeds []domain.Feed

	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		pushFeeds, er = f.repo.FindPushFeed(ctx, uid, maxId, timestamp, limit)
		return er
	})
	eg.Go(func() error {
		// TODO 这里先查2000，后续修改
		followingIds, er := f.followSvc.FollowingIds(ctx, uid, 0, 2000)
		if er != nil {
			return er
		}
		pullFeeds, er = f.repo.FindPullFeed(ctx, followingIds, maxId, timestamp, limit)
		return er
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	feeds := append(pushFeeds, pullFeeds...)
	sort.Slice(feeds, func(i, j int) bool {
		return feeds[i].CreatedAt.After(feeds[j].CreatedAt)
	})

	return feeds[:min(len(feeds), limit)], nil
}
