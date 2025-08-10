package integration

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yourorg/go-api-template/core/auth"
	"github.com/yourorg/go-api-template/core/logger"
	middleware "github.com/yourorg/go-api-template/core/transport/httpserver/middlewares"
)

// AuthMiddlewareTestSuite defines the test suite for authentication middleware
type AuthMiddlewareTestSuite struct {
	suite.Suite
	authService  *auth.AuthService
	authConfig   middleware.AuthConfig
	testServer   *httptest.Server
	validToken   string
	invalidToken string
}

// SetupSuite runs before all tests in the suite
func (suite *AuthMiddlewareTestSuite) SetupSuite() {
	// Initialize logger for tests
	logger.Slog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	
	secretKey := "test-secret-key-for-middleware-testing"
	suite.authService = auth.NewAuthService(secretKey)
	
	suite.authConfig = middleware.AuthConfig{
		JWTSecretKey: secretKey,
		SkipPaths: []string{
			"/public",
			"/health",
		},
	}

	// Generate a valid token for testing
	tokenPair, err := suite.authService.GenerateTokens("test-user", "test@example.com", []string{"user"})
	assert.NoError(suite.T(), err)
	suite.validToken = tokenPair.AccessToken
	suite.invalidToken = "invalid.jwt.token"

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, exists := middleware.GetUserIDFromContext(r.Context())
		if exists {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Hello %s", userID)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("User ID not found in context"))
		}
	})

	// Wrap handler with auth middleware
	authHandler := middleware.AuthMiddleware(suite.authConfig)(handler)
	suite.testServer = httptest.NewServer(authHandler)
}

// TearDownSuite runs after all tests in the suite
func (suite *AuthMiddlewareTestSuite) TearDownSuite() {
	if suite.testServer != nil {
		suite.testServer.Close()
	}
}

// TestAuthMiddleware_ValidToken tests middleware with valid token
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_ValidToken() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/api/test", nil)
	assert.NoError(suite.T(), err)
	
	req.Header.Set("Authorization", "Bearer "+suite.validToken)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	
	resp.Body.Close()
}

// TestAuthMiddleware_MissingToken tests middleware without token
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_MissingToken() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/api/test", nil)
	assert.NoError(suite.T(), err)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
	
	resp.Body.Close()
}

// TestAuthMiddleware_InvalidToken tests middleware with invalid token
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/api/test", nil)
	assert.NoError(suite.T(), err)
	
	req.Header.Set("Authorization", "Bearer "+suite.invalidToken)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
	
	resp.Body.Close()
}

// TestAuthMiddleware_InvalidAuthHeaderFormat tests invalid auth header format
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidAuthHeaderFormat() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/api/test", nil)
	assert.NoError(suite.T(), err)
	
	req.Header.Set("Authorization", "InvalidFormat "+suite.validToken)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
	
	resp.Body.Close()
}

// TestAuthMiddleware_SkipPaths tests that certain paths skip authentication
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_SkipPaths() {
	// Test public path that should skip auth
	req, err := http.NewRequest("GET", suite.testServer.URL+"/public", nil)
	assert.NoError(suite.T(), err)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	
	assert.NoError(suite.T(), err)
	// Should return 500 because the test handler expects user context, 
	// but it means auth was skipped
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
	
	resp.Body.Close()

	// Test health path that should skip auth
	req, err = http.NewRequest("GET", suite.testServer.URL+"/health", nil)
	assert.NoError(suite.T(), err)
	
	resp, err = client.Do(req)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
	
	resp.Body.Close()
}

// Run the test suite
func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

// TestRequireRoles tests role-based authorization
func TestRequireRoles(t *testing.T) {
	// Initialize logger for test
	logger.Slog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	
	secretKey := "test-secret-for-roles"
	authService := auth.NewAuthService(secretKey)
	
	authConfig := middleware.AuthConfig{
		JWTSecretKey: secretKey,
		SkipPaths:    []string{},
	}

	// Generate tokens for different users
	adminToken, err := authService.GenerateTokens("admin-user", "admin@example.com", []string{"admin"})
	assert.NoError(t, err)
	
	userToken, err := authService.GenerateTokens("regular-user", "user@example.com", []string{"user"})
	assert.NoError(t, err)

	// Create test handler that requires admin role
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Admin access granted"))
	})

	// Chain auth middleware and role middleware
	authHandler := middleware.AuthMiddleware(authConfig)(
		middleware.RequireRoles("admin")(handler),
	)

	testCases := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "admin user should access",
			token:          adminToken.AccessToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "regular user should be forbidden",
			token:          userToken.AccessToken,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/admin", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			
			rr := httptest.NewRecorder()
			authHandler.ServeHTTP(rr, req)
			
			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}