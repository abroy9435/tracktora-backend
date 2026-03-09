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

// UpdateApplication modifies an existing application in the database
func UpdateApplication(userID, appID string, req *models.UpdateApplicationRequest) error {
	query := `
		UPDATE applications 
		SET company_name = $1, role_title = $2, status = $3, job_url = $4, notes = $5, applied_date = NULLIF($6, '')::DATE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND user_id = $8
	`

	// database.DB.Exec is used when we don't need to return any rows, just run a command
	result, err := database.DB.Exec(
		context.Background(),
		query,
		req.CompanyName,
		req.RoleTitle,
		req.Status,
		req.JobURL,
		req.Notes,
		req.AppliedDate,
		appID,
		userID,
	)

	if err != nil {
		return errors.New("failed to update application")
	}

	// Check if any row was actually updated (prevents updating non-existent apps)
	if result.RowsAffected() == 0 {
		return errors.New("application not found or you don't have permission to update it")
	}

	return nil
}
