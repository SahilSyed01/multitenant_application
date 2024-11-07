package handlers

import (
	"encoding/json"
	"multitenant/db"
	"net/http"
	"strings"
)


// CreateUserHandler handles the creation of a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	response:= db.CreateUser(input.Username, input.Password)
	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfully"))

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// CreateGroupHandler handles the creation of a new group
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string `json:"username"`
		GroupName string `json:"group_name"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Trim whitespace and check if Username or GroupName is empty
	if strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.GroupName) == "" {
		http.Error(w, "Username and group name cannot be empty or whitespace", http.StatusBadRequest)
		return
	}

	// Attempt to create the group
	response:= db.CreateGroup(input.Username, input.GroupName)
	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Group created successfully"))

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// AddUserHandler adds an existing user to a group
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string `json:"username"`
		GroupName string `json:"group_name"`
		User      string `json:"user"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	response:= db.AddUserToGroup(input.Username, input.GroupName, input.User)
	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User added to group successfully"))

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// RemoveUserHandler removes a user from a group
func RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string `json:"username"`
		GroupName string `json:"group_name"`
		User      string `json:"user"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	response:= db.RemoveUserFromGroup(input.Username, input.GroupName, input.User)
	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User removed from group successfully"))

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// DeleteUserHandler handles the deletion of a user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
	}

	// Decode the JSON request body to get the username
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Username == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the DeleteUser function to delete the user
	response:= db.DeleteUser(input.Username)
	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func AddBudgetHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GroupName string  `json:"group_name"`
		Budget    float64 `json:"budget"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input
	if strings.TrimSpace(input.GroupName) == "" {
		http.Error(w, "Group name cannot be empty or whitespace", http.StatusBadRequest)
		return
	}

	if input.Budget <= 0 {
		http.Error(w, "Budget must be greater than zero", http.StatusBadRequest)
		return
	}

	// Call the AddBudget function
	response:= db.AddBudget(input.GroupName, input.Budget)
	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Budget added to group successfully"))

	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func ListGroupsHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Username string `json:"username"`
    }
 
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil || input.Username == "" {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
 
    response:=db.ListGroupsByManager(input.Username)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
 
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}