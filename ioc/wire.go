//go:build wireinject

package ioc

import (
	"github.com/KNICEX/InkFlow/internal/action"
	"github.com/KNICEX/InkFlow/internal/ai"
	"github.com/KNICEX/InkFlow/internal/bff"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/email"
	"github.com/KNICEX/InkFlow/internal/feed"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/notification"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/review"
	"github.com/KNICEX/InkFlow/internal/search"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/internal/workflow/inkpub"
	"github.com/KNICEX/InkFlow/internal/workflow/schedule"
	"github.com/google/wire"
)

var thirdPartSet = wire.NewSet(
	InitLogger,
	InitDB,
	InitMeiliSearch,
	InitKafka,
	InitSyncProducer,
	InitRedisUniversalClient,
	InitRedisCmdable,
	InitGeminiClient,
	InitTemporalClient,
	InitGorseCli,
)

var webSet = wire.NewSet(
	InitJwtHandler,
	InitAuthMiddleware,
)

func InitApp() *App {
	wire.Build(
		thirdPartSet,
		webSet,
		user.InitUserService,
		//email.InitService,
		email.InitService,
		code.InitEmailCodeService,

		relation.InitFollowService,

		interactive.InitInteractiveService,
		interactive.InitInteractiveInkReadConsumer,

		ink.InitInkService,
		ink.InitRankingService,

		notification.InitNotificationService,
		notification.InitNotificationConsumer,

		search.InitSyncService,
		search.InitSearchService,
		search.InitSyncConsumer,

		recommend.InitSyncService,
		recommend.InitSyncConsumer,
		recommend.InitService,

		comment.InitCommentService,

		ai.InitLLMService,
		review.InitService,
		review.InitAsyncService,
		review.InitReviewConsumer,
		review.InitFailoverService,

		action.InitService,
		feed.InitService,

		inkpub.NewActivities,
		schedule.NewRankActivities,
		schedule.NewReviewFailoverActivity,

		InitRankTagWorker,
		InitRankInkWorker,
		InitInkPubWorker,
		InitRetryReviewWorker,

		InitRankInkScheduler,
		InitRankTagScheduler,
		InitReviewRetryScheduler,
		InitSchedulers,

		bff.InitBff,
		InitConsumers,
		InitWorkers,
		InitGin,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
