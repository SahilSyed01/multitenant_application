package models

// Manager represents the structure of a manager document in the "managers" collection
type Manager struct {
	Username   string `bson:"username"`
	GroupLimit int    `bson:"group_limit"`
}

// CreateManagerRequest represents the input required to create a new manager
type CreateManagerRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	GroupLimit int    `json:"group_limit"`
}

// CreateManagerResponse represents the structure of the response after creating a manager
type CreateManagerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
