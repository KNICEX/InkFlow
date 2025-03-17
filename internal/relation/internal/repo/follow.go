package repo

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/relation/internal/domain"
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type FollowRepo interface {
	AddFollowRelation(ctx context.Context, c domain.FollowRelation) error
	RemoveFollowRelation(ctx context.Context, c domain.FollowRelation) error
	FindFollowList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error)
	FindFlowerList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error)
	GetFollowStatistic(ctx context.Context, uid, viewUid int64) (domain.FollowStatistic, error)
}

type CachedFollowRepo struct {
	dao   dao.FollowRelationDAO
	cache cache.FollowCache
	l     logx.Logger
}

func NewCachedFollowRepo(dao dao.FollowRelationDAO, cache cache.FollowCache, l logx.Logger) FollowRepo {
	return &CachedFollowRepo{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (repo *CachedFollowRepo) AddFollowRelation(ctx context.Context, c domain.FollowRelation) error {
	err := repo.dao.CreateFollowRelation(ctx, repo.toEntity(c))
	if err != nil {
		return err
	}
	err = repo.cache.Follow(ctx, c.Follower, c.Followee)
	if err != nil {
		repo.l.WithCtx(ctx).Error("add follow cache error", logx.Error(err), logx.Int64("UserId", c.Follower))
	}
	return err
}

func (repo *CachedFollowRepo) RemoveFollowRelation(ctx context.Context, c domain.FollowRelation) error {
	err := repo.dao.CancelFollow(ctx, repo.toEntity(c))
	if err != nil {
		return err
	}
	err = repo.cache.CancelFollow(ctx, c.Follower, c.Followee)
	if err != nil {
		repo.l.WithCtx(ctx).Error("cancel follow cache error", logx.Error(err), logx.Int64("UserId", c.Follower))
	}
	return err
}

func (repo *CachedFollowRepo) FindFollowList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error) {
	res, err := repo.dao.FollowList(ctx, uid, maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(res, func(item dao.UserFollow, index int) domain.FollowInfo {
		return domain.FollowInfo{
			Uid: item.FolloweeId,
		}
	}), nil
}

func (repo *CachedFollowRepo) FindFlowerList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error) {
	res, err := repo.dao.FollowerList(ctx, uid, maxId, limit)
	if err != nil {
		return nil, err
	}
	return lo.Map(res, func(item dao.UserFollow, index int) domain.FollowInfo {
		return domain.FollowInfo{
			Uid: item.FollowerId,
		}
	}), nil
}

func (repo *CachedFollowRepo) GetFollowStatistic(ctx context.Context, uid, viewUid int64) (domain.FollowStatistic, error) {
	res, err := repo.cache.GetStatisticInfo(ctx, uid)
	if err == nil {
		if uid == viewUid {
			// 如果是查看自己的关注数，直接返回
			return res, nil
		}
		followed, er := repo.dao.Followed(ctx, viewUid, uid)
		if er != nil {
			repo.l.Error("get follow statistic cache error", logx.Error(er), logx.Int64("UserId", uid))
			return res, er
		}
		res.Followed = followed
		return res, nil
	}

	if !errors.Is(err, cache.ErrKeyNotFound) {
		repo.l.Error("get follow statistic cache error", logx.Error(err), logx.Int64("UserId", uid))
	}

	var followerCount, followeeCount int64
	var followed bool
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		followerCount, er = repo.dao.CntFollower(ctx, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		followeeCount, er = repo.dao.CntFollowee(ctx, uid)
		return er
	})
	if viewUid != uid {
		eg.Go(func() error {
			var er error
			followed, er = repo.dao.Followed(ctx, viewUid, uid)
			return er
		})
	}
	if err = eg.Wait(); err != nil {
		return domain.FollowStatistic{}, err
	}

	go func() {
		if er := repo.cache.SetStatisticInfo(ctx, uid, domain.FollowStatistic{
			Followers: followerCount,
			Following: followeeCount,
		}); er != nil {
			repo.l.Error("set follow statistic cache error", logx.Error(err), logx.Int64("UserId", uid))
		}
	}()
	return domain.FollowStatistic{
		Followers: followerCount,
		Following: followeeCount,
		Followed:  followed,
	}, nil
}

func (repo *CachedFollowRepo) toDomain(follow dao.UserFollow) domain.FollowRelation {
	return domain.FollowRelation{
		Follower:  follow.FollowerId,
		Followee:  follow.FolloweeId,
		CreatedAt: follow.CreatedAt,
	}
}
func (repo *CachedFollowRepo) toEntity(follow domain.FollowRelation) dao.UserFollow {
	return dao.UserFollow{
		FollowerId: follow.Follower,
		FolloweeId: follow.Followee,
		CreatedAt:  follow.CreatedAt,
	}
}
