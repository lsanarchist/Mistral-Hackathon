package llm

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// InsightsGenerator orchestrates LLM insight generation
type InsightsGenerator struct {
	Client         *MistralClient
	Cache          *InsightsCache
	DryRun         bool
	MaxPromptChars int
}

// NewInsightsGenerator creates a new insights generator
func NewInsightsGenerator(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool) *InsightsGenerator {
	client := NewMistralClient(apiKey, model, timeout, maxResponse)
	return &InsightsGenerator{
		Client:         client,
		Cache:          NewInsightsCache(CacheConfig{Enabled: false}), // Disabled by default
		DryRun:         dryRun,
		MaxPromptChars: maxPromptChars,
	}
}

// NewInsightsGeneratorWithRetries creates a new insights generator with retry configuration
func NewInsightsGeneratorWithRetries(apiKey, model string, timeout, maxResponse, maxPromptChars, maxRetries, retryDelaySec int, dryRun bool) *InsightsGenerator {
	client := NewMistralClientWithRetries(apiKey, model, timeout, maxResponse, maxRetries, retryDelaySec)
	return &InsightsGenerator{
		Client:         client,
		Cache:          NewInsightsCache(CacheConfig{Enabled: false}), // Disabled by default
		DryRun:         dryRun,
		MaxPromptChars: maxPromptChars,
	}
}

// NewInsightsGeneratorWithCache creates a new insights generator with caching enabled
func NewInsightsGeneratorWithCache(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool, cacheConfig CacheConfig) *InsightsGenerator {
	client := NewMistralClient(apiKey, model, timeout, maxResponse)
	return &InsightsGenerator{
		Client:         client,
		Cache:          NewInsightsCache(cacheConfig),
		DryRun:         dryRun,
		MaxPromptChars: maxPromptChars,
	}
}

// NewInsightsGeneratorWithCacheAndRetries creates a new insights generator with caching and retries
func NewInsightsGeneratorWithCacheAndRetries(apiKey, model string, timeout, maxResponse, maxPromptChars, maxRetries, retryDelaySec int, dryRun bool, cacheConfig CacheConfig) *InsightsGenerator {
	client := NewMistralClientWithRetries(apiKey, model, timeout, maxResponse, maxRetries, retryDelaySec)
	return &InsightsGenerator{
		Client:         client,
		Cache:          NewInsightsCache(cacheConfig),
		DryRun:         dryRun,
		MaxPromptChars: maxPromptChars,
	}
}

// GenerateInsights creates LLM insights from bundle and findings
func (g *InsightsGenerator) GenerateInsights(ctx context.Context, 
	bundle *model.ProfileBundle, findings *model.FindingsBundle) (*model.InsightsBundle, error) {

	// Check cache first if enabled
	if g.Cache != nil && g.Cache.config.Enabled {
		if cachedInsights, found := g.Cache.GetCachedInsights(ctx, bundle, findings); found {
			log.Printf("Using cached insights for profile %s", bundle.Metadata.Service)
			return cachedInsights, nil
		}
	}

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

	// Cache the insights if caching is enabled
	if g.Cache != nil && g.Cache.config.Enabled {
		if err := g.Cache.CacheInsights(ctx, bundle, findings, insights); err != nil {
			log.Printf("Warning: failed to cache insights: %v", err)
		}
	}

	return insights, nil
}

