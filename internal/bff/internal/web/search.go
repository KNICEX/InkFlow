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

func (s *SearchHandler) RegisterRoutes(server *gin.RouterGroup) {
	searchGroup := server.Group("/search")
	{
		searchGroup.GET("/user", ginx.WrapBody(s.l, s.SearchUser))
	}
}

func (s *SearchHandler) SearchUser(ctx *gin.Context, req SearchReq) (ginx.Result, error) {
	u, _ := jwt.GetUserClaims(ctx)
	users, err := s.svc.SearchUser(ctx, req.Keyword, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	uids := lo.Map(users, func(item search.User, index int) int64 {
		return item.Id
	})

	followInfos, err := s.followSvc.FindFollowStatsBatch(ctx, uids, u.UserId)
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
