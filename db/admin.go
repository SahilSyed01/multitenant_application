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

// AddManager adds a new manager and user to the respective collections
func AddManager(username, password string, groupLimit int) (models.CreateManagerResponse) {
	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		// Return JSON formatted error
		return models.CreateManagerResponse{
			Success: false,
			Message: fmt.Sprintf("could not connect to MongoDB: %v", err),
		}
	}
	defer client.Disconnect(context.TODO())

	// Define collections
	managerCollection := client.Database("mydatabase").Collection("managers")
	userCollection := client.Database("mydatabase").Collection("users")

	// Check if a manager with the same username already exists
	var existingManager models.Manager
	err = managerCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&existingManager)
	if err == nil {
		// Manager with the same username exists
		return models.CreateManagerResponse{
			Success: false,
			Message: fmt.Sprintf("manager with username '%s' already exists", username),
		}
	} else if err != mongo.ErrNoDocuments {
		// Error while checking for existing manager
		return models.CreateManagerResponse{
			Success: false,
			Message: fmt.Sprintf("error checking for existing manager: %v", err),
		}
	}

	// Insert only the username and group limit into the managers collection
	manager := models.Manager{
		Username:   username,
		GroupLimit: groupLimit,
	}
	_, err = managerCollection.InsertOne(context.TODO(), manager)
	if err != nil {
		// Error inserting manager
		return models.CreateManagerResponse{
			Success: false,
			Message: fmt.Sprintf("could not insert manager: %v", err),
		}
	}

	// Insert username, password, and "manager" tag into the users collection
	user := models.User{
		Username: username,
		Password: password,
		Tag:      "manager",
	}
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		// Error inserting user
		return models.CreateManagerResponse{
			Success: false,
			Message: fmt.Sprintf("manager created, but could not add user: %v", err),
		}
	}

	// Return success response
	return models.CreateManagerResponse{
		Success: true,
		Message: "manager created successfully",
	}
}

func RemoveManager(username string) (bool, string) {
	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return false, fmt.Sprintf("could not connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Define collections
	managerCollection := client.Database("mydatabase").Collection("managers")
	userCollection := client.Database("mydatabase").Collection("users")

	// Check if the manager exists in the managers collection
	var existingManager models.Manager
	err = managerCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&existingManager)
	if err == mongo.ErrNoDocuments {
		return false, fmt.Sprintf("manager with username '%s' does not exist", username)
	} else if err != nil {
		return false, fmt.Sprintf("error checking for existing manager: %v", err)
	}

	// Remove the manager from the managers collection
	_, err = managerCollection.DeleteOne(context.TODO(), bson.M{"username": username})
	if err != nil {
		return false, fmt.Sprintf("could not remove manager: %v", err)
	}

	// Remove the user from the users collection (the user tag is "manager")
	_, err = userCollection.DeleteOne(context.TODO(), bson.M{"username": username, "tag": "manager"})
	if err != nil {
		return false, fmt.Sprintf("could not remove user: %v", err)
	}

	return true, fmt.Sprintf("manager with username '%s' removed successfully", username)
}