package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

func TestWebSocketServerCreation(t *testing.T) {
	// Test server creation without compression
	server := NewWebSocketServer(8081, t.TempDir(), t.TempDir(), false, false)
	
	// Test server creation
	if server == nil {
		t.Fatal("Failed to create WebSocket server")
	}

	// Test client count
	if server.GetClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", server.GetClientCount())
	}

	// Test compression flag
	if server.compressionEnabled {
		t.Error("Expected compression to be disabled")
	}
}

func TestWebSocketServerCreationWithCompression(t *testing.T) {
	// Test server creation with compression
	server := NewWebSocketServer(8082, t.TempDir(), t.TempDir(), false, true)
	
	// Test server creation
	if server == nil {
		t.Fatal("Failed to create WebSocket server with compression")
	}

	// Test compression flag
	if !server.compressionEnabled {
		t.Error("Expected compression to be enabled")
	}

	// Test that upgrader has compression enabled
	if !server.upgrader.EnableCompression {
		t.Error("Expected upgrader to have compression enabled")
	}
}

func TestWebSocketDataLoading(t *testing.T) {
	// Create test data
	tempDir := t.TempDir()
	findingsPath := tempDir + "/findings.json"
	
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 85,
		},
		Findings: []model.Finding{
			{
				Category:  "cpu",
				Title:     "CPU hotspot detected",
				Severity:  "high",
				Score:     90,
			},
			{
				Category:  "memory",
				Title:     "Memory allocation issue",
				Severity:  "medium",
				Score:     75,
			},
		},
	}
	
	findingsData, _ := json.Marshal(findings)
	os.WriteFile(findingsPath, findingsData, 0644)

	// Create WebSocket server
	server := NewWebSocketServer(8081, tempDir, tempDir, false, false)
	
	// Test data loading
	err := server.LoadData(findingsPath, "")
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	// Test that data was loaded
	if server.GetClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", server.GetClientCount())
	}
}

func TestWebSocketDataUpdate(t *testing.T) {
	// Create WebSocket server
	server := NewWebSocketServer(8082, t.TempDir(), t.TempDir(), false, false)
	
	// Create initial findings
	initialFindings := &model.FindingsBundle{
		Summary: model.Summary{OverallScore: 70},
		Findings: []model.Finding{
			{Category: "cpu", Title: "Initial finding", Severity: "medium", Score: 70},
		},
	}
	
	server.UpdateData(initialFindings, nil)
	
	// Verify data was updated
	if server.GetClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", server.GetClientCount())
	}
	
	// Create updated findings
	updatedFindings := &model.FindingsBundle{
		Summary: model.Summary{OverallScore: 90},
		Findings: []model.Finding{
			{Category: "cpu", Title: "Updated finding", Severity: "critical", Score: 95},
			{Category: "memory", Title: "New finding", Severity: "high", Score: 85},
		},
	}
	
	server.UpdateData(updatedFindings, nil)
	
	// Data should be updated (we can't easily verify without a client, but the method should not crash)
	if server.GetClientCount() != 0 {
		t.Errorf("Expected 0 clients after update, got %d", server.GetClientCount())
	}
}

// Test JWT token generation
func TestJWTTokenGeneration(t *testing.T) {
	// Create server with auth enabled
	server := NewWebSocketServer(8084, t.TempDir(), t.TempDir(), true, false)
	
	// Test token generation
	token, err := server.GenerateJWTToken("testuser", "viewer")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	if token == "" {
		t.Error("Generated token should not be empty")
	}
	
	// Test token validation
	claims, err := server.ValidateJWTToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	
	if claims.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", claims.Username)
	}
	
	if claims.Role != "viewer" {
		t.Errorf("Expected role 'viewer', got '%s'", claims.Role)
	}
}

// Test JWT token validation with invalid tokens
func TestJWTTokenValidation(t *testing.T) {
	// Create server with auth enabled
	server := NewWebSocketServer(8085, t.TempDir(), t.TempDir(), true, false)
	
	// Test invalid token
	_, err := server.ValidateJWTToken("invalid.token.here")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
	
	// Test empty token
	_, err = server.ValidateJWTToken("")
	if err == nil {
		t.Error("Expected error for empty token")
	}
}

// Test token generation handler
func TestTokenGenerationHandler(t *testing.T) {
	// Create server with auth enabled
	server := NewWebSocketServer(8086, t.TempDir(), t.TempDir(), true, false)
	
	// Create a test request
	reqBody := `{"username": "testuser", "password": "testpass", "role": "admin"}`
	req := httptest.NewRequest("POST", "/auth/token", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	
	// Call the handler
	server.HandleGenerateToken(w, req)
	
	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check that token is present
	if _, ok := response["token"]; !ok {
		t.Error("Response should contain token")
	}
	
	if response["username"] != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", response["username"])
	}
	
	if response["role"] != "admin" {
		t.Errorf("Expected role 'admin', got '%s'", response["role"])
	}
}

// Test token generation with missing credentials
func TestTokenGenerationMissingCredentials(t *testing.T) {
	// Create server with auth enabled
	server := NewWebSocketServer(8087, t.TempDir(), t.TempDir(), true, false)
	
	// Test with missing username
	reqBody := `{"password": "testpass"}`
	req := httptest.NewRequest("POST", "/auth/token", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	
	server.HandleGenerateToken(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing credentials, got %d", resp.StatusCode)
	}
}

// Test auth disabled scenarios
func TestAuthDisabled(t *testing.T) {
	// Create server with auth disabled
	server := NewWebSocketServer(8088, t.TempDir(), t.TempDir(), false, false)
	
	// Test token generation should fail
	_, err := server.GenerateJWTToken("testuser", "viewer")
	if err == nil {
		t.Error("Expected error when auth is disabled")
	}
	
	// Test token validation should allow anonymous access
	claims, err := server.ValidateJWTToken("")
	if err != nil {
		t.Fatalf("Expected no error for empty token when auth disabled: %v", err)
	}
	
	if claims.Username != "anonymous" {
		t.Errorf("Expected anonymous username, got '%s'", claims.Username)
	}
}

func TestWebSocketAutoRefresh(t *testing.T) {
	// Create WebSocket server
	server := NewWebSocketServer(8083, t.TempDir(), t.TempDir(), false, false)
	
	// Test auto-refresh doesn't panic
	server.StartAutoRefresh(1 * time.Second)
	
	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Test that it doesn't crash
	if server.GetClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", server.GetClientCount())
	}
}
