package database

import (
	"context"
	"fmt"

	"github.com/sukantamajhi/go_rest_api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func Connect_to_db() {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.AppConfig.MongoDBURI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	var err error
	client, err = mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
	return client.Database("go_rest_api").Collection(collectionName)
}

func CloseDB() {
	if client != nil {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}
