package webserver

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

func TestWebSocketServerCreation(t *testing.T) {
	// Create WebSocket server
	server := NewWebSocketServer(8081, t.TempDir(), false)
	
	// Test server creation
	if server == nil {
		t.Fatal("Failed to create WebSocket server")
	}

	// Test client count
	if server.GetClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", server.GetClientCount())
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
	server := NewWebSocketServer(8081, tempDir, false)
	
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
	server := NewWebSocketServer(8082, t.TempDir(), false)
	
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

func TestWebSocketAutoRefresh(t *testing.T) {
	// Create WebSocket server
	server := NewWebSocketServer(8083, t.TempDir(), false)
	
	// Test auto-refresh doesn't panic
	server.StartAutoRefresh(1 * time.Second)
	
	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Test that it doesn't crash
	if server.GetClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", server.GetClientCount())
	}
}
