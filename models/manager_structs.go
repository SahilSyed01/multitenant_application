package models

// Structs for managers and groups
type Group struct {
    Manager   string   `json:"manager" bson:"manager"`
    GroupName string   `json:"group_name" bson:"group_name"`
    Members   []string `json:"members" bson:"members"`
    Budget    float64  `json:"budget" bson:"budget"` // Add Budget field here
}
// ManagerResponse standardizes the response structure
type UserResponse struct {
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}