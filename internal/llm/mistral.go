package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// MistralProvider implements the Mistral API client
type MistralProvider struct {
	APIKey      string
	modelName   string
	Timeout     time.Duration
	MaxResponse int
	HTTPClient  *http.Client
	DryRun      bool
}

// NewMistralProvider creates a new Mistral provider
func NewMistralProvider(config ProviderConfig) (*MistralProvider, error) {
	// Allow empty API key at construction; GenerateInsights will return disabled insights at call time.

	if config.Model == "" {
		config.Model = "devstral-small-latest"
	}

	if config.Timeout == 0 {
		config.Timeout = 20 * time.Second
	}

	if config.MaxResponse == 0 {
		config.MaxResponse = 4096
	}

	return &MistralProvider{
		APIKey:      config.APIKey,
		modelName:   config.Model,
		Timeout:     config.Timeout,
		MaxResponse: config.MaxResponse,
		HTTPClient: &http.Client{
			Timeout: config.Timeout,
		},
		DryRun: config.DryRun,
	}, nil
}

// Name returns the provider name
func (p *MistralProvider) Name() string {
	return "mistral"
}

// Model returns the current model name
func (p *MistralProvider) Model() string {
	return p.modelName
}

// GenerateInsights calls Mistral API to generate insights
func (p *MistralProvider) GenerateInsights(ctx context.Context, prompt string) (*model.InsightsBundle, error) {
	if p.DryRun {
		return createDisabledInsights("dry-run mode enabled - no API call made"), nil
	}

	if p.APIKey == "" {
		return createDisabledInsights("MISTRAL_API_KEY not configured"), nil
	}

	// Prepare request
	url := "https://api.mistral.ai/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model":       p.modelName,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens":  p.MaxResponse,
		"temperature": 0.2,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return nil, NewLLMError("no insights generated")
	}

	rawContent := result.Choices[0].Message.Content

	// Parse insights from response
	insights, err := parseInsightsResponse(rawContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse insights: %w", err)
	}

	insights.Model = p.modelName
	insights.RequestID = fmt.Sprintf("mistral-%d", time.Now().Unix())

	return insights, nil
}

// stripProblematicFields removes fields that LLMs commonly return with wrong types
// to prevent json.Unmarshal failures on the main struct.
func stripProblematicFields(raw string) string {
	// Parse into a generic map, drop known-problematic optional fields, re-serialize
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return raw // can't fix it, return as-is
	}
	// Drop roi_analysis if it exists (LLMs often return it as an object instead of array)
	delete(m, "roi_analysis")
	// Drop technical_deep_dive if present (complex nested object, optional)
	delete(m, "technical_deep_dive")
	result, err := json.Marshal(m)
	if err != nil {
		return raw
	}
	return string(result)
}

// parseInsightsResponse parses the LLM response into structured insights with strict validation
func parseInsightsResponse(response string) (*model.InsightsBundle, error) {
	// Strip markdown code fences if present (```json ... ``` or ``` ... ```)
	trimmed := strings.TrimSpace(response)
	if strings.HasPrefix(trimmed, "```") {
		// Remove opening fence (```json or ```)
		end := strings.Index(trimmed, "\n")
		if end >= 0 {
			trimmed = trimmed[end+1:]
		}
		// Remove closing fence
		if idx := strings.LastIndex(trimmed, "```"); idx >= 0 {
			trimmed = trimmed[:idx]
		}
		trimmed = strings.TrimSpace(trimmed)
	}

	// First try strict JSON parsing
	var insights model.InsightsBundle
	insights.GeneratedAt = time.Now()
	insights.SchemaVersion = "2.0"

	// Try to parse as JSON first (preferred format)
	// Use a two-pass approach: first strip fields that could have wrong types (e.g. roi_analysis)
	cleanJSON := stripProblematicFields(trimmed)
	if err := json.Unmarshal([]byte(cleanJSON), &insights); err == nil {
		// Apply guardrails (truncation) BEFORE validation so limits are enforced
		applyInsightsGuardrails(&insights)
		// JSON parsing succeeded, validate the structure
		if err := validateParsedInsights(&insights); err != nil {
			return nil, fmt.Errorf("insights validation failed: %v", err)
		}
		return &insights, nil
	}

	// Fallback to legacy text parsing for backward compatibility
	return parseLegacyTextResponse(response)
}

// validateParsedInsights validates the structure of parsed JSON insights
func validateParsedInsights(insights *model.InsightsBundle) error {
	// Validate executive summary
	if insights.ExecutiveSummary.Overview == "" {
		return fmt.Errorf("executive summary overview is required")
	}

	if insights.ExecutiveSummary.OverallSeverity == "" {
		return fmt.Errorf("executive summary severity is required")
	}

	// Validate confidence range
	if insights.ExecutiveSummary.Confidence < 0 || insights.ExecutiveSummary.Confidence > 100 {
		return fmt.Errorf("confidence must be between 0-100")
	}

	// Validate per-finding insights
	for i, finding := range insights.PerFinding {
		if finding.FindingID == "" {
			return fmt.Errorf("finding %d has empty finding_id", i)
		}
		if finding.Narrative == "" {
			return fmt.Errorf("finding %d has empty narrative", i)
		}
		if finding.Confidence < 0 || finding.Confidence > 100 {
			return fmt.Errorf("finding %d confidence must be between 0-100", i)
		}
	}

	return nil
}

// parseLegacyTextResponse parses legacy text format responses (backward compatibility)
func parseLegacyTextResponse(response string) (*model.InsightsBundle, error) {
	var insights model.InsightsBundle
	insights.GeneratedAt = time.Now()
	insights.SchemaVersion = "1.0"

	// Extract sections from response
	lines := strings.Split(response, "\n")
	var currentSection string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "## ") {
			currentSection = strings.TrimPrefix(line, "## ")
			continue
		}

		switch currentSection {
		case "Executive Summary":
			if insights.ExecutiveSummary.Overview == "" {
				insights.ExecutiveSummary.Overview = line
			} else if strings.HasPrefix(line, "Severity: ") {
				insights.ExecutiveSummary.OverallSeverity = strings.TrimPrefix(line, "Severity: ")
			} else if strings.HasPrefix(line, "Confidence: ") {
				fmt.Sscanf(line, "Confidence: %d", &insights.ExecutiveSummary.Confidence)
			}
		case "Top Risks":
			if !strings.HasPrefix(line, "- ") {
				continue
			}
			risk := strings.TrimPrefix(line, "- ")
			insights.TopRisks = append(insights.TopRisks, model.RiskItem{
				Description: risk,
				Severity:    "medium",
			})
		case "Top Actions":
			if !strings.HasPrefix(line, "- ") {
				continue
			}
			action := strings.TrimPrefix(line, "- ")
			insights.TopActions = append(insights.TopActions, model.ActionItem{
				Description: action,
				Priority:    "medium",
			})
		}
	}

	return &insights, nil
}

// createDisabledInsights creates an insights bundle with disabled reason
func createDisabledInsights(reason string) *model.InsightsBundle {
	return &model.InsightsBundle{
		SchemaVersion:  "1.0",
		GeneratedAt:    time.Now(),
		DisabledReason: reason,
		ExecutiveSummary: model.ExecutiveSummary{
			Overview:        "LLM insights disabled",
			OverallSeverity: "none",
			Confidence:      0,
		},
	}
}