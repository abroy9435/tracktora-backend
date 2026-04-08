package handlers

import (
	"context"
	"tracktora-backend/internal/database"
	"tracktora-backend/internal/models"
	"tracktora-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// GetProfile returns the logged-in user's details
func GetProfile(c *fiber.Ctx) error {
	// 1. Grab the secure user_id from the JWT bouncer
	userID := c.Locals("user_id").(string)

	// 2. Fetch the user from the database
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 3. Return the user object
	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateProfile handles changing user details dynamically
func UpdateProfile(c *fiber.Ctx) error {
	// 1. Grab the secure user_id
	userID := c.Locals("user_id").(string)

	// 2. Parse the JSON body
	req := new(models.UpdateProfileRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// 3. Send to the repository to update PostgreSQL
	if err := repository.UpdateUser(userID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 4. Return success!
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

// UpdatePrivacySettings toggles the share_stats flag
func UpdatePrivacySettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	type Request struct {
		ShareStats bool `json:"share_stats"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := database.DB.Exec(context.Background(), "UPDATE users SET share_stats = $1 WHERE id = $2", req.ShareStats, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update settings"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Privacy settings updated!"})
}
