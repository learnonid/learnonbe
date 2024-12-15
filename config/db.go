package config

import (
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading environtment")
	}
}

func GetDB() *gorm.DB {

	// Load environtment
	LoadEnv()

	// Load string connection
	dbConf := os.Getenv("SQLSTRING")

	// Create connection to database
	DB, err := gorm.Open(mysql.Open(dbConf), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("Failed to connect to database! " + err.Error())
	}

	// Set connection to database
	db0, err := DB.DB()
	if err != nil {
		panic("Failed to connect to database! " + err.Error())
	}
	db0.SetConnMaxIdleTime(time.Duration(1) * time.Minute)
	db0.SetConnMaxLifetime(time.Duration(1) * time.Minute)
	db0.SetMaxIdleConns(2)

	// Show log
	DB.Statement.RaiseErrorOnNotFound = true // Raise error on not found

	return DB
}