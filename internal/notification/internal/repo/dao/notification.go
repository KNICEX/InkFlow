package dao

import "time"

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
