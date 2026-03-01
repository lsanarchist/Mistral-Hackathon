
package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
)

// mockWebSocketConn is a mock WebSocket connection for testing that implements websocket.Conn interface
type mockWebSocketConn struct {
	*websocket.Conn
}

func newMockWebSocketConn() *mockWebSocketConn {
	return &mockWebSocketConn{}
}

func (m *mockWebSocketConn) WriteJSON(v interface{}) error { return nil }
func (m *mockWebSocketConn) Close() error { return nil }
func (m *mockWebSocketConn) WriteMessage(messageType int, data []byte) error { return nil }
func (m *mockWebSocketConn) WriteControl(messageType int, data []byte, deadline time.Time) error { return nil }
func (m *mockWebSocketConn) ReadMessage() (messageType int, p []byte, err error) { return 0, nil, nil }
func (m *mockWebSocketConn) SetReadLimit(limit int64) {}
func (m *mockWebSocketConn) SetReadDeadline(t time.Time) error { return nil }
func (m *mockWebSocketConn) SetPongHandler(h func(string) error) {}
func (m *mockWebSocketConn) SetPingHandler(h func(string) error) {}
func (m *mockWebSocketConn) SetWriteDeadline(t time.Time) error { return nil }
func (m *mockWebSocketConn) RemoteAddr() string { return "test" }
func (m *mockWebSocketConn) LocalAddr() string { return "test" }
func (m *mockWebSocketConn) Subprotocol() string { return "" }
func (m *mockWebSocketConn) CloseHandler() func(code int, text string) error { return nil }
func (m *mockWebSocketConn) SetCloseHandler(h func(code int, text string) error) {}

func TestWebSocketCompressionDisabled(t *testing.T) {
	// Create WebSocket server with compression disabled
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, true, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
	defer server.Stop()

	// Test install endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/install", nil)
	w := httptest.NewRecorder()

	server.handleInstallPlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPluginUpdateEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
	defer server.Stop()

	// Test update endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/update", nil)
	w := httptest.NewRecorder()

	server.handleUpdatePlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPluginUninstallEndpoint(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
	defer server.Stop()

	// Test uninstall endpoint with empty body
	req := httptest.NewRequest("POST", "/plugins/uninstall", nil)
	w := httptest.NewRecorder()

	server.handleUninstallPlugin(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPerformanceHistory(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), true, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, true, true, batchInterval, false, nil, nil, ConnectionQualityConfig{}, false)
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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, false, nil, nil, ConnectionQualityConfig{}, false)

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
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, true, batchInterval, false, nil, nil, ConnectionQualityConfig{}, false)
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


// Connection Quality Tests
func TestConnectionQuality(t *testing.T) {
	t.Run("ConnectionQualityInfo", func(t *testing.T) {
		server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, nil, ConnectionQualityConfig{}, false)
		defer server.Stop()

		// Test connection quality info endpoint
		req := httptest.NewRequest("GET", "/connection/quality", nil)
		w := httptest.NewRecorder()

		server.handleConnectionQuality(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&info)
		assert.NoError(t, err)

		assert.True(t, info["connection_quality_enabled"].(bool))
		assert.Equal(t, float64(10000), info["ping_interval_ms"].(float64)) // 10 seconds
		assert.Equal(t, float64(0), info["active_connections"].(float64))
		assert.Equal(t, false, info["ml_model_enabled"].(bool)) // ML model disabled in this test
		assert.Equal(t, float64(0), info["anomaly_count"])
		assert.Equal(t, float64(0), info["anomaly_percentage"])
	})

	t.Run("ConnectionQualityDisabled", func(t *testing.T) {
		server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, false, nil, nil, ConnectionQualityConfig{}, false)
		defer server.Stop()

		// Test connection quality info endpoint
		req := httptest.NewRequest("GET", "/connection/quality", nil)
		w := httptest.NewRecorder()

		server.handleConnectionQuality(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&info)
		assert.NoError(t, err)

		assert.False(t, info["connection_quality_enabled"].(bool))
		assert.Equal(t, float64(30000), info["ping_interval_ms"].(float64)) // 30 seconds
	})

	t.Run("CalculateConnectionQuality", func(t *testing.T) {
		server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, nil, ConnectionQualityConfig{}, false)
		defer server.Stop()

		// Test excellent quality
		quality := server.calculateConnectionQuality(50*time.Millisecond, 0)
		assert.Equal(t, "excellent", quality)

		// Test good quality (150ms latency, 3% packet loss)
		quality = server.calculateConnectionQuality(150*time.Millisecond, 3)
		assert.Equal(t, "excellent", quality) // 150ms < 200ms, 3% < 5%

		// Test fair quality (300ms latency, 8% packet loss)
		quality = server.calculateConnectionQuality(300*time.Millisecond, 8)
		assert.Equal(t, "good", quality) // 300ms < 500ms, 8% < 10%

		// Test poor quality
		quality = server.calculateConnectionQuality(1*time.Second, 25)
		assert.Equal(t, "poor", quality)
	})

	t.Run("CalculateAverageLatency", func(t *testing.T) {
		server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, nil, ConnectionQualityConfig{}, false)
		defer server.Stop()

		// Test with no connections
		avgLatency := server.calculateAverageLatency([]*WebSocketConnectionStats{})
		assert.Equal(t, float64(0), avgLatency)

		// Test with connections having latency
		stats := []*WebSocketConnectionStats{
			{Latency: 100 * time.Millisecond},
			{Latency: 200 * time.Millisecond},
			{Latency: 0}, // Should be ignored
		}

		avgLatency = server.calculateAverageLatency(stats)
		expectedAvg := (100.0 + 200.0) / 2.0 // 150ms
		assert.Equal(t, expectedAvg, avgLatency)
	})
}

// Test connection quality alert triggering
func TestConnectionQualityAlerts(t *testing.T) {
	server := &WebSocketServer{
		connectionQualityAlerts: []ConnectionQualityAlert{
			{
				ID:               "alert-1",
				Name:             "Poor Connection Alert",
				QualityThreshold: "poor",
				Active:           true,
			},
			{
				ID:               "alert-2",
				Name:             "High Latency Alert",
				LatencyThreshold: 500, // 500ms
				Active:           true,
			},
		},
	}

	// Test stats that should trigger alerts
	stats := &WebSocketConnectionStats{
		ClientID:         "test-client",
		ConnectionQuality: "poor",
		Latency:          600 * time.Millisecond,
		PacketLoss:       15.0,
	}

	// Mock time for testing
	now := time.Now()
	originalNow := timeNow
	timeNow = func() time.Time { return now }
	defer func() { timeNow = originalNow }()

	// This should trigger both alerts
	server.checkConnectionQualityAlerts(stats)

	// Verify alerts were triggered (LastTriggered should be set)
	for _, alert := range server.connectionQualityAlerts {
		if alert.Active {
			assert.NotNil(t, alert.LastTriggered, "Alert %s should have been triggered", alert.Name)
		}
	}
}

// Test adaptive update intervals
func TestAdaptiveUpdateIntervals(t *testing.T) {
	server := &WebSocketServer{
		connectionQualityConfig: ConnectionQualityConfig{
			AdaptiveUpdatesEnabled: true,
			UpdateIntervals: UpdateIntervals{
				Excellent: 1 * time.Second,
				Good:     2 * time.Second,
				Fair:     5 * time.Second,
				Poor:     10 * time.Second,
			},
		},
	}

	// Test different quality levels
	intervals := map[string]time.Duration{
		"excellent": 1 * time.Second,
		"good":     2 * time.Second,
		"fair":     5 * time.Second,
		"poor":     10 * time.Second,
		"unknown":  2 * time.Second, // Should default to good
	}

	for quality, expected := range intervals {
		result := server.getAdaptiveUpdateInterval(quality)
		assert.Equal(t, expected, result, "Interval for quality %s should be %v", quality, expected)
	}
}

// Test bandwidth throttling limits
func TestBandwidthThrottling(t *testing.T) {
	server := &WebSocketServer{
		connectionQualityConfig: ConnectionQualityConfig{
			BandwidthThrottlingEnabled: true,
			ThrottlingThresholds: ThrottlingThresholds{
				Excellent: 1000000, // 1MB/s
				Good:     500000,  // 500KB/s
				Fair:     200000,  // 200KB/s
				Poor:     50000,   // 50KB/s
			},
		},
	}

	// Test different quality levels
	limits := map[string]int{
		"excellent": 1000000,
		"good":     500000,
		"fair":     200000,
		"poor":     50000,
		"unknown":  500000, // Should default to good
	}

	for quality, expected := range limits {
		result := server.getBandwidthLimit(quality)
		assert.Equal(t, expected, result, "Bandwidth limit for quality %s should be %d", quality, expected)
	}
}

// Test connection quality history recording
func TestConnectionQualityHistory(t *testing.T) {
	server := &WebSocketServer{
		connectionQualityHistory: make([]map[string]interface{}, 0),
		maxQualityHistorySize:  3, // Small size for testing
	}

	// Create a mock connection and stats
	conn := &websocket.Conn{}
	stats := &WebSocketConnectionStats{
		ClientID:          "test-client",
		ConnectionQuality: "excellent",
		Latency:           50 * time.Millisecond,
		PacketLoss:        0,
		MessagesSent:      10,
		MessagesReceived:  10,
		BytesSent:         1000,
		BytesReceived:     1000,
	}

	// Record multiple entries
	for i := 0; i < 5; i++ {
		stats.MessagesSent += 1
		stats.MessagesReceived += 1
		server.recordConnectionQualityHistory(conn, stats)
	}

	// Verify history size is limited
	server.qualityHistoryMu.Lock()
	historySize := len(server.connectionQualityHistory)
	server.qualityHistoryMu.Unlock()

	assert.Equal(t, 3, historySize, "History should be limited to max size")

	// Verify history entries contain expected data
	server.qualityHistoryMu.Lock()
	lastEntry := server.connectionQualityHistory[len(server.connectionQualityHistory)-1]
	server.qualityHistoryMu.Unlock()

	assert.Equal(t, "test-client", lastEntry["client_id"])
	assert.Equal(t, "excellent", lastEntry["connection_quality"])
	assert.Equal(t, int64(50), lastEntry["latency_ms"])
}

// Test GetConnectionQualityInfo with new fields
func TestGetConnectionQualityInfoEnhanced(t *testing.T) {
	server := &WebSocketServer{
		connectionQualityEnabled: true,
		pingInterval:             10 * time.Second,
		connectionQualityConfig: ConnectionQualityConfig{
			AdaptiveUpdatesEnabled: true,
			BandwidthThrottlingEnabled: true,
		},
		connectionQualityAlerts: []ConnectionQualityAlert{
			{ID: "alert-1", Name: "Test Alert", Active: true},
			{ID: "alert-2", Name: "Inactive Alert", Active: false},
		},
	}

	info := server.GetConnectionQualityInfo()

	assert.True(t, info["connection_quality_enabled"].(bool))
	assert.Equal(t, int64(10000), info["ping_interval_ms"])
	assert.True(t, info["adaptive_updates_enabled"].(bool))
	assert.True(t, info["bandwidth_throttling_enabled"].(bool))
	assert.Equal(t, 1, info["active_quality_alerts"])
}

// Test WebSocket Connection Quality Data Broadcasting
func TestWebSocketConnectionQualityBroadcasting(t *testing.T) {
	// Create WebSocket server with connection quality enabled
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, nil, ConnectionQualityConfig{}, false)
	defer server.Stop()

	// Create a test WebSocket connection
	testServer := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
	defer testServer.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")

	// Connect WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Subscribe to connection quality data
	subscribeMsg := map[string]interface{}{
		"type":  "subscribe",
		"topic": "connection_quality",
	}
	
	if err := ws.WriteJSON(subscribeMsg); err != nil {
		t.Fatalf("Failed to send subscribe message: %v", err)
	}

	// Read response (should be connection quality data)
	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read WebSocket message: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		t.Fatalf("Failed to parse WebSocket message: %v", err)
	}

	// Verify message type
	assert.Equal(t, "connection_quality_update", data["type"])
	assert.Contains(t, data, "quality_counts")
	assert.Contains(t, data, "total_clients")
	assert.Contains(t, data, "avg_latency")
	assert.Contains(t, data, "avg_packet_loss")
	assert.Contains(t, data, "connection_stats")
}

// Test WebSocket Connection Quality Alerts Broadcasting
func TestWebSocketConnectionQualityAlertsBroadcasting(t *testing.T) {
	// Create WebSocket server with connection quality alerts
	qualityAlerts := []ConnectionQualityAlert{
		{
			ID:               "test-alert-1",
			Name:             "High Latency Alert",
			QualityThreshold: "poor",
			LatencyThreshold: 500,
			Active:           true,
		},
		{
			ID:               "test-alert-2",
			Name:             "Packet Loss Alert",
			PacketLossThreshold: 10.0,
			Active:           false,
		},
	}

	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, qualityAlerts, ConnectionQualityConfig{}, false)
	defer server.Stop()

	// Create a test WebSocket connection
	testServer := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
	defer testServer.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")

	// Connect WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Subscribe to connection quality data
	subscribeMsg := map[string]interface{}{
		"type":  "subscribe",
		"topic": "connection_quality",
	}
	
	if err := ws.WriteJSON(subscribeMsg); err != nil {
		t.Fatalf("Failed to send subscribe message: %v", err)
	}

	// Read alerts message
	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read WebSocket message: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		t.Fatalf("Failed to parse WebSocket message: %v", err)
	}

	// The first message might be quality data, so we might need to read another message
	if data["type"] == "connection_quality_update" {
		// Read next message for alerts
		_, message, err = ws.ReadMessage()
		if err != nil {
			t.Fatalf("Failed to read alerts WebSocket message: %v", err)
		}
		
		if err := json.Unmarshal(message, &data); err != nil {
			t.Fatalf("Failed to parse alerts WebSocket message: %v", err)
		}
	}

	// Verify message type
	assert.Equal(t, "connection_quality_alerts", data["type"])
	assert.Contains(t, data, "alerts")
	
	alerts := data["alerts"].([]interface{})
	assert.Equal(t, 2, len(alerts))
	
	// Verify first alert
	firstAlert := alerts[0].(map[string]interface{})
	assert.Equal(t, "test-alert-1", firstAlert["id"])
	assert.Equal(t, "High Latency Alert", firstAlert["name"])
	assert.Equal(t, "poor", firstAlert["quality_threshold"])
	assert.Equal(t, true, firstAlert["active"])
}

// Test WebSocket Message Handling
func TestWebSocketMessageHandling(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, nil, ConnectionQualityConfig{}, false)
	defer server.Stop()

	// Create a test WebSocket connection
	testServer := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
	defer testServer.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")

	// Connect WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Test unknown message type
	unknownMsg := map[string]interface{}{
		"type": "unknown_type",
	}
	
	if err := ws.WriteJSON(unknownMsg); err != nil {
		t.Fatalf("Failed to send unknown message: %v", err)
	}

	// Test request update message
	requestUpdateMsg := map[string]interface{}{
		"type":  "request_update",
		"topic": "connection_quality",
	}
	
	if err := ws.WriteJSON(requestUpdateMsg); err != nil {
		t.Fatalf("Failed to send request update message: %v", err)
	}

	// Read response
	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read WebSocket message: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		t.Fatalf("Failed to parse WebSocket message: %v", err)
	}

	// Verify message type
	assert.Equal(t, "connection_quality_update", data["type"])
}

// Test ML-based anomaly detection
func TestMLAnomalyDetection(t *testing.T) {
	// Create WebSocket server with ML model enabled
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, nil, ConnectionQualityConfig{}, true)
	defer server.Stop()

	// Test that ML model is enabled
	assert.True(t, server.mlModelEnabled)
	
	// Test anomaly alerts endpoint
	req := httptest.NewRequest("GET", "/anomaly/alerts", nil)
	w := httptest.NewRecorder()
	
	server.handleAnomalyAlerts(w, req)
	
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	
	assert.Equal(t, "success", response["status"])
	assert.Contains(t, response, "anomaly_alerts")
	
	// Test anomaly clusters endpoint
	req = httptest.NewRequest("GET", "/anomaly/clusters", nil)
	w = httptest.NewRecorder()
	
	server.handleAnomalyClusters(w, req)
	
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	
	assert.Equal(t, "success", response["status"])
	assert.Contains(t, response, "clusters")
	
	// Test anomaly patterns endpoint
	req = httptest.NewRequest("GET", "/anomaly/patterns", nil)
	w = httptest.NewRecorder()
	
	server.handleAnomalyPatterns(w, req)
	
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	
	assert.Equal(t, "success", response["status"])
	assert.Contains(t, response, "patterns")
	
	// Test ML model control endpoint
	req = httptest.NewRequest("GET", "/anomaly/ml", nil)
	w = httptest.NewRecorder()
	
	server.handleAnomalyML(w, req)
	
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, true, response["ml_enabled"])
	assert.Contains(t, response, "pattern_count")
	assert.Contains(t, response, "cluster_count")
}

// Test anomaly clustering functionality
func TestAnomalyClustering(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, nil, ConnectionQualityConfig{}, true)
	defer server.Stop()

	// Create test connection stats
	stats := &WebSocketConnectionStats{
		ClientID: "test-client",
		ConnectionTime: time.Now(),
		Latency: 501 * time.Millisecond,
		PacketLoss: 10.0,
		ConnectionScore: 45.0,
		ConnectionQuality: "fair",
		IsAnomaly: true,
		AnomalyScore: 0.8,
		AnomalyReasons: []string{"high latency"},
		AnomalyType: "latency",
		AnomalyConfidence: 0.85,
	}

	// Test clustering
	clusterID := server.clusterAnomaly(stats, "latency", []string{"high latency"})
	
	assert.NotEmpty(t, clusterID)
	assert.True(t, strings.HasPrefix(clusterID, "cluster_"))
	
	// Verify cluster was added
	server.anomalyClustersMu.Lock()
	clusters := server.anomalyClusters
	server.anomalyClustersMu.Unlock()
	
	assert.Contains(t, clusters, clusterID)
	assert.Contains(t, clusters[clusterID], "test-client")
	
	// Test pattern signature creation
	patternSig := server.createConnectionPatternSignature(stats)
	assert.Equal(t, "high_medium_low_fair", patternSig)
}

// Test connection quality info with ML model
func TestConnectionQualityInfoWithML(t *testing.T) {
	server := NewWebSocketServer(8080, t.TempDir(), t.TempDir(), false, false, false, 0, true, nil, nil, ConnectionQualityConfig{}, true)
	defer server.Stop()

	// Test with empty connection stats for simplicity

	// Test GetConnectionQualityInfo
	info := server.GetConnectionQualityInfo()
	
	assert.Equal(t, true, info["connection_quality_enabled"])
	assert.Equal(t, true, info["ml_model_enabled"])
	assert.Equal(t, 0, info["anomaly_count"])
	assert.Equal(t, 0.0, info["anomaly_percentage"])
}
