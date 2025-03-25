package service

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/relation/internal/domain"
	"github.com/KNICEX/InkFlow/internal/relation/internal/event"
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"time"
)

// FollowService TODO 增加关注分组
type FollowService interface {
	// Follow TODO 考虑动态返回关注数
	Follow(ctx context.Context, uid, followeeId int64) error
	CancelFollow(ctx context.Context, uid, followeeId int64) error
	// FollowList TODO 查看他人关注列表时，考虑查询我是否也关注
	FollowList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error)
	FollowerList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error)
	FindFollowStats(ctx context.Context, uid, viewUid int64) (domain.FollowStatistic, error)
	FindFollowStatsBatch(ctx context.Context, uids []int64, viewUid int64) (map[int64]domain.FollowStatistic, error)
}

type followService struct {
	repo     repo.FollowRepo
	producer event.FollowProducer
	l        logx.Logger
}

func NewFollowService(repo repo.FollowRepo, producer event.FollowProducer, l logx.Logger) FollowService {
	return &followService{
		repo: repo,
	}
}

func (svc *followService) Follow(ctx context.Context, uid, followeeId int64) error {
	err := svc.repo.AddFollowRelation(ctx, domain.FollowRelation{
		FollowerId: uid,
		FolloweeId: followeeId,
	})
	if err == nil || errors.Is(err, repo.ErrAlreadyFollowed) {
		er := svc.producer.Produce(ctx, event.FollowEvt{
			FollowerId: uid,
			FolloweeId: followeeId,
			CreatedAt:  time.Now(),
		})
		if er != nil {
			svc.l.Error("produce follow event error", logx.Error(er),
				logx.Int64("followerId", uid),
				logx.Int64("followeeId", followeeId))
		}
	}
	return err
}

func (svc *followService) CancelFollow(ctx context.Context, uid, followeeId int64) error {
	return svc.repo.RemoveFollowRelation(ctx, domain.FollowRelation{
		FollowerId: uid,
		FolloweeId: followeeId,
	})
}

func (svc *followService) FollowList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error) {
	return svc.repo.FindFollowingList(ctx, uid, viewUid, maxId, limit)
}

func (svc *followService) FollowerList(ctx context.Context, uid, viewUid int64, maxId int64, limit int) ([]domain.FollowInfo, error) {
	return svc.repo.FindFlowerList(ctx, uid, viewUid, maxId, limit)
}

func (svc *followService) FindFollowStats(ctx context.Context, uid, viewUid int64) (domain.FollowStatistic, error) {
	return svc.repo.GetFollowStats(ctx, uid, viewUid)
}
func (svc *followService) FindFollowStatsBatch(ctx context.Context, uids []int64, viewUid int64) (map[int64]domain.FollowStatistic, error) {
	return svc.repo.GetFollowStatBatch(ctx, uids, viewUid)
}
