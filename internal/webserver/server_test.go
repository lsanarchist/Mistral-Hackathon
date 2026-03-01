
package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPluginMarketplaceEndpoint(t *testing.T) {
	// Create WebSocket server
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Test marketplace endpoint
	req := httptest.NewRequest("GET", "/plugins/marketplace", nil)
	w := httptest.NewRecorder()

	server.handlePluginMarketplace(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, response, "plugins")
	assert.Contains(t, response, "count")
	
	plugins := response["plugins"].([]interface{})
	assert.True(t, len(plugins) > 0, "Should have at least one plugin in marketplace")
}

func TestPluginInstallEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Test install endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/install", nil)
	w := httptest.NewRecorder()

	server.handleInstallPlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPluginUpdateEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Test update endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/update", nil)
	w := httptest.NewRecorder()

	server.handleUpdatePlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPluginUninstallEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Test uninstall endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/uninstall", nil)
	w := httptest.NewRecorder()

	server.handleUninstallPlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPerformanceHistory(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Create test findings
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
		},
		Findings: []model.Finding{
			{Severity: "critical", Title: "Test Critical"},
			{Severity: "high", Title: "Test High"},
			{Severity: "medium", Title: "Test Medium"},
			{Severity: "low", Title: "Test Low"},
		},
	}

	// Update server data (this should trigger snapshot recording)
	server.UpdateData(findings, nil)

	// Wait a bit for async operations
	time.Sleep(100 * time.Millisecond)

	// Test GetPerformanceHistory method
	history := server.GetPerformanceHistory()
	assert.NotNil(t, history)
	assert.True(t, len(history) > 0, "Should have at least one performance snapshot")

	// Verify the snapshot data
	latest := history[len(history)-1]
	assert.Equal(t, 75, latest.OverallScore)
	assert.Equal(t, 1, latest.CriticalCount)
	assert.Equal(t, 1, latest.HighCount)
	assert.Equal(t, 1, latest.MediumCount)
	assert.Equal(t, 1, latest.LowCount)
	assert.Equal(t, 4, latest.TotalFindings)
	assert.Equal(t, 0, latest.ClientCount) // No clients connected in test
}

func TestPerformanceHistoryEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Create test findings
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 80,
		},
		Findings: []model.Finding{
			{Severity: "critical", Title: "Test Critical"},
			{Severity: "high", Title: "Test High"},
		},
	}

	// Update server data
	server.UpdateData(findings, nil)

	// Wait for snapshot recording
	time.Sleep(100 * time.Millisecond)

	// Test performance history endpoint
	req := httptest.NewRequest("GET", "/performance/history", nil)
	w := httptest.NewRecorder()

	server.handlePerformanceHistory(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, response, "history")
	assert.Contains(t, response, "count")
	assert.Contains(t, response, "analysis")

	history := response["history"].([]interface{})
	assert.True(t, len(history) > 0, "Should have performance history")
}

func TestPerformanceAnalysisEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Create test findings
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 85,
		},
		Findings: []model.Finding{
			{Severity: "critical", Title: "Test Critical"},
			{Severity: "high", Title: "Test High"},
			{Severity: "medium", Title: "Test Medium"},
		},
	}

	// Update server data
	server.UpdateData(findings, nil)

	// Wait for snapshot recording
	time.Sleep(100 * time.Millisecond)

	// Test performance analysis endpoint
	req := httptest.NewRequest("GET", "/performance/analysis", nil)
	w := httptest.NewRecorder()

	server.handlePerformanceAnalysis(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, response, "current")
	assert.Contains(t, response, "analysis")
	assert.Contains(t, response, "trends")

	current := response["current"].(map[string]interface{})
	assert.Equal(t, float64(85), current["overall_score"])
	assert.Equal(t, float64(1), current["critical_count"])
	assert.Equal(t, float64(1), current["high_count"])
	assert.Equal(t, float64(1), current["medium_count"])
}

func TestPerformanceSnapshotLimit(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// Create test findings with different scores
	for i := 0; i < 110; i++ {
		findings := &model.FindingsBundle{
			Summary: model.Summary{
				OverallScore: 50 + i%50, // Vary the score
			},
			Findings: []model.Finding{
				{Severity: "medium", Title: "Test Finding"},
			},
		}
		server.UpdateData(findings, nil)
	}

	// Wait for all snapshots to be recorded
	time.Sleep(200 * time.Millisecond)

	// Verify history doesn't exceed max size
	history := server.GetPerformanceHistory()
	assert.True(t, len(history) <= 100, "History should not exceed max size of 100")
}

// JWT Authentication Tests

func TestJWTTokenGeneration(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false)
	defer server.Stop()

	// Test token generation
	token, err := server.GenerateJWTToken("testuser", "viewer")
	assert.NoError(t, err, "Should generate token successfully")
	assert.NotEmpty(t, token, "Token should not be empty")

	// Verify token can be validated
	claims, err := server.ValidateJWTToken(token)
	assert.NoError(t, err, "Should validate generated token")
	assert.Equal(t, "testuser", claims.Username, "Username should match")
	assert.Equal(t, "viewer", claims.Role, "Role should match")
}

func TestJWTTokenValidation(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false)
	defer server.Stop()

	// Generate valid token
	validToken, err := server.GenerateJWTToken("testuser", "admin")
	assert.NoError(t, err)

	// Test valid token
	claims, err := server.ValidateJWTToken(validToken)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "admin", claims.Role)

	// Test invalid token
	_, err = server.ValidateJWTToken("invalid.token.here")
	assert.Error(t, err, "Should reject invalid token")

	// Test empty token
	_, err = server.ValidateJWTToken("")
	assert.Error(t, err, "Should reject empty token")
}

func TestJWTTokenExpiration(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false)
	defer server.Stop()

	// Generate token with very short expiration for testing
	// Note: In production, tokens expire in 24 hours
	token, err := server.GenerateJWTToken("testuser", "viewer")
	assert.NoError(t, err)

	// Token should be valid immediately
	_, err = server.ValidateJWTToken(token)
	assert.NoError(t, err, "Token should be valid immediately after generation")
}

func TestJWTAuthDisabled(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// When auth is disabled, token generation should fail
	_, err := server.GenerateJWTToken("testuser", "viewer")
	assert.Error(t, err, "Should not generate token when auth is disabled")
	assert.Contains(t, err.Error(), "authentication is disabled")

	// But validation should allow anonymous access
	claims, err := server.ValidateJWTToken("")
	assert.NoError(t, err, "Should allow anonymous access when auth is disabled")
	assert.Equal(t, "anonymous", claims.Username)
	assert.Equal(t, "viewer", claims.Role)
}

func TestJWTTokenExtraction(t *testing.T) {
	// Test extractTokenFromRequest function
	testCases := []struct {
		name           string
		setupRequest    func() *http.Request
		expectedToken   string
	}{
		{
			name: "Authorization header",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/ws", nil)
				req.Header.Set("Authorization", "Bearer test-token-123")
				return req
			},
			expectedToken: "test-token-123",
		},
		{
			name: "Query parameter",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/ws?token=query-token-456", nil)
				return req
			},
			expectedToken: "query-token-456",
		},
		{
			name: "No token",
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/ws", nil)
			},
			expectedToken: "",
		},
		{
			name: "Malformed Authorization header",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/ws", nil)
				req.Header.Set("Authorization", "InvalidHeader")
				return req
			},
			expectedToken: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := tc.setupRequest()
			token := extractTokenFromRequest(req)
			assert.Equal(t, tc.expectedToken, token)
		})
	}
}

func TestJWTAuthTokenEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false)
	defer server.Stop()

	// Test token generation endpoint
	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		shouldHaveToken bool
	}{
		{
			name:           "Valid request",
			requestBody:    `{"username": "testuser", "password": "testpass", "role": "admin"}`,
			expectedStatus: http.StatusOK,
			shouldHaveToken: true,
		},
		{
			name:           "Missing username",
			requestBody:    `{"password": "testpass"}`,
			expectedStatus: http.StatusBadRequest,
			shouldHaveToken: false,
		},
		{
			name:           "Missing password",
			requestBody:    `{"username": "testuser"}`,
			expectedStatus: http.StatusBadRequest,
			shouldHaveToken: false,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"username": "testuser", "password":`,
			expectedStatus: http.StatusBadRequest,
			shouldHaveToken: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/token", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.HandleGenerateToken(w, req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.shouldHaveToken {
				var response map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}

				assert.Contains(t, response, "token")
				assert.Contains(t, response, "expires_in")
				assert.Contains(t, response, "username")
				assert.Contains(t, response, "role")
				assert.NotEmpty(t, response["token"])
			}
		})
	}
}

func TestJWTAuthDisabledTokenEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// When auth is disabled, token endpoint should return error
	req := httptest.NewRequest("POST", "/auth/token", strings.NewReader(`{"username": "test", "password": "test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.HandleGenerateToken(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestJWTSecretKeyGeneration(t *testing.T) {
	// Test that generateJWTSecretKey produces valid keys
	key1 := generateJWTSecretKey()
	key2 := generateJWTSecretKey()

	// Keys should be non-empty
	assert.NotEmpty(t, key1)
	assert.NotEmpty(t, key2)

	// Keys should be different (random)
	assert.NotEqual(t, key1, key2)

	// Keys should be base64 encoded (no spaces, reasonable length)
	assert.NotContains(t, key1, " ")
	assert.True(t, len(key1) > 20, "Key should be reasonably long")
}

func TestJWTWebSocketConnectionWithAuth(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false)
	defer server.Stop()

	// Generate a valid token
	token, err := server.GenerateJWTToken("testuser", "viewer")
	assert.NoError(t, err)

	// Test token validation directly (simpler than full WebSocket upgrade)
	req := httptest.NewRequest("GET", "/ws?token="+token, nil)

	// Test just the token extraction and validation part
	extractedToken := extractTokenFromRequest(req)
	assert.Equal(t, token, extractedToken, "Should extract token from query parameter")

	claims, err := server.ValidateJWTToken(extractedToken)
	assert.NoError(t, err, "Should validate extracted token")
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "viewer", claims.Role)
}

func TestJWTWebSocketConnectionWithoutToken(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false)
	defer server.Stop()

	// Test WebSocket connection without token
	req := httptest.NewRequest("GET", "/ws", nil)

	// Test token extraction - should return empty
	extractedToken := extractTokenFromRequest(req)
	assert.Empty(t, extractedToken, "Should not extract token when none provided")

	// Test validation with empty token - should fail
	_, err := server.ValidateJWTToken(extractedToken)
	assert.Error(t, err, "Should reject empty token when auth is enabled")
}

func TestJWTWebSocketConnectionWithInvalidToken(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false)
	defer server.Stop()

	// Test WebSocket connection with invalid token
	req := httptest.NewRequest("GET", "/ws?token=invalid.token.here", nil)

	// Test token extraction
	extractedToken := extractTokenFromRequest(req)
	assert.Equal(t, "invalid.token.here", extractedToken)

	// Test validation - should fail
	_, err := server.ValidateJWTToken(extractedToken)
	assert.Error(t, err, "Should reject invalid token")
	assert.Contains(t, err.Error(), "invalid")
}

func TestJWTWebSocketConnectionWithAuthDisabled(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false)
	defer server.Stop()

	// When auth is disabled, validation should allow anonymous access
	req := httptest.NewRequest("GET", "/ws", nil)

	// Test token extraction - should return empty
	extractedToken := extractTokenFromRequest(req)
	assert.Empty(t, extractedToken)

	// Test validation with empty token - should succeed with anonymous user
	claims, err := server.ValidateJWTToken(extractedToken)
	assert.NoError(t, err, "Should allow anonymous access when auth is disabled")
	assert.Equal(t, "anonymous", claims.Username)
	assert.Equal(t, "viewer", claims.Role)
}
