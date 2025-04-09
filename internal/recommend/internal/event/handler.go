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
	return h.svc.InputUser(ctx, domain.User{
		Id:        evt.UserId,
		Account:   evt.Account,
		CreatedAt: evt.CreatedAt,
	})
}

type InkViewHandler struct {
	svc service.SyncService
}

func NewInkViewHandler(svc service.SyncService) Handler {
	return &InkViewHandler{
		svc: svc,
	}
}

func (h *InkViewHandler) Topic() string {
	return topicInkView
}

func (h *InkViewHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt InkViewEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	return h.svc.InputFeedback(ctx, domain.Feedback{
		UserId:       evt.UserId,
		InkId:        evt.InkId,
		FeedbackType: domain.FeedbackTypeView,
		CreatedAt:    evt.CreatedAt,
	})
}

type InkLiKeHandler struct {
	svc service.SyncService
}

func NewInkLikeHandler(svc service.SyncService) Handler {
	return &InkLiKeHandler{
		svc: svc,
	}
}

func (h *InkLiKeHandler) Topic() string {
	return topicInkLike
}

func (h *InkLiKeHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt InkLikeEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	return h.svc.InputFeedback(ctx, domain.Feedback{
		UserId:       evt.UserId,
		InkId:        evt.InkId,
		FeedbackType: domain.FeedbackTypeLike,
		CreatedAt:    evt.CreatedAt,
	})
}

type InkCancelLikeHandler struct {
	svc service.SyncService
}

func NewInkCancelLikeHandler(svc service.SyncService) Handler {
	return &InkCancelLikeHandler{
		svc: svc,
	}
}

func (h *InkCancelLikeHandler) Topic() string {
	return topicInkCancelLike
}

func (h *InkCancelLikeHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt InkCancelLikeEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	return h.svc.DeleteFeedback(ctx, domain.Feedback{
		UserId:       evt.UserId,
		InkId:        evt.InkId,
		FeedbackType: domain.FeedbackTypeLike,
	})
}
