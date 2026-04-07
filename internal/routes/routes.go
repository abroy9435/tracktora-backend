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
	authGroup.Post("/verify-email", handlers.VerifyEmail)
	authGroup.Post("/login", handlers.Login)
	authGroup.Post("/forgot-password", handlers.ForgotPassword)
	authGroup.Post("/reset-password", handlers.ResetPassword)
	authGroup.Post("/resend-verification", handlers.ResendVerification)

	// 3. Protected Routes (Requires JWT)
	protectedGroup := app.Group("/api", middleware.RequireAuth)

	// Explicit Application Routes
	protectedGroup.Put("/auth/update-password", handlers.UpdatePassword)
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
	protectedGroup.Get("/connect/search", handlers.SearchUsers)

	// --- NEW: RESUME BUILDER & MASTER VAULT ---
	resumeGroup := protectedGroup.Group("/resume")

	// Vault: Projects
	resumeGroup.Post("/vault/project", handlers.AddProject)
	resumeGroup.Get("/vault/project", handlers.GetProjects)
	resumeGroup.Put("/vault/project", handlers.UpdateProject)
	resumeGroup.Delete("/vault/project", handlers.DeleteProject)

	// Vault: Experiences
	resumeGroup.Post("/vault/experience", handlers.AddExperience)
	resumeGroup.Get("/vault/experience", handlers.GetExperiences)
	resumeGroup.Put("/vault/experience", handlers.UpdateExperience)
	resumeGroup.Delete("/vault/experience", handlers.DeleteExperience)

	// Vault: Education
	resumeGroup.Post("/vault/education", handlers.AddEducation)
	resumeGroup.Get("/vault/education", handlers.GetEducations)
	resumeGroup.Put("/vault/education", handlers.UpdateEducation)
	resumeGroup.Delete("/vault/education", handlers.DeleteEducation)

	// Vault: Skills
	resumeGroup.Post("/vault/skill", handlers.AddSkill)
	resumeGroup.Get("/vault/skill", handlers.GetSkills)
	resumeGroup.Put("/vault/skill", handlers.UpdateSkill)
	resumeGroup.Delete("/vault/skill", handlers.DeleteSkill)

	// Vault: Certifications (NEW TWEAK)
	resumeGroup.Post("/vault/certification", handlers.AddCertification)
	resumeGroup.Get("/vault/certification", handlers.GetCertifications)
	resumeGroup.Put("/vault/certification", handlers.UpdateCertification)
	resumeGroup.Delete("/vault/certification", handlers.DeleteCertification)

	// Resume Compilation & Management
	resumeGroup.Post("/build", handlers.SaveResumeBlueprint)    // Saves the selected array of IDs + Summary
	resumeGroup.Get("/list", handlers.GetSavedResumes)          // Lists all generated resumes
	resumeGroup.Get("/compile/:id", handlers.GetCompiledResume) // The Magic Endpoint: Fetches fully assembled JSON
	resumeGroup.Delete("/delete/:id", handlers.DeleteResume)
}
