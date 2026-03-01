package llm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// InsightsGenerator orchestrates LLM insight generation
type InsightsGenerator struct {
	Provider       Provider
	Cache          *InsightsCache
	MaxPromptChars int
}

// NewInsightsGenerator creates a new insights generator with default Mistral provider
func NewInsightsGenerator(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool) (*InsightsGenerator, error) {
	config := ProviderConfig{
		ProviderName: "mistral",
		APIKey:       apiKey,
		Model:        model,
		Timeout:      time.Duration(timeout) * time.Second,
		MaxResponse:  maxResponse,
		DryRun:       dryRun,
	}
	
	provider, err := NewProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	return &InsightsGenerator{
		Provider:       provider,
		Cache:          NewInsightsCache(CacheConfig{Enabled: true}), // Enabled by default
		MaxPromptChars: maxPromptChars,
	}, nil
}

// NewInsightsGeneratorWithProvider creates a new insights generator with a specific provider
func NewInsightsGeneratorWithProvider(config ProviderConfig) (*InsightsGenerator, error) {
	provider, err := NewProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	return &InsightsGenerator{
		Provider:       provider,
		Cache:          NewInsightsCache(CacheConfig{Enabled: true}), // Enabled by default
		MaxPromptChars: config.MaxPrompt,
	}, nil
}

// GenerateInsights creates LLM insights from bundle and findings with guardrails
func (g *InsightsGenerator) GenerateInsights(ctx context.Context, 
	bundle *model.ProfileBundle, findings *model.FindingsBundle) (*model.InsightsBundle, error) {

	// Validate inputs
	if bundle == nil || findings == nil {
		return &model.InsightsBundle{
			DisabledReason: "bundle or findings is nil",
		}, nil
	}

	// Check cache first
	if cachedInsights, found := g.Cache.GetCachedInsights(ctx, bundle, findings); found {
		return cachedInsights, nil
	}

	// Create prompt builder
	builder := NewPromptBuilder(bundle, findings, g.MaxPromptChars)

	// Build secure prompt
	prompt, err := builder.Build()
	if err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("failed to build prompt: %v", err),
		}, nil
	}

	// Generate insights using the provider
	insights, err := g.Provider.GenerateInsights(ctx, prompt)
	if err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("LLM generation failed: %v", err),
		}, nil
	}

	// Validate insights with strict guardrails
	if err := validateInsights(insights, findings); err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("insights validation failed: %v", err),
		}, nil
	}

	// Apply post-processing guardrails
	applyInsightsGuardrails(insights)

	// Cache the validated insights
	if err := g.Cache.CacheInsights(ctx, bundle, findings, insights); err != nil {
		log.Printf("Warning: failed to cache insights: %v", err)
	}

	return insights, nil
}

// validateInsights performs strict validation of LLM-generated insights
func validateInsights(insights *model.InsightsBundle, findings *model.FindingsBundle) error {
	if insights == nil {
		return fmt.Errorf("insights bundle is nil")
	}

	// Check for disabled insights
	if insights.DisabledReason != "" {
		return nil // Already disabled, no need to validate
	}

	// Validate schema version
	if insights.SchemaVersion == "" {
		insights.SchemaVersion = "2.0"
	}

	// Validate executive summary
	if insights.ExecutiveSummary.Overview == "" {
		return fmt.Errorf("executive summary overview is empty")
	}

	if insights.ExecutiveSummary.OverallSeverity == "" {
		return fmt.Errorf("executive summary severity is empty")
	}

	// Validate confidence range
	if insights.ExecutiveSummary.Confidence < 0 || insights.ExecutiveSummary.Confidence > 100 {
		return fmt.Errorf("confidence must be between 0-100, got %d", insights.ExecutiveSummary.Confidence)
	}

	// Validate per-finding insights
	for i, findingInsight := range insights.PerFinding {
		if findingInsight.FindingID == "" {
			return fmt.Errorf("finding insight %d has empty finding_id", i)
		}

		if findingInsight.Narrative == "" {
			return fmt.Errorf("finding insight %d has empty narrative", i)
		}

		// Validate confidence range
		if findingInsight.Confidence < 0 || findingInsight.Confidence > 100 {
			return fmt.Errorf("finding insight %d confidence must be between 0-100, got %d", i, findingInsight.Confidence)
		}

		// Check that finding ID exists in actual findings
		if !findingExists(findings, findingInsight.FindingID) {
			return fmt.Errorf("finding insight %d references non-existent finding_id: %s", i, findingInsight.FindingID)
		}

		// Validate evidence citations - narrative should reference the finding
		if !containsFindingReference(findingInsight.Narrative, findingInsight.FindingID) {
			return fmt.Errorf("finding insight %d narrative must contain evidence reference to finding_id: %s", i, findingInsight.FindingID)
		}

		// Validate code example length limits
		for j, codeExample := range findingInsight.CodeExamples {
			if len(codeExample) > 200 {
				return fmt.Errorf("finding insight %d code example %d exceeds 200 character limit: %d chars", i, j, len(codeExample))
			}
		}
	}

	// Validate top risks
	for i, risk := range insights.TopRisks {
		if risk.Description == "" {
			return fmt.Errorf("top risk %d has empty description", i)
		}

		if risk.Severity == "" {
			return fmt.Errorf("top risk %d has empty severity", i)
		}

		// Validate severity values
		validSeverities := map[string]bool{"low": true, "medium": true, "high": true, "critical": true}
		if !validSeverities[strings.ToLower(risk.Severity)] {
			return fmt.Errorf("top risk %d has invalid severity: %s", i, risk.Severity)
		}
	}

	// Validate top actions
	for i, action := range insights.TopActions {
		if action.Description == "" {
			return fmt.Errorf("top action %d has empty description", i)
		}

		if action.Priority == "" {
			return fmt.Errorf("top action %d has empty priority", i)
		}

		// Validate priority values
		validPriorities := map[string]bool{"low": true, "medium": true, "high": true, "critical": true}
		if !validPriorities[strings.ToLower(action.Priority)] {
			return fmt.Errorf("top action %d has invalid priority: %s", i, action.Priority)
		}
	}

	return nil
}

// findingExists checks if a finding ID exists in the findings bundle
func findingExists(findings *model.FindingsBundle, findingID string) bool {
	for _, finding := range findings.Findings {
		if finding.ID == findingID {
			return true
		}
	}
	return false
}

// containsFindingReference checks if text contains a reference to a finding ID
func containsFindingReference(text, findingID string) bool {
	// Check for direct reference
	if strings.Contains(text, findingID) {
		return true
	}
	
	// Check for common reference patterns
	patterns := []string{
		"Finding " + findingID,
		"finding " + findingID,
		"Issue " + findingID,
		"issue " + findingID,
		"ID " + findingID,
		"id " + findingID,
	}
	
	for _, pattern := range patterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	
	return false
}

// applyInsightsGuardrails applies post-processing guardrails to insights
func applyInsightsGuardrails(insights *model.InsightsBundle) {
	// Set generated timestamp if not set
	if insights.GeneratedAt.IsZero() {
		insights.GeneratedAt = time.Now()
	}

	// Set schema version if not set
	if insights.SchemaVersion == "" {
		insights.SchemaVersion = "2.0"
	}

	// Ensure model is set
	if insights.Model == "" {
		insights.Model = "unknown"
	}

	// Truncate long text fields to prevent excessive output
	truncateField := func(s *string, maxLength int) {
		if len(*s) > maxLength {
			*s = (*s)[:maxLength] + "..."
		}
	}

	// Truncate executive summary
	truncateField(&insights.ExecutiveSummary.Overview, 1000)

	// Truncate per-finding insights
	for i := range insights.PerFinding {
		truncateField(&insights.PerFinding[i].Narrative, 2000)

		// Truncate code examples to 200 chars as per requirements
		for j := range insights.PerFinding[i].CodeExamples {
			truncateField(&insights.PerFinding[i].CodeExamples[j], 200)
		}

		// Truncate suggestions and root causes
		for j := range insights.PerFinding[i].Suggestions {
			truncateField(&insights.PerFinding[i].Suggestions[j], 500)
		}

		for j := range insights.PerFinding[i].LikelyRootCauses {
			truncateField(&insights.PerFinding[i].LikelyRootCauses[j], 500)
		}
	}

	// Truncate top risks and actions
	for i := range insights.TopRisks {
		truncateField(&insights.TopRisks[i].Description, 500)
		truncateField(&insights.TopRisks[i].Impact, 500)
		truncateField(&insights.TopRisks[i].PotentialImpact, 200)
	}

	for i := range insights.TopActions {
		truncateField(&insights.TopActions[i].Description, 500)
		truncateField(&insights.TopActions[i].ExpectedImpact, 200)

		// Truncate code examples in actions
		for j := range insights.TopActions[i].CodeExamples {
			truncateField(&insights.TopActions[i].CodeExamples[j], 200)
		}
	}

	// Apply redaction to sensitive information
	applyRedactionToInsights(insights)
}

// applyRedactionToInsights applies redaction to sensitive information in insights
func applyRedactionToInsights(insights *model.InsightsBundle) {
	replacePatterns := []struct {
		pattern   *regexp.Regexp
		replacement string
	}{
		{regexp.MustCompile(`(?i)(token|secret|key|password)=[^\s]+`), "$1=[REDACTED]"},
		{regexp.MustCompile(`(?i)(localhost|127\.0\.0\.1|\d+\.\d+\.\d+\.\d+)`), "[REDACTED_HOSTNAME]"},
		{regexp.MustCompile(`[A-Za-z0-9]{32,}`), "[REDACTED_TOKEN]"},
		{regexp.MustCompile(`(?i)(http|https)://[^\s]+`), "[REDACTED_URL]"},
	}

	redactText := func(text string) string {
		result := text
		for _, p := range replacePatterns {
			result = p.pattern.ReplaceAllString(result, p.replacement)
		}
		return result
	}

	// Redact executive summary
	insights.ExecutiveSummary.Overview = redactText(insights.ExecutiveSummary.Overview)
	for i := range insights.ExecutiveSummary.KeyThemes {
		insights.ExecutiveSummary.KeyThemes[i] = redactText(insights.ExecutiveSummary.KeyThemes[i])
	}

	// Redact per-finding insights
	for i := range insights.PerFinding {
		insights.PerFinding[i].Narrative = redactText(insights.PerFinding[i].Narrative)

		for j := range insights.PerFinding[i].LikelyRootCauses {
			insights.PerFinding[i].LikelyRootCauses[j] = redactText(insights.PerFinding[i].LikelyRootCauses[j])
		}

		for j := range insights.PerFinding[i].Suggestions {
			insights.PerFinding[i].Suggestions[j] = redactText(insights.PerFinding[i].Suggestions[j])
		}

		for j := range insights.PerFinding[i].CodeExamples {
			insights.PerFinding[i].CodeExamples[j] = redactText(insights.PerFinding[i].CodeExamples[j])
		}
	}

	// Redact top risks and actions
	for i := range insights.TopRisks {
		insights.TopRisks[i].Description = redactText(insights.TopRisks[i].Description)
		insights.TopRisks[i].Impact = redactText(insights.TopRisks[i].Impact)
		insights.TopRisks[i].PotentialImpact = redactText(insights.TopRisks[i].PotentialImpact)
	}

	for i := range insights.TopActions {
		insights.TopActions[i].Description = redactText(insights.TopActions[i].Description)
		insights.TopActions[i].ExpectedImpact = redactText(insights.TopActions[i].ExpectedImpact)

		for j := range insights.TopActions[i].CodeExamples {
			insights.TopActions[i].CodeExamples[j] = redactText(insights.TopActions[i].CodeExamples[j])
		}
	}
}

// GenerateRemediations creates automated code fix suggestions from findings and insights
	func (g *InsightsGenerator) GenerateRemediations(ctx context.Context, findings *model.FindingsBundle, insights *model.InsightsBundle, config model.RemediationConfig) (*model.RemediationBundle, error) {

		// Validate inputs
		if findings == nil || len(findings.Findings) == 0 {
			return nil, fmt.Errorf("no findings available for remediation")
		}

		// Create remediation generator
		remediationGen := NewRemediationGenerator(g.Provider, config, findings, insights)

		// Generate remediations
		remediations, err := remediationGen.GenerateRemediations(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate remediations: %w", err)
		}

		return remediations, nil
	}