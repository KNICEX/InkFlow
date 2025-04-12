package web

import (
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RecommendHandler struct {
	svc           recommend.Service
	inkAggregate  *inkAggregate
	userAggregate *userAggregate
	auth          middleware.Authentication
	l             logx.Logger
}

func NewRecommendHandler(svc recommend.Service, inkSvc ink.Service, userSvc user.Service,
	followSvc relation.FollowService, intrSvc interactive.Service, commentSvc comment.Service,
	auth middleware.Authentication, l logx.Logger) *RecommendHandler {
	return &RecommendHandler{
		svc:           svc,
		inkAggregate:  newInkAggregate(inkSvc, userSvc, followSvc, intrSvc, commentSvc),
		userAggregate: newUserAggregate(userSvc, followSvc),
		auth:          auth,
		l:             l,
	}
}

func (h *RecommendHandler) RegisterRoutes(server *gin.RouterGroup) {
	recommendGroup := server.Group("/recommend", h.auth.ExtractPayload())
	{
		recommendGroup.GET("/ink/similar/:id", ginx.WrapBody(h.l, h.ListSimilarInk))
		recommendGroup.GET("/author", ginx.WrapBody(h.l, h.ListRecommendAuthor))
		recommendGroup.GET("/author/similar/:id", ginx.WrapBody(h.l, h.ListSimilarAuthor))
	}
}

func (h *RecommendHandler) ListSimilarInk(ctx *gin.Context, req OffsetPagedReq) (ginx.Result, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}

	inkIds, err := h.svc.FindSimilarInk(ctx, id, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	uc, _ := jwt.GetUserClaims(ctx)
	inkVoMap, err := h.inkAggregate.GetInkList(ctx, inkIds, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	voList := make([]InkVO, 0, len(inkVoMap))
	for _, inkId := range inkIds {
		if inkVo, ok := inkVoMap[inkId]; ok {
			voList = append(voList, inkVo)
		}
	}
	return ginx.SuccessWithData(voList), nil
}

func (h *RecommendHandler) ListRecommendAuthor(ctx *gin.Context, req OffsetPagedReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	userIds, err := h.svc.FindRecommendAuthor(ctx, uc.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	userVoMap, err := h.userAggregate.GetUserList(ctx, userIds, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	voList := make([]UserVO, 0, len(userVoMap))
	for _, userId := range userIds {
		if userVo, ok := userVoMap[userId]; ok {
			voList = append(voList, userVo)
		}
	}
	return ginx.SuccessWithData(voList), nil
}

func (h *RecommendHandler) ListSimilarAuthor(ctx *gin.Context, req OffsetPagedReq) (ginx.Result, error) {
	authorId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	userIds, err := h.svc.FindSimilarAuthor(ctx, authorId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	uc, _ := jwt.GetUserClaims(ctx)
	userVoMap, err := h.userAggregate.GetUserList(ctx, userIds, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	voList := make([]UserVO, 0, len(userVoMap))
	for _, userId := range userIds {
		if userVo, ok := userVoMap[userId]; ok {
			voList = append(voList, userVo)
		}
	}
	return ginx.SuccessWithData(voList), nil
}
