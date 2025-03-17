//go:build wireinject

package bff

import (
	"github.com/KNICEX/InkFlow/internal/bff/internal/web"
	"github.com/KNICEX/InkFlow/internal/code"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/user"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/google/wire"
)

func InitHandlers(uh *web.UserHandler, ih *web.InkHandler, fh *web.FileHandler) []ginx.Handler {
	return []ginx.Handler{uh, ih, fh}
}

func InitBff(userSvc user.Service, codeSvc code.Service, inkService ink.Service,
	interactiveSvc interactive.Service,
	jwtHandler jwt.Handler, auth middleware.Authentication, log logx.Logger) []ginx.Handler {
	wire.Build(
		web.NewUserHandler,
		web.NewInkHandler,
		initCloudinary,
		web.NewFileHandler,
		InitHandlers,
	)
	return []ginx.Handler{}
}
