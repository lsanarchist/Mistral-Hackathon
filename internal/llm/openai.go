package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// OpenAIProvider implements the OpenAI API client
type OpenAIProvider struct {
	APIKey      string
	modelName   string
	Timeout     time.Duration
	MaxResponse int
	HTTPClient  *http.Client
	DryRun      bool
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config ProviderConfig) (*OpenAIProvider, error) {
	// Allow empty API key at construction; GenerateInsights will return disabled insights at call time.

	if config.Model == "" {
		config.Model = "gpt-3.5-turbo"
	}

	if config.Timeout == 0 {
		config.Timeout = 20 * time.Second
	}

	if config.MaxResponse == 0 {
		config.MaxResponse = 4096
	}

	return &OpenAIProvider{
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
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Model returns the current model name
func (p *OpenAIProvider) Model() string {
	return p.modelName
}

// GenerateInsights calls OpenAI API to generate insights
func (p *OpenAIProvider) GenerateInsights(ctx context.Context, prompt string) (*model.InsightsBundle, error) {
	if p.DryRun {
		return createDisabledInsights("dry-run mode enabled"), nil
	}

	if p.APIKey == "" {
		return createDisabledInsights("openai API key not configured"), nil
	}

	// Prepare request
	url := "https://api.openai.com/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model":       p.Model,
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

	// Parse insights from response
	insights, err := parseInsightsResponse(result.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse insights: %w", err)
	}

	insights.Model = p.modelName
	insights.RequestID = fmt.Sprintf("openai-%d", time.Now().Unix())

	return insights, nil
}