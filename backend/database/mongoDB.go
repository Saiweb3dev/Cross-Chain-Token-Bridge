package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variable to hold the MongoDB client connection
var Client *mongo.Client

// Connects to the MongoDB server
func ConnectToMongoDB() {
	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up client options with the connection string
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Attempt to connect to MongoDB
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to check the connection
	err = Client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	// Log successful connection
	log.Println("Connected to MongoDB!")
}

// Returns the main database instance
func GetDatabase() *mongo.Database {
	return Client.Database("go_ccip_server")
}
