package web

import (
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/search"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type SearchHandler struct {
	auth      middleware.Authentication
	svc       search.Service
	followSvc relation.FollowService
	l         logx.Logger
}

func NewSearchHandler(auth middleware.Authentication, svc search.Service, followSvc relation.FollowService, l logx.Logger) *SearchHandler {
	return &SearchHandler{
		auth:      auth,
		svc:       svc,
		followSvc: followSvc,
		l:         l,
	}
}

func (h *SearchHandler) RegisterRoutes(server *gin.RouterGroup) {
	searchGroup := server.Group("/search")
	{
		searchGroup.GET("/user", ginx.WrapBody(h.l, h.SearchUser))
		searchGroup.GET("/ink", ginx.WrapBody(h.l, h.SearchInK))
		searchGroup.GET("/comment", ginx.WrapBody(h.l, h.SearchComment))
	}
}

func (h *SearchHandler) SearchUser(ctx *gin.Context, req SearchReq) (ginx.Result, error) {
	u, _ := jwt.GetUserClaims(ctx)
	users, err := h.svc.SearchUser(ctx, req.Keyword, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	uids := lo.Map(users, func(item search.User, index int) int64 {
		return item.Id
	})

	followInfos, err := h.followSvc.FindFollowStatsBatch(ctx, uids, u.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]UserVO, 0, len(users))
	for _, user := range users {
		var followStats relation.FollowStatistic
		if followInfo, ok := followInfos[user.Id]; ok {
			followStats = followInfo
		} else {
			followStats = relation.FollowStatistic{}
		}
		res = append(res, UserVO{
			Id:        user.Id,
			Username:  user.Username,
			Account:   user.Account,
			Avatar:    user.Avatar,
			AboutMe:   user.AboutMe,
			CreatedAt: user.CreatedAt,
			Following: followStats.Following,
			Followers: followStats.Followers,
			Followed:  followStats.Followed,
		})
	}
	return ginx.SuccessWithData(res), nil
}

func (h *SearchHandler) SearchInK(ctx *gin.Context, req SearchReq) (ginx.Result, error) {
	inks, err := h.svc.SearchInk(ctx, req.Keyword, req.Offset, req.Limit)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	return ginx.SuccessWithData(lo.Map(inks, func(item search.Ink, index int) InkVO {
		return searchInkToInkVO(item)
	})), nil
}

func (h *SearchHandler) SearchComment(ctx *gin.Context, req SearchReq) (ginx.Result, error) {
	comments, err := h.svc.SearchComment(ctx, req.Keyword, req.Offset, req.Limit)
	if err != nil {
		return ginx.InvalidParam(), err
	}
	return ginx.SuccessWithData(lo.Map(comments, func(item search.Comment, index int) CommentVO {
		return searchCommentToCommentVO(item)
	})), nil
}
