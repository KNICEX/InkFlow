package ioc

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	Server *gin.Engine
	//Consumers []saramax.Consumer

}
