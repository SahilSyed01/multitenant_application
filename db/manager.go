package db

import (
	"context"
	"fmt"
	"multitenant/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


// Create user logic
func CreateUser(username, password string) error {
	// Check if user already exists
	var existingUser bson.M
	err := models.GetUsersCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return fmt.Errorf("user already exists")
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	// Insert new user
	user := bson.M{
		"username": username,
		"password": password, 
		"Tag":"user",
	}
	_, err = models.GetUsersCollection().InsertOne(context.Background(), user)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

// Create group logic
func CreateGroup(username, groupName string) error {
	var manager models.Manager
	err := models.GetManagerCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&manager)
	if err != nil {
		return fmt.Errorf("manager not found: %v", err)
	}

	groupCount, err := models.GetGroupsCollection().CountDocuments(context.Background(), bson.M{"manager": username})
	if err != nil {
		return fmt.Errorf("error counting groups: %v", err)
	}

	if groupCount >= int64(manager.GroupLimit) {
		return fmt.Errorf("cannot create more groups, group limit reached")
	}

	group := bson.M{
		"manager":    username,
		"group_name": groupName,
		"members":    []string{},
	}

	_, err = models.GetGroupsCollection().InsertOne(context.Background(), group)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	return nil
}

// Add user to group logic
func AddUserToGroup(manager, groupName, username string) error {
	// Check if user exists in the users collection
	var user bson.M
	err := models.GetUsersCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("user does not exist")
	} else if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	// Add user to group
	update := bson.M{"$push": bson.M{"members": username}}
	_, err = models.GetGroupsCollection().UpdateOne(context.Background(), bson.M{
		"manager":    manager,
		"group_name": groupName,
	}, update)
	if err != nil {
		return fmt.Errorf("error adding user to group: %v", err)
	}

	return nil
}

// Remove user from group logic
func RemoveUserFromGroup(manager, groupName, username string) error {
	update := bson.M{"$pull": bson.M{"members": username}}
	_, err := models.GetGroupsCollection().UpdateOne(context.Background(), bson.M{
		"manager":    manager,
		"group_name": groupName,
	}, update)
	if err != nil {
		return fmt.Errorf("error removing user from group: %v", err)
	}

	return nil
}

// DeleteUser deletes a user from the "users" collection
func DeleteUser(username string) error {
	_, err := models.GetUsersCollection().DeleteOne(context.Background(), bson.M{"username": username})
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
