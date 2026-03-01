// Copyright 2026 Mistral AI. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license.

package model

// Remediation represents an automated code fix suggestion for a performance finding
type Remediation struct {
	FindingID       string   `json:"finding_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	CodeChanges     []CodeChange `json:"code_changes"`
	Confidence      float64  `json:"confidence"`
	ImpactEstimate  string   `json:"impact_estimate"`
	Tags            []string  `json:"tags"`
	EvidenceRefs    []string  `json:"evidence_refs"`
}

// CodeChange represents a specific code modification suggestion
type CodeChange struct {
	FilePath    string `json:"file_path"`
	LineNumber  int    `json:"line_number"`
	CurrentCode string `json:"current_code"`
	SuggestedCode string `json:"suggested_code"`
	ChangeType  string `json:"change_type"` // add, modify, delete, refactor
	Explanation  string `json:"explanation"`
}

// RemediationBundle contains all remediation suggestions for a set of findings
type RemediationBundle struct {
	Version         string       `json:"version"`
	GeneratedBy     string       `json:"generated_by"`
	Timestamp       string       `json:"timestamp"`
	Remediations    []Remediation `json:"remediations"`
	Summary         RemediationSummary `json:"summary"`
}

// RemediationSummary provides an overview of remediation suggestions
type RemediationSummary struct {
	TotalRemediations int     `json:"total_remediations"`
	HighImpact        int     `json:"high_impact"`
	MediumImpact      int     `json:"medium_impact"`
	LowImpact         int     `json:"low_impact"`
	EstimatedTotalGain string `json:"estimated_total_gain"`
	ConfidenceScore   float64 `json:"confidence_score"`
}

