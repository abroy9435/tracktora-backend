package models

import "time"

// RegisterRequest defines the JSON payload for signing up
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest defines the JSON payload for logging in
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents the database user and what is returned to the frontend
type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	ShareStats bool      `json:"share_stats"`
	CreatedAt  time.Time `json:"created_at"`

	// --- NEW: Contact & Links (Pointers handle SQL NULLs gracefully) ---
	Phone         *string `json:"phone"`
	City          *string `json:"city"`
	State         *string `json:"state"`
	LinkedinURL   *string `json:"linkedin_url"`
	GithubURL     *string `json:"github_url"`
	PortfolioURL  *string `json:"portfolio_url"`
	OtherLinkName *string `json:"other_link_name"`
	OtherLinkURL  *string `json:"other_link_url"`
}

// UpdateProfileRequest defines the payload for updating user details
// Using pointers ensures we only update fields the user actually sent.
type UpdateProfileRequest struct {
	Username      *string `json:"username"`
	Phone         *string `json:"phone"`
	City          *string `json:"city"`
	State         *string `json:"state"`
	LinkedinURL   *string `json:"linkedin_url"`
	GithubURL     *string `json:"github_url"`
	PortfolioURL  *string `json:"portfolio_url"`
	OtherLinkName *string `json:"other_link_name"`
	OtherLinkURL  *string `json:"other_link_url"`
}
