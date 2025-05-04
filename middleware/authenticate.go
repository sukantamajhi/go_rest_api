package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sukantamajhi/go_rest_api/config"
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
			c.Set("user", claims["sub"])
		} else {
			fmt.Println(err)
		}

		c.Next()
	}
}
