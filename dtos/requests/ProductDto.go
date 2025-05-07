package requests

import "github.com/go-playground/validator/v10"

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Sku         string `json:"sku" binding:"required"`
}

func (c *CreateProductRequest) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}

type ProductResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sku         string `json:"sku"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ProductCreateResponse struct {
	Error   bool            `json:"error"`
	Message string          `json:"message"`
	Product ProductResponse `json:"product"`
}
