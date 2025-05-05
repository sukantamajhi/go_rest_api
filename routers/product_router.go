package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sukantamajhi/go_rest_api/controllers"
	"github.com/sukantamajhi/go_rest_api/middleware"
)

func ProductRouter(router *gin.RouterGroup) *gin.RouterGroup {
	productRouter := router.Group("/products")

	productRouter.POST("/", middleware.Authenticate(), controllers.CreateProduct)
	productRouter.GET("/", controllers.GetProducts)
	// productRouter.GET("/:id", controllers.GetProductById)
	// productRouter.PUT("/:id", controllers.UpdateProduct)
	// productRouter.DELETE("/:id", controllers.DeleteProduct)
	return productRouter
}
