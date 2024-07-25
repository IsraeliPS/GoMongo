package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IsraeliPS/GoMongo/config"
	"github.com/IsraeliPS/GoMongo/handlers"
	"github.com/IsraeliPS/GoMongo/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    router := mux.NewRouter()
    
    // Define routes
    router.HandleFunc("/users", handlers.GetUsers).Methods("GET")
    router.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
    router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
    router.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
    router.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
    router.HandleFunc("/login", handlers.Login).Methods("POST")

    // Apply middleware
    router.Use(middleware.ErrorHandler)
    
     // MongoDB connection
     config.ConnectDatabase()

     // Set up CORS
    handler := cors.Default().Handler(router)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    fmt.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":"+port, handler))
}