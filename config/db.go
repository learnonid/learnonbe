package config

import (
	"log"
	"os"
	"context"

	"github.com/joho/godotenv"
	// "github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func Init() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MongoDB URI from .env
	mongoURI := os.Getenv("MONGO_URI")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Set global MongoClient
	MongoClient = client
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}