package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Sku         string             `json:"sku" bson:"sku"`
	CreatedBy   primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (p *Product) CollectionName() string {
	return "products"
}

func NewProduct(name, description, sku string, createdBy primitive.ObjectID) *Product {
	now := time.Now()
	return &Product{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		Sku:         sku,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
