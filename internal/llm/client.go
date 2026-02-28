package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
	var sections []string

	// Metadata section
	metadata := p.buildMetadataSection()
	if metadata != "" {
		sections = append(sections, metadata)
	}

	// Findings summary section
	findingsSummary := p.buildFindingsSummary()
	if findingsSummary != "" {
		sections = append(sections, findingsSummary)
	}

	// Combine sections
	prompt := strings.Join(sections, "\n\n")

	// Validate size
	if len(prompt) > p.MaxSize {
		return "", fmt.Errorf("prompt exceeds maximum size: %d > %d", len(prompt), p.MaxSize)
	}

	return prompt, nil
}

// buildMetadataSection creates redacted metadata section
func (p *PromptBuilder) buildMetadataSection() string {
	if p.Bundle == nil {
		return ""
	}

	var lines []string
	lines = append(lines, "=== PROFILE METADATA ===")
	
	// Redact service name
	service := p.redactSensitiveInfo(p.Bundle.Metadata.Service)
	if service != "" {
		lines = append(lines, fmt.Sprintf("Service: %s", service))
	}

	// Redact scenario
	scenario := p.redactSensitiveInfo(p.Bundle.Metadata.Scenario)
	if scenario != "" {
		lines = append(lines, fmt.Sprintf("Scenario: %s", scenario))
	}

	// Duration
	if p.Bundle.Metadata.DurationSec > 0 {
		lines = append(lines, fmt.Sprintf("Duration: %d seconds", p.Bundle.Metadata.DurationSec))
	}

	// Redact git SHA (keep only first 7 chars)
	if p.Bundle.Metadata.GitSha != "" {
		gitSha := p.Bundle.Metadata.GitSha
		if len(gitSha) > 7 {
			gitSha = gitSha[:7]
		}
		lines = append(lines, fmt.Sprintf("Git SHA: %s", gitSha))
	}

	// Target info (redacted)
	if p.Bundle.Target.BaseURL != "" {
		redactedURL := p.redactURL(p.Bundle.Target.BaseURL)
		lines = append(lines, fmt.Sprintf("Target: %s", redactedURL))
	}

	return strings.Join(lines, "\n")
}

// buildFindingsSummary creates redacted findings summary
func (p *PromptBuilder) buildFindingsSummary() string {
	if p.Findings == nil || len(p.Findings.Findings) == 0 {
		return ""
	}

	var sections []string
	sections = append(sections, "=== FINDINGS SUMMARY ===")

	// Overall summary
	if p.Findings.Summary.OverallScore > 0 {
		severity := "low"
		if p.Findings.Summary.OverallScore > 70 {
			severity = "medium"
		}
		if p.Findings.Summary.OverallScore > 90 {
			severity = "high"
		}
		sections = append(sections, fmt.Sprintf("Overall Score: %d/100 (%s)", 
			p.Findings.Summary.OverallScore, severity))
	}

	if len(p.Findings.Summary.TopIssueTags) > 0 {
		sections = append(sections, fmt.Sprintf("Top Issues: %s", 
			strings.Join(p.Findings.Summary.TopIssueTags, ", ")))
	}

	// Per-finding details (limit to top 5)
	findingsCount := len(p.Findings.Findings)
	if findingsCount > 5 {
		findingsCount = 5
	}

	for i := 0; i < findingsCount; i++ {
		finding := p.Findings.Findings[i]
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("---"))
		sections = append(sections, fmt.Sprintf("Finding: %s", finding.Title))
		sections = append(sections, fmt.Sprintf("Category: %s", finding.Category))
		sections = append(sections, fmt.Sprintf("Severity: %s", finding.Severity))
		sections = append(sections, fmt.Sprintf("Score: %d", finding.Score))

		// Top stack frames (redacted and limited)
		if len(finding.Top) > 0 {
			sections = append(sections, "Top Hotspots:")
			hotspotsCount := len(finding.Top)
			if hotspotsCount > 10 {
				hotspotsCount = 10
			}
			for j := 0; j < hotspotsCount; j++ {
				frame := finding.Top[j]
				redactedFunc := p.redactFunctionName(frame.Function)
				redactedFile := p.redactPath(frame.File)
				sections = append(sections, fmt.Sprintf("  - %s (%s:%d) - %.2f%%", 
					redactedFunc, redactedFile, frame.Line, frame.Flat))
			}
		}

		// Evidence (redacted)
		if finding.Evidence.ProfileType != "" {
			sections = append(sections, fmt.Sprintf("Evidence: %s profile", finding.Evidence.ProfileType))
		}
	}

	return strings.Join(sections, "\n")
}

// redactSensitiveInfo removes sensitive information from strings
func (p *PromptBuilder) redactSensitiveInfo(text string) string {
	if text == "" {
		return ""
	}

	// Redact common sensitive patterns
	patterns := []struct {
		regex    *regexp.Regexp
		replace string
	}{
		{regexp.MustCompile(`(?i)(token|secret|key|password)=[^&\s]+`), "$1=[REDACTED]"},
		{regexp.MustCompile(`(?i)(localhost|127\.0\.0\.1|\d+\.\d+\.\d+\.\d+)`), "[REDACTED_HOSTNAME]"},
		{regexp.MustCompile(`[A-Za-z0-9]{32,}`), "[REDACTED_TOKEN]"},
	}

	result := text
	for _, pattern := range patterns {
		result = pattern.regex.ReplaceAllString(result, pattern.replace)
	}

	// Limit length
	if len(result) > 200 {
		result = result[:200] + "..."
	}

	return result
}

// redactURL redacts sensitive parts of URLs
func (p *PromptBuilder) redactURL(url string) string {
	// Simple URL redaction - keep scheme and path, redact host
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return strings.Replace(url, url[strings.Index(url, "://")+3:], "[REDACTED_HOSTNAME]", 1)
	}
	return "[REDACTED_URL]"
}

// redactPath redacts absolute paths, keeping only filename
func (p *PromptBuilder) redactPath(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}

// redactFunctionName redacts sensitive info from function names
func (p *PromptBuilder) redactFunctionName(funcName string) string {
	if funcName == "" {
		return ""
	}

	// Redact long tokens that might be secrets
	patterns := []struct {
		regex    *regexp.Regexp
		replace string
	}{
		{regexp.MustCompile(`[A-Za-z0-9]{32,}`), "[REDACTED_TOKEN]"},
	}

	result := funcName
	for _, pattern := range patterns {
		result = pattern.regex.ReplaceAllString(result, pattern.replace)
	}

	// Limit length
	if len(result) > 100 {
		result = result[:100] + "..."
	}

	return result
}

// MistralClient handles communication with Mistral API
type MistralClient struct {
	APIKey      string
	Model       string
	Timeout     time.Duration
	MaxResponse int
	HTTPClient  *http.Client
}

// NewMistralClient creates a new Mistral API client
func NewMistralClient(apiKey, model string, timeout, maxResponse int) *MistralClient {
	return &MistralClient{
		APIKey:      apiKey,
		Model:       model,
		Timeout:     time.Duration(timeout) * time.Second,
		MaxResponse: maxResponse,
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// GenerateInsights calls Mistral API to generate insights
func (c *MistralClient) GenerateInsights(ctx context.Context, prompt string) (*model.InsightsBundle, error) {
	// Validate API key
	if c.APIKey == "" {
		return &model.InsightsBundle{
			DisabledReason: "MISTRAL_API_KEY environment variable not set",
		}, nil
	}

	// Validate prompt size
	if len(prompt) > 12000 {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("prompt too large: %d characters (max 12000)", len(prompt)),
		}, nil
	}

	// Prepare request
	apiURL := "https://api.mistral.ai/v1/chat/completions"
	
	requestBody := map[string]interface{}{
		"model":  c.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": c.MaxResponse,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("failed to marshal request: %v", err),
		}, nil
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("failed to create request: %v", err),
		}, nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("API request failed: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)),
		}, nil
	}

	// Parse response
	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("failed to parse response: %v", err),
		}, nil
	}

	if len(apiResponse.Choices) == 0 {
		return &model.InsightsBundle{
			DisabledReason: "no choices returned from API",
		}, nil
	}

	// Parse the insights from the response
	var insights model.InsightsBundle
	if err := json.Unmarshal([]byte(apiResponse.Choices[0].Message.Content), &insights); err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("failed to parse insights: %v", err),
		}, nil
	}

	// Set metadata
	insights.GeneratedAt = time.Now()
	insights.Model = c.Model
	insights.SchemaVersion = "1.0"

	return &insights, nil
}

