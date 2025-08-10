package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yourorg/go-api-template/core/auth"
)

// AuthServiceTestSuite defines the test suite for AuthService
type AuthServiceTestSuite struct {
	suite.Suite
	authService *auth.AuthService
	secretKey   string
}

// SetupSuite runs before all tests in the suite
func (suite *AuthServiceTestSuite) SetupSuite() {
	suite.secretKey = "test-secret-key-for-jwt-signing"
	suite.authService = auth.NewAuthService(suite.secretKey)
}

// TestGenerateTokens tests token generation
func (suite *AuthServiceTestSuite) TestGenerateTokens() {
	userID := "user-123"
	email := "test@example.com"
	roles := []string{"user", "admin"}

	tokenPair, err := suite.authService.GenerateTokens(userID, email, roles)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenPair.AccessToken)
	assert.NotEmpty(suite.T(), tokenPair.RefreshToken)
	assert.NotEqual(suite.T(), tokenPair.AccessToken, tokenPair.RefreshToken)
}

// TestValidateRefreshToken tests refresh token validation
func (suite *AuthServiceTestSuite) TestValidateRefreshToken() {
	userID := "user-123"
	email := "test@example.com"
	roles := []string{"user"}

	// Generate tokens
	tokenPair, err := suite.authService.GenerateTokens(userID, email, roles)
	assert.NoError(suite.T(), err)

	// Validate refresh token
	extractedUserID, err := suite.authService.ValidateRefreshToken(tokenPair.RefreshToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)
}

// TestValidateRefreshTokenWithInvalidToken tests validation with invalid token
func (suite *AuthServiceTestSuite) TestValidateRefreshTokenWithInvalidToken() {
	invalidToken := "invalid.token.here"

	_, err := suite.authService.ValidateRefreshToken(invalidToken)
	assert.Error(suite.T(), err)
}

// TestHashPassword tests password hashing
func (suite *AuthServiceTestSuite) TestHashPassword() {
	password := "mySecretPassword123"

	hashedPassword, err := suite.authService.HashPassword(password)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
	assert.NotEqual(suite.T(), password, hashedPassword)
}

// TestVerifyPassword tests password verification
func (suite *AuthServiceTestSuite) TestVerifyPassword() {
	password := "mySecretPassword123"

	// Hash the password
	hashedPassword, err := suite.authService.HashPassword(password)
	assert.NoError(suite.T(), err)

	// Verify correct password
	isValid := suite.authService.VerifyPassword(hashedPassword, password)
	assert.True(suite.T(), isValid)

	// Verify incorrect password
	isValid = suite.authService.VerifyPassword(hashedPassword, "wrongPassword")
	assert.False(suite.T(), isValid)
}

// TestGenerateSecretKey tests secret key generation
func (suite *AuthServiceTestSuite) TestGenerateSecretKey() {
	key1, err := auth.GenerateSecretKey()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), key1)

	key2, err := auth.GenerateSecretKey()
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), key2)

	// Keys should be different
	assert.NotEqual(suite.T(), key1, key2)

	// Keys should be 64 characters (32 bytes hex encoded)
	assert.Equal(suite.T(), 64, len(key1))
	assert.Equal(suite.T(), 64, len(key2))
}

// Run the test suite
func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

// Additional table-driven tests for edge cases
func TestAuthService_EdgeCases(t *testing.T) {
	authService := auth.NewAuthService("test-secret")

	tests := []struct {
		name    string
		userID  string
		email   string
		roles   []string
		wantErr bool
	}{
		{
			name:    "empty user ID",
			userID:  "",
			email:   "test@example.com",
			roles:   []string{"user"},
			wantErr: false, // Should still work
		},
		{
			name:    "empty email",
			userID:  "user-123",
			email:   "",
			roles:   []string{"user"},
			wantErr: false, // Should still work
		},
		{
			name:    "empty roles",
			userID:  "user-123",
			email:   "test@example.com",
			roles:   []string{},
			wantErr: false, // Should still work
		},
		{
			name:    "nil roles",
			userID:  "user-123",
			email:   "test@example.com",
			roles:   nil,
			wantErr: false, // Should still work
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := authService.GenerateTokens(tt.userID, tt.email, tt.roles)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tokens)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
		})
	}
}
