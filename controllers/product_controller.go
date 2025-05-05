package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sukantamajhi/go_rest_api/database"
	"github.com/sukantamajhi/go_rest_api/models"
	"github.com/sukantamajhi/go_rest_api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Sku         string `json:"sku" binding:"required"`
}

func CreateProduct(c *gin.Context) {
	request, err := utils.ParseRequest[CreateProductRequest](c)

	user, exists := c.Get("user")
	if !exists {
		c.SecureJSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Unauthorized",
		})
		return
	}

	log.Printf("user: %+v", user)

	if err != nil {
		c.SecureJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	productCollection := database.GetCollection("products")

	// Check if sku already exists or not
	product, err := models.GetProductBySku(request.Sku)

	if err == nil && product.ID != "" {
		c.SecureJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Product already exists with this sku",
		})
		return
	}

	// creating a product
	_, err = productCollection.InsertOne(context.Background(), bson.M{
		"name":        request.Name,
		"description": request.Description,
		"sku":         request.Sku,
		"createdAt":   time.Now().Format(time.RFC3339),
		"updatedAt":   time.Now().Format(time.RFC3339),
	})

	if err != nil {
		c.SecureJSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

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
