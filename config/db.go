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
	// Hanya muat file .env jika aplikasi berjalan secara lokal
	if os.Getenv("HEROKU") == "" {
		// Muat file .env untuk lingkungan lokal
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Dapatkan MongoDB URI dari variabel lingkungan
	mongoURI := os.Getenv("MONGO_URI")

	// Hubungkan ke MongoDB
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
