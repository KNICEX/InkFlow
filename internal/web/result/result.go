package result

import (
	"github.com/KNICEX/InkFlow/pkg/ginx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result = ginx.Result

func Success(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ginx.Success())
}

func SuccessWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, ginx.SuccessWithMsg(msg))
}

func SuccessWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, ginx.SuccessWithData(data))
}

func InvalidToken(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ginx.InvalidToken())
}

func InvalidParam(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ginx.InvalidParam())
}

func InvalidParamWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, ginx.InvalidParamWithMsg(msg))
}

func BizError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, ginx.BizError(msg))
}

func InternalError(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ginx.InternalError())
}
