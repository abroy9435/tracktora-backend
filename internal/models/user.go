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

// User represents the database user
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
