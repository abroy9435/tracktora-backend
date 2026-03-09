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

// UpdateApplicationRequest defines what the frontend will send to update an app
type UpdateApplicationRequest struct {
	ID          string `json:"id"`
	CompanyName string `json:"company_name"`
	RoleTitle   string `json:"role_title"`
	Status      string `json:"status"`
	JobURL      string `json:"job_url"`
	Notes       string `json:"notes"`
	AppliedDate string `json:"applied_date"`
}

// DeleteApplicationRequest defines what the frontend will send to delete an app
type DeleteApplicationRequest struct {
	ID string `json:"id"`
}

// ApplicationStats represents the user's dashboard statistics
type ApplicationStats struct {
	Total        int `json:"total"`
	Wishlist     int `json:"wishlist"`
	Applied      int `json:"applied"`
	Interviewing int `json:"interviewing"`
	Offer        int `json:"offer"`
	Rejected     int `json:"rejected"`
}
