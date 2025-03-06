//go:build wireinject

package ioc

import (
	"github.com/KNICEX/InkFlow/internal/bff"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/email"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/google/wire"
)

var thirdPartSet = wire.NewSet(
	InitLogger,
	InitDB,
	InitEs,
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
		user.InitGithubOAuth2Service,
		email.InitService,
		code.InitEmailCodeService,
		bff.InitBff,
		InitGin,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
