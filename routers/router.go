package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Health check route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	apiRouter := router.Group("/api")

	AuthRouter(apiRouter)
	ProductRouter(apiRouter)

	return router
}
