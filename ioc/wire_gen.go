// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func InitApp() *App {
	logger := InitLogger()
	db := InitDB(logger)
	universalClient := InitRedisUniversalClient()
	cmdable := InitRedisCmdable(universalClient)
	client := InitKafka()
	syncProducer := InitSyncProducer(client)
	userService := user.InitUserService(db, cmdable, syncProducer, logger)
	service := email.InitService(logger)
	serviceService := code.InitEmailCodeService(cmdable, service)
	inkService := ink.InitInkService(cmdable, db, logger)
	interactiveService := interactive.InitInteractiveService(cmdable, syncProducer, db, logger)
	rankingService := ink.InitRankingService(cmdable, db, logger, interactiveService)
	followService := relation.InitFollowService(cmdable, db, syncProducer, logger)
	commentService := comment.InitCommentService(db, cmdable, inkService, syncProducer, logger)
	notificationService := notification.InitNotificationService(db)
	gorsexClient := InitGorseCli()
	recommendService := recommend.InitService(gorsexClient, followService, interactiveService, logger)
	actionService := action.InitService()
	feedService := feed.InitService(db, followService, actionService, logger)
	serviceManager := InitMeiliSearch()
	searchService := search.InitSearchService(serviceManager)
	clientClient := InitTemporalClient()
	handler := InitJwtHandler(cmdable)
	authentication := InitAuthMiddleware(handler, logger)
	v := bff.InitBff(userService, serviceService, inkService, rankingService, followService, interactiveService, commentService, notificationService, recommendService, feedService, searchService, clientClient, handler, authentication, logger)
	engine := InitGin(v, logger)
	inkViewConsumer := interactive.InitInteractiveInkReadConsumer(client, logger)
	v2 := InitGeminiClient()
	llmService := ai.InitLLMService(v2)
	service2 := review.InitService(llmService)
	reviewConsumer := review.InitReviewConsumer(clientClient, client, service2, logger)
	syncService := search.InitSyncService(serviceManager)
	syncConsumer := search.InitSyncConsumer(syncService, client, logger)
	notificationConsumer := notification.InitNotificationConsumer(client, notificationService, inkService, commentService, logger)
	recommendSyncService := recommend.InitSyncService(gorsexClient)
	eventSyncConsumer := recommend.InitSyncConsumer(client, recommendSyncService, logger)
	v3 := InitConsumers(inkViewConsumer, reviewConsumer, syncConsumer, notificationConsumer, eventSyncConsumer)
	asyncService := review.InitAsyncService(syncProducer, logger)
	activities := inkpub.NewActivities(inkService, interactiveService, asyncService, syncService, recommendSyncService, notificationService, feedService)
	inkPubWorker := InitInkPubWorker(clientClient, activities)
	rankActivities := schedule.NewRankActivities(rankingService)
	rankTagWorker := InitRankTagWorker(clientClient, rankActivities)
	rankInkWorker := InitRankInkWorker(clientClient, rankActivities)
	failoverService := review.InitFailoverService(clientClient, service2, db, logger)
	reviewFailoverActivity := schedule.NewReviewFailoverActivity(failoverService)
	retryReviewWorker := InitRetryReviewWorker(clientClient, reviewFailoverActivity)
	v4 := InitWorkers(inkPubWorker, rankTagWorker, rankInkWorker, retryReviewWorker)
	rankInkScheduler := InitRankInkScheduler(clientClient)
	rankTagScheduler := InitRankTagScheduler(clientClient)
	reviewFailRetryScheduler := InitReviewRetryScheduler(clientClient)
	v5 := InitSchedulers(rankInkScheduler, rankTagScheduler, reviewFailRetryScheduler)
	app := &App{
		Server:     engine,
		Consumers:  v3,
		Workers:    v4,
		Schedulers: v5,
	}
	return app
}

// wire.go:

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
