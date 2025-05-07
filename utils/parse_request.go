package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// TODO: I have to create a middleware for parsing request
func ParseRequest[T any](c *gin.Context) (T, error) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}

			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
			c.Abort()
			return request, errors.New("validation error")
		}
	}

	return request, nil
}
