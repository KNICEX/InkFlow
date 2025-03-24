package dao

import (
	"context"
	"time"
)

// Notification
// 取消点赞后，又重新点赞
// 删除评论后，又重新评论...
type Notification struct {
	Id               int64
	RecipientId      int64
	SenderId         int64
	NotificationType string
	SubjectType      string
	SubjectId        int64
	Content          string
	CreatedAt        time.Time `gorm:"index:created_read"`
	Read             bool      `gorm:"index:created_read"`
}

type NotificationDAO interface {
	Insert(ctx context.Context, no Notification) error
	BatchInsert(ctx context.Context, notifications []Notification) error
	FindByType(ctx context.Context, uid int64, notificationType string, maxId int64, limit int) ([]Notification, error)
	ReadAll(ctx context.Context, userId int64, notificationType ...string) error
	CountTotalUnread(ctx context.Context, userId int64) (int64, error)
	CountUnreadByType(ctx context.Context, userId int64, notificationType ...string) (map[string]int64, error)
}
