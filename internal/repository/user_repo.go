package repository

import (
	"context"
	"errors"
	"tracktora-backend/internal/database" // Update if your go.mod module name is different
)

// CreateUser inserts a new user into the database and returns the generated UUID
func CreateUser(username, email, passwordHash string) (string, error) {
	var userID string

	// The SQL query to insert a user and return their new ID
	query := `
		INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`

	// Execute the query
	err := database.DB.QueryRow(context.Background(), query, username, email, passwordHash).Scan(&userID)
	if err != nil {
		return "", errors.New("failed to create user, email or username might already exist")
	}

	return userID, nil
}
