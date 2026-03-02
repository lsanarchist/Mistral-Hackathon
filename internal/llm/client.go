package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"math/rand"

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

	// Add enhanced analysis context for better LLM insights
	sections = append(sections, "")
	sections = append(sections, "=== ANALYSIS CONTEXT ===")
	sections = append(sections, "You are an expert performance engineer analyzing profiling data.")
	sections = append(sections, "Provide deep technical analysis with actionable insights.")
	sections = append(sections, "Focus on production-grade recommendations with concrete implementation details.")
	sections = append(sections, "")
	sections = append(sections, "=== ANALYSIS REQUIREMENTS ===")
	sections = append(sections, "For each finding, provide:")
	sections = append(sections, "1. Narrative explanation: Clear technical explanation of the root cause")
	sections = append(sections, "2. Likely root causes: 2-4 specific technical reasons with evidence")
	sections = append(sections, "3. Concrete suggestions: Actionable recommendations with code examples")
	sections = append(sections, "4. Next measurements: Specific metrics to validate fixes")
	sections = append(sections, "5. Caveats: Limitations and assumptions of the analysis")
	sections = append(sections, "6. Confidence score: 0-100 based on evidence quality")
	sections = append(sections, "7. Performance impact: Quantitative estimate of improvement potential")
	sections = append(sections, "8. Implementation complexity: Low/Medium/High with justification")
	sections = append(sections, "")
	sections = append(sections, "=== EXECUTIVE SUMMARY REQUIREMENTS ===")
	sections = append(sections, "Also provide:")
	sections = append(sections, "- Executive summary: Concise overview with overall severity assessment")
	sections = append(sections, "- Top 3 risks: Most critical issues with impact analysis (high/medium/low)")
	sections = append(sections, "- Top 3 action items: Prioritized recommendations with effort estimates")
	sections = append(sections, "- Key themes: Patterns and common issues across findings")
	sections = append(sections, "- Performance categories: Distribution of issues by type")
	sections = append(sections, "- ROI analysis: Cost-benefit assessment of proposed fixes")
	sections = append(sections, "")
	sections = append(sections, "=== OUTPUT FORMAT REQUIREMENTS ===")
	sections = append(sections, "- Use JSON format with the exact schema provided")
	sections = append(sections, "- Be specific and technical in explanations")
	sections = append(sections, "- Provide code examples where applicable")
	sections = append(sections, "- Reference specific functions and files from the evidence")
	sections = append(sections, "- Use appropriate confidence scores based on evidence quality")
	sections = append(sections, "- Include quantitative metrics and benchmarks where possible")
	sections = append(sections, "- Prioritize recommendations based on impact vs effort")
	sections = append(sections, "- CRITICAL: Include evidence_refs for all insights citing specific finding IDs")
	sections = append(sections, "- CRITICAL: All code examples must be limited to 200 characters maximum")
	sections = append(sections, "- CRITICAL: Use schema_version 2.0 for all responses")

	// Add technical deep dive section
	sections = append(sections, "")
	sections = append(sections, "=== TECHNICAL DEEP DIVE ===")
	sections = append(sections, "Provide advanced analysis including:")
	sections = append(sections, "- Memory allocation patterns and optimization opportunities")
	sections = append(sections, "- CPU utilization breakdown by function and goroutine")
	sections = append(sections, "- Blocking operations and synchronization bottlenecks")
	sections = append(sections, "- Cache efficiency and data locality analysis")
	sections = append(sections, "- Algorithm complexity analysis and optimization suggestions")
	sections = append(sections, "- Concurrency patterns and parallelization opportunities")
	sections = append(sections, "- I/O patterns and optimization strategies")
	sections = append(sections, "- Garbage collection pressure and memory management")

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

		// Add finding context
		sections = append(sections, "")
		sections = append(sections, "Context:")
		sections = append(sections, "This finding represents a performance bottleneck that requires attention.")
		sections = append(sections, "Analyze the technical details and provide specific, actionable recommendations.")

		// Top stack frames (redacted and limited)
		if len(finding.Top) > 0 {
			sections = append(sections, "")
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

		// Evidence (redacted) - handle both new and legacy formats
		evidenceDesc := getEvidenceDescription(finding)
		if evidenceDesc != "" {
			sections = append(sections, fmt.Sprintf("Evidence: %s", evidenceDesc))
		}

		// Add callgraph analysis if available
		if len(finding.Callgraph) > 0 {
			sections = append(sections, "")
			sections = append(sections, "Callgraph Analysis:")
			callgraphCount := len(finding.Callgraph)
			if callgraphCount > 5 {
				callgraphCount = 5
			}
			for j := 0; j < callgraphCount; j++ {
				node := finding.Callgraph[j]
				p.addCallgraphNode(&sections, node, 0)
			}
		}

		// Add allocation analysis if available
		if finding.AllocationAnalysis != nil {
			sections = append(sections, "")
			sections = append(sections, "Allocation Analysis:")
			sections = append(sections, fmt.Sprintf("  Total Allocations: %.2f", finding.AllocationAnalysis.TotalAllocations))
			sections = append(sections, fmt.Sprintf("  Top Concentration: %.2f%%", finding.AllocationAnalysis.TopConcentration))
			sections = append(sections, fmt.Sprintf("  Severity: %s", finding.AllocationAnalysis.Severity))
			
			if len(finding.AllocationAnalysis.Hotspots) > 0 {
				hotspotsCount := len(finding.AllocationAnalysis.Hotspots)
				if hotspotsCount > 5 {
					hotspotsCount = 5
				}
				for j := 0; j < hotspotsCount; j++ {
					hotspot := finding.AllocationAnalysis.Hotspots[j]
					redactedFunc := p.redactFunctionName(hotspot.Function)
					redactedFile := p.redactPath(hotspot.File)
					sections = append(sections, fmt.Sprintf("  - %s (%s:%d) - count: %.2f, percent: %.2f%%", 
						redactedFunc, redactedFile, hotspot.Line, hotspot.Count, hotspot.Percent))
				}
			}
		}

		// Add regression analysis if available
		if finding.Regression != nil {
			sections = append(sections, "")
			sections = append(sections, "Regression Analysis:")
			sections = append(sections, fmt.Sprintf("  Baseline Score: %d", finding.Regression.BaselineScore))
			sections = append(sections, fmt.Sprintf("  Current Score: %d", finding.Regression.CurrentScore))
			sections = append(sections, fmt.Sprintf("  Delta: %d (%.2f%%)", finding.Regression.Delta, finding.Regression.Percentage))
			sections = append(sections, fmt.Sprintf("  Severity: %s", finding.Regression.Severity))
			sections = append(sections, fmt.Sprintf("  Confidence: %d%%", finding.Regression.Confidence))
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

// addCallgraphNode adds a callgraph node and its children recursively
func (p *PromptBuilder) addCallgraphNode(sections *[]string, node model.CallgraphNode, depth int) {
	indent := ""
	for i := 0; i <= depth; i++ {
		indent += "  "
	}

	redactedFunc := p.redactFunctionName(node.Function)
	redactedFile := p.redactPath(node.File)
	*sections = append(*sections, fmt.Sprintf("%s- %s (%s:%d) - cum: %.2f%%, flat: %.2f%%", 
		indent, redactedFunc, redactedFile, node.Line, node.Cum, node.Flat))

	// Add children recursively (limit depth to avoid excessive output)
	if depth < 3 && len(node.Children) > 0 {
		childCount := len(node.Children)
		if childCount > 5 {
			childCount = 5
		}
		for i := 0; i < childCount; i++ {
			p.addCallgraphNode(sections, node.Children[i], depth+1)
		}
	}
}

// getEvidenceDescription extracts evidence description from finding (handles both new and legacy formats)
func getEvidenceDescription(finding model.Finding) string {
	// Handle new evidence format
	if len(finding.Evidence) > 0 {
		evidenceTypes := []string{}
		for _, evidence := range finding.Evidence {
			evidenceTypes = append(evidenceTypes, evidence.Type)
		}
		if len(evidenceTypes) > 0 {
			return fmt.Sprintf("%s evidence", strings.Join(evidenceTypes, ", "))
		}
	}
	
	// Handle legacy evidence format
	if finding.EvidenceLegacy.ProfileType != "" {
		return fmt.Sprintf("%s profile", finding.EvidenceLegacy.ProfileType)
	}
	
	return ""
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
	MaxRetries  int
	RetryDelay  time.Duration
}

// NewMistralClient creates a new Mistral API client
func NewMistralClient(apiKey, model string, timeout, maxResponse int) *MistralClient {
	return NewMistralClientWithRetries(apiKey, model, timeout, maxResponse, 3, 1)
}

// NewMistralClientWithRetries creates a new Mistral API client with retry configuration
func NewMistralClientWithRetries(apiKey, model string, timeout, maxResponse, maxRetries, retryDelaySec int) *MistralClient {
	return &MistralClient{
		APIKey:      apiKey,
		Model:       model,
		Timeout:     time.Duration(timeout) * time.Second,
		MaxResponse: maxResponse,
		MaxRetries:  maxRetries,
		RetryDelay:  time.Duration(retryDelaySec) * time.Second,
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
			DisabledReason: "MISTRAL_API_KEY not configured",
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
		"temperature": 0.2, // More deterministic responses for technical analysis
		"top_p": 0.9,
		"frequency_penalty": 0.0,
		"presence_penalty": 0.0,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return &model.InsightsBundle{
			DisabledReason: fmt.Sprintf("failed to marshal request: %v", err),
		}, nil
	}

	// Retry loop with exponential backoff
	var lastError error
	var insights *model.InsightsBundle
	
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		if attempt > 0 {
			// Apply exponential backoff with jitter
			delay := c.RetryDelay * time.Duration(1<<(attempt-1))
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			totalDelay := delay + jitter
			
			select {
			case <-time.After(totalDelay):
			case <-ctx.Done():
				return &model.InsightsBundle{
					DisabledReason: fmt.Sprintf("context canceled during retry delay: %v", ctx.Err()),
				}, nil
			}
		}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		lastError = fmt.Errorf("failed to create request (attempt %d): %v", attempt+1, err)
		continue
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("User-Agent", "TriageProf/1.0")

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		lastError = fmt.Errorf("API request failed (attempt %d): %v", attempt+1, err)
		continue
	}
	
	// Handle response
	insights, err = c.processAPIResponse(resp)
	if err == nil {
		break // Success
	}
	lastError = err
	
	// Close response body
	resp.Body.Close()
	
	// Log retry attempt
	log.Printf("Mistral API attempt %d failed: %v, retrying...", attempt+1, err)
	}

	if insights != nil {
		// Set metadata
		insights.GeneratedAt = time.Now()
		insights.Model = c.Model
		insights.SchemaVersion = "2.0" // Updated schema version
		return insights, nil
	}

	// All attempts failed
	return &model.InsightsBundle{
		DisabledReason: fmt.Sprintf("all retry attempts failed: %v", lastError),
	}, nil
}

// processAPIResponse handles the API response parsing and error checking
func (c *MistralClient) processAPIResponse(resp *http.Response) (*model.InsightsBundle, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
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
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(apiResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from API")
	}

	// Parse the insights from the response
	var insights model.InsightsBundle
	if err := json.Unmarshal([]byte(apiResponse.Choices[0].Message.Content), &insights); err != nil {
		return nil, fmt.Errorf("failed to parse insights: %v", err)
	}

	return &insights, nil
}

