package webserver

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestPhase4AnomalyDetection tests Phase 4 anomaly detection with simulated connections
func TestPhase4AnomalyDetection(t *testing.T) {
	// Create a test WebSocket server with Phase 4 features
	server := NewWebSocketServer(8083, "./testdata", "./plugins", false, false, false, 0, 
		true, nil, nil, ConnectionQualityConfig{}, true, true, true, false, false)

	// Create a mock connection with anomalous behavior
	mockConn := &websocket.Conn{}

	// Add a connection with high latency and packet loss to trigger anomaly detection
	stats := &WebSocketConnectionStats{
		ClientID:         "test-client-1",
		ConnectionTime:  time.Now(),
		LastPingTime:    time.Now(),
		LastPongTime:    time.Now(),
		Latency:         1500 * time.Millisecond, // High latency
		PacketLoss:      25.0,                   // High packet loss
		MessagesSent:    100,
		MessagesReceived: 85,                    // Simulate some packet loss
		BytesSent:       10000,
		BytesReceived:   8500,
		ConnectionQuality: "poor",
		Geolocation:     "Test Region",
	}

	// Manually add to connection stats for testing
	server.statsMu.Lock()
	server.connectionStats[mockConn] = stats

	// Add two more connections to meet the minimum requirement of 3
	mockConn2 := &websocket.Conn{}
	stats2 := &WebSocketConnectionStats{
		ClientID:         "test-client-2",
		ConnectionTime:  time.Now(),
		LastPingTime:    time.Now(),
		LastPongTime:    time.Now(),
		Latency:         50 * time.Millisecond, // Normal latency
		PacketLoss:      1.0,                   // Normal packet loss
		MessagesSent:    100,
		MessagesReceived: 99,
		BytesSent:       10000,
		BytesReceived:   9900,
		ConnectionQuality: "excellent",
		Geolocation:     "Test Region",
	}
	server.connectionStats[mockConn2] = stats2

	mockConn3 := &websocket.Conn{}
	stats3 := &WebSocketConnectionStats{
		ClientID:         "test-client-3",
		ConnectionTime:  time.Now(),
		LastPingTime:    time.Now(),
		LastPongTime:    time.Now(),
		Latency:         100 * time.Millisecond, // Slightly elevated latency
		PacketLoss:      2.0,                   // Slightly elevated packet loss
		MessagesSent:    100,
		MessagesReceived: 98,
		BytesSent:       10000,
		BytesReceived:   9800,
		ConnectionQuality: "good",
		Geolocation:     "Test Region",
	}
	server.connectionStats[mockConn3] = stats3
	server.statsMu.Unlock()

	// Update connection stats to trigger anomaly detection
	server.updateConnectionStats(mockConn, 1000, 850)

	// Verify anomaly was detected
	if !stats.IsAnomaly {
		t.Error("Anomaly should have been detected for high latency and packet loss")
	}

	if stats.AnomalyScore <= 0.5 {
		t.Errorf("Expected high anomaly score for poor connection, got %.2f", stats.AnomalyScore)
	}

	if stats.AnomalyType == "" {
		t.Error("Anomaly type should be set")
	}

	// Test anomaly detection with Phase 4 algorithms
	anomalyDetection := server.detectAdvancedConnectionQualityAnomaliesPhase4()

	if anomalyDetection["anomaly_count"] != float64(1) {
		t.Errorf("Expected 1 anomaly, got %v", anomalyDetection["anomaly_count"])
	}

	// Verify anomaly has comprehensive information
	anomalies := anomalyDetection["anomalies"].([]map[string]interface{})
	if len(anomalies) != 1 {
		t.Errorf("Expected 1 anomaly in list, got %d", len(anomalies))
	}

	anomaly := anomalies[0]
	if anomaly["client_id"] != "test-client-1" {
		t.Errorf("Expected client ID test-client-1, got %s", anomaly["client_id"])
	}

	if anomaly["anomaly_score"] == nil {
		t.Error("Anomaly should have a score")
	}

	if anomaly["root_cause"] == nil {
		t.Error("Anomaly should have root cause analysis")
	}

	if anomaly["prediction"] == nil {
		t.Error("Anomaly should have prediction")
	}

	// Clean up
	server.statsMu.Lock()
	delete(server.connectionStats, mockConn)
	delete(server.connectionStats, mockConn2)
	delete(server.connectionStats, mockConn3)
	server.statsMu.Unlock()

	err := server.Stop()
	if err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}
}