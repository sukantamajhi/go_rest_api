package requests

import (
	"log"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (c *RegisterRequest) Validate() error {
	log.Printf("RegisterRequest: %+v", c)
	validate := validator.New()

	return validate.Struct(c)
}

type RegisterResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
