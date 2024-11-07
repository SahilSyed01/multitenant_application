package db

import (
	"context"
	"multitenant/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthenticateUser checks if the user exists with the correct credentials and returns the tag
func AuthenticateUser(username, password string) (bool, string, error) {
	collection := client.Database("mydatabase").Collection("users")

	// Check if the user exists with the given username and password
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"username": username, "password": password}).Decode(&user)
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
