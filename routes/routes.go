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
	router.HandleFunc("/delete-manager", handlers.RemoveManagerHandler).Methods("DELETE")

	return router
}
