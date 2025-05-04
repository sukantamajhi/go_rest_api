package utils

import "github.com/gin-gonic/gin"

func ParseRequest[T any](c *gin.Context) (T, error) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		return request, err
	}
	return request, nil
}
