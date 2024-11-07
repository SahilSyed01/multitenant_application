package main

import (
	"fmt"
	"log"
	"multitenant/db"
	"multitenant/routes"
	"net/http"
	"github.com/rs/cors"
	"multitenant/models"
)

func main() {

	models.InitMongoDB()
	// Initialize MongoDB connection
	err := db.ConnectMongoDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer db.DisconnectMongoDB()

	// Initialize routes
	router := routes.InitializeRoutes()

	// Setup CORS with the allowed origin for Angular
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"}, // Allow requests from Angular
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Wrap the router with CORS middleware
	handler := c.Handler(router)

	// Start the server
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
