package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/IsraeliPS/GoMongo/db"
	"github.com/IsraeliPS/GoMongo/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// GetUsers godoc
// @Summary Get all users
// @Description Get all users with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
    var users []models.User
    collection := db.DB.Collection("users")
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        panic(err)
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            panic(err)
        }
        users = append(users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// GetUser godoc
// @Summary Get a user
// @Description Get a user by ID
// @Tags users
// @Security Bearer
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        panic("Invalid ID")
    }

    var user models.User
    collection := db.DB.Collection("users")
    err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        } else {
            panic(err)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}


// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Security Bearer
// @Accept  json
// @Produce  json
// @Param   user     body    models.User     true        "User data"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate input
    if err := validate.Struct(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user.ID = primitive.NewObjectID()
    collection := db.DB.Collection("users")
    // result, err := collection.InsertOne(context.Background(), user)
    _, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        panic(err)
    }

    userResponse := user
    userResponse.Password = "" 

    response := map[string]interface{}{
        "message": "User created successfully",
        "user":    userResponse,
    }

    w.Header().Set("Content-Type", "application/json")
    // json.NewEncoder(w).Encode(result.InsertedID)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

// UpdateUser godoc
// @Summary  Update an existing user
// @Description Update a user's information
// @Tags users
// @Security Bearer
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "User ID"
// @Param   user     body    models.User     true        "User data"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate input
    if err := validate.Struct(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user.ID = id

    collection := db.DB.Collection("users")
    result, err := collection.ReplaceOne(context.Background(), bson.M{"_id": id}, user)
    if err != nil {
        panic(err)
    }
    if result.MatchedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    userResponse := user
    userResponse.Password = ""

    response := map[string]interface{}{
        "message": "User updated successfully",
        "user":    userResponse,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Security Bearer
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    collection := db.DB.Collection("users")
    result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
    if err != nil {
        panic(err)
    }
    if result.DeletedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    response := map[string]string{"message": "User deleted successfully"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNoContent)
    json.NewEncoder(w).Encode(response)
}

// Login godoc
// @Summary Login 
// @Description Login a user and get a token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user     body    models.Credentials     true        "User credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
    var creds models.Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate input
    if err := validate.Struct(creds); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var user models.User
    collection := db.DB.Collection("users")
    err := collection.FindOne(context.Background(), bson.M{"email": creds.Email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        } else {
            panic(err)
        }
    }

    if user.Password != creds.Password {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &models.Claims{
        Email: creds.Email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        panic(err)
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   tokenString,
        Expires: expirationTime,
    })

    response:=map[string]interface{}{
        "message": "Login successful",
        "token": tokenString,
        "expires": expirationTime.Local(),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
