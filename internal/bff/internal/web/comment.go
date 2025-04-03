package web

import (
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/mapx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"strconv"
)

type CommentHandler struct {
	svc     comment.Service
	userSvc user.Service
	*userAggregate
	auth middleware.Authentication

	l logx.Logger
}

func NewCommentHandler(svc comment.Service, followSvc relation.FollowService, userSvc user.Service, auth middleware.Authentication, l logx.Logger) *CommentHandler {
	return &CommentHandler{
		svc:           svc,
		userSvc:       userSvc,
		userAggregate: newUserAggregate(userSvc, followSvc),
		auth:          auth,
		l:             l,
	}
}

func (h *CommentHandler) RegisterRoutes(server *gin.RouterGroup) {
	commentGroup := server.Group("/comment")
	{
		commentGroup.GET("", h.auth.ExtractPayload(), ginx.WrapBody(h.l, h.List))
		commentGroup.GET("/child/:rid", h.auth.ExtractPayload(), ginx.WrapBody(h.l, h.LoadMoreChild))
		commentGroup.POST("/reply", h.auth.CheckLogin(), ginx.WrapBody(h.l, h.Reply))
		commentGroup.DELETE("/:id", h.auth.CheckLogin(), ginx.Wrap(h.l, h.DelComment))

		commentGroup.POST("/like/:id", h.auth.CheckLogin(), ginx.Wrap(h.l, h.Like))
		commentGroup.DELETE("/like/:id", h.auth.CheckLogin(), ginx.Wrap(h.l, h.CancelLike))
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

func (h *CommentHandler) DelComment(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	if err = h.svc.Delete(ctx, id, uc.UserId); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}

func (h *CommentHandler) List(ctx *gin.Context, req BizCommentReq) (ginx.Result, error) {
	uc, _ := jwt.GetUserClaims(ctx)
	coms, err := h.svc.LoadLastedList(ctx, req.Biz, req.BizId, uc.UserId, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}

	if len(coms) == 0 {
		return ginx.SuccessWithData([]CommentVO{}), nil
	}

	uidSet := mapx.NewSet[int64]()
	for _, com := range coms {
		uidSet.Add(com.Commentator.Id)
		for _, child := range com.Children {
			uidSet.Add(child.Commentator.Id)
		}
	}
	users, err := h.userAggregate.GetUserList(ctx, uidSet.Values(), uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}

	res := make([]CommentVO, 0, len(coms))
	for _, com := range coms {
		vo := commentToVO(com)
		vo.Commentator = users[com.Commentator.Id]
		if len(com.Children) > 0 {
			vo.Children = make([]CommentVO, len(com.Children))
			for i, child := range com.Children {
				vo.Children[i] = commentToVO(child)
				vo.Children[i].Commentator = users[child.Commentator.Id]
			}
		}
		res = append(res, vo)
	}
	return ginx.SuccessWithData(res), nil
}

func (h *CommentHandler) LoadMoreChild(ctx *gin.Context, req ChildCommentReq) (ginx.Result, error) {
	uc, _ := jwt.GetUserClaims(ctx)
	rid, err := strconv.ParseInt(ctx.Param("rid"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	coms, err := h.svc.LoadMoreRepliesByRid(ctx, rid, uc.UserId, req.MaxId, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}

	uids := lo.UniqMap(coms, func(item comment.Comment, index int) int64 {
		return item.Commentator.Id
	})
	users, err := h.userAggregate.GetUserList(ctx, uids, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}

	res := make([]CommentVO, 0, len(coms))
	for _, com := range coms {
		vo := commentToVO(com)
		vo.Commentator = users[com.Commentator.Id]
		res = append(res, vo)
	}
	return ginx.SuccessWithData(res), nil
}

func (h *CommentHandler) Like(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	if err = h.svc.Like(ctx, uc.UserId, id); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}
func (h *CommentHandler) CancelLike(ctx *gin.Context) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	if err = h.svc.CancelLike(ctx, uc.UserId, id); err != nil {
		return ginx.InternalError(), err
	}
	return ginx.Success(), nil
}
