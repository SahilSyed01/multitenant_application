package db

import (
	"context"
	"fmt"
	"multitenant/config"
	"multitenant/models"

	"go.mongodb.org/mongo-driver/bson"
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

// AuthenticateUser checks if the user exists with the correct credentials and returns the tag
func AuthenticateUser(username, password string) (bool, string, error) {
	collection := client.Database("mydatabase").Collection("users")

	// Check if the user exists with the given username and password
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": username, "password": password}).Decode(&user)
	if err != nil {
		// If the error is not nil, it might mean the user is not found
		if err == mongo.ErrNoDocuments {
			return false, "", nil // User not found
		}
		return false, "", err // Other error
	}

	// Return authentication success and the user's tag
	return true, user.Tag, nil
}

func AddManager(username, password string, groupLimit int) (bool, string) {
	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return false, fmt.Sprintf("could not connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Define collections
	managerCollection := client.Database("mydatabase").Collection("managers")
	userCollection := client.Database("mydatabase").Collection("users")

	// Check if a manager with the same username already exists
	var existingManager models.Manager
	err = managerCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&existingManager)
	if err == nil {
		return false, fmt.Sprintf("manager with username '%s' already exists", username)
	} else if err != mongo.ErrNoDocuments {
		return false, fmt.Sprintf("error checking for existing manager: %v", err)
	}

	// Insert only the username and group limit into the managers collection
	manager := models.Manager{
		Username:   username,
		GroupLimit: groupLimit,
	}
	_, err = managerCollection.InsertOne(context.TODO(), manager)
	if err != nil {
		return false, fmt.Sprintf("could not insert manager: %v", err)
	}

	// Insert username, password, and "manager" tag into the users collection
	user := models.User{
		Username: username,
		Password: password,
		Tag:      "manager",
	}
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return false, fmt.Sprintf("manager created, but could not add user: %v", err)
	}

	return true, "manager created successfully"
}