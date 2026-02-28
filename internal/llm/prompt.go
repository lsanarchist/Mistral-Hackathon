package llm

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// PromptBuilder creates structured prompts with redaction
type PromptBuilder struct {
	Bundle    *model.ProfileBundle
	Findings  *model.FindingsBundle
	MaxSize   int
}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder(bundle *model.ProfileBundle, findings *model.FindingsBundle, maxSize int) *PromptBuilder {
	return &PromptBuilder{
		Bundle:   bundle,
		Findings: findings,
		MaxSize:  maxSize,
	}
}

// Build creates final prompt with redaction and size limiting
func (p *PromptBuilder) Build() (string, error) {
	if p.Bundle == nil || p.Findings == nil {
		return "", fmt.Errorf("bundle and findings are required")
	}

	var sb strings.Builder

	// Header
	sb.WriteString("Analyze the following performance findings and provide insights.\n")
	sb.WriteString("Focus on root causes, actionable recommendations, and confidence levels.\n")
	sb.WriteString("Be concise and technical.\n\n")

	// Metadata (redacted)
	sb.WriteString("=== METADATA ===\n")
	sb.WriteString(fmt.Sprintf("Service: %s\n", redactSensitiveInfo(p.Bundle.Metadata.Service)))
	sb.WriteString(fmt.Sprintf("Scenario: %s\n", redactSensitiveInfo(p.Bundle.Metadata.Scenario)))
	sb.WriteString(fmt.Sprintf("Duration: %d seconds\n", p.Bundle.Metadata.DurationSec))
	sb.WriteString(fmt.Sprintf("Timestamp: %s\n\n", p.Bundle.Metadata.Timestamp.Format("2006-01-02 15:04:05")))

	// Findings summary
	sb.WriteString("=== FINDINGS SUMMARY ===\n")
	sb.WriteString(fmt.Sprintf("Overall Score: %d/100\n", p.Findings.Summary.OverallScore))
	sb.WriteString(fmt.Sprintf("Top Issues: %s\n\n", strings.Join(p.Findings.Summary.TopIssueTags, ", ")))

	// Per-finding details (limited to top 10)
	for i, finding := range p.Findings.Findings {
		if i >= 10 {
			break
		}

		sb.WriteString(fmt.Sprintf("=== FINDING %d: %s ===\n", i+1, finding.Category))
		sb.WriteString(fmt.Sprintf("Title: %s\n", finding.Title))
		sb.WriteString(fmt.Sprintf("Severity: %s\n", finding.Severity))
		sb.WriteString(fmt.Sprintf("Score: %d\n", finding.Score))
		sb.WriteString(fmt.Sprintf("Profile Type: %s\n", finding.Evidence.ProfileType))
		sb.WriteString(fmt.Sprintf("Artifact: %s\n\n", redactPath(finding.Evidence.ArtifactPath)))

		// Top stack frames (limited to 5)
		if len(finding.Top) > 0 {
			sb.WriteString("Top Hotspots:\n")
			for j, frame := range finding.Top {
				if j >= 5 {
					break
				}
				sb.WriteString(fmt.Sprintf("  %d. %s (%s:%d) - Cum: %.2f, Flat: %.2f\n",
					j+1,
					reactStackFrame(frame.Function),
					reactPath(frame.File),
					frame.Line,
					frame.Cum,
					frame.Flat))
			}
			sb.WriteString("\n")
		}
	}

	// Instructions
	sb.WriteString("=== INSTRUCTIONS ===\n")
	sb.WriteString("Provide insights in JSON format matching the InsightsBundle schema.\n")
	sb.WriteString("Include executive summary, top risks, top actions, and per-finding analysis.\n")
	sb.WriteString("Use confidence scores (0-100) for all insights.\n")
	sb.WriteString("Be specific about root causes and actionable recommendations.\n")

	prompt := sb.String()

	// Validate size
	if len(prompt) > p.MaxSize {
		return "", fmt.Errorf("prompt exceeds maximum size of %d characters (actual: %d)", p.MaxSize, len(prompt))
	}

	return prompt, nil
}

// redactSensitiveInfo removes sensitive information from strings
func redactSensitiveInfo(input string) string {
	if input == "" {
		return input
	}

	// Common sensitive patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(token|secret|key|password|credential)[=: ]*[A-Za-z0-9]{8,}`),
		regexp.MustCompile(`(?i)(localhost|127\.0\.0\.1|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`),
		regexp.MustCompile(`(?i)(http|https)://[^\s]+`),
	}

	result := input
	for _, pattern := range patterns {
		result = pattern.ReplaceAllString(result, "[REDACTED]")
	}

	// Limit length
	if len(result) > 200 {
		result = result[:200] + "..."
	}

	return result
}

// redactPath keeps only filename and removes directory paths
func redactPath(path string) string {
	if path == "" {
		return path
	}
	return filepath.Base(path)
}

// redactStackFrame sanitizes function names and files
func redactStackFrame(function string) string {
	if function == "" {
		return function
	}

	// Remove sensitive function names
	sensitivePrefixes := []string{"auth", "Auth", "token", "Token", "secret", "Secret"}
	for _, prefix := range sensitivePrefixes {
		if strings.HasPrefix(function, prefix) {
			return "[REDACTED_FUNCTION]"
		}
	}

	// Limit length
	if len(function) > 100 {
		return function[:100] + "..."
	}

	return function
}