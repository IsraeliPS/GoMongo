package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IsraeliPS/GoMongo/config"
	"github.com/IsraeliPS/GoMongo/db"
	_ "github.com/IsraeliPS/GoMongo/docs"
	"github.com/IsraeliPS/GoMongo/handlers"
	"github.com/IsraeliPS/GoMongo/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Proyect Go MongoDB API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Bearer token for authorization

func main() {
    config.LoadEnv()
    router := mux.NewRouter()
    
    // Define routes
    api:= router.PathPrefix("/api").Subrouter()
    api.HandleFunc("/users", handlers.GetUsers).Methods("GET")
    api.HandleFunc("/login", handlers.Login).Methods("POST")



   // Protected routes
    protectedRoutes := api.PathPrefix("").Subrouter()
    protectedRoutes.Use(middleware.JWTAuthentication)
    protectedRoutes.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
    protectedRoutes.HandleFunc("/users", handlers.CreateUser).Methods("POST")
    protectedRoutes.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
    protectedRoutes.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

    // Apply middleware
    router.Use(middleware.ErrorHandler)
    
     // MongoDB connection
     db.ConnectDatabase()

     // Set up CORS
    handler := cors.Default().Handler(router)

    // Swagger documentation
    router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
    log.Println("Swagger started at http://localhost:8080/swagger/index.html")

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    fmt.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":"+port, handler))
}