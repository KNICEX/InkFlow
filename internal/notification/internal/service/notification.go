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
	ReadAll(ctx context.Context, recipientId int64, types ...domain.NotificationType) error
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

func (svc *notificationService) ReadAll(ctx context.Context, recipientId int64, types ...domain.NotificationType) error {
	return svc.repo.MarkAllRead(ctx, recipientId, types...)
}
