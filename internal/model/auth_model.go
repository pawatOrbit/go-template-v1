package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	User         *UserInfo `json:"user"`
}

type UserInfo struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FirstName    *string   `db:"first_name"`
	LastName     *string   `db:"last_name"`
	Roles        []string  `db:"roles"`
	IsActive     bool      `db:"is_active"`
	EmailVerified bool     `db:"email_verified"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}