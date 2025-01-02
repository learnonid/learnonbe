package config

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading environment")
	}
}

func InitDB() {
	// Load environment
	LoadEnv()

	// Load MongoDB connection string
	dbConf := os.Getenv("MONGOSTRING")
	if dbConf == "" {
		panic("MONGOSTRING is not set in environment variables")
	}

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(dbConf)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		panic("Failed to create MongoDB client: " + err.Error())
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	// Ping MongoDB to ensure connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic("Failed to ping MongoDB: " + err.Error())
	}

	// Set the database instance
	db = client.Database("nama_database") // Ganti dengan nama database yang sesuai
}

func GetDB() *mongo.Database {
	return db
}