package service

import (
	"context"
	"time"

	"github.com/yourorg/go-api-template/core/auth"
	"github.com/yourorg/go-api-template/core/exception"
	"github.com/yourorg/go-api-template/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
}

type authService struct {
	authCore *auth.AuthService
	errors   *exception.MockDataServiceErrors
}

func NewAuthService(authCore *auth.AuthService, errors *exception.MockDataServiceErrors) AuthService {
	return &authService{
		authCore: authCore,
		errors:   errors,
	}
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// Mock user data - in production, you would fetch this from database
	// For demonstration, we'll use a hardcoded user
	mockUsers := map[string]struct {
		ID           string
		PasswordHash string
		FirstName    string
		LastName     string
		Roles        []string
	}{
		"user@example.com": {
			ID:           "550e8400-e29b-41d4-a716-446655440001",
			PasswordHash: "$2a$10$YourHashedPasswordHere", // Password: "password123"
			FirstName:    "John",
			LastName:     "Doe",
			Roles:        []string{"user"},
		},
		"admin@example.com": {
			ID:           "550e8400-e29b-41d4-a716-446655440002",
			PasswordHash: "$2a$10$AnotherHashedPassword", // Password: "admin123"
			FirstName:    "Admin",
			LastName:     "User",
			Roles:        []string{"admin", "user"},
		},
	}

	// Validate request fields
	if req.Email == "" || req.Password == "" {
		fields := []string{}
		if req.Email == "" {
			fields = append(fields, "email")
		}
		if req.Password == "" {
			fields = append(fields, "password")
		}
		return nil, s.errors.ErrInvalidRequest.
			WithMessage("Missing required fields").
			WithFields(fields).
			WithDebugMessage("Email and password are required")
	}

	// Find user by email
	user, exists := mockUsers[req.Email]
	if !exists {
		return nil, s.errors.ErrNotFound.WithDebugMessage("User not found")
	}

	// For demo purposes, check a simple password match
	// In production, you would verify the password hash using bcrypt
	validPasswords := map[string]string{
		"user@example.com":  "password123",
		"admin@example.com": "admin123",
	}

	if expectedPassword, ok := validPasswords[req.Email]; !ok || req.Password != expectedPassword {
		return nil, s.errors.ErrUnauthorized.
			WithMessage("Authentication failed").
			WithDatas(map[string]string{
				"email":  req.Email,
				"reason": "Invalid credentials",
			}).
			WithDebugMessage("Invalid password for user: " + req.Email)
	}

	// Generate tokens
	tokenPair, err := s.authCore.GenerateTokens(user.ID, req.Email, user.Roles)
	if err != nil {
		return nil, s.errors.ErrUnauthorized.WithDebugMessage(err.Error())
	}

	// Calculate expiration time
	expiresIn := int64(24 * time.Hour / time.Second) // 24 hours in seconds

	return &model.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		User: &model.UserInfo{
			ID:        user.ID,
			Email:     req.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     user.Roles,
		},
	}, nil
}
