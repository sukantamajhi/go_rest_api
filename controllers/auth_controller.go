package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sukantamajhi/go_rest_api/config"
	"github.com/sukantamajhi/go_rest_api/database"
	"github.com/sukantamajhi/go_rest_api/models"
	"github.com/sukantamajhi/go_rest_api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type ResponseUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterResponse struct {
	Message string       `json:"message"`
	User    ResponseUser `json:"user"`
	UserId  any          `json:"userId"`
}

func Register(c *gin.Context) {
	request, err := utils.ParseRequest[RegisterRequest](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userCollection := database.GetCollection("users")

	// Find the existing user
	var existingUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&existingUser)

	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check user existence",
		})
		return
	}

	if existingUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		ID:        primitive.NewObjectID().Hex(),
		Name:      request.Name,
		Email:     request.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertedUser, err := userCollection.InsertOne(context.TODO(), user)

	log.Println("error", err)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to insert user",
		})
		return
	}

	log.Println("insertedUser", insertedUser)

	response := RegisterResponse{
		Message: "User registered successfully",
		User: ResponseUser{
			Name:  user.Name,
			Email: user.Email,
		},
		UserId: insertedUser.InsertedID,
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": response,
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func Login(c *gin.Context) {
	request, err := utils.ParseRequest[LoginRequest](c)

	if err != nil {
		log.Println("Error in login validation", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"code":    "LOGIN_FAILED",
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}

	userCollection := database.GetCollection("users")

	// Find the existing user
	var existingUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&existingUser)

	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"code":    "SERVER_ERROR",
			"message": "Failed to check user existence",
		})
		return
	}

	if existingUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"code":    "USER_NOT_FOUND",
			"message": "User does not exist",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(request.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"code":    "LOGIN_FAILED",
			"message": "Invalid password",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	mySigningKey := []byte(config.AppConfig.JwtSecret)

	// Create the Claims
	claims := &jwt.StandardClaims{
		Subject:   existingUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    existingUser.Email,
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Println("Error: ", err)
		var message string

		if config.AppConfig.GinMode == "release" {
			message = "Something went wrong"
		} else {
			message = err.Error()
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"code":    "SERVER_ERROR",
			"message": message,
		})
		c.Abort()
		return
	}

	c.SecureJSON(http.StatusOK, gin.H{
		"status":  true,
		"token":   tokenString,
		"message": "User logged in successfully",
	})
}
