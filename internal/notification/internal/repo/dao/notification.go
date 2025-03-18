package dao

import "time"

// Notification
// 可以
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
