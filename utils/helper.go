package utils

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sukantamajhi/go_rest_api/dtos/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SuccessResponse(c *gin.Context, message string, data any, status ...int) {
	response := requests.Response{}

	if len(status) > 0 {
		c.SecureJSON(status[0], response.Success(message, data))
	} else {
		c.SecureJSON(http.StatusOK, response.Success(message, data))
	}
}

func ErrorResponse(c *gin.Context, message string, data any, status ...int) {
	response := requests.Response{}

	if len(status) > 0 {
		c.SecureJSON(status[0], response.ErrorWithData(message, data))
	} else {
		c.SecureJSON(http.StatusInternalServerError, response.ErrorWithData(message, data))
	}
}

func TrimmedString(str *string) string {
	if str == nil {
		return ""
	}
	return strings.TrimSpace(*str)
}

func ObjectIDFromHex(hex string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Fatal(err)
	}
	return objectID
}
