package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	middleware "github.com/yourorg/go-api-template/core/transport/httpserver/middlewares"
	"golang.org/x/crypto/bcrypt"
)

// AuthService provides authentication services
type AuthService struct {
	jwtSecretKey    string
	tokenExpiration time.Duration
	refreshTokenExp time.Duration
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the authentication response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSecretKey string) *AuthService {
	return &AuthService{
		jwtSecretKey:    jwtSecretKey,
		tokenExpiration: 24 * time.Hour,     // 24 hours
		refreshTokenExp: 7 * 24 * time.Hour, // 7 days
	}
}

// GenerateTokens creates JWT access and refresh tokens for a user
func (s *AuthService) GenerateTokens(userID, email string, roles []string) (*TokenPair, error) {
	// Create access token
	accessToken, err := s.generateAccessToken(userID, email, roles)
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateAccessToken creates a JWT access token
func (s *AuthService) generateAccessToken(userID, email string, roles []string) (string, error) {
	now := time.Now()
	claims := &middleware.UserClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenExpiration)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "go-api-template",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecretKey))
}

// generateRefreshToken creates a refresh token
func (s *AuthService) generateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenExp)),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "go-api-template",
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecretKey))
}

// ValidateRefreshToken validates and extracts user ID from refresh token
func (s *AuthService) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", jwt.ErrTokenInvalidClaims
	}

	return claims.Subject, nil
}

// HashPassword hashes a plain text password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword compares a plain text password with its hash
func (s *AuthService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateSecretKey generates a secure random secret key for JWT signing
func GenerateSecretKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
