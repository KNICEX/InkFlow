package web

import (
	"github.com/KNICEX/InkFlow/internal/notification"
	"time"
)

type NotificationVO struct {
	Id               int64     `json:"id"`
	User             *UserVO   `json:"sender"`
	SubjectType      string    `json:"subject_type"`
	SubjectId        int64     `json:"subject_id"`
	Subject          any       `json:"subject"`
	NotificationType string    `json:"notification_type"`
	Content          any       `json:"content"`
	Read             bool      `json:"read"`
	CreatedAt        time.Time `json:"created_at"`
}

func notificationToVO(n notification.Notification) NotificationVO {
	return NotificationVO{
		Id:               n.Id,
		User:             nil,
		SubjectType:      n.SubjectType.ToString(),
		SubjectId:        n.SubjectId,
		Subject:          nil,
		NotificationType: n.NotificationType.ToString(),
		Content:          n.Content,
		Read:             n.Read,
		CreatedAt:        n.CreatedAt,
	}
}

type MergedLikeVO struct {
	Users       []UserVO  `json:"users"`
	Total       int64     `json:"total"`
	SubjectType string    `json:"subject_type"`
	Subject     any       `json:"subject"`
	Read        bool      `json:"read"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func mergedLikeToVO(ml notification.MergedLike) MergedLikeVO {
	return MergedLikeVO{
		Users:       nil,
		Total:       ml.Total,
		SubjectType: ml.SubjectType.ToString(),
		Read:        ml.Read,
		UpdatedAt:   ml.UpdatedAt,
	}
}

type PagedReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type MaxIdPagedReq struct {
	MaxId int64 `json:"maxId"`
	Limit int   `json:"limit"`
}
