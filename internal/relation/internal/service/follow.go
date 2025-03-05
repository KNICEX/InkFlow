package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/relation/internal/domain"
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo"
)

// FollowService TODO 增加关注分组
type FollowService interface {
	// Follow TODO 考虑动态返回关注数
	Follow(ctx context.Context, uid, followeeId int64) error
	CancelFollow(ctx context.Context, uid, followeeId int64) error
	// FollowList TODO 查看他人关注列表时，考虑查询我是否也关注
	FollowList(ctx context.Context, uid int64, maxId int64, limit int) ([]domain.FollowRelation, error)
	FollowStatistic(ctx context.Context, uid int64) (domain.FollowStatistic, error)
}

type followService struct {
	repo repo.FollowRepo
}

func (svc *followService) Follow(ctx context.Context, uid, followeeId int64) error {
	return svc.repo.AddFollowRelation(ctx, domain.FollowRelation{
		Follower: uid,
		Followee: followeeId,
	})
}

func (svc *followService) CancelFollow(ctx context.Context, uid, followeeId int64) error {
	return svc.repo.RemoveFollowRelation(ctx, domain.FollowRelation{
		Follower: uid,
		Followee: followeeId,
	})
}

func (svc *followService) FollowList(ctx context.Context, uid int64, maxId int64, limit int) ([]domain.FollowRelation, error) {
	return svc.repo.FindFollowList(ctx, uid, maxId, limit)
}

func (svc *followService) FollowStatistic(ctx context.Context, uid int64) (domain.FollowStatistic, error) {
	return svc.repo.GetFollowStatistic(ctx, uid)
}
