package models

import "time"

// CreateApplicationRequest defines what the frontend will send us
type CreateApplicationRequest struct {
	CompanyName string `json:"company_name"`
	RoleTitle   string `json:"role_title"`
	Status      string `json:"status"` // e.g., "Wishlist", "Applied", "Interviewing", "Rejected", "Offer"
	JobURL      string `json:"job_url"`
	Notes       string `json:"notes"`
	AppliedDate string `json:"applied_date"` // YYYY-MM-DD format
}

// Application represents the full database record
type Application struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CompanyName string    `json:"company_name"`
	RoleTitle   string    `json:"role_title"`
	Status      string    `json:"status"`
	JobURL      string    `json:"job_url"`
	Notes       string    `json:"notes"`
	AppliedDate string    `json:"applied_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
