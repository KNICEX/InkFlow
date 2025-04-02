package inkpub

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/feed"
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
	feedSvc          feed.Service
}

func NewActivities(
	inkSvc ink.Service,
	intrSvc interactive.Service,
	reviewSvc review.AsyncService,
	searchSyncSvc search.SyncService,
	recommendSyncSvc recommend.SyncService,
	notificationSvc notification.Service,
	feedSvc feed.Service,
) *Activities {
	return &Activities{
		inkSvc:           inkSvc,
		intrSvc:          intrSvc,
		reviewSvc:        reviewSvc,
		searchSyncSvc:    searchSyncSvc,
		recommendSyncSvc: recommendSyncSvc,
		notificationSvc:  notificationSvc,
		feedSvc:          feedSvc,
	}
}
func (a *Activities) FindInkInfo(ctx context.Context, inkId, uid int64) (ink.Ink, error) {
	return a.inkSvc.FindPendingInk(ctx, inkId, uid)
}

func (a *Activities) SubmitReview(ctx context.Context, ink review.Ink) error {
	return a.reviewSvc.SubmitInk(ctx, ink)
}

func (a *Activities) CreateIntr(ctx context.Context, biz string, bizId int64) error {
	return a.intrSvc.CreateInteractive(ctx, biz, bizId)
}

func (a *Activities) UpdateToPublished(ctx context.Context, inkId, uid int64) error {
	return a.inkSvc.UpdateDraftStatus(ctx, inkId, uid, ink.StatusPublished)
}

func (a *Activities) UpdateInkToRejected(ctx context.Context, inkId int64, uid int64) error {
	return a.inkSvc.UpdateDraftStatus(ctx, inkId, uid, ink.StatusRejected)
}

func (a *Activities) SyncToLive(ctx context.Context, inkInfo ink.Ink) error {
	inkInfo.Status = ink.StatusPublished
	return a.inkSvc.SyncToLive(ctx, inkInfo)
}

func (a *Activities) SyncToSearch(ctx context.Context, ink ink.Ink) error {
	return a.searchSyncSvc.InputInk(ctx, []search.Ink{
		{
			Id:    ink.Id,
			Title: ink.Title,
			Author: search.User{
				Id: ink.Author.Id,
			},
			Content:   ink.ContentHtml,
			Cover:     ink.Cover,
			Tags:      ink.Tags,
			AiTags:    ink.AiTags,
			CreatedAt: ink.CreatedAt,
			UpdatedAt: ink.UpdatedAt,
		},
	})
}

func (a *Activities) SyncToRecommend(ctx context.Context, ink ink.Ink) error {
	tagMap := make(map[string]struct{})
	mergedTags := make([]string, 0)
	for _, tag := range ink.Tags {
		if _, ok := tagMap[tag]; !ok {
			tagMap[tag] = struct{}{}
			mergedTags = append(mergedTags, tag)
		}
	}
	for _, tag := range ink.AiTags {
		if _, ok := tagMap[tag]; !ok {
			tagMap[tag] = struct{}{}
			mergedTags = append(mergedTags, tag)
		}
	}
	recommendInk := recommend.Ink{
		Id:        ink.Id,
		AuthorId:  ink.Author.Id,
		Title:     ink.Title,
		Tags:      mergedTags,
		CreatedAt: ink.CreatedAt,
	}
	return a.recommendSyncSvc.InputInk(ctx, recommendInk)
}

func (a *Activities) SyncToFeed(ctx context.Context, ink ink.Ink) error {
	return a.feedSvc.CreateFeed(ctx, feed.Feed{
		UserId: ink.Author.Id,
		Biz:    bizInk,
		BizId:  ink.Id,
		Content: feed.Ink{
			InkId:     ink.Id,
			Title:     ink.Title,
			AuthorId:  ink.Author.Id,
			Cover:     ink.Cover,
			Abstract:  "",
			CreatedAt: ink.CreatedAt,
		},
	})
}

func (a *Activities) NotifyRejected(ctx context.Context, ink ink.Ink, reason string) error {
	return a.notificationSvc.SendNotification(ctx, notification.Notification{
		RecipientId:      ink.Author.Id,
		NotificationType: notification.TypeSystem,
		SubjectType:      notification.SubjectTypeInk,
		SubjectId:        ink.Id,
		Content:          fmt.Sprintf("您的稿件《%s》未通过审核，原因：%s", ink.Title, reason),
	})
}
