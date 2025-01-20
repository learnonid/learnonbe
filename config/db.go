package config

import (
	"log"
	"os"
	"context"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func Init() {
	// Only load .env file if running locally (not on Heroku)
	if os.Getenv("DYNO") == "" {
		// Load .env for local environment
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Get MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Set the global MongoClient
	MongoClient = client
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}
