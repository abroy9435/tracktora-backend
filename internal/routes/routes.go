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
	authGroup.Post("/forgot-password", handlers.ForgotPassword)
	authGroup.Post("/reset-password", handlers.ResetPassword)

	// 3. Protected Routes (Requires JWT)
	protectedGroup := app.Group("/api", middleware.RequireAuth)

	// Explicit Application Routes
	protectedGroup.Post("/applications/add", handlers.CreateApplication)
	protectedGroup.Get("/applications/list", handlers.GetApplications)
	protectedGroup.Get("/applications/stats", handlers.GetApplicationStats)
	protectedGroup.Put("/applications/update", handlers.UpdateApplication)
	protectedGroup.Delete("/applications/delete", handlers.DeleteApplication)

	// User Profile Routes
	protectedGroup.Get("/profile", handlers.GetProfile)
	protectedGroup.Put("/profile/update", handlers.UpdateProfile)
	protectedGroup.Put("/profile/privacy", handlers.UpdatePrivacySettings)

	// Explore Feed (Adzuna)
	protectedGroup.Get("/explore", handlers.GetExplorePage)
	protectedGroup.Post("/explore/save", handlers.SaveExternalJob)

	// Connection / Friends System
	protectedGroup.Post("/connect/invite", handlers.SendFriendRequest)       // User A invites User B
	protectedGroup.Put("/connect/respond", handlers.RespondToFriendRequest)  // User B accepts/rejects
	protectedGroup.Get("/connect/stats/:friend_id", handlers.GetFriendStats) // View mutual friend stats
	protectedGroup.Get("/connect/requests", handlers.GetPendingRequests)
	protectedGroup.Get("/connect/list", handlers.GetFriendList)
}
