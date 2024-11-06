package main

import (
    "fmt"
    "log"
    "net/http"

    "multitenant/db"
    "multitenant/routes"
)

func main() {
    // Initialize MongoDB connection
    err := db.ConnectMongoDB()
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    defer db.DisconnectMongoDB()

    // Initialize routes
    router := routes.InitializeRoutes()
    // Start server
    fmt.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
