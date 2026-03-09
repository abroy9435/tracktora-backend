package handlers

import (
	"tracktora-backend/internal/models"
	"tracktora-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// POST /api/connect/invite
func SendFriendRequest(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	req := new(models.ConnectionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := repository.SendInviteByID(userID, req.FriendID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Friend request sent!"})
}

// GET /api/connect/requests
func GetPendingRequests(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	requests, err := repository.GetPendingRequests(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch requests"})
	}
	return c.Status(200).JSON(fiber.Map{"pending_requests": requests})
}

// PUT /api/connect/respond
func RespondToFriendRequest(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string) // The Receiver
	req := new(models.UpdateConnectionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	// We use req.FriendID here as the 'SenderID' we want to accept
	if err := repository.UpdateFriendStatus(userID, req.FriendID, req.Status); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Friend request " + req.Status})
}

// GET /api/connect/stats/:friend_id
func GetFriendStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	friendID := c.Params("friend_id")
	stats, err := repository.GetFriendStats(userID, friendID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(stats)
}
