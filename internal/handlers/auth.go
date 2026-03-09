package handlers

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"tracktora-backend/internal/models"
	"tracktora-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user sign-up
func Register(c *fiber.Ctx) error {
	req := new(models.RegisterRequest)

	// 1. Parse the incoming JSON body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// 2. Validate basic input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email, and password are required",
		})
	}

	// 3. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process password",
		})
	}

	// 4. Save to the database
	userID, err := repository.CreateUser(req.Username, req.Email, string(hashedPassword))
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 5. Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user_id": userID,
	})
}

// Login handles user authentication and returns a JWT
func Login(c *fiber.Ctx) error {
	req := new(models.LoginRequest)

	// 1. Parse the incoming JSON body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// 2. Find the user by email
	user, hashedPassword, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// 3. Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// 4. Generate the JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	})

	// 5. Sign the token with our secret key
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	// 6. Return the user info and the token!
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
		"user":    user,
	})
}

// Helper to generate a 6-digit random code
func generateRandomCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

// ForgotPassword generates a dynamic random code
func ForgotPassword(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}

	// Generate the random 6-digit code
	resetCode := generateRandomCode()

	if err := repository.StoreResetToken(req.Email, resetCode); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Returning the dynamic code in the response for testing
	return c.Status(200).JSON(fiber.Map{
		"message": "Reset code generated. Check your email (simulated).",
		"code":    resetCode,
	})
}

// internal/handlers/auth.go

// ResetPassword handles the final password change using the token
func ResetPassword(c *fiber.Ctx) error {
	type Request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to secure new password"})
	}

	// Call repository to verify token and update DB
	if err := repository.ResetPassword(req.Token, string(hashedPassword)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Password updated successfully. You can now login.",
	})
}
