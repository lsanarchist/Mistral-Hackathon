package plugin

import (
	"testing"
	"time"
)

func TestPluginPerformance(t *testing.T) {
	// Create a plugin manager
	manager := NewPluginManager("test-plugins")

	// Test initial state
	if len(manager.GetPluginPerformance()) != 0 {
		t.Errorf("Expected empty performance list initially, got %d entries", len(manager.GetPluginPerformance()))
	}

	// Test recording performance
	performance := PluginPerformance{
		PluginName:      "test-plugin",
		ExecutionTime:   100 * time.Millisecond,
		MemoryUsageMB:   50.5,
		CPUUsagePercent: 10.2,
		Timestamp:       time.Now(),
		Success:         true,
	}

	manager.RecordPluginPerformance(performance)

	// Test that performance was recorded
	recordedPerformance := manager.GetPluginPerformance()
	if len(recordedPerformance) != 1 {
		t.Errorf("Expected 1 performance record, got %d", len(recordedPerformance))
	}

	if recordedPerformance[0].PluginName != "test-plugin" {
		t.Errorf("Expected plugin name 'test-plugin', got '%s'", recordedPerformance[0].PluginName)
	}

	if recordedPerformance[0].Success != true {
		t.Errorf("Expected success true, got %v", recordedPerformance[0].Success)
	}

	// Test that performance list is limited to 100 entries
	for i := 0; i < 110; i++ {
		testPerf := PluginPerformance{
			PluginName:      "test-plugin",
			ExecutionTime:   time.Duration(i) * time.Millisecond,
			MemoryUsageMB:   float64(i),
			CPUUsagePercent: float64(i) * 0.1,
			Timestamp:       time.Now(),
			Success:         i%2 == 0,
		}
		manager.RecordPluginPerformance(testPerf)
	}

	finalPerformance := manager.GetPluginPerformance()
	if len(finalPerformance) != 100 {
		t.Errorf("Expected performance list to be limited to 100 entries, got %d", len(finalPerformance))
	}
}

func TestGetMemoryUsageMB(t *testing.T) {
	// This is a basic test to ensure the function runs without error
	memoryUsage := getMemoryUsageMB()
	
	if memoryUsage < 0 {
		t.Errorf("Expected non-negative memory usage, got %f", memoryUsage)
	}
}

func TestGetCPUUsagePercent(t *testing.T) {
	// This is a basic test to ensure the function runs without error
	cpuUsage := getCPUUsagePercent()
	
	if cpuUsage < 0 {
		t.Errorf("Expected non-negative CPU usage, got %f", cpuUsage)
	}
}