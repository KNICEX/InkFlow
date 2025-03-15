package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/domain"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/events"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	View(ctx context.Context, biz string, bizId int64, uid int64) error
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	ListView(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.ViewRecord, error)
	ListLike(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.LikeRecord, error)
	Collect(ctx context.Context, biz string, bizId, cid, uid int64)
	Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error)
	GetMulti(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	repo     repo.InteractiveRepo
	producer events.InteractiveProducer
	l        logx.Logger
}

func NewInteractiveService(repo repo.InteractiveRepo, producer events.InteractiveProducer, l logx.Logger) InteractiveService {
	return &interactiveService{
		repo:     repo,
		producer: producer,
		l:        l,
	}
}

func (svc *interactiveService) View(ctx context.Context, biz string, bizId int64, uid int64) error {
	if biz == domain.BizInk {
		return svc.producer.ProduceInkView(ctx, events.InkViewEvent{
			InkId:  bizId,
			UserId: uid,
		})
	}
	return svc.repo.IncrView(ctx, biz, bizId, uid)
}

func (svc *interactiveService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	if err := svc.repo.IncrLike(ctx, biz, bizId, uid); err != nil {
		return err
	}

	if biz == domain.BizInk {
		if err := svc.producer.ProduceInkLike(ctx, events.InkLikeEvent{
			InkId:  bizId,
			UserId: uid,
		}); err != nil {
			svc.l.WithCtx(ctx).Error("produce ink like event error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId),
				logx.Int64("uid", uid))
		}
	}
	return nil
}

func (svc *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	if err := svc.repo.DecrLike(ctx, biz, bizId, uid); err != nil {
		return err
	}

	if biz == domain.BizInk {
		if err := svc.producer.ProduceInkCancelLike(ctx, events.InkCancelLikeEvent{
			InkId:  bizId,
			UserId: uid,
		}); err != nil {
			svc.l.WithCtx(ctx).Error("produce ink cancel like event error", logx.Error(err),
				logx.String("biz", biz),
				logx.Int64("bizId", bizId),
				logx.Int64("uid", uid))
		}
	}
	return nil
}
func (svc *interactiveService) ListView(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.ViewRecord, error) {
	return svc.repo.ListView(ctx, biz, uid, maxId, limit)
}

func (svc *interactiveService) ListLike(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.LikeRecord, error) {
	return svc.repo.ListLike(ctx, biz, uid, maxId, limit)
}
func (svc *interactiveService) Collect(ctx context.Context, biz string, bizId, cid, uid int64) {
	//TODO implement me
	panic("implement me")
}

func (svc *interactiveService) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	eg := errgroup.Group{}
	var intr domain.Interactive
	eg.Go(func() error {
		var er error
		intr, er = svc.repo.Get(ctx, biz, bizId)
		return er
	})
	var liked bool
	if uid != 0 {
		// 聚合是否点赞
		// TODO 后续还有收藏或更多其他操作
		eg.Go(func() error {
			var er error
			liked, er = svc.repo.Liked(ctx, biz, bizId, uid)
			return er
		})
	}
	if err := eg.Wait(); err != nil {
		return domain.Interactive{}, err
	}
	intr.Liked = liked
	return intr, nil
}

func (svc *interactiveService) GetMulti(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]domain.Interactive, error) {
	// TODO 还需要处理uid是否点赞等
	return svc.repo.GetMulti(ctx, biz, bizIds)
}
