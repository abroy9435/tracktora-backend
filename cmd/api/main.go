package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/yourusername/tracktora-backend/internal/database" // Update with your actual module path
)

func main() {
	// 1. Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or error reading it")
	}

	// 2. Connect to PostgreSQL
	database.ConnectDB()
	defer database.DB.Close() // Ensure the database connection closes when the app stops

	// 3. Initialize Fiber app
	app := fiber.New()

	// 4. Simple Health Check Route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "TrackTora Backend is up and running! 🚀",
		})
	})

	// 5. Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
