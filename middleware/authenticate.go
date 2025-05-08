package middleware

import (
	"context"
	"errors"
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

// Common error definitions
var (
	ErrUserNotFound     = errors.New("user not found in context")
	ErrInvalidUserType  = errors.New("invalid user type in context")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidToken     = errors.New("invalid token")
	ErrUserDoesNotExist = errors.New("user does not exist")
)

// Authenticate middleware validates JWT tokens and adds user to context
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		token := extractToken(c)
		if token == "" {
			abortWithError(c, ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Validate token and extract claims
		claims, err := validateToken(token)
		if err != nil {
			abortWithError(c, err, http.StatusUnauthorized)
			return
		}

		// Get user from database
		userID := claims["sub"].(string)
		user, err := getUserByID(userID)
		if err != nil {
			abortWithError(c, ErrUserDoesNotExist, http.StatusUnauthorized)
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}

// GetUser retrieves the user from the context
func GetUser(c *gin.Context) (*models.User, error) {
	userValue, exists := c.Get("user")
	if !exists {
		return nil, ErrUserNotFound
	}

	user, ok := userValue.(models.User)
	if !ok {
		return nil, ErrInvalidUserType
	}

	return &user, nil
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) (string, error) {
	user, err := GetUser(c)
	if err != nil {
		return "", err
	}

	return user.ID.Hex(), nil
}

// Helper functions

// extractToken extracts the token from the Authorization header
func extractToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		return ""
	}
	return strings.TrimPrefix(token, "Bearer ")
}

// validateToken validates the JWT token and returns the claims
func validateToken(tokenString string) (jwt.MapClaims, error) {
	mySigningKey := []byte(config.Env.JwtSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// getUserByID retrieves a user from the database by ID
func getUserByID(userID string) (models.User, error) {
	userCollection := database.GetCollection("users")
	var user models.User

	err := userCollection.FindOne(
		context.Background(),
		bson.M{"_id": utils.ObjectIDFromHex(userID)},
	).Decode(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// abortWithError aborts the request with an error response
func abortWithError(c *gin.Context, err error, statusCode int) {
	utils.ErrorResponse(c, err.Error(), nil, statusCode)
	c.Abort()
}
