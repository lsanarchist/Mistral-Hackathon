package model

import "time"

// InsightsBundle contains LLM-generated insights about performance findings
type InsightsBundle struct {
    SchemaVersion    string            `json:"schema_version"`
    GeneratedAt      time.Time         `json:"generated_at"`
    Model            string            `json:"model,omitempty"`
    RequestID        string            `json:"request_id,omitempty"`
    DisabledReason   string            `json:"disabled_reason,omitempty"`
    ExecutiveSummary ExecutiveSummary `json:"executive_summary"`
    TopRisks        []RiskItem        `json:"top_risks,omitempty"`
    TopActions      []ActionItem      `json:"top_actions,omitempty"`
    PerformanceCategories map[string]int `json:"performance_categories,omitempty"`
    PerFinding      []FindingInsight   `json:"per_finding,omitempty"`
}

// ExecutiveSummary provides high-level overview
type ExecutiveSummary struct {
    Overview        string   `json:"overview"`
    OverallSeverity string   `json:"overall_severity"`
    KeyThemes       []string `json:"key_themes,omitempty"`
    Confidence      int      `json:"confidence"` // 0-100
}

// RiskItem represents a significant risk identified by LLM
type RiskItem struct {
    Description string `json:"description"`
    Severity    string `json:"severity"`
    Impact      string `json:"impact"`
    Likelihood  string `json:"likelihood"`
}

// ActionItem represents a recommended action
type ActionItem struct {
    Description    string   `json:"description"`
    Priority       string   `json:"priority"`
    EstimatedEffort string   `json:"estimated_effort"`
    Categories     []string `json:"categories,omitempty"`
}

// FindingInsight provides per-finding analysis
type FindingInsight struct {
    FindingID        string   `json:"finding_id"`
    Narrative        string   `json:"narrative"`
    LikelyRootCauses []string `json:"likely_root_causes,omitempty"`
    Suggestions      []string `json:"suggestions,omitempty"`
    NextMeasurements  []string `json:"next_measurements,omitempty"`
    Caveats          []string `json:"caveats,omitempty"`
    Confidence       int      `json:"confidence"` // 0-100
}