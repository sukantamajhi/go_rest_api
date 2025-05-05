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
			c.SecureJSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"code":    "UNAUTHORIZED",
				"message": "Unauthorized",
			})
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
			userID := claims["sub"]

			log.Printf("userID: %+v", userID)

			userCollection := database.GetCollection("users")

			var user models.User
			err := userCollection.FindOne(context.Background(), bson.M{"_id": userID.(string)}).Decode(&user)

			if err != nil {
				log.Println("Error fetching user", err)
				c.SecureJSON(http.StatusUnauthorized, gin.H{
					"status":  false,
					"code":    "UNAUTHORIZED",
					"message": "Unauthorized",
				})
				c.Abort()
			}

			c.Set("user", user)
		} else {
			log.Println("Error parsing token", err)
			c.SecureJSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"code":    "UNAUTHORIZED",
				"message": "Unauthorized",
			})
			c.Abort()
		}

		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (models.User, error) {
	user, ok := c.Get("user")

	log.Printf("user: %+v", user)

	if !ok {
		utils.ErrorResponse(c, "User not found in context")
		return models.User{}, errors.New("user not found in context")
	}

	return user.(models.User), nil
}
