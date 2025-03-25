package inkpub

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/notification"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/review"
	"github.com/KNICEX/InkFlow/internal/search"
)

type Activities struct {
	inkSvc           ink.Service
	intrSvc          interactive.Service
	reviewSvc        review.AsyncService
	searchSyncSvc    search.SyncService
	notificationSvc  notification.Service
	recommendSyncSvc recommend.SyncService
}

func NewActivities(
	inkSvc ink.Service,
	intrSvc interactive.Service,
	reviewSvc review.AsyncService,
	searchSyncSvc search.SyncService,
	recommendSyncSvc recommend.SyncService,
	notificationSvc notification.Service,
) *Activities {
	return &Activities{
		inkSvc:           inkSvc,
		intrSvc:          intrSvc,
		reviewSvc:        reviewSvc,
		searchSyncSvc:    searchSyncSvc,
		recommendSyncSvc: recommendSyncSvc,
		notificationSvc:  notificationSvc,
	}
}
func (a *Activities) FindInkInfo(ctx context.Context, inkId int64) (ink.Ink, error) {
	return a.inkSvc.FindById(ctx, inkId)
}

func (a *Activities) SubmitReview(ctx context.Context, ink review.Ink) error {
	return a.reviewSvc.SubmitInk(ctx, ink)
}

func (a *Activities) CreateIntr(ctx context.Context, biz string, bizId int64) error {
	return a.intrSvc.CreateInteractive(ctx, biz, bizId)
}

func (a *Activities) UpdateToPublished(ctx context.Context, inkId, uid int64) error {
	return a.inkSvc.UpdateInkStatus(ctx, inkId, uid, ink.StatusPublished)
}

func (a *Activities) UpdateInkToRejected(ctx context.Context, inkId int64) error {
	return a.inkSvc.UpdateInkStatus(ctx, inkId, 0, ink.StatusRejected)
}

func (a *Activities) SyncToSearch(ctx context.Context, ink search.Ink) error {
	return a.searchSyncSvc.InputInk(ctx, []search.Ink{ink})
}

func (a *Activities) SyncToRecommend(ctx context.Context, ink recommend.Ink) error {
	return a.recommendSyncSvc.InputInk(ctx, ink)
}

func (a *Activities) NotifyRejected(ctx context.Context, inkId int64, uid int64, reason string) error {
	// TODO 构建系统通知
	return nil
}
