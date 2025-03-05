package ioc

import (
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

func InitGin(handlers []ginx.Handler) *gin.Engine {
	r := gin.New()
	ginx.InitErrCodeMetrics(prometheus.CounterOpts{
		Namespace: "ink-flow",
		Subsystem: "web",
		Name:      "http_response_err_code",
		Help:      "http response err code",
	})
	for _, h := range handlers {
		h.RegisterRoutes(r)
	}
	return r
}
