package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
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

func sendEmail(to, subject, htmlBody string) {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// RFC 822 format for Gmail SMTP
	message := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		htmlBody + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Background execution
	go func() {
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
		if err != nil {
			log.Printf("!!! SMTP ERROR to %s: %v", to, err)
			return
		}
		log.Printf("Email successfully sent to %s via Gmail SMTP", to)
	}()
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

	sendEmail(req.Email, "Verify Your TrackTora Account",
		fmt.Sprintf("<h1>Welcome!</h1><p>Your verification code is: <strong>%s</strong></p>", code))

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

	sendEmail(req.Email, "Password Reset Code",
		fmt.Sprintf("<h1>Reset Code: %s</h1>", token))

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

	// 1. Check if user exists and their current status
	_, _, isVerified, err := repository.GetUserByEmailWithVerification(req.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found"})
	}

	if isVerified {
		return c.Status(400).JSON(fiber.Map{"error": "Account is already verified"})
	}

	// 2. Generate and store new code (this overrides the old expired one)
	code := generateRandomCode()
	err = repository.StoreVerificationToken(req.Email, code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate new code"})
	}

	// 3. Send the email
	sendEmail(req.Email, "New Verification Code",
		fmt.Sprintf("<h1>TrackTora</h1><p>Your new verification code is: <strong>%s</strong></p>", code))

	return c.Status(200).JSON(fiber.Map{"message": "New verification code sent!"})
}
