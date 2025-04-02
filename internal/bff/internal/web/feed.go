package web

import (
	"github.com/KNICEX/InkFlow/internal/feed"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type FeedHandler struct {
	svc          feed.Service
	recommendSvc recommend.Service
	*inkAggregate
	auth middleware.Authentication
	l    logx.Logger
}

func NewFeedHandler(
	svc feed.Service,
	inkSvc ink.Service,
	intrSvc interactive.Service,
	userSvc user.Service,
	recommendSvc recommend.Service,
	auth middleware.Authentication,
	l logx.Logger,
) *FeedHandler {
	return &FeedHandler{
		svc:          svc,
		recommendSvc: recommendSvc,
		inkAggregate: newInkAggregate(inkSvc, userSvc, intrSvc),
		auth:         auth,
		l:            l,
	}
}

func (h *FeedHandler) RegisterRoutes(server *gin.RouterGroup) {
	feedGroup := server.Group("/feed")
	{
		feedGroup.GET("/ink/follow", h.auth.CheckLogin(), ginx.WrapBody(h.l, h.Follow))
		feedGroup.GET("/ink/recommend", h.auth.CheckLogin(), ginx.WrapBody(h.l, h.Recommend))

		feedGroup.GET("/ink/hot", ginx.WrapBody(h.l, h.Hot))
	}
}

func (h *FeedHandler) Follow(ctx *gin.Context, req FeedFollowReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	feeds, err := h.svc.FollowFeedInkList(ctx, uc.UserId, req.MaxId, req.Timestamp, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(feeds) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
	}

	inkIds := lo.Map(feeds, func(item feed.Feed, index int) int64 {
		return item.BizId
	})
	authorIds := lo.Map(feeds, func(item feed.Feed, index int) int64 {
		return item.Content.(feed.Ink).AuthorId
	})
	var intrs map[int64]interactive.Interactive
	var authors map[int64]user.User
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		intrs, er = h.intrSvc.GetMulti(ctx, bizInk, inkIds, uc.UserId)
		return er
	})
	eg.Go(func() error {
		var er error
		authors, er = h.userSvc.FindByIds(ctx, authorIds)
		return er
	})
	if err = eg.Wait(); err != nil {
		return ginx.InternalError(), err
	}

	inkVos := make([]InkVO, 0, len(feeds))
	for _, f := range feeds {
		vo := feedInkToVO(f.Content.(feed.Ink))
		vo.Interactive = intrToVo(intrs[f.BizId])
		vo.Author = userToVO(authors[f.Content.(feed.Ink).AuthorId])
		inkVos = append(inkVos, vo)
	}
	return ginx.SuccessWithData(inkVos), nil
}

func (h *FeedHandler) Recommend(ctx *gin.Context, req FeedRecommendReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	inkIds, err := h.recommendSvc.FindRecommendInk(ctx, uc.UserId, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(inkIds) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
	}

	inkVos, err := h.inkAggregate.GetInkList(ctx, inkIds, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkVO, 0, len(inkVos))
	for _, id := range inkIds {
		vo, ok := inkVos[id]
		if !ok {
			continue
		}
		res = append(res, vo)
	}
	return ginx.SuccessWithData(res), nil
}

func (h *FeedHandler) Hot(ctx *gin.Context, req FeedRecommendReq) (ginx.Result, error) {
	uc := jwt.MustGetUserClaims(ctx)
	inkIds, err := h.recommendSvc.FindPopular(ctx, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(inkIds) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
	}

	inkVos, err := h.inkAggregate.GetInkList(ctx, inkIds, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}
	res := make([]InkVO, 0, len(inkVos))
	for _, id := range inkIds {
		vo, ok := inkVos[id]
		if !ok {
			continue
		}
		res = append(res, vo)
	}
	return ginx.SuccessWithData(res), nil
}
