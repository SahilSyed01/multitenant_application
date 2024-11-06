package handlers

import (
	"encoding/json"
	"fmt"
	"multitenant/db"
	"multitenant/models"
	"net/http"
)

// LoginHandler handles the login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body into LoginRequest
	var loginRequest models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Authenticate the user using the provided credentials and fetch the user's tag
	isAuthenticated, tag, err := db.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during authentication: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	var response models.LoginResponse
	if isAuthenticated {
		// Authentication succeeded, now route based on the tag
		response = models.LoginResponse{
			Success: true,
			Message: "Login successful",
		}

		// Send the login success response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to send response", http.StatusInternalServerError)
			return
		}

		// Route to the appropriate page based on the tag
		switch tag {
		case "admin":
			// Route to admin page
			//http.Redirect(w, r, "/admin-dashboard", http.StatusSeeOther)
           // fmt.Print("got admin tag")
			return
		case "manager":
			// Route to manager page
			//http.Redirect(w, r, "/manager-dashboard", http.StatusSeeOther)
			return
		case "user":
			// Route to user page
			//http.Redirect(w, r, "/user-dashboard", http.StatusSeeOther)
			return
		default:
			// If the tag is invalid, deny access
			response = models.LoginResponse{
				Success: false,
				Message: "Invalid user tag",
			}
			// Send the error response
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, "Failed to send response", http.StatusInternalServerError)
			}
		}
	} else {
		// Invalid credentials
		response = models.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}

		// Send the error response
		w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to send response", http.StatusInternalServerError)
		}
	}
}
