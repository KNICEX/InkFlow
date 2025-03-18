package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/notification/internal/domain"
)

type NotificationService interface {
	SendNotification(ctx context.Context, n domain.Notification) error
	DeleteNotification(ctx context.Context, id int64, recipientId int64) error
	ListNotification(ctx context.Context, recipientId int64, maxId int64, limit int) ([]domain.Notification, error)
	Read(ctx context.Context, id int64, recipientId int64) error
	ReadAll(ctx context.Context, recipientId int64) error
}
