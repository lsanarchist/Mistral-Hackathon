package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// AuditLogger handles audit logging for enterprise features
type AuditLogger struct {
	mu               sync.Mutex
	logEntries       []model.AuditLogEntry
	logFilePath      string
	enabled          bool
	maxEntries       int
}

// NewAuditLogger creates a new audit logger
type AuditLoggerConfig struct {
	Enabled      bool
	LogDir       string
	MaxEntries   int
}

func NewAuditLogger(config AuditLoggerConfig) *AuditLogger {
	logger := &AuditLogger{
		logEntries: make([]model.AuditLogEntry, 0),
		enabled:    config.Enabled,
		maxEntries: config.MaxEntries,
	}

	if config.Enabled && config.LogDir != "" {
		logger.logFilePath = filepath.Join(config.LogDir, "audit.log")
		logger.loadExistingLogs()
	}

	return logger
}

// loadExistingLogs loads existing audit logs from file if it exists
func (a *AuditLogger) loadExistingLogs() {
	if a.logFilePath == "" {
		return
	}

	data, err := os.ReadFile(a.logFilePath)
	if err != nil {
		// File doesn't exist or can't be read, start fresh
		return
	}

	var entries []model.AuditLogEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		// Invalid JSON, start fresh
		return
	}

	a.logEntries = entries
	// Trim to max entries if needed
	if len(a.logEntries) > a.maxEntries && a.maxEntries > 0 {
		a.logEntries = a.logEntries[len(a.logEntries)-a.maxEntries:]
	}
}

// LogAction logs an audit action
func (a *AuditLogger) LogAction(userID, action, resource, details, status string) {
	if !a.enabled {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	entry := model.AuditLogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Details:   details,
		Status:    status,
	}

	a.logEntries = append(a.logEntries, entry)
	// Trim to max entries if needed
	if len(a.logEntries) > a.maxEntries && a.maxEntries > 0 {
		a.logEntries = a.logEntries[len(a.logEntries)-a.maxEntries:]
	}

	// Write to file if path is configured
	if a.logFilePath != "" {
		a.writeToFile()
	}
}

// writeToFile writes the current log entries to file
func (a *AuditLogger) writeToFile() {
	data, err := json.MarshalIndent(a.logEntries, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal audit logs: %v\n", err)
		return
	}

	err = os.WriteFile(a.logFilePath, data, 0644)
	if err != nil {
		fmt.Printf("Failed to write audit logs to file: %v\n", err)
	}
}

// GetLogs returns all audit log entries
func (a *AuditLogger) GetLogs() []model.AuditLogEntry {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Return a copy to prevent external modification
	entries := make([]model.AuditLogEntry, len(a.logEntries))
	copy(entries, a.logEntries)
	return entries
}

// GetLogsByUser returns audit log entries for a specific user
func (a *AuditLogger) GetLogsByUser(userID string) []model.AuditLogEntry {
	a.mu.Lock()
	defer a.mu.Unlock()

	var result []model.AuditLogEntry
	for _, entry := range a.logEntries {
		if entry.UserID == userID {
			result = append(result, entry)
		}
	}
	return result
}

// GetLogsByAction returns audit log entries for a specific action
func (a *AuditLogger) GetLogsByAction(action string) []model.AuditLogEntry {
	a.mu.Lock()
	defer a.mu.Unlock()

	var result []model.AuditLogEntry
	for _, entry := range a.logEntries {
		if entry.Action == action {
			result = append(result, entry)
		}
	}
	return result
}

// GetLogsByResource returns audit log entries for a specific resource
func (a *AuditLogger) GetLogsByResource(resource string) []model.AuditLogEntry {
	a.mu.Lock()
	defer a.mu.Unlock()

	var result []model.AuditLogEntry
	for _, entry := range a.logEntries {
		if entry.Resource == resource {
			result = append(result, entry)
		}
	}
	return result
}

// GetLogsByTimeRange returns audit log entries within a time range
func (a *AuditLogger) GetLogsByTimeRange(startTime, endTime time.Time) []model.AuditLogEntry {
	a.mu.Lock()
	defer a.mu.Unlock()

	var result []model.AuditLogEntry
	for _, entry := range a.logEntries {
		entryTime, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err != nil {
			continue
		}

		if (startTime.IsZero() || !entryTime.Before(startTime)) && 
		   (endTime.IsZero() || !entryTime.After(endTime)) {
			result = append(result, entry)
		}
	}
	return result
}

// ClearLogs clears all audit log entries
func (a *AuditLogger) ClearLogs() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.logEntries = []model.AuditLogEntry{}
	if a.logFilePath != "" {
		a.writeToFile()
	}
}

// GetLogCount returns the number of log entries
func (a *AuditLogger) GetLogCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	return len(a.logEntries)
}

// IsEnabled returns whether audit logging is enabled
func (a *AuditLogger) IsEnabled() bool {
	return a.enabled
}

// Enable enables audit logging
func (a *AuditLogger) Enable() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.enabled = true
}

// Disable disables audit logging
func (a *AuditLogger) Disable() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.enabled = false
}

// SetMaxEntries sets the maximum number of log entries to keep
func (a *AuditLogger) SetMaxEntries(max int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.maxEntries = max
	if max > 0 && len(a.logEntries) > max {
		a.logEntries = a.logEntries[len(a.logEntries)-max:]
	}
	if a.logFilePath != "" {
		a.writeToFile()
	}
}

// GetAuditSummary returns a summary of audit log activity
func (a *AuditLogger) GetAuditSummary() map[string]interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	summary := map[string]interface{}{
		"total_entries":    len(a.logEntries),
		"enabled":         a.enabled,
		"max_entries":     a.maxEntries,
		"actions_by_type": make(map[string]int),
		"users":           make(map[string]int),
		"resources":       make(map[string]int),
		"status_counts":   make(map[string]int),
	}

	if len(a.logEntries) == 0 {
		return summary
	}

	for _, entry := range a.logEntries {
		summary["actions_by_type"].(map[string]int)[entry.Action]++
		summary["users"].(map[string]int)[entry.UserID]++
		summary["resources"].(map[string]int)[entry.Resource]++
		summary["status_counts"].(map[string]int)[entry.Status]++
	}

	return summary
}

// String representation of audit logger status
func (a *AuditLogger) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.enabled {
		return "Audit Logger: Disabled"
	}

	return fmt.Sprintf("Audit Logger: Enabled (Entries: %d, Max: %d, File: %s)", 
		len(a.logEntries), a.maxEntries, a.logFilePath)
}