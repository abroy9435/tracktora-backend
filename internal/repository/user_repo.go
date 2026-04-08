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

	// Fetch all the new columns
	query := `SELECT id, username, email, is_verified, share_stats, created_at, 
	          phone, city, state, linkedin_url, github_url, portfolio_url, other_link_name, other_link_url 
	          FROM users WHERE id = $1`

	err := database.DB.QueryRow(context.Background(), query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.IsVerified, &user.ShareStats, &user.CreatedAt,
		&user.Phone, &user.City, &user.State, &user.LinkedinURL, &user.GithubURL, &user.PortfolioURL, &user.OtherLinkName, &user.OtherLinkURL,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func UpdateUser(userID string, req *models.UpdateProfileRequest) error {
	// COALESCE ensures that if a pointer is nil (meaning Flutter didn't send that field),
	// it safely keeps the existing database value instead of erasing it.
	query := `UPDATE users SET 
		username = COALESCE($1, username),
		phone = COALESCE($2, phone),
		city = COALESCE($3, city),
		state = COALESCE($4, state),
		linkedin_url = COALESCE($5, linkedin_url),
		github_url = COALESCE($6, github_url),
		portfolio_url = COALESCE($7, portfolio_url),
		other_link_name = COALESCE($8, other_link_name),
		other_link_url = COALESCE($9, other_link_url)
		WHERE id = $10`

	result, err := database.DB.Exec(context.Background(), query,
		req.Username, req.Phone, req.City, req.State,
		req.LinkedinURL, req.GithubURL, req.PortfolioURL,
		req.OtherLinkName, req.OtherLinkURL,
		userID,
	)
	if err != nil {
		return errors.New("failed to update profile")
	}
	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
