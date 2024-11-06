package models

type Manager struct {
	Username   string `bson:"username"`
	GroupLimit int    `bson:"group_limit"`
}
type CreateManagerRequest struct {
	Username   string `json:"username"`
	GroupLimit int    `json:"group_limit"`
}

// CreateManagerResponse represents the structure of the response body
type CreateManagerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
