package main

import (
	"log"

	"github.com/learnonid/learnonbe/config"
	"github.com/learnonid/learnonbe/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Inisialisasi koneksi ke database
	config.Init()

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new Fiber app
	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Auth",
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	app.Use(logger.New(logger.Config{
		Format: "${status} - ${method} ${path}\n",
	}))

	// Save the database connection in the app
	app.Use(func(c *fiber.Ctx) error {
		config.Init()
		return c.Next()
	})

	// Routes
	routes.SetupRoutes(app)

	// Listen to port 3000
	log.Fatal(app.Listen(":3000"))
}
