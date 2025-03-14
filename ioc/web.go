package ioc

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/KNICEX/InkFlow/pkg/ginx/jwt"
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

func InitAuthMiddleware(h jwt.Handler, l logx.Logger) middleware.Authentication {
	return middleware.NewJwtLoginBuilder(h, l)
}

func InitJwtHandler(cmd redis.Cmdable) jwt.Handler {
	return jwt.NewRedisHandler(cmd)
}

func InitGin(handlers []ginx.Handler, l logx.Logger) *gin.Engine {
	r := gin.New()
	r.Use(middleware.NewLoggerBuilder(func(ctx context.Context, al *middleware.AccessLog) {
		l.WithCtx(ctx).Debug("gin access log", logx.Any("content", al))
	}).AllowRespBody().AllowReqBody().Build())
	ginx.InitErrCodeMetrics(prometheus.CounterOpts{
		Namespace: "ink-flow",
		Subsystem: "web",
		Name:      "http_response_err_code",
		Help:      "http response err code",
	})
	g := r.Group("/api/v1")
	for _, h := range handlers {
		h.RegisterRoutes(g)
	}
	return r
}
