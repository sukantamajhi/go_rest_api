package utils

import "go.mongodb.org/mongo-driver/bson"

var ProductProjection = bson.M{
	"_id":         1,
	"name":        1,
	"description": 1,
	"sku":         1,
	"createdBy":   1,
	"createdAt":   1,
	"updatedAt":   1,
	"creator":     UserProjection,
}

var UserProjection = bson.M{
	"_id":      1,
	"name":     1,
	"email":    1,
	"username": 1,
}
