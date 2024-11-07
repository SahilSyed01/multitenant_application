package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"multitenant/config"
)

var Client *mongo.Client

// Initialize MongoDB connection
func InitMongoDB() {
	var err error
	Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
}

func GetManagerCollection() *mongo.Collection {
	return Client.Database("mydatabase").Collection("managers")
}

func GetGroupsCollection() *mongo.Collection {
	return Client.Database("groups").Collection("groupnames") 
}

func GetUsersCollection() *mongo.Collection {
	return Client.Database("mydatabase").Collection("users") 
}

func DisconnectMongoDB() error {
	if Client != nil {
		return Client.Disconnect(context.Background())
	}
	return nil
}

// Structs for managers and groups
type Group struct {
	Manager   string   `bson:"manager"`
	GroupName string   `bson:"group_name"`
	Members   []string `bson:"members"`
}

// type User struct {
// 	Username string `bson:"username"`
// 	Password string `bson:"password"` // In a real app, consider hashing passwords
// }
