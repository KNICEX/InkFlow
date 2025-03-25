package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/feed/internal/domain"
	"github.com/KNICEX/InkFlow/internal/feed/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/feed/internal/repo/dao"
	"github.com/samber/lo"
)

type FeedRepo interface {
	CreatePushFeed(ctx context.Context, push []domain.Feed) error
	CreatePullFeed(ctx context.Context, pull domain.Feed) error
	DeleteFeed(ctx context.Context, feed domain.Feed) error

	FindPullFeed(ctx context.Context, uids []int64, maxId, timestamp int64, limit int) ([]domain.Feed, error)
	FindPushFeed(ctx context.Context, uids int64, maxId, timestamp int64, limit int) ([]domain.Feed, error)
	FindPullByType(ctx context.Context, uids []int64, feedType domain.FeedType, maxId, timestamp int64, limit int, ) ([]domain.Feed, error)
	FindPushByType(ctx context.Context, uids int64, feedType domain.FeedType, maxId, timestamp int64, limit int) ([]domain.Feed, error)
}

type CachedFeedRepo struct {
	pushDAO dao.PushFeedDAO
	pullDAO dao.PullFeedDAO
	cache   cache.FeedCache
}

func (repo *CachedFeedRepo) CreatePullFeed(ctx context.Context, pull domain.Feed) error {
	return repo.pullDAO.CreatePull(ctx, repo.domainToPullFeed(pull))
}

func (repo *CachedFeedRepo) CreatePushFeed(ctx context.Context, push []domain.Feed) error {
	return repo.pushDAO.CreatePush(ctx, lo.Map(push, func(item domain.Feed, index int) dao.PushFeed {
		return repo.domainToPushFeed(item)
	}))
}

func (repo *CachedFeedRepo) DeleteFeed(ctx context.Context, feed domain.Feed) error {
	if err := repo.pushDAO.UpdateStatus(ctx, dao.PushFeed{
		FeedId:   feed.FeedId,
		UserId:   feed.UserId,
		FeedType: feed.FeedType,
		Status:   dao.FeedStatusHidden,
	}); err != nil {
		return err
	}

	if err := repo.pullDAO.UpdateStatus(ctx, dao.PullFeed{
		FeedId:   feed.FeedId,
		UserId:   feed.UserId,
		FeedType: feed.FeedType,
		Status:   dao.FeedStatusHidden,
	}); err != nil {
		return err
	}
	return nil
}

func (repo *CachedFeedRepo) FindPullFeed(ctx context.Context, uids []int64, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	pullFeeds, err := repo.pullDAO.FindPull(ctx, uids, maxId, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(pullFeeds, func(item dao.PullFeed, index int) domain.Feed {
		return repo.pullToDomain(item)
	}), nil
}

func (repo *CachedFeedRepo) FindPushFeed(ctx context.Context, uids int64, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	pushFeeds, err := repo.pushDAO.FindPush(ctx, uids, maxId, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(pushFeeds, func(item dao.PushFeed, index int) domain.Feed {
		return repo.pushToDomain(item)
	}), nil
}

func (repo *CachedFeedRepo) FindPullByType(ctx context.Context, uids []int64, feedType domain.FeedType, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	pullFeeds, err := repo.pullDAO.FindPullByType(ctx, uids, feedType, maxId, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(pullFeeds, func(item dao.PullFeed, index int) domain.Feed {
		return repo.pullToDomain(item)
	}), nil
}

func (repo *CachedFeedRepo) FindPushByType(ctx context.Context, uids int64, feedType domain.FeedType, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	pushFeeds, err := repo.pushDAO.FindPushByType(ctx, uids, feedType, maxId, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(pushFeeds, func(item dao.PushFeed, index int) domain.Feed {
		return repo.pushToDomain(item)
	}), nil
}

func (repo *CachedFeedRepo) pullToDomain(pull dao.PullFeed) domain.Feed {
	return domain.Feed{
		Id:        pull.Id,
		UserId:    pull.UserId,
		FeedType:  pull.FeedType,
		CreatedAt: pull.CreatedAt,
		UpdatedAt: pull.UpdatedAt,
		FeedId:    pull.FeedId,
		Content:   pull.Content,
	}
}

func (repo *CachedFeedRepo) pushToDomain(push dao.PushFeed) domain.Feed {
	return domain.Feed{
		Id:        push.Id,
		UserId:    push.UserId,
		FeedType:  push.FeedType,
		CreatedAt: push.CreatedAt,
		UpdatedAt: push.UpdatedAt,
		FeedId:    push.FeedId,
		Content:   push.Content,
	}
}

func (repo *CachedFeedRepo) domainToPullFeed(feed domain.Feed) dao.PullFeed {
	return dao.PullFeed{
		Id:        feed.Id,
		UserId:    feed.UserId,
		FeedType:  feed.FeedType,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		FeedId:    feed.FeedId,
		Content:   feed.Content,
		Status:    dao.FeedStatusNormal,
	}
}

func (repo *CachedFeedRepo) domainToPushFeed(feed domain.Feed) dao.PushFeed {
	return dao.PushFeed{
		Id:        feed.Id,
		UserId:    feed.UserId,
		FeedType:  feed.FeedType,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		FeedId:    feed.FeedId,
		Content:   feed.Content,
		Status:    dao.FeedStatusNormal,
	}
}
