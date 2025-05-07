package controllers

import (
	"context"
	"errors"
	"net/http"

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
		utils.ErrorResponse(c, err.Error(), nil, http.StatusUnprocessableEntity)
		return
	}

	// Get user context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	productCollection := database.GetCollection("products")

	// Check if sku already exists or not
	var product *models.Product

	err = productCollection.FindOne(context.Background(), bson.M{"sku": request.Sku}).Decode(&product)

	// Check if product already exists with the same SKU
	if err == nil && product.ID != primitive.NilObjectID {
		utils.ErrorResponse(c, "Product already exists with this sku", nil, http.StatusBadRequest)
		return
	}

	newProduct := models.NewProduct(request.Name, request.Description, request.Sku, utils.ObjectIDFromHex(userID))
	_, err = productCollection.InsertOne(context.TODO(), newProduct)

	if err != nil {
		utils.ErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Product created successfully", newProduct, http.StatusCreated)
}

func GetProducts(c *gin.Context) {
	productCollection := database.GetCollection("products")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		utils.ErrorResponse(c, "Unauthorized", nil)
		return
	}

	// Get pagination parameters from query
	page := utils.GetIntQuery(c, "page", 1)
	limit := utils.GetIntQuery(c, "limit", 10)
	skip := (page - 1) * limit

	// Create a more efficient pipeline
	pipeline := []bson.M{
		{
			"$match": bson.M{"createdBy": utils.ObjectIDFromHex(userID)},
		},
		{
			"$sort": bson.M{"createdAt": -1},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "createdBy",
				"foreignField": "_id",
				"as":           "creator",
			},
		},
		{
			"$skip": skip,
		},
		{
			"$limit": limit,
		},
		{
			"$project": utils.ProductProjection,
		},
	}

	// Get total count for pagination
	countPipeline := []bson.M{
		{
			"$match": bson.M{"createdBy": utils.ObjectIDFromHex(userID)},
		},
		{
			"$count": "total",
		},
	}

	// Execute count query
	var countResult []bson.M
	countCursor, err := productCollection.Aggregate(context.Background(), countPipeline)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), nil)
		return
	}
	if err = countCursor.All(context.Background(), &countResult); err != nil {
		utils.ErrorResponse(c, err.Error(), nil)
		return
	}

	total := 0
	if len(countResult) > 0 {
		total = int(countResult[0]["total"].(int32))
	}

	// Execute main query
	cursor, err := productCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), nil)
		return
	}

	var products []bson.M
	if err = cursor.All(context.Background(), &products); err != nil {
		utils.ErrorResponse(c, err.Error(), nil)
		return
	}

	// Return paginated response
	utils.SuccessResponse(c, "Products fetched successfully", gin.H{
		"products": products,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + limit - 1) / limit,
		},
	})
}
