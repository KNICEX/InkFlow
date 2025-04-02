package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/notification/internal/domain"
	"github.com/KNICEX/InkFlow/internal/notification/internal/repo"
)

type NotificationService interface {
	SendNotification(ctx context.Context, n domain.Notification) error
	DeleteByType(ctx context.Context, recipientId int64, types ...domain.NotificationType) error
	ListNotification(ctx context.Context, recipientId int64, types []domain.NotificationType, maxId int64, limit int) ([]domain.Notification, error)
	ListMergedLike(ctx context.Context, recipientId int64, offset, limit int) ([]domain.MergedLikeNotification, error)
	ReadAll(ctx context.Context, recipientId int64, types ...domain.NotificationType) error
	UnreadCount(ctx context.Context, recipientId int64) (map[domain.NotificationType]int64, error)
}

type notificationService struct {
	repo repo.NotificationRepo
}

func NewNotificationService(repo repo.NotificationRepo) NotificationService {
	return &notificationService{
		repo: repo,
	}
}

func (svc *notificationService) SendNotification(ctx context.Context, n domain.Notification) error {
	return svc.repo.CreateNotification(ctx, n)
}

func (svc *notificationService) DeleteByType(ctx context.Context, recipientId int64, types ...domain.NotificationType) error {
	return svc.repo.DeleteByType(ctx, recipientId, types...)
}

func (svc *notificationService) ListNotification(ctx context.Context, recipientId int64, types []domain.NotificationType, maxId int64, limit int) ([]domain.Notification, error) {
	return svc.repo.FindByType(ctx, recipientId, types, maxId, limit)
}

func (svc *notificationService) ListMergedLike(ctx context.Context, recipientId int64, offset, limit int) ([]domain.MergedLikeNotification, error) {
	likes, err := svc.repo.FindMergedLike(ctx, recipientId, offset, limit)
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (svc *notificationService) ReadAll(ctx context.Context, recipientId int64, types ...domain.NotificationType) error {
	return svc.repo.MarkAllRead(ctx, recipientId, types...)
}

func (svc *notificationService) UnreadCount(ctx context.Context, recipientId int64) (map[domain.NotificationType]int64, error) {
	return svc.repo.CountUnreadByType(ctx, recipientId, []domain.NotificationType{
		domain.NotificationTypeLike,
		domain.NotificationTypeReply,
		domain.NotificationTypeFollow,
		domain.NotificationTypeSystem,
		domain.NotificationTypeMention,
		domain.NotificationTypeSubscribe,
	})
}
