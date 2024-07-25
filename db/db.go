package db

import (
	"context"
	"log"
	"time"

	"github.com/IsraeliPS/GoMongo/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var DB *mongo.Database

func ConnectDatabase() {
    mongoURI := config.LoadEnv().MONGO_URI
    mongoDB := config.LoadEnv().MONGO_DB
    

    client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    DB = client.Database(mongoDB)
    log.Println("Connected to MongoDB!")
}