package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth is the security guard that checks for a valid JWT
func RequireAuth(c *fiber.Ctx) error {
	// 1. Get the "Authorization" header from the incoming request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Access denied. No authorization header provided.",
		})
	}

	// 2. Check if it is formatted correctly as "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization format. Expected 'Bearer <token>'",
		})
	}

	tokenString := parts[1]
	secret := os.Getenv("JWT_SECRET")

	// 3. Parse and validate the token using our secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is exactly what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	// 4. If the token is fake, altered, or expired, reject them
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// 5. If it IS valid, extract the user's ID and save it into the request context
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// c.Locals allows us to pass data to the next function down the line
		c.Locals("user_id", claims["user_id"])

		// Let them pass through to the actual route!
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Failed to parse token claims",
	})
}
