package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    MONGO_URI     string
    MONGO_DB    string
    JWT_SECRET string
    PORT string
}

var EnvConfig *Config

func LoadEnv() *Config {
    if EnvConfig != nil {
        return EnvConfig
    }

    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    EnvConfig = &Config{
        MONGO_URI:     os.Getenv("MONGO_URI"),
        MONGO_DB:    os.Getenv("MONGO_DB"),
        JWT_SECRET: os.Getenv("JWT_SECRET"),
        PORT: os.Getenv("PORT"),
    }

    return EnvConfig
}

