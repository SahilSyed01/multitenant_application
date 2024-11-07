package db

import (
	"context"
	"fmt"
	"log"
	"multitenant/config"
	"multitenant/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func GetManagerCollection() *mongo.Collection {
	return Client.Database("mydatabase").Collection("managers")
}

func GetGroupsCollection() *mongo.Collection {
	return Client.Database("groups").Collection("groupnames")
}

func GetUsersCollection() *mongo.Collection {
	return Client.Database("mydatabase").Collection("users")
}

//var Client *mongo.Client

// Initialize MongoDB connection
func InitMongoDB() {
	var err error
	Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
}

// func DisconnectMongoDB() error {
// 	if Client != nil {
// 		return Client.Disconnect(context.Background())
// 	}
// 	return nil
// }

// Create user logic
func CreateUser(username, password string) error {
	// Check if user already exists
	var existingUser bson.M
	err := GetUsersCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return fmt.Errorf("user already exists")
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	// Insert new user
	user := bson.M{
		"username": username,
		"password": password,
		"Tag":      "user",
	}
	_, err = GetUsersCollection().InsertOne(context.Background(), user)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

// Create group logic
func CreateGroup(username, groupName string) error {
	var manager models.Manager

	// Check if manager exists
	err := GetManagerCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&manager)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("manager not found")
		}
		return fmt.Errorf("error finding manager: %v", err)
	}

	// Check if group with the same name already exists for the manager
	var existingGroup bson.M
	err = GetGroupsCollection().FindOne(context.Background(), bson.M{"manager": username, "group_name": groupName}).Decode(&existingGroup)
	if err == nil {
		return fmt.Errorf("group with this name already exists for the manager")
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking existing groups: %v", err)
	}

	// Check if group limit has been reached
	groupCount, err := GetGroupsCollection().CountDocuments(context.Background(), bson.M{"manager": username})
	if err != nil {
		return fmt.Errorf("error counting groups: %v", err)
	}

	if groupCount >= int64(manager.GroupLimit) {
		return fmt.Errorf("cannot create more groups, group limit reached")
	}

	// Create new group
	group := bson.M{
		"manager":    username,
		"group_name": groupName,
		"members":    []string{},
	}

	_, err = GetGroupsCollection().InsertOne(context.Background(), group)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	return nil
}

// Add user to group logic
func AddUserToGroup(manager, groupName, username string) error {
	// Check if user exists in the users collection
	var user bson.M
	err := GetUsersCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("user does not exist")
	} else if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	// Check if the user is already a member of any group
	filter := bson.M{
		"members": username,
	}
	var existingGroup bson.M
	err = GetGroupsCollection().FindOne(context.Background(), filter).Decode(&existingGroup)
	if err == nil {
		return fmt.Errorf("user is already a member of another group and cannot be added to multiple groups")
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking user's group membership: %v", err)
	}

	// Add user to the specified group
	update := bson.M{"$push": bson.M{"members": username}}
	_, err = GetGroupsCollection().UpdateOne(context.Background(), bson.M{
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
	_, err := GetGroupsCollection().UpdateOne(context.Background(), bson.M{
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
	_, err := GetUsersCollection().DeleteOne(context.Background(), bson.M{"username": username})
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
