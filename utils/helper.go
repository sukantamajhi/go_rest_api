package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SuccessResponse(c *gin.Context, message string, data any) {
	c.SecureJSON(http.StatusOK, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, message string) {
	c.SecureJSON(http.StatusBadRequest, Response{
		Status:  false,
		Message: message,
		Data:    nil,
	})
}

func TrimmedString(str *string) string {
	if str == nil {
		return ""
	}
	return strings.TrimSpace(*str)
}
