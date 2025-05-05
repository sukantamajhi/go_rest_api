package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	

	c.SecureJSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Product created successfully",
	})
}

func GetProducts(c *gin.Context) {
	fmt.Println("Get Products")
	c.SecureJSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}
