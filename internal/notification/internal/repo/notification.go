package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/notification/internal/domain"
)

type NotificationRepo interface {
	CreateNotification(ctx context.Context, n domain.Notification) error
	DelNotification(ctx context.Context, id int64, recipientId int64) error
	FindByRecipientId(ctx context.Context, recipientId int64, maxId int64, limit int) ([]domain.Notification, error)
	FindByIds(ctx context.Context, ids []int64) (map[int64]domain.Notification, error)
	CountNoRead(ctx context.Context, recipientId int64) (int64, error)
	MarkRead(ctx context.Context, recipientId int64, ids []int64) error
	MarkAllRead(ctx context.Context, recipientId int64) error
	DelAll(ctx context.Context, recipientId int64) error
}
