// Copyright 2026 Mistral AI. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license.

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// RemediationGenerator generates automated code fix suggestions
type RemediationGenerator struct {
	provider       Provider
	config         model.RemediationConfig
	findingsBundle *model.FindingsBundle
	insightsBundle *model.InsightsBundle
}

func NewRemediationGenerator(provider Provider, config model.RemediationConfig, 
	findingsBundle *model.FindingsBundle, insightsBundle *model.InsightsBundle) *RemediationGenerator {
	return &RemediationGenerator{
		provider:       provider,
		config:         config,
		findingsBundle: findingsBundle,
		insightsBundle: insightsBundle,
	}
}

func (rg *RemediationGenerator) GenerateRemediations(ctx context.Context) (*model.RemediationBundle, error) {
	if !rg.config.Enabled {
		return nil, fmt.Errorf("remediation is disabled in config")
	}

	if rg.findingsBundle == nil || len(rg.findingsBundle.Findings) == 0 {
		return nil, fmt.Errorf("no findings available for remediation")
	}

	// Build prompt for remediation generation
	prompt := rg.buildRemediationPrompt()

	// Generate remediation using LLM
	// Note: We need to adapt the Provider interface to support remediation generation
	// For now, we'll use a simple approach
	response, err := rg.generateRemediationWithProvider(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate remediations: %w", err)
	}

	// Parse and validate remediation response
	bundle, err := rg.parseAndValidateRemediationResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remediation response: %w", err)
	}

	return bundle, nil
}

func (rg *RemediationGenerator) buildRemediationPrompt() string {
	var sb strings.Builder

	sb.WriteString("Generate automated code remediation suggestions for the following performance findings:\n\n")

	// Add findings context
	for _, finding := range rg.findingsBundle.Findings {
		sb.WriteString(fmt.Sprintf("Finding ID: %s\n", finding.ID))
		sb.WriteString(fmt.Sprintf("Title: %s\n", finding.Title))
		sb.WriteString(fmt.Sprintf("Category: %s\n", finding.Category))
		sb.WriteString(fmt.Sprintf("Severity: %s\n", finding.Severity))
		sb.WriteString(fmt.Sprintf("Impact Summary: %s\n", finding.ImpactSummary))
		sb.WriteString("\n")
	}

	// Add insights context if available
	if rg.insightsBundle != nil {
		sb.WriteString("\nAdditional context from LLM insights:\n")
		for _, insight := range rg.insightsBundle.PerFinding {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", insight.FindingID, insight.Narrative))
		}
	}

	sb.WriteString("\nPlease provide specific code changes with the following format:\n")
	sb.WriteString("1. File path and line number\n")
	sb.WriteString("2. Current code (max 200 chars)\n")
	sb.WriteString("3. Suggested code (max 200 chars)\n")
	sb.WriteString("4. Change type (add, modify, delete, refactor)\n")
	sb.WriteString("5. Explanation of why this change improves performance\n")
	sb.WriteString("6. Confidence score (0.0-1.0)\n")
	sb.WriteString("7. Impact estimate (high, medium, low)\n")
	sb.WriteString("8. Evidence references from the findings\n")

	sb.WriteString("\nConstraints:\n")
	sb.WriteString("- Maximum of 3 code changes per finding\n")
	sb.WriteString("- Each code change must be under 200 characters\n")
	sb.WriteString("- Only suggest changes with confidence > 0.7\n")
	sb.WriteString("- Reference specific finding IDs in evidence_refs\n")

	return sb.String()
}

func (rg *RemediationGenerator) generateRemediationWithProvider(ctx context.Context, prompt string) (string, error) {
	// TextGenerator is an optional extended interface that providers may implement
	// to return raw text responses (used by mocks and future providers).
	type textGenerator interface {
		Generate(ctx context.Context, prompt string) (string, error)
	}

	if tg, ok := rg.provider.(textGenerator); ok {
		return tg.Generate(ctx, prompt)
	}

	// Fallback: use GenerateInsights and re-serialise to JSON so the caller
	// can still parse it as a RemediationBundle when the provider returns a
	// pre-built InsightsBundle (unlikely in production, but keeps the path
	// consistent).  For the standard Mistral/OpenAI providers we build a
	// plain prompt that asks for a JSON remediation response.
	bundle, err := rg.provider.GenerateInsights(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("provider call failed: %w", err)
	}
	// If the provider returned a disabled reason we treat it as an empty result.
	if bundle.DisabledReason != "" {
		emptyResp := `{"version":"1.0","generated_by":"triageprof-remediation","timestamp":"` +
			time.Now().UTC().Format(time.RFC3339) +
			`","remediations":[],"summary":{"total_remediations":0,"high_impact":0,"medium_impact":0,"low_impact":0,"estimated_total_gain":"none","confidence_score":0}}`
		return emptyResp, nil
	}
	// In the happy path the model should have returned valid JSON as the
	// content; we just re-serialise whatever the provider gave us.
	data, err := json.Marshal(bundle)
	if err != nil {
		return "", fmt.Errorf("failed to serialise provider response: %w", err)
	}
	return string(data), nil
}

func (rg *RemediationGenerator) parseAndValidateRemediationResponse(response string) (*model.RemediationBundle, error) {
	// Try to parse as JSON first
	var bundle model.RemediationBundle
	if err := json.Unmarshal([]byte(response), &bundle); err == nil {
		// Validate the parsed bundle
		if err := rg.validateRemediationBundle(&bundle); err != nil {
			return nil, fmt.Errorf("remediation validation failed: %w", err)
		}
		return &bundle, nil
	}

	// Fallback to legacy text parsing if JSON fails
	return rg.parseLegacyRemediationResponse(response)
}

func (rg *RemediationGenerator) validateRemediationBundle(bundle *model.RemediationBundle) error {
	// Set default values
	if bundle.Version == "" {
		bundle.Version = "1.0"
	}
	if bundle.GeneratedBy == "" {
		bundle.GeneratedBy = "triageprof-remediation"
	}
	if bundle.Timestamp == "" {
		bundle.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	// Validate each remediation
	for i, remediation := range bundle.Remediations {
		// Check if finding exists
		if !rg.findingExists(remediation.FindingID) {
			return fmt.Errorf("remediation %d references non-existent finding: %s", i, remediation.FindingID)
		}

		// Validate confidence range
		if remediation.Confidence < 0 || remediation.Confidence > 1 {
			return fmt.Errorf("remediation %d has invalid confidence: %f", i, remediation.Confidence)
		}

		// Apply minimum confidence filter
		if remediation.Confidence < rg.config.MinConfidence {
			bundle.Remediations[i].Confidence = rg.config.MinConfidence
		}

		// Validate code changes
		if len(remediation.CodeChanges) > rg.config.MaxCodeChanges {
			bundle.Remediations[i].CodeChanges = remediation.CodeChanges[:rg.config.MaxCodeChanges]
		}

		for j, change := range remediation.CodeChanges {
			// Truncate long code snippets
			if len(change.CurrentCode) > rg.config.CodeChangeLimit {
				bundle.Remediations[i].CodeChanges[j].CurrentCode = change.CurrentCode[:rg.config.CodeChangeLimit]
			}
			if len(change.SuggestedCode) > rg.config.CodeChangeLimit {
				bundle.Remediations[i].CodeChanges[j].SuggestedCode = change.SuggestedCode[:rg.config.CodeChangeLimit]
			}

			// Set default change type if empty
			if change.ChangeType == "" {
				bundle.Remediations[i].CodeChanges[j].ChangeType = "modify"
			}
		}

		// Validate evidence references
		for _, evidenceRef := range remediation.EvidenceRefs {
			if !rg.findingExists(evidenceRef) {
				return fmt.Errorf("remediation %d references invalid evidence: %s", i, evidenceRef)
			}
		}
	}

	// Update summary
	rg.updateRemediationSummary(bundle)

	return nil
}

func (rg *RemediationGenerator) findingExists(findingID string) bool {
	for _, finding := range rg.findingsBundle.Findings {
		if finding.ID == findingID {
			return true
		}
	}
	return false
}

func (rg *RemediationGenerator) parseLegacyRemediationResponse(response string) (*model.RemediationBundle, error) {
	// Fallback parsing for non-JSON responses
	log.Printf("Warning: Remediation response is not valid JSON, attempting legacy parsing")

	bundle := &model.RemediationBundle{
		Version:     "1.0",
		GeneratedBy: "triageprof-remediation",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	// Simple parsing - in production this would be more sophisticated
	lines := strings.Split(response, "\n")
	var currentRemediation *model.Remediation

	for _, line := range lines {
		if strings.HasPrefix(line, "Finding:") {
			findingID := strings.TrimPrefix(line, "Finding:")
			findingID = strings.TrimSpace(findingID)
			currentRemediation = &model.Remediation{
				FindingID:      findingID,
				Confidence:     0.7,
				ImpactEstimate: "medium",
			}
			bundle.Remediations = append(bundle.Remediations, *currentRemediation)
		} else if currentRemediation != nil && strings.HasPrefix(line, "Suggestion:") {
			suggestion := strings.TrimPrefix(line, "Suggestion:")
			suggestion = strings.TrimSpace(suggestion)
			currentRemediation.Description = suggestion
		}
	}

	// Validate the parsed bundle
	if err := rg.validateRemediationBundle(bundle); err != nil {
		return nil, fmt.Errorf("legacy remediation validation failed: %w", err)
	}

	return bundle, nil
}

func (rg *RemediationGenerator) updateRemediationSummary(bundle *model.RemediationBundle) {
	summary := model.RemediationSummary{
		TotalRemediations: len(bundle.Remediations),
		ConfidenceScore:   0.0,
		EstimatedTotalGain: "unknown",
	}

	for _, remediation := range bundle.Remediations {
		summary.ConfidenceScore += remediation.Confidence

		switch remediation.ImpactEstimate {
		case "high":
			summary.HighImpact++
		case "medium":
			summary.MediumImpact++
		case "low":
			summary.LowImpact++
		}
	}

	if summary.TotalRemediations > 0 {
		summary.ConfidenceScore /= float64(summary.TotalRemediations)
	}

	// Simple impact estimation
	switch {
	case summary.HighImpact >= 1:
		summary.EstimatedTotalGain = "high (>30% improvement)"
	case summary.MediumImpact >= 3:
		summary.EstimatedTotalGain = "medium (15-30% improvement)"
	case summary.MediumImpact >= 1 || summary.LowImpact >= 3:
		summary.EstimatedTotalGain = "low (5-15% improvement)"
	default:
		summary.EstimatedTotalGain = "minimal (<5% improvement)"
	}

	bundle.Summary = summary
}