package models

import (
	"context"

	"github.com/sukantamajhi/go_rest_api/database"
	"go.mongodb.org/mongo-driver/bson"
)

type Product struct {
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Sku         string `json:"sku" bson:"sku"`
	CreatedAt   string `json:"createdAt" bson:"createdAt"`
	UpdatedAt   string `json:"updatedAt" bson:"updatedAt"`
}

func GetProductBySku(sku string) (*Product, error) {
	ProductCollection := database.GetCollection("products")
	product := &Product{}
	err := ProductCollection.FindOne(context.Background(), bson.M{"sku": sku}).Decode(product)
	return product, err
}

func GetAllProducts() ([]*Product, error) {
	ProductCollection := database.GetCollection("products")

	products := []*Product{}
	cursor, err := ProductCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &products)
	return products, err
}

func GetProductById(id string) (*Product, error) {
	ProductCollection := database.GetCollection("products")

	product := &Product{}
	err := ProductCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(product)
	return product, err
}

func UpdateProduct(id string, product *Product) error {
	ProductCollection := database.GetCollection("products")

	_, err := ProductCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": product})
	return err
}
