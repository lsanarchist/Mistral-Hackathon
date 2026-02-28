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