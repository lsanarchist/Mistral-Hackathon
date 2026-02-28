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
    Client         *MistralClient
    DryRun         bool
    MaxPromptChars int
}

// NewInsightsGenerator creates a new insights generator
func NewInsightsGenerator(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool) *InsightsGenerator {
    client := NewMistralClient(apiKey, model,
        WithTimeout(time.Duration(timeout)*time.Second),
        WithMaxResponse(maxResponse))

    return &InsightsGenerator{
        Client:         client,
        DryRun:         dryRun,
        MaxPromptChars: maxPromptChars,
    }
}

// GenerateInsights creates LLM insights from bundle and findings
func (g *InsightsGenerator) GenerateInsights(ctx context.Context, 
    bundle *model.ProfileBundle, findings *model.FindingsBundle) (*model.InsightsBundle, error) {
    
    // Build secure prompt
    builder := NewPromptBuilder(bundle, findings, g.MaxPromptChars)
    userPrompt, err := builder.BuildUserPrompt()
    if err != nil {
        return nil, fmt.Errorf("failed to build prompt: %w", err)
    }

    // Handle dry-run mode
    if g.DryRun {
        // Save prompt for inspection
        promptData := map[string]interface{}{
            "system_prompt": BuildSystemPrompt(),
            "user_prompt":   userPrompt,
            "timestamp":     time.Now(),
        }
        
        promptJSON, _ := json.MarshalIndent(promptData, "", "  ")
        
        if err := os.WriteFile("llm_prompt.json", promptJSON, 0644); err != nil {
            return nil, fmt.Errorf("failed to save dry-run prompt: %w", err)
        }

        return g.Client.createDisabledBundle("dry-run mode enabled"), nil
    }

    // Call Mistral API
    insights, err := g.Client.GenerateInsights(ctx, userPrompt)
    if err != nil {
        return nil, fmt.Errorf("LLM insights generation failed: %w", err)
    }

    return insights, nil
}

// GenerateInsightsFromFiles standalone file-based generation
func GenerateInsightsFromFiles(ctx context.Context, bundlePath, findingsPath, outputPath string,
    apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool) error {
    
    // Read bundle
    bundleData, err := os.ReadFile(bundlePath)
    if err != nil {
        return fmt.Errorf("failed to read bundle: %w", err)
    }
    
    var bundle model.ProfileBundle
    if err := json.Unmarshal(bundleData, &bundle); err != nil {
        return fmt.Errorf("failed to parse bundle: %w", err)
    }

    // Read findings
    findingsData, err := os.ReadFile(findingsPath)
    if err != nil {
        return fmt.Errorf("failed to read findings: %w", err)
    }
    
    var findings model.FindingsBundle
    if err := json.Unmarshal(findingsData, &findings); err != nil {
        return fmt.Errorf("failed to parse findings: %w", err)
    }

    // Generate insights
    generator := NewInsightsGenerator(apiKey, model, timeout, maxResponse, maxPromptChars, dryRun)
    insights, err := generator.GenerateInsights(ctx, &bundle, &findings)
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