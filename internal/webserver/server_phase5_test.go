package webserver

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mistral-hackathon/triageprof/internal/model"
)

// TestPhase5Features tests the Phase 5 advanced ML features
func TestPhase5Features(t *testing.T) {
	// Create a test WebSocket server with Phase 5 features enabled
	server := NewWebSocketServer(8083, "./testdata", "./plugins", false, false, false, 0, 
		true, nil, nil, ConnectionQualityConfig{}, true, true, true, true, false)

	// Verify Phase 5 features are enabled
	if !server.phase5FeaturesEnabled {
		t.Error("Phase 5 features should be enabled")
	}

	if !server.advancedMLEnabled {
		t.Error("Advanced ML should be enabled when Phase 5 features are enabled")
	}

	if !server.mlModelEnabled {
		t.Error("ML model should be enabled when Phase 5 features are enabled")
	}

	if !server.phase4FeaturesEnabled {
		t.Error("Phase 4 features should be enabled when Phase 5 features are enabled")
	}

	// Test that advanced ML model info is initialized
	if server.advancedMLModelInfo.ModelVersion != "3.0" {
		t.Errorf("Expected model version 3.0, got %s", server.advancedMLModelInfo.ModelVersion)
	}

	if server.advancedMLModelInfo.ModelType != "deep_learning" {
		t.Errorf("Expected model type deep_learning, got %s", server.advancedMLModelInfo.ModelType)
	}

	if server.advancedMLModelInfo.TrainingStatus != "advanced" {
		t.Errorf("Expected training status advanced, got %s", server.advancedMLModelInfo.TrainingStatus)
	}

	// Test Phase 5 specific methods exist and don't panic
	_ = server.getConnectionStats()

	// Test anomaly detection with Phase 5
	anomalyDetection := server.detectAdvancedConnectionQualityAnomaliesPhase5()
	if anomalyDetection["status"] != "insufficient_data" && anomalyDetection["status"] != "analyzed" {
		t.Errorf("Unexpected anomaly detection status: %s", anomalyDetection["status"])
	}

	// Test adaptive learning with Phase 5
	server.performAdaptiveLearningPhase5()

	// Verify model statistics are updated
	if server.advancedMLModelInfo.TrainingSamples != 1 {
		t.Errorf("Expected 1 training sample after adaptive learning, got %d", server.advancedMLModelInfo.TrainingSamples)
	}

	// Test correlation detection with Phase 5
	correlationAnalysis := server.performAnomalyCorrelationDetectionPhase5()
	if correlationAnalysis["status"] != "analyzed" {
		t.Errorf("Unexpected correlation analysis status: %s", correlationAnalysis["status"])
	}

	// Test predictive maintenance with Phase 5
	predictiveMaintenance := server.performPredictiveMaintenancePhase5()
	if predictiveMaintenance["status"] != "analyzed" {
		t.Errorf("Unexpected predictive maintenance status: %s", predictiveMaintenance["status"])
	}

	// Test time series data point addition
	stats := server.getConnectionStats()
	if len(stats) > 0 {
		stat := stats[0]
		server.addTimeSeriesDataPoint(stat, 0.5, "test_anomaly")
	}

	// Verify time series data was added
	server.timeSeriesMu.Lock()
	timeSeriesCount := len(server.anomalyTimeSeries)
	server.timeSeriesMu.Unlock()

	if timeSeriesCount != 1 {
		t.Errorf("Expected 1 time series data point, got %d", timeSeriesCount)
	}

	// Test advanced ML connection quality info with Phase 5
	qualityInfo := server.getAdvancedMLConnectionQualityInfoPhase5()
	if qualityInfo["phase_5_features"] == nil {
		t.Error("Phase 5 features should be present in quality info")
	}

	phase5Features := qualityInfo["phase_5_features"].(map[string]interface{})
	if !phase5Features["enabled"].(bool) {
		t.Error("Phase 5 features should be enabled in quality info")
	}
}

// TestPhase5AnomalyDetection tests Phase 5 anomaly detection with simulated connections
func TestPhase5AnomalyDetection(t *testing.T) {
	server := NewWebSocketServer(8084, "./testdata", "./plugins", false, false, false, 0, 
		true, nil, nil, ConnectionQualityConfig{}, true, true, true, true, false)

	// Create some test connection stats with anomalies
	server.statsMu.Lock()
	server.connectionStats = make(map[*websocket.Conn]*WebSocketConnectionStats)

	// Create mock connections
	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}
	conn3 := &websocket.Conn{}

	// Add connection with high latency and packet loss
	server.connectionStats[conn1] = &WebSocketConnectionStats{
		ClientID:         "client1",
		ConnectionTime:   time.Now(),
		Latency:          800 * time.Millisecond,
		PacketLoss:       15.0,
		ConnectionScore:  45.0,
		ConnectionQuality: "fair",
		AnomalyHistory: []AnomalyEvent{
			{
				Timestamp:    time.Now().Add(-10 * time.Minute),
				AnomalyType:  "latency_packet_loss",
				AnomalyScore: 0.85,
			},
		},
		QualityTrend: "degrading",
	}

	// Add connection with persistent high latency
	server.connectionStats[conn2] = &WebSocketConnectionStats{
		ClientID:         "client2",
		ConnectionTime:   time.Now(),
		Latency:          1200 * time.Millisecond,
		PacketLoss:       5.0,
		ConnectionScore:  35.0,
		ConnectionQuality: "poor",
		AnomalyHistory: []AnomalyEvent{
			{
				Timestamp:    time.Now().Add(-5 * time.Minute),
				AnomalyType:  "high_latency",
				AnomalyScore: 0.90,
			},
		},
		QualityTrend: "degrading",
	}

	// Add normal connection
	server.connectionStats[conn3] = &WebSocketConnectionStats{
		ClientID:         "client3",
		ConnectionTime:   time.Now(),
		Latency:          100 * time.Millisecond,
		PacketLoss:       2.0,
		ConnectionScore:  85.0,
		ConnectionQuality: "excellent",
		AnomalyHistory:   nil,
		QualityTrend:     "stable",
	}

	server.statsMu.Unlock()

	// Test Phase 5 anomaly detection
	anomalyDetection := server.detectAdvancedConnectionQualityAnomaliesPhase5()

	if anomalyDetection["status"] != "analyzed" {
		t.Errorf("Expected analyzed status, got %s", anomalyDetection["status"])
	}

	if anomalyDetection["anomaly_count"] != 2 {
		t.Errorf("Expected 2 anomalies, got %d", int(anomalyDetection["anomaly_count"].(float64)))
	}

	if anomalyDetection["severity"] != "high" {
		t.Errorf("Expected high severity, got %s", anomalyDetection["severity"])
	}

	// Verify anomalies contain expected data
	anomalies := anomalyDetection["anomalies"].([]map[string]interface{})
	if len(anomalies) != 2 {
		t.Errorf("Expected 2 anomaly entries, got %d", len(anomalies))
	}

	// Check that anomalies have Phase 5 features
	for _, anomaly := range anomalies {
		if anomaly["automated_analysis"] == nil {
			t.Error("Anomaly should have automated_analysis")
		}

		if anomaly["ml_insights"] == nil {
			t.Error("Anomaly should have ml_insights")
		}

		if anomaly["prediction"] == nil {
			t.Error("Anomaly should have prediction")
		}

		prediction := anomaly["prediction"].(map[string]interface{})
		if prediction["phase_5_features"] == nil {
			t.Error("Prediction should have phase_5_features")
		}
	}

	// Verify correlation analysis
	correlationAnalysis := anomalyDetection["correlation_analysis"]
	if correlationAnalysis == nil {
		t.Error("Should have correlation_analysis")
	}

	// Verify predictive maintenance
	predictiveMaintenance := anomalyDetection["predictive_maintenance"]
	if predictiveMaintenance == nil {
		t.Error("Should have predictive_maintenance")
	}
}

// TestPhase5RootCauseAnalysis tests Phase 5 root cause analysis
func TestPhase5RootCauseAnalysis(t *testing.T) {
	server := NewWebSocketServer(8085, "./testdata", "./plugins", false, false, false, 0, 
		true, nil, nil, ConnectionQualityConfig{}, true, true, true, true, false)

	// Test root cause analysis for different anomaly types
	testCases := []struct {
		anomalyType string
		expectedRootCause string
	}{
		{
			anomalyType: "latency_packet_loss_repeating_degrading",
			expectedRootCause: "Network congestion with packet loss and degrading quality - likely infrastructure issue",
		},
		{
			anomalyType: "persistent_high_latency_with_history",
			expectedRootCause: "Persistent high latency with historical pattern - likely geographical or routing issue",
		},
		{
			anomalyType: "degrading_packet_loss_low_score",
			expectedRootCause: "Degrading packet loss with low connection score - likely network instability or interference",
		},
		{
			anomalyType: "unknown_anomaly",
			expectedRootCause: "Unknown root cause - requires further investigation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.anomalyType, func(t *testing.T) {
			// Create a test stat
			stat := &WebSocketConnectionStats{
				ClientID:         "test_client",
				ConnectionTime:   time.Now(),
				Latency:          500 * time.Millisecond,
				PacketLoss:       10.0,
				ConnectionScore:  40.0,
				ConnectionQuality: "fair",
				QualityTrend:     "degrading",
			}

			rootCause := server.analyzeAnomalyRootCausePhase5(stat, tc.anomalyType)
			if rootCause != tc.expectedRootCause {
				t.Errorf("For anomaly type %s, expected root cause %s, got %s", tc.anomalyType, tc.expectedRootCause, rootCause)
			}
		})
	}
}

// TestPhase5PatternMatching tests Phase 5 pattern matching functionality
func TestPhase5PatternMatching(t *testing.T) {
	server := NewWebSocketServer(8086, "./testdata", "./plugins", false, false, false, 0, 
		true, nil, nil, ConnectionQualityConfig{}, true, true, true, true, false)

	// Add a test pattern to the server
	pattern := RootCausePattern{
		PatternID:   "test_pattern_1",
		PatternType: "latency",
		Conditions: []Condition{
			{
				Field:    "latency",
				Operator: ">",
				Value:    "500",
				Weight:   0.8,
			},
			{
				Field:    "packet_loss",
				Operator: ">",
				Value:    "10",
				Weight:   0.7,
			},
		},
		RootCause:       "High latency with packet loss detected - network congestion likely",
		Confidence:      0.9,
		OccurrenceCount: 5,
		FirstSeen:       time.Now(),
		LastSeen:        time.Now(),
	}

	server.rootCausePatternsMu.Lock()
	server.anomalyRootCausePatterns["test_pattern_1"] = pattern
	server.rootCausePatternsMu.Unlock()

	// Create a test stat that should match the pattern
	stat := &WebSocketConnectionStats{
		ClientID:         "test_client",
		ConnectionTime:   time.Now(),
		Latency:          600 * time.Millisecond, // > 500
		PacketLoss:       12.0,                  // > 10
		ConnectionScore:  50.0,
		ConnectionQuality: "good",
		QualityTrend:     "stable",
	}

	// Test pattern condition checking
	testConditions := []Condition{
		{
			Field:    "latency",
			Operator: ">",
			Value:    "500",
			Weight:   0.8,
		},
		{
			Field:    "packet_loss",
			Operator: ">",
			Value:    "10",
			Weight:   0.7,
		},
		{
			Field:    "connection_score",
			Operator: ">=",
			Value:    "50",
			Weight:   0.6,
		},
	}

	for _, condition := range testConditions {
		if !server.checkPatternCondition(stat, condition) {
			t.Errorf("Pattern condition should match for %s %s %s", condition.Field, condition.Operator, condition.Value)
		}
	}

	// Test root cause analysis with pattern matching
	rootCause := server.analyzeAnomalyRootCausePhase5(stat, "latency")
	if rootCause != pattern.RootCause {
		t.Errorf("Expected pattern-based root cause %s, got %s", pattern.RootCause, rootCause)
	}

	// Test numeric condition checking
	numericTests := []struct {
		name     string
		value    float64
		operator string
		target   string
		expected bool
	}{
		{"greater than true", 600, ">", "500", true},
		{"greater than false", 400, ">", "500", false},
		{"less than true", 400, "<", "500", true},
		{"less than false", 600, "<", "500", false},
		{"equal true", 500, "==", "500", true},
		{"equal false", 600, "==", "500", false},
		{"greater equal true", 500, ">=", "500", true},
		{"greater equal false", 400, ">=", "500", false},
		{"less equal true", 500, "<=", "500", true},
		{"less equal false", 600, "<=", "500", false},
	}

	for _, tt := range numericTests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkNumericCondition(tt.value, tt.operator, tt.target)
			if result != tt.expected {
				t.Errorf("checkNumericCondition(%f, %s, %s) = %v, want %v", tt.value, tt.operator, tt.target, result, tt.expected)
			}
		})
	}
}