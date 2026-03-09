package routes

import (
	"tracktora-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// Setup takes the Fiber app and registers all the endpoints
func Setup(app *fiber.App) {
	// 1. Simple Health Check Route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "TrackTora Backend is up and running! 🚀",
		})
	})

	// Auth Routes
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/login", handlers.Login) // <-- Add this line!
}
