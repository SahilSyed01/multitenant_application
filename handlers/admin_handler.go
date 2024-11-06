package handlers

import (
	"encoding/json"
	"multitenant/db"
	"multitenant/models"
	"net/http"
)

// CreateManagerHandler handles the request to create a manager
func CreateManagerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request models.CreateManagerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate username and group limit
	if request.Username == "admin" {
		http.Error(w, "Username 'admin' is not allowed", http.StatusBadRequest)
		return
	}

	// Check if group limit is a positive integer
	if request.GroupLimit <= 0 {
		http.Error(w, "Group limit must be a positive integer", http.StatusBadRequest)
		return
	}

	// Call AddManager to add the manager
	success, message := db.AddManager(request.Username, request.GroupLimit)

	// Prepare the response
	response := models.CreateManagerResponse{
		Success: success,
		Message: message,
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
}
