package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// MistralClient handles communication with Mistral API
type MistralClient struct {
	APIKey      string
	Model       string
	Timeout     time.Duration
	MaxResponse int
	HTTPClient  *http.Client
}

// MistralRequest represents a chat completion request
type MistralRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MistralResponse represents a chat completion response
type MistralResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewMistralClient creates a new Mistral API client
func NewMistralClient(apiKey, model string, timeout time.Duration, maxResponse int) *MistralClient {
	return &MistralClient{
		APIKey:      apiKey,
		Model:       model,
		Timeout:     timeout,
		MaxResponse: maxResponse,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GenerateInsights calls Mistral API to generate insights from findings
func (c *MistralClient) GenerateInsights(ctx context.Context, prompt string) (*model.InsightsBundle, error) {
	// Check if API key is available
	if c.APIKey == "" {
		apiKey := os.Getenv("MISTRAL_API_KEY")
		if apiKey == "" {
			return &model.InsightsBundle{
				SchemaVersion:  model.InsightsSchemaVersion,
				GeneratedAt:    time.Now(),
				DisabledReason: "MISTRAL_API_KEY environment variable not set",
				ExecutiveSummary: model.ExecutiveSummary{
					Overview:        "LLM insights disabled: API key not configured",
					OverallSeverity: model.SeverityLow,
					Confidence:      0,
				},
			}, nil
		}
		c.APIKey = apiKey
	}

	// Create request
	request := MistralRequest{
		Model: c.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a performance triage assistant. Output JSON ONLY matching the provided schema. No markdown.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.2, // More deterministic
		MaxTokens:   c.MaxResponse,
	}

	// Marshal request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.mistral.ai/v1/chat/completions", bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Execute request with retry logic
	response, err := c.executeWithRetry(ctx, req, 3)
	if err != nil {
		return nil, fmt.Errorf("Mistral API request failed: %w", err)
	}
	defer response.Body.Close()

	// Check status code
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return nil, fmt.Errorf("Mistral API returned status %d: %s", response.StatusCode, string(body))
	}

	// Read and limit response
	var responseBody []byte
	if int64(c.MaxResponse) > 0 {
		limitedReader := io.LimitReader(response.Body, int64(c.MaxResponse))
		responseBody, err = io.ReadAll(limitedReader)
	} else {
		responseBody, err = io.ReadAll(response.Body)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var mistralResponse MistralResponse
	if err := json.Unmarshal(responseBody, &mistralResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Mistral API response: %w", err)
	}

	// Extract insights from response
	if len(mistralResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from Mistral API")
	}

	// Parse the JSON response from LLM
	var insights model.InsightsBundle
	if err := json.Unmarshal([]byte(mistralResponse.Choices[0].Message.Content), &insights); err != nil {
		return nil, fmt.Errorf("failed to parse LLM insights JSON: %w", err)
	}

	// Set metadata
	insights.SchemaVersion = model.InsightsSchemaVersion
	insights.GeneratedAt = time.Now()
	insights.Model = mistralResponse.Model
	insights.RequestID = mistralResponse.ID

	return &insights, nil
}

// executeWithRetry executes HTTP request with retry logic
func (c *MistralClient) executeWithRetry(ctx context.Context, req *http.Request, maxRetries int) (*http.Response, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		response, err := c.HTTPClient.Do(req)
		if err == nil {
			return response, nil
		}

		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 500 * time.Millisecond) // Exponential backoff
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}
