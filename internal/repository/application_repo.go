package repository

import (
	"context"
	"errors"

	"tracktora-backend/internal/database"
	"tracktora-backend/internal/models"
)

// CreateApplication inserts a new job application into the database
func CreateApplication(userID string, req *models.CreateApplicationRequest) (string, error) {
	var applicationID string

	// If no status is provided, default to "Wishlist"
	status := req.Status
	if status == "" {
		status = "Wishlist"
	}

	query := `
		INSERT INTO applications (user_id, company_name, role_title, status, job_url, notes, applied_date) 
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')::DATE) 
		RETURNING id
	`

	err := database.DB.QueryRow(
		context.Background(),
		query,
		userID,
		req.CompanyName,
		req.RoleTitle,
		status,
		req.JobURL,
		req.Notes,
		req.AppliedDate,
	).Scan(&applicationID)

	if err != nil {
		return "", errors.New("failed to create application record")
	}

	return applicationID, nil
}

// GetUserApplications fetches all applications belonging to a specific user
func GetUserApplications(userID string) ([]models.Application, error) {
	// Initialize an empty slice (array) to hold the applications
	applications := []models.Application{}

	query := `
		SELECT id, user_id, company_name, role_title, status, job_url, notes, COALESCE(applied_date::TEXT, ''), created_at, updated_at 
		FROM applications 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, errors.New("failed to fetch applications")
	}
	defer rows.Close()

	// Loop through the database rows and add them to our slice
	for rows.Next() {
		var app models.Application
		err := rows.Scan(
			&app.ID,
			&app.UserID,
			&app.CompanyName,
			&app.RoleTitle,
			&app.Status,
			&app.JobURL,
			&app.Notes,
			&app.AppliedDate,
			&app.CreatedAt,
			&app.UpdatedAt,
		)
		if err != nil {
			continue // Skip broken rows
		}
		applications = append(applications, app)
	}

	return applications, nil
}
