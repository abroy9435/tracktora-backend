package handlers

import (
	"tracktora-backend/internal/clients"
	"tracktora-backend/internal/models"
	"tracktora-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// CreateApplication handles adding a new job application
func CreateApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	req := new(models.CreateApplicationRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.CompanyName == "" || req.RoleTitle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Company name and role title are required",
		})
	}

	appID, err := repository.CreateApplication(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save application",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":        "Application tracked successfully",
		"application_id": appID,
	})
}

// GetApplications fetches all applications for the logged-in user
func GetApplications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	apps, err := repository.GetUserApplications(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch applications",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"applications": apps,
	})
}

// UpdateApplication handles modifying an existing job application
func UpdateApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	req := new(models.UpdateApplicationRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application ID is required in the body",
		})
	}

	err := repository.UpdateApplication(userID, req.ID, req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Application updated successfully",
	})
}

// DeleteApplication handles removing a job application
func DeleteApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	req := new(models.DeleteApplicationRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application ID is required in the body",
		})
	}

	err := repository.DeleteApplication(userID, req.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Application deleted successfully",
	})
}

// GetApplicationStats handles fetching the user's dashboard statistics
func GetApplicationStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	stats, err := repository.GetApplicationStats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch dashboard statistics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

// GetExplorePage returns live ads from Adzuna with advanced filtering
func GetExplorePage(c *fiber.Ctx) error {
	search := c.Query("search")
	location := c.Query("location")
	page := c.QueryInt("page", 1)
	salaryMin := c.QueryInt("salary", 0)

	jobs, err := clients.FetchLiveJobs(search, location, page, salaryMin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch jobs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"page":         page,
		"salary_query": salaryMin,
		"explore_feed": jobs,
	})
}

// SaveExternalJob handles saving a listing from the Explore feed into the user's tracker
func SaveExternalJob(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	req := new(models.CreateApplicationRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid job data",
		})
	}

	if req.Status == "" {
		req.Status = "Wishlist"
	}

	appID, err := repository.CreateApplication(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save listing to your tracker",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":        "Job saved to your Wishlist!",
		"application_id": appID,
	})
}
