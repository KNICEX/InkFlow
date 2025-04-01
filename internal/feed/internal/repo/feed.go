package repo

import (
	"context"
	"encoding/json"
	"github.com/KNICEX/InkFlow/internal/feed/internal/domain"
	"github.com/KNICEX/InkFlow/internal/feed/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/samber/lo"
)

type FeedRepo interface {
	CreatePushFeed(ctx context.Context, push []domain.Feed) error
	CreatePullFeed(ctx context.Context, pull domain.Feed) error
	DeleteFeed(ctx context.Context, feed domain.Feed) error

	FindPullFeed(ctx context.Context, uids []int64, maxId, timestamp int64, limit int) ([]domain.Feed, error)
	FindPushFeed(ctx context.Context, uids int64, maxId, timestamp int64, limit int) ([]domain.Feed, error)
	FindPullByType(ctx context.Context, uids []int64, biz string, maxId, timestamp int64, limit int, ) ([]domain.Feed, error)
	FindPushByType(ctx context.Context, uids int64, biz string, maxId, timestamp int64, limit int) ([]domain.Feed, error)
}

type NoCacheFeedRepo struct {
	pushDAO dao.PushFeedDAO
	pullDAO dao.PullFeedDAO
	l       logx.Logger
}

func NewNoCacheFeedRepo(push dao.PushFeedDAO, pull dao.PullFeedDAO, l logx.Logger) FeedRepo {
	return &NoCacheFeedRepo{
		pullDAO: pull,
		pushDAO: push,
		l:       l,
	}
}

func (repo *NoCacheFeedRepo) CreatePullFeed(ctx context.Context, pull domain.Feed) error {
	return repo.pullDAO.CreatePull(ctx, repo.domainToPullFeed(pull))
}

func (repo *NoCacheFeedRepo) CreatePushFeed(ctx context.Context, push []domain.Feed) error {
	return repo.pushDAO.CreatePush(ctx, lo.Map(push, func(item domain.Feed, index int) dao.PushFeed {
		return repo.domainToPushFeed(item)
	}))
}

func (repo *NoCacheFeedRepo) DeleteFeed(ctx context.Context, feed domain.Feed) error {
	if err := repo.pushDAO.UpdateStatus(ctx, dao.PushFeed{
		BizId:  feed.BizId,
		UserId: feed.UserId,
		Biz:    feed.Biz,
		Status: dao.FeedStatusHidden,
	}); err != nil {
		return err
	}

	if err := repo.pullDAO.UpdateStatus(ctx, dao.PullFeed{
		BizId:  feed.BizId,
		UserId: feed.UserId,
		Biz:    feed.Biz,
		Status: dao.FeedStatusHidden,
	}); err != nil {
		return err
	}
	return nil
}

func (repo *NoCacheFeedRepo) FindPullFeed(ctx context.Context, uids []int64, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	pullFeeds, err := repo.pullDAO.FindPull(ctx, uids, maxId, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(pullFeeds, func(item dao.PullFeed, index int) domain.Feed {
		return repo.pullToDomain(item)
	}), nil
}

func (repo *NoCacheFeedRepo) FindPushFeed(ctx context.Context, uids int64, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	pushFeeds, err := repo.pushDAO.FindPush(ctx, uids, maxId, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(pushFeeds, func(item dao.PushFeed, index int) domain.Feed {
		return repo.pushToDomain(item)
	}), nil
}

func (repo *NoCacheFeedRepo) FindPullByType(ctx context.Context, uids []int64, biz string, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	pullFeeds, err := repo.pullDAO.FindPullByType(ctx, uids, biz, maxId, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(pullFeeds, func(item dao.PullFeed, index int) domain.Feed {
		return repo.pullToDomain(item)
	}), nil
}

func (repo *NoCacheFeedRepo) FindPushByType(ctx context.Context, uids int64, biz string, maxId, timestamp int64, limit int) ([]domain.Feed, error) {
	pushFeeds, err := repo.pushDAO.FindPushByBiz(ctx, uids, biz, maxId, timestamp, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(pushFeeds, func(item dao.PushFeed, index int) domain.Feed {
		return repo.pushToDomain(item)
	}), nil
}

func (repo *NoCacheFeedRepo) pullToDomain(pull dao.PullFeed) domain.Feed {
	var content any
	switch pull.Biz {
	case domain.BizInk:
		content = domain.FeedInk{}
		if err := json.Unmarshal([]byte(pull.Content), &content); err != nil {
			repo.l.Error("json unmarshal feed content failed", logx.Any("feed", pull), logx.Error(err))
		}
	}

	return domain.Feed{
		Id:        pull.Id,
		UserId:    pull.UserId,
		Biz:       pull.Biz,
		CreatedAt: pull.CreatedAt,
		UpdatedAt: pull.UpdatedAt,
		BizId:     pull.BizId,
		Content:   content,
	}
}

func (repo *NoCacheFeedRepo) pushToDomain(push dao.PushFeed) domain.Feed {
	var content any
	switch push.Biz {
	case domain.BizInk:
		content = domain.FeedInk{}
		if err := json.Unmarshal([]byte(push.Content), &content); err != nil {
			repo.l.Error("json unmarshal feed content failed", logx.Any("feed", push), logx.Error(err))
		}
	}
	return domain.Feed{
		Id:        push.Id,
		UserId:    push.UserId,
		Biz:       push.Biz,
		CreatedAt: push.CreatedAt,
		UpdatedAt: push.UpdatedAt,
		BizId:     push.BizId,
		Content:   content,
	}
}

func (repo *NoCacheFeedRepo) domainToPullFeed(feed domain.Feed) dao.PullFeed {
	bs, err := json.Marshal(feed.Content)
	if err != nil {
		repo.l.Error("json marshal feed content failed", logx.Any("feed", feed), logx.Error(err))
		bs = []byte{}
	}
	return dao.PullFeed{
		Id:        feed.Id,
		UserId:    feed.UserId,
		Biz:       feed.Biz,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		BizId:     feed.BizId,
		Content:   string(bs),
		Status:    dao.FeedStatusNormal,
	}
}

func (repo *NoCacheFeedRepo) domainToPushFeed(feed domain.Feed) dao.PushFeed {
	bs, err := json.Marshal(feed.Content)
	if err != nil {
		repo.l.Error("json marshal feed content failed", logx.Any("feed", feed), logx.Error(err))
		bs = []byte{}
	}
	return dao.PushFeed{
		Id:        feed.Id,
		UserId:    feed.UserId,
		Biz:       feed.Biz,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		BizId:     feed.BizId,
		Content:   string(bs),
		Status:    dao.FeedStatusNormal,
	}
}
