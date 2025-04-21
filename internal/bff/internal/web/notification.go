package web

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/notification"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"strconv"
	"sync"
)

type NotificationHandler struct {
	svc           notification.Service
	inkSvc        ink.Service
	commentSvc    comment.Service
	userAggregate *UserAggregate

	auth middleware.Authentication
	l    logx.Logger
}

func NewNotificationHandler(svc notification.Service, userAggregate *UserAggregate, inkSvc ink.Service,
	commentSvc comment.Service, auth middleware.Authentication, l logx.Logger) *NotificationHandler {
	return &NotificationHandler{
		svc:           svc,
		userAggregate: userAggregate,
		inkSvc:        inkSvc,
		commentSvc:    commentSvc,
		auth:          auth,
		l:             l,
	}
}

func (handler *NotificationHandler) RegisterRoutes(server *gin.RouterGroup) {
	notificationGroup := server.Group("/notification", handler.auth.CheckLogin())
	{
		notificationGroup.GET("/like", ginx.WrapBody(handler.l, handler.ListMergedLike))
		notificationGroup.GET("/mention", ginx.WrapBody(handler.l, handler.ListMention))
		notificationGroup.GET("/reply", ginx.WrapBody(handler.l, handler.ListReply))
		notificationGroup.GET("/follow", ginx.WrapBody(handler.l, handler.ListFollow))
		notificationGroup.GET("/system", ginx.WrapBody(handler.l, handler.ListSystem))
		notificationGroup.GET("/count", ginx.Wrap(handler.l, handler.UnreadCount))

		notificationGroup.POST("/read/:type", ginx.Wrap(handler.l, handler.Read))

		notificationGroup.DELETE("/:id", ginx.Wrap(handler.l, handler.Delete))
		notificationGroup.DELETE("/like", ginx.WrapBody(handler.l, handler.DeleteMergedLike))
	}
}

func (handler *NotificationHandler) readAll(ctx context.Context, userId int64, typ notification.Type) {
	go func() {
		err := handler.svc.ReadAll(ctx, userId, typ)
		if err != nil {
			handler.l.Error("failed to mark notification as read", logx.Int64("userId", userId), logx.String("type", typ.ToString()), logx.Error(err))
		}
	}()
}

func (handler *NotificationHandler) ListMergedLike(ctx *gin.Context, req OffsetPagedReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	likes, err := handler.svc.ListMergedLike(ctx, uc.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(likes) == 0 {
		return ginx.SuccessWithData([]MergedLikeVO{}), nil
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

	var users map[int64]UserVO
	var subjects map[notification.SubjectType]map[int64]any
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		users, er = handler.userAggregate.GetUserList(ctx, lo.Keys(uids), uc.UserId)
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
				vo.Users = append(vo.Users, u)
			}
		}
		vo.Subject = subjects[like.SubjectType][like.SubjectId]
		likesVO = append(likesVO, vo)
	}

	handler.readAll(context.WithoutCancel(ctx), uc.UserId, notification.TypeLike)
	return ginx.SuccessWithData(likesVO), nil
}

func (handler *NotificationHandler) DeleteMergedLike(ctx *gin.Context, req SubjectReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	err := handler.svc.DeleteMergedLike(ctx, uc.UserId, notification.SubjectType(req.SubjectType), req.SubjectId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
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
			return nil, fmt.Errorf("unsupported subject type: %s", subjectType)
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
	if len(notifications) == 0 {
		return ginx.SuccessWithData([]NotificationVO{}), nil
	}

	subjectIds := make(map[notification.SubjectType][]int64)
	for _, no := range notifications {
		subjectIds[no.SubjectType] = append(subjectIds[no.SubjectType], no.SubjectId)
	}

	uids := lo.UniqMap(notifications, func(item notification.Notification, index int) int64 {
		return item.SenderId
	})
	var users map[int64]UserVO
	var subjects map[notification.SubjectType]map[int64]any
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		users, er = handler.userAggregate.GetUserList(ctx, uids, uc.UserId)
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
			vo.User = &u
		}
		vo.Subject = subjects[no.SubjectType][no.SubjectId]
		notificationsVO = append(notificationsVO, vo)
	}

	handler.readAll(context.WithoutCancel(ctx), uc.UserId, notification.TypeReply)
	return ginx.SuccessWithData(notificationsVO), nil
}

func (handler *NotificationHandler) ListFollow(ctx *gin.Context, req MaxIdPagedReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	notifications, err := handler.svc.ListNotification(ctx, uc.UserId, []notification.Type{notification.TypeFollow}, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(notifications) == 0 {
		return ginx.SuccessWithData([]NotificationVO{}), nil
	}

	uids := lo.UniqMap(notifications, func(item notification.Notification, index int) int64 {
		return item.SenderId
	})

	users, err := handler.userAggregate.GetUserList(ctx, uids, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}

	followVOs := make([]NotificationVO, 0, len(notifications))
	for _, no := range notifications {
		vo := notificationToVO(no)
		if u, ok := users[no.SenderId]; ok {
			vo.User = &u
		}
		followVOs = append(followVOs, vo)
	}

	handler.readAll(context.WithoutCancel(ctx), uc.UserId, notification.TypeFollow)
	return ginx.SuccessWithData(followVOs), nil
}

func (handler *NotificationHandler) ListMention(ctx *gin.Context, req MaxIdPagedReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	notifications, err := handler.svc.ListNotification(ctx, uc.UserId, []notification.Type{notification.TypeMention}, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(notifications) == 0 {
		return ginx.SuccessWithData([]NotificationVO{}), nil
	}

	uids := lo.UniqMap(notifications, func(item notification.Notification, index int) int64 {
		return item.SenderId
	})

	subjectIds := handler.subjectIds(notifications)

	var users map[int64]UserVO
	var subjects map[notification.SubjectType]map[int64]any
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		users, er = handler.userAggregate.GetUserList(ctx, uids, uc.UserId)
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

	mentionsVOs := make([]NotificationVO, 0, len(notifications))
	for _, no := range notifications {
		vo := notificationToVO(no)
		if u, ok := users[no.SenderId]; ok {
			vo.User = &u
		}
		vo.Subject = subjects[no.SubjectType][no.SubjectId]
		mentionsVOs = append(mentionsVOs, vo)
	}

	handler.readAll(context.WithoutCancel(ctx), uc.UserId, notification.TypeMention)
	return ginx.SuccessWithData(mentionsVOs), nil
}

func (handler *NotificationHandler) ListSystem(ctx *gin.Context, req MaxIdPagedReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	notifications, err := handler.svc.ListNotification(ctx, uc.UserId, []notification.Type{notification.TypeSystem}, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(notifications) == 0 {
		return ginx.SuccessWithData([]NotificationVO{}), nil
	}

	subjectIds := handler.subjectIds(notifications)

	subjects, err := handler.findSubjectsVO(ctx, subjectIds)
	if err != nil {
		return ginx.InternalError(), err
	}

	systemVOs := make([]NotificationVO, 0, len(notifications))
	for _, no := range notifications {
		vo := notificationToVO(no)
		vo.Subject = subjects[no.SubjectType][no.SubjectId]
		systemVOs = append(systemVOs, vo)
	}

	handler.readAll(context.WithoutCancel(ctx), uc.UserId, notification.TypeSystem)
	return ginx.SuccessWithData(systemVOs), nil
}

func (handler *NotificationHandler) subjectIds(nos []notification.Notification) map[notification.SubjectType][]int64 {
	subjectIds := make(map[notification.SubjectType][]int64)
	for _, no := range nos {
		subjectIds[no.SubjectType] = append(subjectIds[no.SubjectType], no.SubjectId)
	}
	return subjectIds
}

func (handler *NotificationHandler) UnreadCount(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	count, err := handler.svc.UnreadCount(ctx, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(count), nil
}
func (handler *NotificationHandler) Delete(ctx *gin.Context) (ginx.Result, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	uc := jwt.MustGetUserClaims(ctx)
	err = handler.svc.DeleteById(ctx, uc.UserId, id)
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}

func (handler *NotificationHandler) Read(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	tp := ctx.Param("type")
	err := handler.svc.ReadAll(ctx, uc.UserId, notification.Type(tp))
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}
