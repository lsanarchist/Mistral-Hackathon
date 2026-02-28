package model

import "time"

// JSONReport represents the structured JSON report format
type JSONReport struct {
	SchemaVersion string          `json:"schema_version"`
	GeneratedAt   time.Time       `json:"generated_at"`
	Summary       ReportSummary   `json:"summary"`
	Findings      []ReportFinding `json:"findings"`
	Insights      *InsightsBundle `json:"insights,omitempty"`
}

// ReportSummary provides high-level overview
type ReportSummary struct {
	OverallScore int      `json:"overall_score"`
	TopIssueTags []string `json:"top_issue_tags"`
	KeyThemes    []string `json:"key_themes,omitempty"`
	Notes        []string `json:"notes,omitempty"`
	Severity     string   `json:"severity"`
}

// ReportFinding represents a single finding in JSON format
type ReportFinding struct {
	ID               string                  `json:"id"`
	Category         string                  `json:"category"`
	Title            string                  `json:"title"`
	Severity         string                  `json:"severity"`
	Score            int                     `json:"score"`
	TopHotspots      []StackFrame            `json:"top_hotspots"`
	Callgraph        []CallgraphNode         `json:"callgraph,omitempty"`
	Regression       *RegressionAnalysis    `json:"regression,omitempty"`
	AllocationAnalysis *AllocationAnalysis `json:"allocationAnalysis,omitempty"`
	Evidence         Evidence                `json:"evidence"`
}

// JSONReportOptions configure JSON report generation
type JSONReportOptions struct {
	IncludeInsights bool `json:"include_insights"`
	PrettyPrint     bool `json:"pretty_print"`
}
