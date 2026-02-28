package llm

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// InsightsGenerator orchestrates LLM insight generation
type InsightsGenerator struct {
	Client         *MistralClient
	DryRun         bool
	MaxPromptChars int
}

// NewInsightsGenerator creates a new insights generator
func NewInsightsGenerator(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool) *InsightsGenerator {
	client := NewMistralClient(apiKey, model, timeout, maxResponse)
	return &InsightsGenerator{
		Client:         client,
		DryRun:         dryRun,
		MaxPromptChars: maxPromptChars,
	}
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

	// Handle dry-run mode
	if g.DryRun {
		// Save prompt for inspection
		if err := os.WriteFile("llm_prompt.json", []byte(prompt), 0644); err != nil {
			return &model.InsightsBundle{
				DisabledReason: fmt.Sprintf("failed to save prompt in dry-run mode: %v", err),
			}, nil
		}
		
		return &model.InsightsBundle{
			DisabledReason: "dry-run mode enabled - no API call made",
			GeneratedAt:   time.Now(),
		}, nil
	}

	// Call Mistral API
	insights, err := g.Client.GenerateInsights(ctx, prompt)
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

