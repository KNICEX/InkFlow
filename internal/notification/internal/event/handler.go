package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/notification/internal/domain"
	"github.com/KNICEX/InkFlow/internal/notification/internal/service"
)

type Handler interface {
	HandleMessage(msg *sarama.ConsumerMessage) error
}

type FollowHandler struct {
	svc service.NotificationService
}

func (f FollowHandler) HandleMessage(msg *sarama.ConsumerMessage) error {
	var evt FollowEvt
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	return f.svc.SendNotification(context.Background(), domain.Notification{
		RecipientId:      evt.FolloweeId,
		SenderId:         evt.FollowerId,
		NotificationType: domain.NotificationTypeFollow,
		SubjectType:      domain.SubjectTypeUser,
		SubjectId:        evt.FollowerId,
		Content:          nil,
		Read:             false,
		CreatedAt:        evt.CreatedAt,
	})
}
