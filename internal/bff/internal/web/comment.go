package web

import (
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type CommentHandler struct {
	svc     comment.Service
	userSvc user.Service
	auth    middleware.Authentication

	l logx.Logger
}

func NewCommentHandler(svc comment.Service, userSvc user.Service, auth middleware.Authentication, l logx.Logger) *CommentHandler {
	return &CommentHandler{
		svc:     svc,
		userSvc: userSvc,
		auth:    auth,
		l:       l,
	}
}

func (h *CommentHandler) RegisterRoutes(server *gin.RouterGroup) {
	commentGroup := server.Group("/comment")
	{
		commentGroup.GET("", ginx.WrapBody(h.l, h.List))
		commentGroup.POST("", h.auth.CheckLogin(), ginx.WrapBody(h.l, h.Reply))
	}
}

func (h *CommentHandler) Reply(ctx *gin.Context, req PostReplyReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)

	id, err := h.svc.Create(ctx, comment.Comment{
		Biz:   req.Biz,
		BizId: req.BizId,
		Commentator: comment.Commentator{
			Id: uc.UserId,
		},
		Payload: comment.Payload{
			Content: req.Payload.Content,
			Images:  req.Payload.Images,
		},
		Parent: &comment.Comment{
			Id: req.ParentId,
		},
		Root: &comment.Comment{
			Id: req.RootId,
		},
	})
	if err != nil {
		return ginx.InternalError(), err
	}
	return ginx.SuccessWithData(id), nil
}

func (h *CommentHandler) List(ctx *gin.Context, req BizCommentReq) (ginx.Result, error) {
	uc, _ := jwt.GetUserClaims(ctx)
	coms, err := h.svc.LoadLastedList(ctx, req.Biz, req.BizId, uc.UserId, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}

	uids := lo.UniqMap(coms, func(item comment.Comment, index int) int64 {
		return item.Commentator.Id
	})
	users, err := h.userSvc.FindByIds(ctx, uids)
	if err != nil {
		return ginx.InternalError(), err
	}

	res := make([]CommentVO, 0, len(coms))
	for _, com := range coms {
		vo := commentToVO(com)
		vo.Commentator = userToVO(users[com.Commentator.Id])
		res = append(res, vo)
	}
	return ginx.SuccessWithData(res), nil
}

func (h *CommentHandler) LoadMoreChild(ctx *gin.Context, req ChildCommentReq) (ginx.Result, error) {
	uc, _ := jwt.GetUserClaims(ctx)
	coms, err := h.svc.LoadMoreRepliesByRid(ctx, req.RootId, uc.UserId, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}

	uids := lo.UniqMap(coms, func(item comment.Comment, index int) int64 {
		return item.Commentator.Id
	})
	users, err := h.userSvc.FindByIds(ctx, uids)
	if err != nil {
		return ginx.InternalError(), err
	}

	res := make([]CommentVO, 0, len(coms))
	for _, com := range coms {
		vo := commentToVO(com)
		vo.Commentator = userToVO(users[com.Commentator.Id])
		res = append(res, vo)
	}
	return ginx.SuccessWithData(res), nil
}
