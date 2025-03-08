package ginx

import (
	"fmt"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

var vector *prometheus.CounterVec

// InitErrCodeMetrics 初始化错误码统计
// 注意：要使用wrap功能，必须调用该方法初始化vector
func InitErrCodeMetrics(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	prometheus.MustRegister(vector)
}

func WrapBody[T any](l logx.Logger, bizFn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.ShouldBind(&req); err != nil {
			res := InvalidParam()
			vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
			ctx.JSON(http.StatusOK, res)
			return
		}

		res, err := bizFn(ctx, req)
		// 统计错误码
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			l.WithCtx(ctx).Error("handle http error: ",
				logx.Error(err),
				logx.String("path", ctx.Request.URL.Path),
				logx.String("route", fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())),
			)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func Wrap(l logx.Logger, bizFn func(ctx *gin.Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := bizFn(ctx)
		// 统计错误码
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			l.WithCtx(ctx).Error("handle http error: ",
				logx.Error(err),
				logx.String("path", ctx.Request.URL.Path),
				logx.String("route", fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())),
			)
		}
		ctx.JSON(http.StatusOK, res)
	}
}
