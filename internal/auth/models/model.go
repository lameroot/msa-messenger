package auth_models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is not exposed in JSON
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	Info      string    `json:"info"`
}

// Token represents an authentication token pair
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// RegisterRequest represents the request body for registration and login
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required"`
}

// LoginRequest represents the request body for registration and login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// AuthResponse represents the response for a successful authentication
type AuthResponse struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

// RefreshRequest represents the request body for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UserResponse represents the response for user-related operations
type UserResponse struct {
	User User `json:"user"`
}

// TokenResponse represents the response for token-related operations
type TokenResponse struct {
	Token Token `json:"token"`
}

// TokenConfig holds configuration for token generation
type TokenConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type UpdateUserRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	AvatarURL string `json:"avatar_url"`
	Info      string `json:"info"`
}
