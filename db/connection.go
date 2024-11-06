package db

import (
	"context"
	"multitenant/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// ConnectMongoDB initializes the MongoDB client
func ConnectMongoDB() error {
	clientOptions := options.Client().ApplyURI(config.MongoURI)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	// Ping MongoDB to verify connection
	return client.Ping(context.TODO(), nil)
}

// DisconnectMongoDB closes the MongoDB connection
func DisconnectMongoDB() {
	if client != nil {
		client.Disconnect(context.TODO())
	}
}
