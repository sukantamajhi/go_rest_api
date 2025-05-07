package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sukantamajhi/go_rest_api/config"
	"github.com/sukantamajhi/go_rest_api/database"
	"github.com/sukantamajhi/go_rest_api/dtos/requests"
	"github.com/sukantamajhi/go_rest_api/models"
	"github.com/sukantamajhi/go_rest_api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	request, err := utils.ParseRequest[requests.RegisterRequest](c)

	if err != nil {
		utils.ErrorResponse(c, err.Error(), nil, http.StatusUnprocessableEntity)
		return
	}

	username := strings.TrimSpace(request.Username)
	name := strings.TrimSpace(request.Name)
	email := strings.TrimSpace(request.Email)
	password := strings.TrimSpace(request.Password)

	userCollection := database.GetCollection("users")

	// Find the existing user
	var existingUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"email": email}, {"username": username}}}).Decode(&existingUser)

	if err != nil && err != mongo.ErrNoDocuments {
		utils.ErrorResponse(c, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if existingUser.Email != "" {
		utils.ErrorResponse(c, "User already exists", gin.H{
			"email":    email,
			"username": username,
		}, http.StatusBadRequest)
		return
	}

	user := models.NewUser(username, email, name, "1234567890", password)

	insertedUser, err := userCollection.InsertOne(context.TODO(), user)

	log.Println("error", err)

	if err != nil {
		utils.ErrorResponse(c, "Failed to insert user", nil, http.StatusInternalServerError)
		return
	}

	user.ID = insertedUser.InsertedID.(primitive.ObjectID)

	utils.SuccessResponse(c, "User registered successfully", user, http.StatusCreated)
}

type LoginRequest struct {
	Identity string `json:"identity" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func Login(c *gin.Context) {
	request, err := utils.ParseRequest[LoginRequest](c)

	if err != nil {
		log.Println("Error in login validation", err)
		utils.ErrorResponse(c, err.Error(), nil, http.StatusUnprocessableEntity)
		return
	}

	identity := strings.TrimSpace(request.Identity)
	password := strings.TrimSpace(request.Password)

	userCollection := database.GetCollection("users")

	// Find the existing user
	var existingUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"email": identity}, {"username": identity}}}).Decode(&existingUser)

	if err != nil && err != mongo.ErrNoDocuments {
		utils.ErrorResponse(c, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if existingUser.Email == "" {
		utils.ErrorResponse(c, "User does not exist", nil, http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))

	if err != nil {
		utils.ErrorResponse(c, "Invalid password", nil, http.StatusUnauthorized)
		return
	}

	mySigningKey := []byte(config.Env.JwtSecret)

	// Create the Claims
	claims := &jwt.StandardClaims{
		Subject:   existingUser.ID.Hex(),
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

		if config.Env.GinMode == "release" {
			message = "Something went wrong"
		} else {
			message = err.Error()
		}

		utils.ErrorResponse(c, message, nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "User logged in successfully", gin.H{
		"token": tokenString,
	})
}
