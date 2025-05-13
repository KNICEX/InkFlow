package web

import (
	"time"

	"github.com/KNICEX/InkFlow/internal/notification"
)

type NotificationVO struct {
	Id               int64     `json:"id,string"`
	User             *UserVO   `json:"user"`
	SubjectType      string    `json:"subjectType"`
	SubjectId        int64     `json:"subjectId,string"`
	Subject          any       `json:"subject"`
	NotificationType string    `json:"notificationType"`
	Content          any       `json:"content"`
	Read             bool      `json:"read"`
	CreatedAt        time.Time `json:"createdAt"`
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
	SubjectType string    `json:"subjectType"`
	Subject     any       `json:"subject"`
	Read        bool      `json:"read"`
	UpdatedAt   time.Time `json:"updatedAt"`
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

type SubjectReq struct {
	SubjectType string `json:"subjectType" form:"subjectType"`
	SubjectId   int64  `json:"subjectId,string" form:"subjectId"`
}

type ReadBatchReq struct {
	Ids []int64 `json:"ids,string" form:"ids"`
}
