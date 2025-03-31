package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/domain"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/service"
)

type Handler interface {
	Topic() string
	HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type UserCreateHandler struct {
	svc service.SyncService
}

func NewUserCreateHandler(svc service.SyncService) Handler {
	return &UserCreateHandler{
		svc: svc,
	}
}

func (h *UserCreateHandler) Topic() string {
	return topicUserCreate
}

func (h *UserCreateHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt UserCreateEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	return h.svc.InputUser(context.Background(), domain.User{
		Id:        evt.UserId,
		Account:   evt.Account,
		CreatedAt: evt.CreatedAt,
	})
}
