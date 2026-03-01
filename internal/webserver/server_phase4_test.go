package webserver

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mistral-hackathon/triageprof/internal/model"
)

// TestPhase4Features tests the Phase 4 advanced ML features
func TestPhase4Features(t *testing.T) {
	// Create a test WebSocket server with Phase 4 features enabled
	server := NewWebSocketServer(8082, "./testdata", "./plugins", false, false, false, 0, 
		true, nil, nil, ConnectionQualityConfig{}, true, true, true, false, false)
	
	// Verify Phase 4 features are enabled
	if !server.phase4FeaturesEnabled {
		t.Error("Phase 4 features should be enabled")
	}
	
	if !server.advancedMLEnabled {
		t.Error("Advanced ML should be enabled when Phase 4 features are enabled")
	}
	
	if !server.mlModelEnabled {
		t.Error("ML model should be enabled when Phase 4 features are enabled")
	}
	
	// Test that ML model info is initialized
	if server.mlModelInfo.ModelVersion != "2.0" {
		t.Errorf("Expected model version 2.0, got %s", server.mlModelInfo.ModelVersion)
	}
	
	// Test Phase 4 specific methods exist and don't panic
	_ = server.getConnectionStats()
	
	// Test anomaly detection with Phase 4
	anomalyDetection := server.detectAdvancedConnectionQualityAnomaliesPhase4()
	if anomalyDetection["status"] != "insufficient_data" && anomalyDetection["status"] != "analyzed" {
		t.Errorf("Unexpected anomaly detection status: %s", anomalyDetection["status"])
	}
	
	// Test adaptive learning with Phase 4
	server.performAdaptiveLearningPhase4()
	
	// Verify model statistics are updated
	if server.mlModelInfo.TrainingSamples != 0 {
		t.Errorf("Expected 0 training samples initially, got %d", server.mlModelInfo.TrainingSamples)
	}
	
	// Test that we can get advanced ML connection quality info for Phase 4
	qualityInfo := server.getAdvancedMLConnectionQualityInfoPhase4()
	if qualityInfo["status"] != "analyzed" {
		t.Errorf("Expected analyzed status, got %s", qualityInfo["status"])
	}
	
	// Verify Phase 4 features are included in the response
	if _, ok := qualityInfo["phase_4_features"]; !ok {
		t.Error("Phase 4 features should be included in quality info")
	}
	
	phase4Features := qualityInfo["phase_4_features"].([]string)
	expectedFeatures := []string{
		"Deep learning anomaly detection with enhanced feature extraction",
		"Time series forecasting with historical pattern analysis",
		"Automated root cause analysis with confidence scoring",
		"Anomaly correlation detection for systemic issue identification",
		"Adaptive learning with dynamic learning rate adjustment",
		"Comprehensive ML model management and statistics",
	}
	
	if len(phase4Features) != len(expectedFeatures) {
		t.Errorf("Expected %d Phase 4 features, got %d", len(expectedFeatures), len(phase4Features))
	}
	
	// Test that the server can be stopped without error
	err := server.Stop()
	if err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}
}

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
	
	if anomalyDetection["anomaly_count"] != 1 {
		t.Errorf("Expected 1 anomaly, got %d", int(anomalyDetection["anomaly_count"].(float64)))
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
	server.statsMu.Unlock()
	
	err := server.Stop()
	if err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}
}

// TestPhase4AdaptiveLearning tests adaptive learning with Phase 4 enhancements
func TestPhase4AdaptiveLearning(t *testing.T) {
	server := NewWebSocketServer(8084, "./testdata", "./plugins", false, false, false, 0, 
		false, nil, nil, ConnectionQualityConfig{}, true, true, true, false, false)
	
	// Get initial learning rate
	initialLearningRate := server.mlModelInfo.LearningRate
	initialAccuracy := server.mlModelInfo.AccuracyScore
	
	// Perform adaptive learning multiple times
	for i := 0; i < 5; i++ {
		server.performAdaptiveLearningPhase4()
	}
	
	// Verify learning rate decreased (should be multiplied by 0.95 each time)
	expectedLearningRate := initialLearningRate * 0.95 * 0.95 * 0.95 * 0.95 * 0.95
	if server.mlModelInfo.LearningRate != expectedLearningRate {
		t.Errorf("Expected learning rate %.6f, got %.6f", expectedLearningRate, server.mlModelInfo.LearningRate)
	}
	
	// Verify accuracy improved
	if server.mlModelInfo.AccuracyScore <= initialAccuracy {
		t.Error("Accuracy should improve with adaptive learning")
	}
	
	// Verify model version updates when accuracy is high
	if server.mlModelInfo.AccuracyScore > 0.95 {
		if server.mlModelInfo.ModelVersion != "4.0" {
			t.Errorf("Expected model version 4.0 for high accuracy, got %s", server.mlModelInfo.ModelVersion)
		}
	}
	
	err := server.Stop()
	if err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}
}