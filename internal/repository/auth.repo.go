package repository

import (
	"context"
	"errors"
	"time"
	"tracktora-backend/internal/database"
)

// StoreResetToken saves a code for a user, valid for 15 mins
func StoreResetToken(email, token string) error {
	var userID string
	// 1. Find user
	err := database.DB.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return errors.New("no account found with this email")
	}

	// 2. Clear old tokens
	database.DB.Exec(context.Background(), "DELETE FROM password_resets WHERE user_id = $1", userID)

	// 3. Insert new token
	expiresAt := time.Now().Add(15 * time.Minute)
	query := `INSERT INTO password_resets (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err = database.DB.Exec(context.Background(), query, userID, token, expiresAt)

	return err
}

// ResetPassword verifies the token and updates the user's password
func ResetPassword(token, newHashedPassword string) error {
	var userID string
	// Verify token is valid and not expired
	query := `SELECT user_id FROM password_resets WHERE token = $1 AND expires_at > $2`
	err := database.DB.QueryRow(context.Background(), query, token, time.Now()).Scan(&userID)
	if err != nil {
		return errors.New("invalid or expired reset code")
	}

	// Update password
	_, err = database.DB.Exec(context.Background(), "UPDATE users SET password_hash = $1 WHERE id = $2", newHashedPassword, userID)
	if err != nil {
		return err
	}

	// Delete used token
	database.DB.Exec(context.Background(), "DELETE FROM password_resets WHERE user_id = $1", userID)
	return nil
}
