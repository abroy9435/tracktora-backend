package handlers

import (
	"tracktora-backend/internal/models"
	"tracktora-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// CreateApplication handles adding a new job application
func CreateApplication(c *fiber.Ctx) error {
	// 1. Grab the user_id that the JWT middleware securely saved for us
	userID := c.Locals("user_id").(string)

	req := new(models.CreateApplicationRequest)

	// 2. Parse the incoming JSON body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// 3. Quick validation
	if req.CompanyName == "" || req.RoleTitle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Company name and role title are required",
		})
	}

	// 4. Send to the repository to save in PostgreSQL
	appID, err := repository.CreateApplication(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save application",
		})
	}

	// 5. Return success!
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":        "Application tracked successfully",
		"application_id": appID,
	})
}

// GetApplications fetches all applications for the logged-in user
func GetApplications(c *fiber.Ctx) error {
	// 1. Grab the secure user_id
	userID := c.Locals("user_id").(string)

	// 2. Fetch from the database
	apps, err := repository.GetUserApplications(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch applications",
		})
	}

	// 3. Return the list!
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"applications": apps,
	})
}
