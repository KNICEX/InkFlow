package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/search/internal/domain"
	"github.com/KNICEX/InkFlow/internal/search/internal/service"
)

type Handler interface {
	Topic() string
	HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type ReplyHandler struct {
	svc service.SyncService
}

func NewReplyHandler(svc service.SyncService) Handler {
	return &ReplyHandler{
		svc: svc,
	}
}

func (h *ReplyHandler) Topic() string {
	return topicCommentReply
}

func (h *ReplyHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt ReplyEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	return h.svc.InputComment(context.Background(), []domain.Comment{
		{
			Id: evt.CommentId,
			Commentator: domain.User{
				Id: evt.CommentatorId,
			},
			Biz:       evt.Biz,
			BizId:     evt.BizId,
			RootId:    evt.RootId,
			ParentId:  evt.ParentId,
			Content:   evt.Payload.Content,
			Images:    evt.Payload.Images,
			CreatedAt: evt.CreatedAt,
		},
	})
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
	return h.svc.InputUser(context.Background(), []domain.User{
		{
			Id:        evt.UserId,
			Account:   evt.Account,
			Avatar:    evt.Avatar,
			Username:  evt.Username,
			CreatedAt: evt.CreatedAt,
		},
	})
}

type UserUpdateHandler struct {
	svc service.SyncService
}

func NewUserUpdateHandler(svc service.SyncService) Handler {
	return &UserUpdateHandler{
		svc: svc,
	}
}

func (h *UserUpdateHandler) Topic() string {
	return topicUserUpdate
}

func (h *UserUpdateHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt UserUpdateEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	return h.svc.InputUser(context.Background(), []domain.User{
		{
			Id:        evt.UserId,
			Account:   evt.Account,
			Avatar:    evt.Avatar,
			AboutMe:   evt.AboutMe,
			Username:  evt.Username,
			CreatedAt: evt.CreatedAt,
			UpdatedAt: evt.UpdatedAt,
		},
	})
}
