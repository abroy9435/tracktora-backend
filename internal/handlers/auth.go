package handlers

import (
	"crypto/tls"
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

func sendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if from == "" || password == "" {
		log.Println("WARNING: SMTP_EMAIL or SMTP_PASSWORD is not set.")
		return fmt.Errorf("SMTP credentials missing")
	}

	// Hugging Face workaround: Use port 465 and explicit TLS
	smtpHost := "smtp.gmail.com"
	smtpPort := "465"

	msg := "From: TrackTora <" + from + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Build the TLS configuration
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	// Connect directly over TLS
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		log.Printf("!!! TLS Connection Error: %v\n", err)
		return err
	}
	defer conn.Close()

	// Create the SMTP client
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("!!! SMTP Client Error: %v\n", err)
		return err
	}

	// Authenticate
	if err = c.Auth(auth); err != nil {
		log.Printf("!!! SMTP Auth Error: %v\n", err)
		return err
	}

	// Set sender and recipient
	if err = c.Mail(from); err != nil {
		return err
	}
	if err = c.Rcpt(to); err != nil {
		return err
	}

	// Send the email body
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	log.Printf("SUCCESS: Email sent to %s via port 465\n", to)
	return nil
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
