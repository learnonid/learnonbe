package routes

import (
	"learnonbe/controller"
	"learnonbe/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Create a new Fiber app
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to LearnOn.id")
	})

	// Auth routes
	AuthRoutes := app.Group("/auth")
	AuthRoutes.Post("/register", controller.RegisterAkun)
	AuthRoutes.Post("/login", controller.Login)

	// Middleware
	ProtectedRoutes := app.Group("/u")
	ProtectedRoutes.Use(middleware.JWTMiddleware("secret"))
	ProtectedRoutes.Get("/profile", controller.GetProfile)
	// ProtectedRoutes.Get("/user", controller.GetUser)
}