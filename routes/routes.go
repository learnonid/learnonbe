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

	// Roles routes
	RolesRoutes := app.Group("/roles")
	RolesRoutes.Post("/create", controller.CreateRole)

	// Auth routes
	AuthRoutes := app.Group("/auth")
	AuthRoutes.Post("/register", controller.RegisterAkun)
	AuthRoutes.Post("/login", controller.Login)
	AuthRoutes.Post("/login/admin", controller.LoginAdmin)

	// Middleware
	ProtectedRoutes := app.Group("/u")
	ProtectedRoutes.Use(middleware.JWTMiddleware("secret"))
	ProtectedRoutes.Get("/profile", controller.GetProfile)
	// ProtectedRoutes.Get("/user", controller.GetUser)

	// Event routes
	EventRoutes := app.Group("/event")
	EventRoutes.Post("/create", controller.CreateEvent)
	EventRoutes.Post("/upload-image", controller.UploadEventImageHandler)
	EventRoutes.Get("/all", controller.GetEvents)

	// Admin routes
	AdminRoutes := app.Group("/admin")
	AdminRoutes.Use(middleware.JWTMiddleware("secret"))
	AdminRoutes.Put("/update-user/:id", controller.UpdateUser)
	AdminRoutes.Delete("/delete-user/:id", controller.DeleteUser)

	// Static file
	app.Static("/uploads", "./uploads")
}