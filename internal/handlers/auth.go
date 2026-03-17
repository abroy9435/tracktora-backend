package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"tracktora-backend/internal/models"
	"tracktora-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// --- UTILS ---

func generateRandomCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

// sendEmail now bypasses the Hugging Face firewall by using a Google Webhook (Port 443)
func sendEmail(to string, subject string, body string) error {
	webhookURL := os.Getenv("GOOGLE_SCRIPT_URL")

	if webhookURL == "" {
		log.Println("WARNING: GOOGLE_SCRIPT_URL is not set in environment variables.")
		return fmt.Errorf("email webhook configuration missing")
	}

	// Prepare the JSON payload for the Google Script
	payload := map[string]string{
		"to":      to,
		"subject": subject,
		"body":    body,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Make the HTTP request to Google (Standard Web Traffic - Never Blocked)
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("!!! Webhook Request Failed: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// Handle the response from your Google Script
	if resp.StatusCode == 200 {
		log.Printf("SUCCESS: Verification email triggered for %s via Google Apps Script\n", to)
		return nil
	}

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("!!! Webhook Error (Status %d): %s\n", resp.StatusCode, string(respBody))
	return fmt.Errorf("failed to send email through webhook")
}

// --- HANDLERS ---

func Register(c *fiber.Ctx) error {
	req := new(models.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	userID, err := repository.CreateUser(req.Username, req.Email, string(hash))
	if err != nil {
		return c.Status(409).JSON(fiber.Map{"error": err.Error()})
	}

	code := generateRandomCode()
	repository.StoreVerificationToken(req.Email, code)

	// Send the email using our new webhook
	go sendEmail(req.Email, "Verify Your TrackTora Account",
		fmt.Sprintf("<h1>Welcome to TrackTora!</h1><p>Your verification code is: <strong>%s</strong></p>", code))

	return c.Status(201).JSON(fiber.Map{"message": "Verification code sent!", "user_id": userID})
}

func VerifyEmail(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	req := new(Request)
	c.BodyParser(req)

	if err := repository.VerifyAndActivateUser(req.Email, req.Code); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Email verified successfully!"})
}

func Login(c *fiber.Ctx) error {
	req := new(models.LoginRequest)
	c.BodyParser(req)

	user, hash, isVerified, err := repository.GetUserByEmailWithVerification(req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !isVerified {
		return c.Status(403).JSON(fiber.Map{"error": "Please verify your email address before logging in."})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	secret := os.Getenv("JWT_SECRET")
	tokenString, _ := token.SignedString([]byte(secret))

	return c.Status(200).JSON(fiber.Map{"token": tokenString, "user": user})
}

func ForgotPassword(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}
	req := new(Request)
	c.BodyParser(req)

	token := generateRandomCode()
	if err := repository.StoreResetToken(req.Email, token); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	go sendEmail(req.Email, "TrackTora Password Reset",
		fmt.Sprintf("<h1>Reset Code: %s</h1><p>Enter this code in the app to set a new password.</p>", token))

	return c.Status(200).JSON(fiber.Map{"message": "If account exists, code has been sent."})
}

func ResetPassword(c *fiber.Ctx) error {
	type Request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	req := new(Request)
	c.BodyParser(req)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err := repository.ResetPassword(req.Token, string(hashed)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Password updated successfully"})
}

func ResendVerification(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	_, _, isVerified, err := repository.GetUserByEmailWithVerification(req.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found"})
	}

	if isVerified {
		return c.Status(400).JSON(fiber.Map{"error": "Account is already verified"})
	}

	code := generateRandomCode()
	err = repository.StoreVerificationToken(req.Email, code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate new code"})
	}

	go sendEmail(req.Email, "New Verification Code",
		fmt.Sprintf("<h1>TrackTora</h1><p>Your new verification code is: <strong>%s</strong></p>", code))

	return c.Status(200).JSON(fiber.Map{"message": "New verification code sent!"})
}

func UpdatePassword(c *fiber.Ctx) error {
	// 1. Get user_id from JWT locals (Ensure your middleware sets this)
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Missing user identity"})
	}

	type Request struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// 2. Fetch current hash from repository
	currentHash, err := repository.GetUserHashByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User record not found"})
	}

	// 3. Compare current password
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Current password is incorrect"})
	}

	// 4. Hash new password and save
	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err := repository.UpdateUserPassword(userID, string(newHash)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update password in database"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Password updated successfully!"})
}
