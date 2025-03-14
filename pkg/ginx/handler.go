package ginx

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoutes(server *gin.RouterGroup)
}
