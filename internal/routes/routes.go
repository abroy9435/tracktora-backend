package routes

import (
	"tracktora-backend/internal/handlers"
	"tracktora-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Setup takes the Fiber app and registers all the endpoints
func Setup(app *fiber.App) {
	// 1. Public Routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "TrackTora Backend is up and running! 🚀",
		})
	})

	// 2. Auth Routes (Public)
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/login", handlers.Login)

	// 3. Protected Routes (Requires JWT)
	protectedGroup := app.Group("/api", middleware.RequireAuth)

	// Explicit Application Routes
	protectedGroup.Post("/applications/add", handlers.CreateApplication)
	protectedGroup.Get("/applications/list", handlers.GetApplications)
}
