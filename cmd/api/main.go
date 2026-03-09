package main

import (
	"log"
	"os"

	"tracktora-backend/internal/database"
	"tracktora-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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

	// 4. Setup all routes cleanly!
	routes.Setup(app)

	// 5. Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
