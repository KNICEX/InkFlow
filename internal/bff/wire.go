//go:build wireinject

package bff

import (
	"github.com/KNICEX/InkFlow/internal/bff/internal/web"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/feed"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/notification"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/search"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/google/wire"
	"go.temporal.io/sdk/client"
)

func InitHandlers(uh *web.UserHandler, ih *web.InkHandler, fh *web.FileHandler,
	ch *web.CommentHandler, nh *web.NotificationHandler, sh *web.SearchHandler, feedH *web.FeedHandler,
	statsH *web.StatsHandler, intrH *web.InteractiveHandler, rh *web.RecommendHandler) []ginx.Handler {
	return []ginx.Handler{uh, ih, fh, ch, nh, sh, feedH, statsH, intrH, rh}
}

func initUserAggregate(userSvc user.Service, followSvc relation.FollowService) *web.UserAggregate {
	return web.NewUserAggregate(userSvc, followSvc)
}

func initInkAggregate(inkSvc ink.Service, userAggregate *web.UserAggregate,
	interactiveAggregate *web.InteractiveAggregate,
	intrSvc interactive.Service, commentSvc comment.Service) *web.InkAggregate {
	return web.NewInkAggregate(inkSvc, userAggregate, interactiveAggregate)
}
func initInteractiveAggregate(intrSvc interactive.Service, commentSvc comment.Service) *web.InteractiveAggregate {
	return web.NewInteractiveAggregate(intrSvc, commentSvc)
}

func InitBff(userSvc user.Service, codeSvc code.Service, inkService ink.Service,
	inkRankService ink.RankingService,
	followService relation.FollowService,
	interactiveSvc interactive.Service,
	commentSvc comment.Service,
	notificationSvc notification.Service,
	recommendSvc recommend.Service,
	feedSvc feed.Service,
	searchSvc search.Service,
	workflowCli client.Client,
	jwtHandler jwt.Handler, auth middleware.Authentication, log logx.Logger) []ginx.Handler {
	wire.Build(
		web.NewUserAggregate,
		web.NewInkAggregate,
		web.NewInteractiveAggregate,
		web.NewUserHandler,
		web.NewInkHandler,
		web.NewCommentHandler,
		web.NewNotificationHandler,
		web.NewSearchHandler,
		web.NewFeedHandler,
		web.NewInteractiveHandler,
		web.NewStatsHandler,
		initCloudinary,
		web.NewFileHandler,
		web.NewRecommendHandler,
		InitHandlers,
	)
	return []ginx.Handler{}
}
