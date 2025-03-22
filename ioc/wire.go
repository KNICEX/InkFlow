//go:build wireinject

package ioc

import (
	"github.com/KNICEX/InkFlow/internal/bff"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/email"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/internal/user"
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
		bff.InitBff,
		InitConsumers,
		InitGin,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
