package db

import (
    "context"
    "fmt"
    "multitenant/config"
    "multitenant/models"
    "strings"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// AddManager adds a new manager and user to the respective collections
func AddManager(username, password string, groupLimit int) models.ManagerResponse {
    // Input validation with whitespace trimming
    if strings.TrimSpace(username) == "" {
        return models.ManagerResponse{
            Success: false,
            Message: "username cannot be empty",
        }
    }
    if strings.TrimSpace(password) == "" {
        return models.ManagerResponse{
            Success: false,
            Message: "password cannot be empty",
        }
    }
    if groupLimit <= 0 {
        return models.ManagerResponse{
            Success: false,
            Message: "group limit must be greater than zero",
        }
    }

    clientOptions := options.Client().ApplyURI(config.MongoURI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return models.ManagerResponse{
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
    err = managerCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&existingManager)
    if err == nil {
        return models.ManagerResponse{
            Success: false,
            Message: fmt.Sprintf("manager with username '%s' already exists", username),
        }
    } else if err != mongo.ErrNoDocuments {
        return models.ManagerResponse{
            Success: false,
            Message: fmt.Sprintf("error checking for existing manager: %v", err),
        }
    }

    // Insert only the username and group limit into the managers collection
    manager := models.Manager{
        Username:   username,
        GroupLimit: groupLimit,
    }
    _, err = managerCollection.InsertOne(context.Background(), manager)
    if err != nil {
        return models.ManagerResponse{
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
    _, err = userCollection.InsertOne(context.Background(), user)
    if err != nil {
        return models.ManagerResponse{
            Success: false,
            Message: fmt.Sprintf("manager created, but could not add user: %v", err),
        }
    }

    return models.ManagerResponse{
        Success: true,
        Message: "manager created successfully",
    }
}

// RemoveManager removes a manager and corresponding user from the collections
func RemoveManager(username string) models.ManagerResponse {
    // Input validation with whitespace trimming
    if strings.TrimSpace(username) == "" {
        return models.ManagerResponse{
            Success: false,
            Message: "username cannot be empty",
        }
    }

    clientOptions := options.Client().ApplyURI(config.MongoURI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return models.ManagerResponse{
            Success: false,
            Message: fmt.Sprintf("could not connect to MongoDB: %v", err),
        }
    }
    defer client.Disconnect(context.TODO())

    // Define collections
    managerCollection := client.Database("mydatabase").Collection("managers")
    userCollection := client.Database("mydatabase").Collection("users")

    // Check if the manager exists in the managers collection
    var existingManager models.Manager
    err = managerCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&existingManager)
    if err == mongo.ErrNoDocuments {
        return models.ManagerResponse{
            Success: false,
            Message: fmt.Sprintf("manager with username '%s' does not exist", username),
        }
    } else if err != nil {
        return models.ManagerResponse{
            Success: false,
            Message: fmt.Sprintf("error checking for existing manager: %v", err),
        }
    }

    // Remove the manager from the managers collection
    _, err = managerCollection.DeleteOne(context.Background(), bson.M{"username": username})
    if err != nil {
        return models.ManagerResponse{
            Success: false,
            Message: fmt.Sprintf("could not remove manager: %v", err),
        }
    }

    // Remove the user from the users collection (the user tag is "manager")
    _, err = userCollection.DeleteOne(context.Background(), bson.M{"username": username, "tag": "manager"})
    if err != nil {
        return models.ManagerResponse{
            Success: false,
            Message: fmt.Sprintf("could not remove user: %v", err),
        }
    }
    return models.ManagerResponse{
        Success: true,
        Message: fmt.Sprintf("manager with username '%s' removed successfully", username),
    }
}
