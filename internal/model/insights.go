package model

import "time"

// InsightsBundle contains LLM-generated insights about performance findings
// Schema is versioned for stability and backward compatibility
type InsightsBundle struct {
	SchemaVersion    string           `json:"schema_version"`
	GeneratedAt      time.Time        `json:"generated_at"`
	Model            string           `json:"model,omitempty"`
	RequestID        string           `json:"request_id,omitempty"`
	DisabledReason   string           `json:"disabled_reason,omitempty"`
	ExecutiveSummary ExecutiveSummary `json:"executive_summary"`
	TopRisks         []RiskItem       `json:"top_risks,omitempty"`
	TopActions       []ActionItem     `json:"top_actions,omitempty"`
	PerFinding       []FindingInsight `json:"per_finding,omitempty"`
}

// ExecutiveSummary provides a high-level overview of performance issues
type ExecutiveSummary struct {
	Overview        string   `json:"overview"`
	OverallSeverity string   `json:"overall_severity"`
	KeyThemes       []string `json:"key_themes,omitempty"`
	Confidence      int      `json:"confidence"` // 0-100
}

// RiskItem represents a potential performance risk
type RiskItem struct {
	Description string `json:"description"`
	Severity    string `json:"severity"` // low, medium, high
	Impact      string `json:"impact"`
	Likelihood  string `json:"likelihood"` // low, medium, high
}

// ActionItem represents a recommended action
type ActionItem struct {
	Description     string   `json:"description"`
	Priority        string   `json:"priority"`             // low, medium, high
	EstimatedEffort string   `json:"estimated_effort"`     // low, medium, high
	Categories      []string `json:"categories,omitempty"` // e.g., ["code", "configuration"]
}

// FindingInsight provides LLM-generated analysis for a specific finding
type FindingInsight struct {
	FindingID        string   `json:"finding_id"`
	Narrative        string   `json:"narrative"`
	LikelyRootCauses []string `json:"likely_root_causes,omitempty"`
	Suggestions      []string `json:"suggestions,omitempty"`
	NextMeasurements []string `json:"next_measurements,omitempty"`
	Caveats          []string `json:"caveats,omitempty"`
	Confidence       int      `json:"confidence"` // 0-100
}

// Schema constants
const (
	InsightsSchemaVersion = "1.0"
	SeverityLow           = "low"
	SeverityMedium        = "medium"
	SeverityHigh          = "high"
	PriorityLow           = "low"
	PriorityMedium        = "medium"
	PriorityHigh          = "high"
)
