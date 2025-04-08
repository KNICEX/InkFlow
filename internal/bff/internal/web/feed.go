package web

import (
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/feed"
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
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type FeedHandler struct {
	svc            feed.Service
	recommendSvc   recommend.Service
	inkRankService ink.RankingService
	inkAggregate   *inkAggregate
	userAggregate  *userAggregate
	intrAggregate  *interactiveAggregate
	auth           middleware.Authentication
	l              logx.Logger
}

func NewFeedHandler(
	svc feed.Service,
	inkSvc ink.Service,
	inkRankService ink.RankingService,
	intrSvc interactive.Service,
	userSvc user.Service,
	followSvc relation.FollowService,
	recommendSvc recommend.Service,
	auth middleware.Authentication,
	commentSvc comment.Service,
	l logx.Logger,
) *FeedHandler {
	return &FeedHandler{
		svc:            svc,
		inkRankService: inkRankService,
		recommendSvc:   recommendSvc,
		inkAggregate:   newInkAggregate(inkSvc, userSvc, followSvc, intrSvc, commentSvc),
		userAggregate:  newUserAggregate(userSvc, followSvc),
		intrAggregate:  newInteractiveAggregate(intrSvc, commentSvc),
		auth:           auth,
		l:              l,
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

	inksMap, err := h.inkAggregate.GetInkList(ctx, inkIds, uc.UserId)
	if err != nil {
		return ginx.InternalError(), err
	}

	inkVos := make([]InkVO, 0, len(feeds))
	for _, f := range feeds {
		inkVo, ok := inksMap[f.BizId]
		if !ok {
			continue
		}
		inkVos = append(inkVos, inkVo)
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
	inks, err := h.inkRankService.FindTopNInk(ctx, req.Offset, req.Limit)
	if err != nil {
		return ginx.InternalError(), err
	}
	if len(inks) == 0 {
		return ginx.SuccessWithData([]InkVO{}), nil
	}
	inkIds := make([]int64, 0, len(inks))
	authorIds := make([]int64, 0, len(inks))
	for _, i := range inks {
		inkIds = append(inkIds, i.Id)
		authorIds = append(authorIds, i.Author.Id)
	}

	var userVos map[int64]UserVO
	var intrVos map[int64]InteractiveVO
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		userVos, er = h.userAggregate.GetUserList(ctx, authorIds, uc.UserId)
		return er
	})
	eg.Go(func() error {
		var er error
		intrVos, er = h.intrAggregate.GetInteractiveList(ctx, bizInk, inkIds, uc.UserId)
		return er
	})

	vos := make([]InkVO, 0, len(inks))
	for _, i := range inks {
		vo := inkToVO(i)
		vo.Interactive = intrVos[i.Id]
		vo.Author = userVos[i.Author.Id]
		vos = append(vos, vo)
	}

	return ginx.SuccessWithData(vos), nil
}
