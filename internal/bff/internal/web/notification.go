package web

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/notification"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"sync"
)

type NotificationHandler struct {
	svc        notification.Service
	userSvc    user.Service
	inkSvc     ink.Service
	commentSvc comment.Service
	l          logx.Logger
}

func (handler *NotificationHandler) RegisterRoutes(server *gin.RouterGroup) {
	notificationGroup := server.Group("/notification")
	{
		notificationGroup.GET("/like", ginx.WrapBody(handler.l, handler.ListMergedLike))
		notificationGroup.GET("/reply", ginx.WrapBody(handler.l, handler.ListReply))
	}
}

func (handler *NotificationHandler) ListMergedLike(ctx *gin.Context, req PagedReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	likes, err := handler.svc.ListMergedLike(ctx, uc.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	uids := make(map[int64]struct{})
	for _, like := range likes {
		for _, uid := range like.UserIds {
			if _, ok := uids[uid]; !ok {
				uids[uid] = struct{}{}
			}
		}
	}

	subjectIdMap := make(map[notification.SubjectType][]int64)
	for _, like := range likes {
		subjectIdMap[like.SubjectType] = append(subjectIdMap[like.SubjectType], like.SubjectId)
	}

	var users map[int64]user.User
	var subjects map[notification.SubjectType]map[int64]any
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		users, er = handler.userSvc.FindByIds(ctx, lo.Keys(uids))
		return er
	})
	eg.Go(func() error {
		var er error
		subjects, er = handler.findSubjectsVO(ctx, subjectIdMap)
		return er
	})
	if err = eg.Wait(); err != nil {
		return ginx.InternalError(), err
	}

	likesVO := make([]MergedLikeVO, 0, len(likes))
	for _, like := range likes {
		vo := mergedLikeToVO(like)
		for _, uid := range like.UserIds {
			if u, ok := users[uid]; ok {
				vo.Users = append(vo.Users, userToVO(u))
			}
		}
		vo.Subject = subjects[like.SubjectType][like.SubjectId]
	}
	return ginx.SuccessWithData(likesVO), nil
}

func (handler *NotificationHandler) findSubjectsVO(ctx context.Context, subjectIds map[notification.SubjectType][]int64) (map[notification.SubjectType]map[int64]any, error) {
	subjectVOMap := make(map[notification.SubjectType]map[int64]any)
	mu := sync.Mutex{}

	eg := errgroup.Group{}
	for subjectType, ids := range subjectIds {
		switch subjectType {
		case notification.SubjectTypeInk:
			eg.Go(func() error {
				inkMap, err := handler.inkSvc.FindByIds(ctx, ids)
				if err != nil {
					return err
				}
				mu.Lock()
				defer mu.Unlock()
				subjectVOMap[subjectType] = lo.MapEntries(inkMap, func(key int64, value ink.Ink) (int64, any) {
					return key, inkToVO(value)
				})
				return nil
			})
		case notification.SubjectTypeComment:
			eg.Go(func() error {
				commentMap, err := handler.commentSvc.FindByIds(ctx, ids, 0)
				if err != nil {
					return err
				}
				mu.Lock()
				defer mu.Unlock()
				subjectVOMap[subjectType] = lo.MapEntries(commentMap, func(key int64, value comment.Comment) (int64, any) {
					return key, commentToVO(value)
				})
				return nil
			})
		default:
			handler.l.Error("unknown like subject type",
				logx.String("type", subjectType.ToString()))
		}
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return subjectVOMap, nil
}

func (handler *NotificationHandler) ListReply(ctx *gin.Context, req MaxIdPagedReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	notifications, err := handler.svc.ListNotification(ctx, uc.UserId, []notification.Type{notification.TypeReply}, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}

	subjectIds := make(map[notification.SubjectType][]int64)
	for _, no := range notifications {
		subjectIds[no.SubjectType] = append(subjectIds[no.SubjectType], no.SubjectId)
	}

	uids := lo.UniqMap(notifications, func(item notification.Notification, index int) int64 {
		return item.SenderId
	})
	var users map[int64]user.User
	var subjects map[notification.SubjectType]map[int64]any
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		users, er = handler.userSvc.FindByIds(ctx, uids)
		return er
	})
	eg.Go(func() error {
		var er error
		subjects, er = handler.findSubjectsVO(ctx, subjectIds)
		return er
	})
	if err = eg.Wait(); err != nil {
		return ginx.InternalError(), err
	}

	notificationsVO := make([]NotificationVO, 0, len(notifications))
	for _, no := range notifications {
		vo := notificationToVO(no)
		if u, ok := users[no.SenderId]; ok {
			userVO := userToVO(u)
			vo.User = &userVO
		}
		vo.Subject = subjects[no.SubjectType][no.SubjectId]
		notificationsVO = append(notificationsVO, vo)
	}

	return ginx.SuccessWithData(notificationsVO), nil
}
