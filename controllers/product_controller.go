package controllers

import (
	"context"
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

	userID := user.(models.User).ID

	log.Println("userID", userID)

	if err != nil {
		c.SecureJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	productCollection := database.GetCollection("products")

	// Check if sku already exists or not
	var product *models.Product
	err = productCollection.FindOne(context.Background(), bson.M{"sku": request.Sku}).Decode(&product)

	log.Printf("product: %+v", product)

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
		"createdBy":   userID,
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
	productCollection := database.GetCollection("products")

	user, exists := c.Get("user")

	if !exists {
		c.SecureJSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Unauthorized",
		})
		return
	}

	userID := user.(models.User).ID

	cursor, err := productCollection.Find(context.Background(), bson.M{"createdBy": userID})
	if err != nil {
		c.SecureJSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
		})
	}

	products := []*models.Product{}
	err = cursor.All(context.Background(), &products)
	if err != nil {
		c.SecureJSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	c.SecureJSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Products fetched successfully",
		"data":    products,
	})
}
