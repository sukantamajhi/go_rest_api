package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sukantamajhi/go_rest_api/config"
	"github.com/sukantamajhi/go_rest_api/database"
	"github.com/sukantamajhi/go_rest_api/models"
	"github.com/sukantamajhi/go_rest_api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.ErrorResponse(c, "Unauthorized", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		mySigningKey := []byte(config.AppConfig.JwtSecret)
		// Parse the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
			return mySigningKey, nil
		})

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			userID := claims["sub"].(string)

			log.Printf("userID: %+v", userID)

			userCollection := database.GetCollection("users")

			var user models.User
			err := userCollection.FindOne(context.Background(), bson.M{"_id": utils.ObjectIDFromHex(userID)}).Decode(&user)

			if err != nil {
				log.Println("Error fetching user", err)
				utils.ErrorResponse(c, "Unauthorized", nil, http.StatusUnauthorized)
				c.Abort()
			}

			c.Set("user", user)
		} else {
			log.Println("Error parsing token", err)
			utils.ErrorResponse(c, "Unauthorized", nil, http.StatusUnauthorized)
			c.Abort()
		}

		c.Next()
	}
}

func GetUser(c *gin.Context) (*models.User, error) {
	userValue, exists := c.Get("user")

	if !exists {
		return nil, errors.New("user not found in context")
	}

	user, ok := userValue.(models.User)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	log.Printf("user: %+v", user)

	return &user, nil
}

func GetUserID(c *gin.Context) (string, error) {
	user, err := GetUser(c)
	if err != nil {
		return "", err
	}

	return user.ID.Hex(), nil
}
