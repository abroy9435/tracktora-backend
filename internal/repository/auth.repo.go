package repository

import (
	"context"
	"errors"
	"time"
	"tracktora-backend/internal/database"
)

// --- PASSWORD RESET (Table: password_resets) ---

func StoreResetToken(email, token string) error {
	var userID string
	err := database.DB.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return errors.New("no account found with this email")
	}

	database.DB.Exec(context.Background(), "DELETE FROM password_resets WHERE user_id = $1", userID)

	expiresAt := time.Now().Add(15 * time.Minute)
	query := `INSERT INTO password_resets (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err = database.DB.Exec(context.Background(), query, userID, token, expiresAt)
	return err
}

func ResetPassword(token, newHashedPassword string) error {
	var userID string
	query := `SELECT user_id FROM password_resets WHERE token = $1 AND expires_at > $2`
	err := database.DB.QueryRow(context.Background(), query, token, time.Now()).Scan(&userID)
	if err != nil {
		return errors.New("invalid or expired reset code")
	}

	_, err = database.DB.Exec(context.Background(), "UPDATE users SET password_hash = $1 WHERE id = $2", newHashedPassword, userID)
	if err != nil {
		return err
	}

	database.DB.Exec(context.Background(), "DELETE FROM password_resets WHERE user_id = $1", userID)
	return nil
}

// --- SIGNUP VERIFICATION (Table: verification_tokens) ---

func StoreVerificationToken(email string, code string) error {
	expiresAt := time.Now().Add(15 * time.Minute)
	query := `
		INSERT INTO verification_tokens (email, code, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE 
		SET code = EXCLUDED.code, expires_at = EXCLUDED.expires_at`
	_, err := database.DB.Exec(context.Background(), query, email, code, expiresAt)
	return err
}

func VerifyAndActivateUser(email, code string) error {
	var expiresAt time.Time
	err := database.DB.QueryRow(context.Background(),
		"SELECT expires_at FROM verification_tokens WHERE email = $1 AND code = $2",
		email, code).Scan(&expiresAt)

	if err != nil {
		return errors.New("invalid verification code")
	}

	if time.Now().After(expiresAt) {
		return errors.New("code has expired")
	}

	_, err = database.DB.Exec(context.Background(), "UPDATE users SET is_verified = TRUE WHERE email = $1", email)
	if err != nil {
		return err
	}

	_, _ = database.DB.Exec(context.Background(), "DELETE FROM verification_tokens WHERE email = $1", email)
	return nil
}
