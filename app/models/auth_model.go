// app/models/auth_model.go

package models

// LoginRequest represents the request body for the login endpoint
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest defines the structure of the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
