package webserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformanceAlerts(t *testing.T) {
	server := NewWebSocketServer(8080, "./testdata", "./testdata", false, false, false, 0, false, nil, nil, ConnectionQualityConfig{})
	defer server.Stop()

	// Test GET empty alerts
	req, _ := http.NewRequest("GET", "/performance/alerts", nil)
	rr := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var alerts []PerformanceAlert
	json.Unmarshal(rr.Body.Bytes(), &alerts)
	assert.Empty(t, alerts)

	// Test POST new alert
	newAlert := PerformanceAlert{
		Name:       "Test Alert",
		Metric:     "critical",
		Threshold:  5,
		Comparator: ">",
		Active:     true,
	}

	body, _ := json.Marshal(newAlert)
	req, _ = http.NewRequest("POST", "/performance/alerts", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var createdAlert PerformanceAlert
	json.Unmarshal(rr.Body.Bytes(), &createdAlert)
	assert.NotEmpty(t, createdAlert.ID)
	assert.Equal(t, "Test Alert", createdAlert.Name)

	// Test GET alerts after adding
	req, _ = http.NewRequest("GET", "/performance/alerts", nil)
	rr = httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &alerts)
	assert.Len(t, alerts, 1)
	assert.Equal(t, "Test Alert", alerts[0].Name)

	// Test DELETE alert
	req, _ = http.NewRequest("DELETE", "/performance/alerts?id="+createdAlert.ID, nil)
	rr = httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify deletion
	req, _ = http.NewRequest("GET", "/performance/alerts", nil)
	rr = httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &alerts)
	assert.Empty(t, alerts)
}

func TestPerformanceAnnotations(t *testing.T) {
	server := NewWebSocketServer(8080, "./testdata", "./testdata", false, false, false, 0, false, nil, nil, ConnectionQualityConfig{})
	defer server.Stop()

	// Test GET empty annotations
	req, _ := http.NewRequest("GET", "/performance/annotations", nil)
	rr := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var annotations []PerformanceAnnotation
	json.Unmarshal(rr.Body.Bytes(), &annotations)
	assert.Empty(t, annotations)

	// Test POST new annotation
	newAnnotation := PerformanceAnnotation{
		Title:   "Deployment v1.0.0",
		Content: "Deployed new version with performance improvements",
		Type:    "deployment",
	}

	body, _ := json.Marshal(newAnnotation)
	req, _ = http.NewRequest("POST", "/performance/annotations", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var createdAnnotation PerformanceAnnotation
	json.Unmarshal(rr.Body.Bytes(), &createdAnnotation)
	assert.NotEmpty(t, createdAnnotation.ID)
	assert.NotZero(t, createdAnnotation.Timestamp)
	assert.Equal(t, "Deployment v1.0.0", createdAnnotation.Title)

	// Test GET annotations after adding
	req, _ = http.NewRequest("GET", "/performance/annotations", nil)
	rr = httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &annotations)
	assert.Len(t, annotations, 1)
	assert.Equal(t, "Deployment v1.0.0", annotations[0].Title)

	// Test DELETE annotation
	req, _ = http.NewRequest("DELETE", "/performance/annotations?id="+createdAnnotation.ID, nil)
	rr = httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify deletion
	req, _ = http.NewRequest("GET", "/performance/annotations", nil)
	rr = httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &annotations)
	assert.Empty(t, annotations)
}

func TestPerformanceExport(t *testing.T) {
	server := NewWebSocketServer(8080, "./testdata", "./testdata", false, false, false, 0, false, nil, nil, ConnectionQualityConfig{})
	defer server.Stop()

	// Load some test data
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 85,
		},
		Findings: []model.Finding{
			{Severity: "critical"},
			{Severity: "high"},
			{Severity: "medium"},
		},
	}
	server.UpdateData(findings, nil)

	// Test JSON export
	req, _ := http.NewRequest("GET", "/performance/export?format=json", nil)
	rr := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, rr.Header().Get("Content-Disposition"), "performance_export.json")

	var exportData []PerformanceSnapshot
	json.Unmarshal(rr.Body.Bytes(), &exportData)
	assert.NotEmpty(t, exportData)
	assert.Equal(t, 85, exportData[0].OverallScore)
	assert.Equal(t, 1, exportData[0].CriticalCount)
	assert.Equal(t, 1, exportData[0].HighCount)
}

func TestPerformanceCompare(t *testing.T) {
	server := NewWebSocketServer(8080, "./testdata", "./testdata", false, false, false, 0, false, nil, nil, ConnectionQualityConfig{})
	defer server.Stop()

	compareRequest := map[string]interface{}{
		"applications": []map[string]interface{}{
			{
				"name": "App A",
				"data": []map[string]interface{}{
					{
						"overall_score":    90,
						"critical_count":   1,
						"high_count":       2,
						"total_findings":   10,
					},
				},
			},
			{
				"name": "App B",
				"data": []map[string]interface{}{
					{
						"overall_score":    75,
						"critical_count":   3,
						"high_count":       5,
						"total_findings":   15,
					},
				},
			},
		},
		"metrics": []string{"overall_score", "critical_count"},
	}

	body, _ := json.Marshal(compareRequest)
	req, _ := http.NewRequest("POST", "/performance/compare", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	
	assert.Contains(t, response, "applications")
	assert.Contains(t, response, "analysis")
	
	apps := response["applications"].(map[string]interface{})
	assert.Contains(t, apps, "App A")
	assert.Contains(t, apps, "App B")
	
	analysis := response["analysis"].(map[string]interface{})
	assert.Contains(t, analysis, "overall_score")
	assert.Contains(t, analysis, "critical_count")
}

func TestLoadPerformanceAlertsFromFile(t *testing.T) {
	// Create a temporary alerts file
	alerts := []PerformanceAlert{
		{
			ID:         "alert-1",
			Name:       "High Critical Alert",
			Metric:     "critical",
			Threshold:  10,
			Comparator: ">",
			Active:     true,
		},
		{
			ID:         "alert-2",
			Name:       "Low Score Alert",
			Metric:     "score",
			Threshold:  50,
			Comparator: "<",
			Active:     true,
		},
	}

	fileContent, _ := json.Marshal(alerts)
	tmpFile := t.TempDir() + "/alerts.json"
	testFile, _ := os.Create(tmpFile)
	defer testFile.Close()
	defer os.Remove(tmpFile)

	_, err := testFile.Write(fileContent)
	require.NoError(t, err)

	// Test loading from file
	loadedAlerts, err := LoadPerformanceAlertsFromFile(tmpFile)
	require.NoError(t, err)
	assert.Len(t, loadedAlerts, 2)
	assert.Equal(t, "High Critical Alert", loadedAlerts[0].Name)
	assert.Equal(t, "Low Score Alert", loadedAlerts[1].Name)
}

func TestAlertTriggering(t *testing.T) {
	alerts := []PerformanceAlert{
		{
			Name:       "Critical Alert",
			Metric:     "critical",
			Threshold:  1,
			Comparator: ">",
			Active:     true,
		},
	}

	server := NewWebSocketServer(8080, "./testdata", "./testdata", false, false, false, 0, false, alerts, nil, ConnectionQualityConfig{})
	defer server.Stop()

	// Load findings that should trigger the alert
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
		},
		Findings: []model.Finding{
			{Severity: "critical"},
			{Severity: "critical"}, // This should trigger the alert (> 1 critical)
			{Severity: "high"},
		},
	}
	server.UpdateData(findings, nil)

	// Check that the alert was triggered
	time.Sleep(100 * time.Millisecond) // Give time for alert processing
	
	server.alertsMu.Lock()
	defer server.alertsMu.Unlock()
	
	assert.Len(t, server.performanceAlerts, 1)
	assert.NotNil(t, server.performanceAlerts[0].LastTriggered)
}

// Test WebSocket connection quality monitoring
func TestWebSocketConnectionQualityMonitoring(t *testing.T) {
	server := NewWebSocketServer(8080, "./testdata", "./testdata", false, false, false, 0, true, nil, nil, ConnectionQualityConfig{})
	defer server.Stop()

	// Test that connection quality is enabled
	assert.True(t, server.connectionQualityEnabled)
	assert.Equal(t, 10*time.Second, server.pingInterval)

	// Test connection quality info
	info := server.GetConnectionQualityInfo()
	assert.True(t, info["connection_quality_enabled"].(bool))
	assert.Equal(t, int64(10000), info["ping_interval_ms"])
}
