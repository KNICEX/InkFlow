//go:build wireinject

package ioc

import (
	"github.com/KNICEX/InkFlow/internal/ai"
	"github.com/KNICEX/InkFlow/internal/bff"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/email"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/notification"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/review"
	"github.com/KNICEX/InkFlow/internal/search"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/internal/workflow/inkpub"
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
		email.InitMemoryService,
		code.InitEmailCodeService,
		ink.InitInkService,
		relation.InitFollowService,

		interactive.InitInteractiveService,
		interactive.InitInteractiveInkReadConsumer,

		notification.InitNotificationService,
		notification.InitNotificationConsumer,

		search.InitSyncService,
		search.InitSearchService,
		search.InitSyncConsumer,

		recommend.InitSyncService,
		recommend.InitSyncConsumer,

		comment.InitCommentService,

		inkpub.NewActivities,

		ai.InitLLMService,
		review.InitAsyncService,
		review.InitReviewConsumer,

		InitInkPubWorker,

		bff.InitBff,
		InitConsumers,
		InitWorkers,
		InitGin,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
