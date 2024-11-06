package models

// User represents a user document in MongoDB
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Tag      string `bson:"tag"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Tag      string `bson:"tag"`
}

// LoginResponse represents the structure of the login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
