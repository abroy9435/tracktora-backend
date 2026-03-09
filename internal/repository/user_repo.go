package repository

import (
	"context"
	"errors"

	"tracktora-backend/internal/database"
	"tracktora-backend/internal/models"
)

// CreateUser inserts a new user into the database and returns the generated UUID
func CreateUser(username, email, passwordHash string) (string, error) {
	var userID string

	query := `
		INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`

	err := database.DB.QueryRow(context.Background(), query, username, email, passwordHash).Scan(&userID)
	if err != nil {
		return "", errors.New("failed to create user, email or username might already exist")
	}

	return userID, nil
}

// GetUserByEmail finds a user and returns their data along with the hashed password
func GetUserByEmail(email string) (*models.User, string, error) {
	var user models.User
	var passwordHash string

	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1`

	// Scan directly into user.CreatedAt
	err := database.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	return &user, passwordHash, nil
}
