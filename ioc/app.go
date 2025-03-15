package ioc

import (
	"github.com/KNICEX/InkFlow/pkg/saramax"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server    *gin.Engine
	Consumers []saramax.Consumer
}
