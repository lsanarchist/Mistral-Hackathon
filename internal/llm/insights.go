package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// InsightsGenerator orchestrates LLM insight generation
type InsightsGenerator struct {
	Client          *MistralClient
	DryRun          bool
	MaxPromptChars  int
}

// NewInsightsGenerator creates a new insights generator
func NewInsightsGenerator(apiKey, model string, timeout int, maxResponse, maxPromptChars int, dryRun bool) *InsightsGenerator {
	timeoutDur := time.Duration(timeout) * time.Second
	return &InsightsGenerator{
		Client:         NewMistralClient(apiKey, model, timeoutDur, maxResponse),
		DryRun:          dryRun,
		MaxPromptChars:  maxPromptChars,
	}
}

// GenerateInsights creates LLM insights from bundle and findings
func (g *InsightsGenerator) GenerateInsights(ctx context.Context, bundle *model.ProfileBundle, findings *model.FindingsBundle) (*model.InsightsBundle, error) {
	// Build prompt
	builder := NewPromptBuilder(bundle, findings)
	builder.MaxSize = g.MaxPromptChars
	prompt, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	// Dry run mode - save prompt to file and return disabled insights
	if g.DryRun {
		if err := os.WriteFile("llm_prompt.json", []byte(prompt), 0644); err != nil {
			return nil, fmt.Errorf("failed to write dry-run prompt: %w", err)
		}
		return &model.InsightsBundle{
			SchemaVersion:  model.InsightsSchemaVersion,
			GeneratedAt:    time.Now(),
			DisabledReason: "Dry run mode - prompt saved to llm_prompt.json",
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:        "LLM insights disabled: dry run mode",
				OverallSeverity: model.SeverityLow,
				Confidence:      0,
			},
		}, nil
	}

	// Generate insights via Mistral API
	insights, err := g.Client.GenerateInsights(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate insights: %w", err)
	}

	// Validate insights structure
	if insights.ExecutiveSummary.Overview == "" {
		insights.ExecutiveSummary.Overview = "No specific insights generated"
	}
	if insights.ExecutiveSummary.Confidence == 0 {
		insights.ExecutiveSummary.Confidence = 50 // Default confidence
	}

	return insights, nil
}

// GenerateInsightsFromFiles loads bundle and findings from files and generates insights
func GenerateInsightsFromFiles(ctx context.Context, bundlePath, findingsPath, outputPath string, 
	apiKey, llmModel string, timeout, maxResponse, maxPromptChars int, dryRun bool) error {

	// Load bundle
	bundleData, err := os.ReadFile(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to read bundle: %w", err)
	}
	var profileBundle model.ProfileBundle
	if err := json.Unmarshal(bundleData, &profileBundle); err != nil {
		return fmt.Errorf("failed to parse bundle: %w", err)
	}

	// Load findings
	findingsData, err := os.ReadFile(findingsPath)
	if err != nil {
		return fmt.Errorf("failed to read findings: %w", err)
	}
	var findingsBundle model.FindingsBundle
	if err := json.Unmarshal(findingsData, &findingsBundle); err != nil {
		return fmt.Errorf("failed to parse findings: %w", err)
	}

	// Generate insights
	generator := NewInsightsGenerator(apiKey, llmModel, timeout, maxResponse, maxPromptChars, dryRun)
	insights, err := generator.GenerateInsights(ctx, &profileBundle, &findingsBundle)
	if err != nil {
		return fmt.Errorf("failed to generate insights: %w", err)
	}

	// Save insights
	insightsData, err := json.MarshalIndent(insights, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal insights: %w", err)
	}

	if err := os.WriteFile(outputPath, insightsData, 0644); err != nil {
		return fmt.Errorf("failed to write insights: %w", err)
	}

	return nil
}