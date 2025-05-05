package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sukantamajhi/go_rest_api/database"
	"github.com/sukantamajhi/go_rest_api/dtos/requests"
	"github.com/sukantamajhi/go_rest_api/middleware"
	"github.com/sukantamajhi/go_rest_api/models"
	"github.com/sukantamajhi/go_rest_api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrProductExists = errors.New("product already exists with this SKU")
	ErrInvalidUser   = errors.New("invalid user context")
)

func CreateProduct(c *gin.Context) {
	// Parse and validate request
	request, err := utils.ParseRequest[requests.CreateProductRequest](c)
	if err != nil {
		c.SecureJSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// Get user context
	user, _ := middleware.GetUserFromContext(c)

	userID := user.ID
	productCollection := database.GetCollection("products")

	// Check if sku already exists or not
	var product *models.Product

	err = productCollection.FindOne(context.Background(), bson.M{"sku": request.Sku}).Decode(&product)

	// Check if product already exists with the same SKU
	if err == nil && product.ID != primitive.NilObjectID {
		utils.ErrorResponse(c, "Product already exists with this sku")
		return
	}

	// Create new product
	now := time.Now()
	newProduct := models.Product{
		ID:          primitive.NewObjectID(),
		Name:        utils.TrimmedString(&request.Name),
		Description: utils.TrimmedString(&request.Description),
		Sku:         utils.TrimmedString(&request.Sku),
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	insertedProduct, err := productCollection.InsertOne(context.TODO(), newProduct)

	if err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	newProduct.ID = insertedProduct.InsertedID.(primitive.ObjectID)

	utils.SuccessResponse(c, "Product created successfully", newProduct)
}

func GetProducts(c *gin.Context) {
	productCollection := database.GetCollection("products")

	user, exists := c.Get("user")

	if !exists {
		utils.ErrorResponse(c, "Unauthorized")
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
	}

	c.SecureJSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Products fetched successfully",
		"data":    products,
	})
}
