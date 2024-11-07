package routes

import (
	"multitenant/handlers"

	"github.com/gorilla/mux"
)

// InitializeRoutes initializes all the routes for the application
func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	// Define routes for login and other endpoints
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
    router.HandleFunc("/create-manager", handlers.CreateManagerHandler).Methods("POST")

	// Set up routes
	router.HandleFunc("/create_group", handlers.CreateGroupHandler).Methods("POST")
	router.HandleFunc("/add_user", handlers.AddUserHandler)
	router.HandleFunc("/remove_user", handlers.RemoveUserHandler)
	router.HandleFunc("/create_user", handlers.CreateUserHandler)

	return router
}
