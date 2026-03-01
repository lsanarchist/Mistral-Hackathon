
package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
