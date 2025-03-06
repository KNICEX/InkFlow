//go:build wireinject

package bff

import (
	"github.com/KNICEX/InkFlow/internal/bff/internal/web"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/google/wire"
)

var handlers = wire.NewSet(
	web.NewUserHandler,
)

func InitHandlers(uh *web.UserHandler) []ginx.Handler {
	return []ginx.Handler{uh}
}

func InitBff(userSvc user.Service, codeSvc code.Service, jwtHandler jwt.Handler, auth middleware.Authentication, log logx.Logger) []ginx.Handler {
	wire.Build(
		web.NewUserHandler,
		InitHandlers,
	)
	return []ginx.Handler{}
}
