package web

import (
	"github.com/KNICEX/InkFlow/pkg/ginx/middleware"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	auth middleware.Authentication
	l    logx.Logger
}

func (handler *FileHandler) RegisterRoutes(server *gin.RouterGroup) {

}
