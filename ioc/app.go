package ioc

import (
	"github.com/KNICEX/InkFlow/pkg/saramax"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/worker"
)

type App struct {
	Server    *gin.Engine
	Consumers []saramax.Consumer
	Workers   []worker.Worker
}
