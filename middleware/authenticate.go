package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sukantamajhi/go_rest_api/config"
	"github.com/sukantamajhi/go_rest_api/database"
	"github.com/sukantamajhi/go_rest_api/models"
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
			fmt.Println(claims)
			userID := claims["sub"]

			userCollection := database.GetCollection("users")

			var user models.User
			err := userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)

			if err != nil {
				c.SecureJSON(http.StatusUnauthorized, gin.H{
					"status":  false,
					"code":    "UNAUTHORIZED",
					"message": "Unauthorized",
				})
				c.Abort()
			}

			c.Set("user", user)
		} else {
			fmt.Println(err)
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
