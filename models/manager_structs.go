package models

// Structs for managers and groups
type Group struct {
	Manager   string   `bson:"manager"`
	GroupName string   `bson:"group_name"`
	Members   []string `bson:"members"`
}
