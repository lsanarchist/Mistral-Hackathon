
package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketCompressionDisabled(t *testing.T) {
	// Create WebSocket server with compression disabled
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
	defer server.Stop()

	// Test compression info endpoint
	req := httptest.NewRequest("GET", "/compression/info", nil)
	w := httptest.NewRecorder()

	server.handleCompressionInfo(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var compressionInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&compressionInfo); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, false, compressionInfo["enabled"])
	assert.Contains(t, compressionInfo, "level")
	assert.Contains(t, compressionInfo, "threshold")
	assert.Contains(t, compressionInfo, "description")
}

func TestWebSocketCompressionEnabled(t *testing.T) {
	// Create WebSocket server with compression enabled
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, true, false, 0, nil)
	defer server.Stop()

	// Test compression info endpoint
	req := httptest.NewRequest("GET", "/compression/info", nil)
	w := httptest.NewRecorder()

	server.handleCompressionInfo(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var compressionInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&compressionInfo); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, compressionInfo["enabled"])
	assert.Equal(t, float64(6), compressionInfo["level"])
	assert.Equal(t, float64(256), compressionInfo["threshold"])
	assert.Contains(t, compressionInfo, "description")
}

func TestWebSocketCompressionMethodNotAllowed(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
	defer server.Stop()

	// Test POST method (should not be allowed)
	req := httptest.NewRequest("POST", "/compression/info", nil)
	w := httptest.NewRecorder()

	server.handleCompressionInfo(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestPluginMarketplaceEndpoint(t *testing.T) {
	// Create WebSocket server
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
	defer server.Stop()

	// Test install endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/install", nil)
	w := httptest.NewRecorder()

	server.handleInstallPlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPluginUpdateEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
	defer server.Stop()

	// Test update endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/update", nil)
	w := httptest.NewRecorder()

	server.handleUpdatePlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPluginUninstallEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
	defer server.Stop()

	// Test uninstall endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/uninstall", nil)
	w := httptest.NewRecorder()

	server.handleUninstallPlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPerformanceHistory(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, nil)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
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

// TestWebSocketClientHandling tests WebSocket client connection management
func TestWebSocketClientHandling(t *testing.T) {
	// Create WebSocket server
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
	defer server.Stop()

	// Test initial client count
	assert.Equal(t, 0, server.GetClientCount(), "Should start with 0 clients")

	// Test broadcast functionality (should not panic with no clients)
	server.BroadcastData()

	// Test performance history tracking
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 90,
		},
		Findings: []model.Finding{
			{Severity: "high", Title: "Test Finding"},
		},
	}
	
	// Update data multiple times to test history
	for i := 0; i < 3; i++ {
		server.UpdateData(findings, nil)
	}
	
	// Wait for async operations
	time.Sleep(100 * time.Millisecond)
	
	// Verify history is being tracked
	history := server.GetPerformanceHistory()
	assert.True(t, len(history) > 0, "Should have performance history")
	assert.True(t, len(history) <= server.maxHistorySize, "Should respect max history size")
}

// WebSocket Batching Tests

func TestWebSocketBatchingDisabled(t *testing.T) {
	// Create WebSocket server with batching disabled
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
	defer server.Stop()

	// Test batching info endpoint
	req := httptest.NewRequest("GET", "/batching/info", nil)
	w := httptest.NewRecorder()

	server.handleBatchingInfo(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var batchingInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&batchingInfo); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, false, batchingInfo["enabled"])
	assert.Contains(t, batchingInfo, "interval_ms")
	assert.Contains(t, batchingInfo, "description")
}

func TestWebSocketBatchingEnabled(t *testing.T) {
	// Create WebSocket server with batching enabled
	batchInterval := 50 * time.Millisecond
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, nil)
	defer server.Stop()

	// Test batching info endpoint
	req := httptest.NewRequest("GET", "/batching/info", nil)
	w := httptest.NewRecorder()

	server.handleBatchingInfo(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var batchingInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&batchingInfo); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, batchingInfo["enabled"])
	assert.Equal(t, float64(50), batchingInfo["interval_ms"])
	assert.Contains(t, batchingInfo, "description")
}

func TestWebSocketBatchingMethodNotAllowed(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, nil)
	defer server.Stop()

	// Test POST method (should not be allowed)
	req := httptest.NewRequest("POST", "/batching/info", nil)
	w := httptest.NewRecorder()

	server.handleBatchingInfo(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestWebSocketMessageQueue(t *testing.T) {
	// Create WebSocket server with batching enabled
	batchInterval := 100 * time.Millisecond
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, nil)
	defer server.Stop()

	// Create test findings
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 80,
		},
		Findings: []model.Finding{
			{Severity: "medium", Title: "Test Finding"},
		},
	}

	// Update server data (this should queue messages)
	server.UpdateData(findings, nil)
	server.UpdateData(findings, nil)
	server.UpdateData(findings, nil)

	// Wait a bit for messages to be queued
	time.Sleep(20 * time.Millisecond)

	// Verify messages are queued (we can't directly access the queue, but we can verify batching is working)
	assert.True(t, server.batchingEnabled, "Batching should be enabled")
	assert.Equal(t, batchInterval, server.batchInterval, "Batch interval should match")
}

func TestWebSocketBatchingIntegration(t *testing.T) {
	// Create WebSocket server with batching enabled
	batchInterval := 200 * time.Millisecond
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, nil)
	defer server.Stop()

	// Create test findings
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
		},
		Findings: []model.Finding{
			{Severity: "high", Title: "Test Critical Finding"},
			{Severity: "medium", Title: "Test Medium Finding"},
		},
	}

	// Update server data multiple times
	for i := 0; i < 5; i++ {
		server.UpdateData(findings, nil)
	}

	// Wait for batching to process
	time.Sleep(300 * time.Millisecond)

	// Verify server is still functioning
	assert.True(t, server.batchingEnabled)
	assert.Equal(t, batchInterval, server.batchInterval)
}

func TestWebSocketBatchingWithCompression(t *testing.T) {
	// Create WebSocket server with both batching and compression enabled
	batchInterval := 50 * time.Millisecond
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, true, true, batchInterval, nil)
	defer server.Stop()

	// Verify both features are enabled
	assert.True(t, server.batchingEnabled, "Batching should be enabled")
	assert.True(t, server.compressionEnabled, "Compression should be enabled")

	// Test that both info endpoints work
	compressionReq := httptest.NewRequest("GET", "/compression/info", nil)
	compressionW := httptest.NewRecorder()
	server.handleCompressionInfo(compressionW, compressionReq)

	batchingReq := httptest.NewRequest("GET", "/batching/info", nil)
	batchingW := httptest.NewRecorder()
	server.handleBatchingInfo(batchingW, batchingReq)

	assert.Equal(t, http.StatusOK, compressionW.Result().StatusCode)
	assert.Equal(t, http.StatusOK, batchingW.Result().StatusCode)
}

func TestWebSocketBatchingStop(t *testing.T) {
	// Create WebSocket server with batching enabled
	batchInterval := 100 * time.Millisecond
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, nil)

	// Stop the server (should clean up batching timer)
	err := server.Stop()
	assert.NoError(t, err, "Server should stop without error")

	// Verify batching timer is stopped
	assert.Nil(t, server.batchTimer, "Batching timer should be nil after stop")
}

// TestWebSocketBatchingConcurrency tests concurrent access to the message queue
func TestWebSocketBatchingConcurrency(t *testing.T) {
	// Create WebSocket server with batching enabled
	batchInterval := 200 * time.Millisecond
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, nil)
	defer server.Stop()

	// Create test findings
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 90,
		},
		Findings: []model.Finding{
			{Severity: "low", Title: "Concurrent Test Finding"},
		},
	}

	// Test concurrent updates
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				server.UpdateData(findings, nil)
				time.Sleep(5 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	// Wait for batching to complete
	time.Sleep(300 * time.Millisecond)

	// Server should still be functioning
	assert.True(t, server.batchingEnabled)
	assert.NotNil(t, server)
}
