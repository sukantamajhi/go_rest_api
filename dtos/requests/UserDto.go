package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type ResponseUser struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
}

type RegisterResponse struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    ResponseUser `json:"data"`
}
