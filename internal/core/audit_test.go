package core

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuditLogger(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir := t.TempDir()

	// Test disabled audit logger
	logger := NewAuditLogger(AuditLoggerConfig{
		Enabled: false,
	})
	assert.False(t, logger.IsEnabled())
	logger.LogAction("user1", "test_action", "test_resource", "test_details", "success")
	assert.Equal(t, 0, logger.GetLogCount())

	// Test enabled audit logger
	logger = NewAuditLogger(AuditLoggerConfig{
		Enabled:    true,
		LogDir:     tempDir,
		MaxEntries: 100,
	})
	assert.True(t, logger.IsEnabled())

	// Test logging actions
	logger.LogAction("user1", "test_action1", "test_resource1", "test_details1", "success")
	logger.LogAction("user2", "test_action2", "test_resource2", "test_details2", "failed")

	// Test log count
	assert.Equal(t, 2, logger.GetLogCount())

	// Test getting all logs
	logs := logger.GetLogs()
	assert.Equal(t, 2, len(logs))
	assert.Equal(t, "user1", logs[0].UserID)
	assert.Equal(t, "test_action1", logs[0].Action)
	assert.Equal(t, "test_resource1", logs[0].Resource)
	assert.Equal(t, "success", logs[0].Status)

	// Test getting logs by user
	user1Logs := logger.GetLogsByUser("user1")
	assert.Equal(t, 1, len(user1Logs))
	assert.Equal(t, "user1", user1Logs[0].UserID)

	// Test getting logs by action
	actionLogs := logger.GetLogsByAction("test_action1")
	assert.Equal(t, 1, len(actionLogs))

	// Test getting logs by resource
	resourceLogs := logger.GetLogsByResource("test_resource2")
	assert.Equal(t, 1, len(resourceLogs))

	// Test getting logs by time range
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(1 * time.Hour)
	timeLogs := logger.GetLogsByTimeRange(startTime, endTime)
	assert.Equal(t, 2, len(timeLogs))

	// Test clearing logs
	logger.ClearLogs()
	assert.Equal(t, 0, logger.GetLogCount())

	// Test max entries
	logger.SetMaxEntries(1)
	logger.LogAction("user1", "action1", "resource1", "details1", "success")
	logger.LogAction("user2", "action2", "resource2", "details2", "success")
	assert.Equal(t, 1, logger.GetLogCount()) // Should only keep the last entry

	// Test audit summary
	summary := logger.GetAuditSummary()
	assert.Equal(t, 1, summary["total_entries"])
	assert.True(t, summary["enabled"].(bool))
	assert.Equal(t, 1, summary["max_entries"])
}

func TestAuditLoggerPersistence(t *testing.T) {
	tempDir := t.TempDir()

	// Create first logger and add some entries
	logger1 := NewAuditLogger(AuditLoggerConfig{
		Enabled:    true,
		LogDir:     tempDir,
		MaxEntries: 10,
	})

	logger1.LogAction("user1", "action1", "resource1", "details1", "success")
	logger1.LogAction("user2", "action2", "resource2", "details2", "failed")

	// Create second logger from same directory
	logger2 := NewAuditLogger(AuditLoggerConfig{
		Enabled:    true,
		LogDir:     tempDir,
		MaxEntries: 10,
	})

	// Should load existing logs
	assert.Equal(t, 2, logger2.GetLogCount())
	logs := logger2.GetLogs()
	assert.Equal(t, "user1", logs[0].UserID)
	assert.Equal(t, "action1", logs[0].Action)

	// Add more entries
	logger2.LogAction("user3", "action3", "resource3", "details3", "success")
	assert.Equal(t, 3, logger2.GetLogCount())

	// Verify file was written
	logFile := tempDir + "/audit.log"
	_, err := os.Stat(logFile)
	assert.NoError(t, err)
}

func TestAuditLoggerEdgeCases(t *testing.T) {
	// Test with invalid log directory
	logger := NewAuditLogger(AuditLoggerConfig{
		Enabled:    true,
		LogDir:     "/invalid/directory/path",
		MaxEntries: 10,
	})

	// Should still work, just won't persist to file
	logger.LogAction("user1", "action1", "resource1", "details1", "success")
	assert.Equal(t, 1, logger.GetLogCount())

	// Test with empty user ID
	logger.LogAction("", "action2", "resource2", "details2", "success")
	assert.Equal(t, 2, logger.GetLogCount())

	// Test with very long details
	longDetails := ""
	for i := 0; i < 1000; i++ {
		longDetails += "very long details "
	}
	logger.LogAction("user2", "action3", "resource3", longDetails, "success")
	assert.Equal(t, 3, logger.GetLogCount())

	// Test enable/disable
	logger.Disable()
	logger.LogAction("user3", "action4", "resource4", "details4", "success")
	assert.Equal(t, 3, logger.GetLogCount()) // Should not be logged

	logger.Enable()
	logger.LogAction("user4", "action5", "resource5", "details5", "success")
	assert.Equal(t, 4, logger.GetLogCount()) // Should be logged
}