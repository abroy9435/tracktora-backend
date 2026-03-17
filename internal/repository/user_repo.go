package repository

import (
	"context"
	"errors"

	"tracktora-backend/internal/database"
	"tracktora-backend/internal/models"
)

func CreateUser(username, email, passwordHash string) (string, error) {
	var userID string
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`

	err := database.DB.QueryRow(context.Background(), query, username, email, passwordHash).Scan(&userID)
	if err != nil {
		return "", errors.New("failed to create user, email or username might already exist")
	}
	return userID, nil
}

func GetUserByEmailWithVerification(email string) (*models.User, string, bool, error) {
	var user models.User
	var passwordHash string
	var isVerified bool

	query := `SELECT id, username, email, password_hash, is_verified, share_stats, created_at FROM users WHERE email = $1`

	err := database.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &isVerified, &user.ShareStats, &user.CreatedAt,
	)
	if err != nil {
		return nil, "", false, errors.New("invalid email or password")
	}

	return &user, passwordHash, isVerified, nil
}

func GetUserByID(userID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, is_verified, share_stats, created_at FROM users WHERE id = $1`

	err := database.DB.QueryRow(context.Background(), query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.IsVerified, &user.ShareStats, &user.CreatedAt,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func UpdateUser(userID string, req *models.UpdateProfileRequest) error {
	query := `UPDATE users SET username = $1 WHERE id = $2`
	result, err := database.DB.Exec(context.Background(), query, req.Username, userID)
	if err != nil {
		return errors.New("failed to update profile")
	}
	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
