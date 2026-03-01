package llm

import (
	"context"
	"fmt"
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
		Cache:          NewInsightsCache(CacheConfig{Enabled: false}), // Disabled by default
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
		Cache:          NewInsightsCache(CacheConfig{Enabled: false}), // Disabled by default
		MaxPromptChars: config.MaxPrompt,
	}, nil
}

// GenerateInsights creates LLM insights from bundle and findings
func (g *InsightsGenerator) GenerateInsights(ctx context.Context, 
	bundle *model.ProfileBundle, findings *model.FindingsBundle) (*model.InsightsBundle, error) {

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

	// Validate and return insights
	if insights == nil {
		return &model.InsightsBundle{
			DisabledReason: "no insights generated",
		}, nil
	}

	return insights, nil
}