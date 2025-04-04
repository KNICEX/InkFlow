package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/notification/internal/domain"
	"github.com/KNICEX/InkFlow/internal/notification/internal/service"
)

const (
	bizInk = "ink"
)

var (
	ErrUnknownBiz = errors.New("unknown biz")
)

type Handler interface {
	Topic() string
	HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type FollowHandler struct {
	svc service.NotificationService
}

func NewFollowHandler(svc service.NotificationService) Handler {
	return &FollowHandler{
		svc: svc,
	}
}

func (f *FollowHandler) Topic() string {
	return topicFollow
}

func (f *FollowHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt FollowEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	return f.svc.SendNotification(ctx, domain.Notification{
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

type ReplyHandler struct {
	svc        service.NotificationService
	commentSvc comment.Service
	inkSvc     ink.Service
}

func NewReplyHandler(svc service.NotificationService, commentSvc comment.Service, inkSvc ink.Service) Handler {
	return &ReplyHandler{
		svc:        svc,
		commentSvc: commentSvc,
		inkSvc:     inkSvc,
	}
}

func (r *ReplyHandler) Topic() string {
	return topicCommentReply
}

func (r *ReplyHandler) bizToSubjectType(biz string) (domain.SubjectType, error) {
	switch biz {
	case bizInk:
		return domain.SubjectTypeInk, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownBiz, biz)
	}
}

func (r *ReplyHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt ReplyEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}
	if evt.ParentId > 0 {
		// 回复的是评论
		parent, err := r.commentSvc.FindById(ctx, evt.ParentId, 0)
		if err != nil {
			return err
		}

		if parent.Commentator.Id == evt.CommentatorId {
			// 自己评论自己的
			return nil
		}

		subjectType, err := r.bizToSubjectType(evt.Biz)
		if err != nil {
			return err
		}

		return r.svc.SendNotification(ctx, domain.Notification{
			RecipientId:      parent.Commentator.Id,
			SenderId:         evt.CommentatorId,
			NotificationType: domain.NotificationTypeReply,
			SubjectType:      subjectType,
			SubjectId:        evt.BizId,
			Content: domain.ReplyContent{
				CommentId: evt.CommentId,
				SourceContent: domain.ReplyPayload{
					Content: evt.Payload.Content,
					Images:  evt.Payload.Images,
				},
				TargetContent: domain.ReplyPayload{
					Content: parent.Payload.Content,
					Images:  parent.Payload.Images,
				},
			},
			Read:      false,
			CreatedAt: evt.CreatedAt,
		})
	} else {
		// 一级回复
		inkInfo, err := r.inkSvc.FindLiveInk(ctx, evt.BizId)
		if err != nil {
			return err
		}

		if inkInfo.Author.Id == evt.CommentatorId {
			// 自己评论自己的
			return nil
		}

		subjectType, err := r.bizToSubjectType(evt.Biz)
		if err != nil {
			return err
		}

		return r.svc.SendNotification(ctx, domain.Notification{
			RecipientId:      inkInfo.Author.Id,
			SenderId:         evt.CommentatorId,
			NotificationType: domain.NotificationTypeReply,
			SubjectType:      subjectType,
			SubjectId:        evt.BizId,
			Content: domain.ReplyContent{
				CommentId: evt.CommentId,
				SourceContent: domain.ReplyPayload{
					Content: evt.Payload.Content,
					Images:  evt.Payload.Images,
				},
			},
			Read:      false,
			CreatedAt: evt.CreatedAt,
		})
	}
}

type CommentLikeHandler struct {
	svc        service.NotificationService
	commentSvc comment.Service
}

func NewCommentLikeHandler(svc service.NotificationService, commentSvc comment.Service) Handler {
	return &CommentLikeHandler{
		commentSvc: commentSvc,
		svc:        svc,
	}
}

func (c *CommentLikeHandler) Topic() string {
	return topicCommentLike
}

func (c *CommentLikeHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt CommentLikeEvt
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}

	comm, err := c.commentSvc.FindById(ctx, evt.CommentId, 0)
	if err != nil {
		return err
	}

	if evt.LikeUid == comm.Commentator.Id {
		// 自己点赞
		return nil
	}

	return c.svc.SendNotification(ctx, domain.Notification{
		RecipientId:      comm.Commentator.Id,
		SenderId:         evt.LikeUid,
		NotificationType: domain.NotificationTypeLike,
		SubjectType:      domain.SubjectTypeComment,
		SubjectId:        comm.Id,
		Content:          nil,
		Read:             false,
		CreatedAt:        evt.CreatedAt,
	})
}

type InkLikeHandler struct {
	svc    service.NotificationService
	inkSvc ink.Service
}

func NewInkLikeHandler(svc service.NotificationService, inkSvc ink.Service) Handler {
	return &InkLikeHandler{
		svc:    svc,
		inkSvc: inkSvc,
	}
}

func (i *InkLikeHandler) Topic() string {
	return topicInkLike
}

func (i *InkLikeHandler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var evt InkLikeEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}

	inkInfo, err := i.inkSvc.FindLiveInk(ctx, evt.InkId)
	if err != nil {
		return err
	}
	if evt.UserId == inkInfo.Author.Id {
		// 自己点赞
		return nil
	}

	return i.svc.SendNotification(ctx, domain.Notification{
		RecipientId:      inkInfo.Author.Id,
		SenderId:         evt.UserId,
		NotificationType: domain.NotificationTypeLike,
		SubjectType:      domain.SubjectTypeInk,
		SubjectId:        inkInfo.Id,
		Content:          nil,
		Read:             false,
		CreatedAt:        evt.CreatedAt,
	})
}
