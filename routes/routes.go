package routes

import (
	"github.com/learnonid/learnonbe/controller"
	"github.com/learnonid/learnonbe/middleware"

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

	// User routes
	UserRoutes := app.Group("/user")
	UserRoutes.Get("/all", controller.GetAllUsers)
	UserRoutes.Get("/detail/:id", controller.GetUserByID)
	UserRoutes.Get("/detail/email/:email", controller.GetUserByEmail)
	UserRoutes.Put("/update/:id", controller.UpdateUser)
	UserRoutes.Delete("/delete/:id", controller.DeleteUser)

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
	EventRoutes.Get("/detail/:id", controller.GetEventByID)
	EventRoutes.Get("/type", controller.GetEventsByType)
	EventRoutes.Get("/type/online", controller.GetEventByTypeOnline)
	EventRoutes.Get("/type/offline", controller.GetEventByTypeOffline)
	EventRoutes.Put("/update/:id", controller.EditEvent)
	EventRoutes.Delete("/delete/:id", controller.DeleteEvent)

	// Book routes
	BookRoutes := app.Group("/book")
	BookRoutes.Post("/create", controller.CreateBook)
	BookRoutes.Get("/all", controller.GetBooks)
	BookRoutes.Get("/detail/:id", controller.GetBookByID)
	BookRoutes.Put("/update/:id", controller.EditBook)
	BookRoutes.Delete("/delete/:id", controller.DeleteBook)

	// Admin routes
	AdminRoutes := app.Group("/admin")
	AdminRoutes.Use(middleware.JWTMiddleware("secret"))
	AdminRoutes.Put("/update-user/:id", controller.UpdateUser)
	AdminRoutes.Delete("/delete-user/:id", controller.DeleteUser)

	// File upload routes
	FileRoutes := app.Group("/file")
	FileRoutes.Post("/upload", controller.PostUploadPayment)

	// Event registration routes
	UERRoutes := app.Group("/uer")
	UERRoutes.Get("/all", controller.GetAllUERegistration)
	UERRoutes.Get("/user/:id", controller.GetUERegistrationByUserID)
	UERRoutes.Get("/:id", controller.GetUERegistrationByID)
	UERRoutes.Put("/update/:id", controller.UpdateEventRegistration)

	// Book Payment routes
	BookPaymentRoutes := app.Group("/book")
	BookPaymentRoutes.Post("/pay", controller.PostBookPayment)
	BookPaymentRoutes.Get("/all/pay", controller.GetAllBookPayment)
	BookPaymentRoutes.Get("/pay/user/:id", controller.GetBookPaymentByUserID)
	BookPaymentRoutes.Get("/pay/:id", controller.GetBookPaymentByID)

}
