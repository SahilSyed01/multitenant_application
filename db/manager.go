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
	return Client.Database("mydatabase").Collection("groups")
}

func GetUsersCollection() *mongo.Collection {
	return Client.Database("mydatabase").Collection("users")
}

// Initialize MongoDB connection
func InitMongoDB() {
	var err error
	Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
}

// Create user logic
func CreateUser(username, password string) models.UserResponse  {
	var existingUser bson.M
	err := GetUsersCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return models.UserResponse {
			Message: "user already exists",
			Status:  "error",
		}
	} else if err != mongo.ErrNoDocuments {
		return models.UserResponse {
			Message: fmt.Sprintf("error checking user existence: %v", err),
			Status:  "error",
		}
	}

	user := bson.M{
		"username": username,
		"password": password,
		"Tag":      "user",
	}
	_, err = GetUsersCollection().InsertOne(context.Background(), user)
	if err != nil {
		return models.UserResponse {
			Message: fmt.Sprintf("error creating user: %v", err),
			Status:  "error",
		}
	}

	return models.UserResponse {Message: "User created successfully"}
}

// Create group logic
func CreateGroup(username, groupName string) models.UserResponse  {
	var manager models.Manager

	err := GetManagerCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&manager)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.UserResponse {
				Message: "manager not found",
				Status:  "error",
			}
		}
		return models.UserResponse {
			Message: fmt.Sprintf("error finding manager: %v", err),
			Status:  "error",
		}
	}

	var existingGroup bson.M
	err = GetGroupsCollection().FindOne(context.Background(), bson.M{"manager": username, "group_name": groupName}).Decode(&existingGroup)
	if err == nil {
		return models.UserResponse {
			Message: "group with this name already exists for the manager",
			Status:  "error",
		}
	} else if err != mongo.ErrNoDocuments {
		return models.UserResponse {
			Message: fmt.Sprintf("error checking existing groups: %v", err),
			Status:  "error",
		}
	}

	groupCount, err := GetGroupsCollection().CountDocuments(context.Background(), bson.M{"manager": username})
	if err != nil {
		return models.UserResponse {
			Message: fmt.Sprintf("error counting groups: %v", err),
			Status:  "error",
		}
	}

	if groupCount >= int64(manager.GroupLimit) {
		return models.UserResponse {
			Message: "cannot create more groups, group limit reached",
			Status:  "error",
		}
	}

	group := bson.M{
		"manager":    username,
		"group_name": groupName,
		"members":    []string{},
	}

	_, err = GetGroupsCollection().InsertOne(context.Background(), group)
	if err != nil {
		return models.UserResponse {
			Message: fmt.Sprintf("error creating group: %v", err),
			Status:  "error",
		}
	}

	return models.UserResponse {Message: "Group created successfully"}
}

// Add user to group logic
func AddUserToGroup(manager, groupName, username string) models.UserResponse  {
	var user bson.M
	err := GetUsersCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return models.UserResponse {
			Message: "user does not exist",
			Status:  "error",
		}
	} else if err != nil {
		return models.UserResponse {
			Message: fmt.Sprintf("error checking user existence: %v", err),
			Status:  "error",
		}
	}

	filter := bson.M{"members": username}
	var existingGroup bson.M
	err = GetGroupsCollection().FindOne(context.Background(), filter).Decode(&existingGroup)
	if err == nil {
		return models.UserResponse {
			Message: "user is already a member of another group and cannot be added to multiple groups",
			Status:  "error",
		}
	} else if err != mongo.ErrNoDocuments {
		return models.UserResponse {
			Message: fmt.Sprintf("error checking user's group membership: %v", err),
			Status:  "error",
		}
	}

	update := bson.M{"$push": bson.M{"members": username}}
	_, err = GetGroupsCollection().UpdateOne(context.Background(), bson.M{
		"manager":    manager,
		"group_name": groupName,
	}, update)
	if err != nil {
		return models.UserResponse {
			Message: fmt.Sprintf("error adding user to group: %v", err),
			Status:  "error",
		}
	}

	return models.UserResponse {Message: "User added to group successfully"}
}

// Remove user from group logic
func RemoveUserFromGroup(manager, groupName, username string) models.UserResponse  {
	update := bson.M{"$pull": bson.M{"members": username}}
	_, err := GetGroupsCollection().UpdateOne(context.Background(), bson.M{
		"manager":    manager,
		"group_name": groupName,
	}, update)
	if err != nil {
		return models.UserResponse {
			Message: fmt.Sprintf("error removing user from group: %v", err),
			Status:  "error",
		}
	}

	return models.UserResponse {Message: "User removed from group successfully"}
}

// DeleteUser deletes a user from the "users" collection
func DeleteUser(username string) models.UserResponse  {
	_, err := GetUsersCollection().DeleteOne(context.Background(), bson.M{"username": username})
	if err != nil {
		return models.UserResponse {
			Message: fmt.Sprintf("error deleting user: %v", err),
			Status:  "error",
		}
	}
	return models.UserResponse {Message: "User deleted successfully"}
}
func ListGroupsByManager(manager string) models.UserResponse {
    filter := bson.M{"manager": manager}
    
    // Check if the manager exists
    count, err := GetGroupsCollection().CountDocuments(context.Background(), filter)
    if err != nil {
        return models.UserResponse{
            Message: fmt.Sprintf("error checking manager existence: %v", err),
            Status:  "error",
        }
    }
    
    if count == 0 {
        return models.UserResponse{
            Message: fmt.Sprintf("manager '%s' not found", manager),
            Status:  "error",
        }
    }

    // Retrieve groups if the manager exists
    cursor, err := GetGroupsCollection().Find(context.Background(), filter)
    if err != nil {
        return models.UserResponse{
            Message: fmt.Sprintf("error retrieving groups: %v", err),
            Status:  "error",
        }
    }
    defer cursor.Close(context.Background())

    var groups []models.Group
    if err = cursor.All(context.Background(), &groups); err != nil {
        return models.UserResponse{
            Message: fmt.Sprintf("error decoding groups: %v", err),
            Status:  "error",
        }
    }

    return models.UserResponse{
        Message: "Groups retrieved successfully",
        Status:  "success",
        Data:    groups,
    }
}

func AddBudget(manager, groupName string, budget float64) models.UserResponse {
    if budget <= 0 {
        return models.UserResponse{
            Message: "budget must be greater than zero",
            Status:  "error",
        }
    }

    // Check if the group exists for the given manager
    filter := bson.M{"manager": manager, "group_name": groupName}
    var group models.Group
    err := GetGroupsCollection().FindOne(context.Background(), filter).Decode(&group)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return models.UserResponse{
                Message: fmt.Sprintf("group '%s' not found for manager '%s'", groupName, manager),
                Status:  "error",
            }
        }
        return models.UserResponse{
            Message: fmt.Sprintf("error finding group: %v", err),
            Status:  "error",
        }
    }

    // Check if the budget is already assigned to the group
    if group.Budget > 0 {
        return models.UserResponse{
            Message: fmt.Sprintf("budget is already allocated for the group '%s'", groupName),
            Status:  "error",
        }
    }

    // Assign the new budget to the group
    update := bson.M{"$set": bson.M{"budget": budget}}
    result, err := GetGroupsCollection().UpdateOne(context.Background(), filter, update)
    if err != nil {
        return models.UserResponse{
            Message: fmt.Sprintf("error updating budget: %v", err),
            Status:  "error",
        }
    }

    if result.MatchedCount == 0 {
        return models.UserResponse{
            Message: fmt.Sprintf("group '%s' not found for manager '%s'", groupName, manager),
            Status:  "error",
        }
    }

    return models.UserResponse{
        Message: fmt.Sprintf("Budget successfully allocated to group '%s'", groupName),
        Status:  "success",
    }
}

func UpdateBudget(manager, groupName string, budget float64) models.UserResponse {
    if budget <= 0 {
        return models.UserResponse{
            Message: "Budget must be greater than zero",
            Status:  "error",
        }
    }

    // Check if the group exists for the given manager
    filter := bson.M{"manager": manager, "group_name": groupName}
    var group models.Group
    err := GetGroupsCollection().FindOne(context.Background(), filter).Decode(&group)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return models.UserResponse{
                Message: fmt.Sprintf("Group '%s' not found for manager '%s'", groupName, manager),
                Status:  "error",
            }
        }
        return models.UserResponse{
            Message: fmt.Sprintf("Error finding group: %v", err),
            Status:  "error",
        }
    }

    // Check if the group already has a budget assigned
    if group.Budget == 0 {
        return models.UserResponse{
            Message: fmt.Sprintf("No budget allocated for group '%s'. Cannot update budget.", groupName),
            Status:  "error",
        }
    }

    // Update the existing budget in the group
    update := bson.M{"$set": bson.M{"budget": budget}}
    result, err := GetGroupsCollection().UpdateOne(context.Background(), filter, update)
    if err != nil {
        return models.UserResponse{
            Message: fmt.Sprintf("Error updating budget: %v", err),
            Status:  "error",
        }
    }

    if result.MatchedCount == 0 {
        return models.UserResponse{
            Message: fmt.Sprintf("Group '%s' not found for manager '%s'", groupName, manager),
            Status:  "error",
        }
    }

    return models.UserResponse{
        Message: fmt.Sprintf("Budget successfully updated for group '%s'", groupName),
        Status:  "success",
    }
}