package routes

import (
	handler "server/services/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) *gin.RouterGroup {
	userRouter := router.Group("/user")
	userRouter.POST("/create", handler.CreateUser)
	userRouter.GET("/:id", handler.GetUser)

	return userRouter
}
