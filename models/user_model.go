// Generate User Model
package models

import (
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username,omitempty,unique,required,index"`
	Email     string             `json:"email" bson:"email,omitempty,unique,required,index"`
	Name      string             `json:"name" bson:"name,omitempty,required"`
	Phone     string             `json:"phone" bson:"phone,omitempty,unique,required,index"`
	Role      string             `json:"role" bson:"role,omitempty"`
	Password  string             `json:"-" bson:"password,omitempty,required"`
	Status    bool               `json:"status" bson:"status,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func (u *User) CollectionName() string {
	return "users"
}

func NewUser(username, email, name, phone, password string) *User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(password)), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	return &User{
		ID:        primitive.NewObjectID(),
		Username:  strings.TrimSpace(username),
		Email:     strings.TrimSpace(email),
		Phone:     strings.TrimSpace(phone),
		Name:      strings.TrimSpace(name),
		Password:  string(hashedPassword),
		Status:    true,
		Role:      "customer",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
