package models

// Structs for managers and groups
type Group struct {
	Manager   string   `bson:"manager"`
	GroupName string   `bson:"group_name"`
	Members   []string `bson:"members"`
}
// ManagerResponse standardizes the response structure
type UserResponse struct {
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}