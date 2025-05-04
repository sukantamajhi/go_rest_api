package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sukantamajhi/go_rest_api/controllers"
)

func AuthRouter(router *gin.RouterGroup) *gin.RouterGroup {
	authRouter := router.Group("/auth")
	authRouter.POST("/register", controllers.Register)
	authRouter.POST("/login", controllers.Login)

	return authRouter
}
