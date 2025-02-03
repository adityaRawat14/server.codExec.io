package routes

import (
	handler "server/services/handler"

	"github.com/gin-gonic/gin"
) 

func CodeRoutes(router *gin.Engine) *gin.RouterGroup {
	codeRouter := router.Group("/:language")
	codeRouter.POST("/",handler.HandleCodeExecution)

	return codeRouter
}
