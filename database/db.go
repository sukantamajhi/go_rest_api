package database

import (
	"context"
	"fmt"
	"time"

	"github.com/sukantamajhi/go_rest_api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client   *mongo.Client
	database *mongo.Database
)

func Connect_to_db() {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// Configure connection pool and other performance settings
	opts := options.Client().
		ApplyURI(config.Env.MongoDBURI).
		SetServerAPIOptions(serverAPI).
		SetMaxPoolSize(100).                       // Maximum number of connections in the pool
		SetMinPoolSize(10).                        // Minimum number of connections in the pool
		SetMaxConnIdleTime(5 * time.Minute).       // Maximum time a connection can remain idle
		SetRetryWrites(true).                      // Enable retryable writes
		SetRetryReads(true).                       // Enable retryable reads
		SetSocketTimeout(10 * time.Second).        // Socket timeout
		SetConnectTimeout(10 * time.Second).       // Connection timeout
		SetServerSelectionTimeout(5 * time.Second) // Server selection timeout

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

	// Get database instance
	database = client.Database(config.Env.Database_Name)

	defer fmt.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}

func CloseDB() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}
}
