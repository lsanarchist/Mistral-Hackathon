package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

	// Validate prompt
	if len(prompt) == 0 {
		return nil, fmt.Errorf("empty prompt")
	}

	// Create API request
	apiURL := "https://api.mistral.ai/v1/chat/completions"
	
	requestBody := map[string]interface{}{
		"model":    c.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": c.MaxResponse,
		"temperature": 0.7,
	}

	// Marshal request
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
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
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResponse.Error.Message != "" {
		return nil, fmt.Errorf("API error: %s", apiResponse.Error.Message)
	}

	// Extract content
	if len(apiResponse.Choices) == 0 || apiResponse.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("no insights generated")
	}

	content := apiResponse.Choices[0].Message.Content

	// Parse insights
	var insights model.InsightsBundle
	if err := json.Unmarshal([]byte(content), &insights); err != nil {
		return nil, fmt.Errorf("failed to parse insights: %w", err)
	}

	// Set metadata
	insights.GeneratedAt = time.Now()
	insights.Model = c.Model

	return &insights, nil
}

// GetAPIKeyFromEnv retrieves API key from environment variable
func GetAPIKeyFromEnv() string {
	return os.Getenv("MISTRAL_API_KEY")
}