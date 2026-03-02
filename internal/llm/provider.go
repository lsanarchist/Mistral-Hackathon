// Package llm provides LLM integration for performance insights
package llm

import (
	"context"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// Provider defines the interface for LLM providers
type Provider interface {
	// GenerateInsights generates performance insights from findings
	GenerateInsights(ctx context.Context, prompt string) (*model.InsightsBundle, error)
	
	// Name returns the provider name
	Name() string
	
	// Model returns the current model being used
	Model() string
}

// ProviderConfig holds configuration for LLM providers
type ProviderConfig struct {
	ProviderName string            `json:"provider" yaml:"provider"`
	Model        string            `json:"model" yaml:"model"`
	APIKey       string            `json:"apiKey" yaml:"apiKey"`
	Timeout      time.Duration     `json:"timeout" yaml:"timeout"`
	MaxResponse  int               `json:"maxResponse" yaml:"maxResponse"`
	MaxPrompt    int               `json:"maxPrompt" yaml:"maxPrompt"`
	DryRun       bool              `json:"dryRun" yaml:"dryRun"`
	ProviderConfig map[string]string `json:"providerConfig" yaml:"providerConfig"`
}

// NewProvider creates a new LLM provider based on configuration
func NewProvider(config ProviderConfig) (Provider, error) {
	switch config.ProviderName {
	case "mistral", "":
		return NewMistralProvider(config)
	case "openai":
		return NewOpenAIProvider(config)
	default:
		return nil, ErrUnknownProvider
	}
}

// DefaultConfig returns default configuration for Mistral provider
func DefaultConfig() ProviderConfig {
	return ProviderConfig{
		ProviderName: "mistral",
		Model:        "mistral-large-latest",
		Timeout:      20 * time.Second,
		MaxResponse:  8192,
		MaxPrompt:    12000,
		DryRun:       false,
		ProviderConfig: map[string]string{},
	}
}

// ErrUnknownProvider is returned when an unknown provider is requested
var ErrUnknownProvider = NewLLMError("unknown LLM provider")

// LLMError represents LLM-specific errors
type LLMError struct {
	Message string
}

func NewLLMError(msg string) *LLMError {
	return &LLMError{Message: msg}
}

func (e *LLMError) Error() string {
	return "llm: " + e.Message
}

func (e *LLMError) Is(target error) bool {
	_, ok := target.(*LLMError)
	return ok
}