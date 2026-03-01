package llm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/mistral-hackathon/triageprof/internal/model"
)

func TestInsightsValidation(t *testing.T) {
	// Create test findings bundle
	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 85,
			TopIssueTags:  []string{"cpu", "alloc"},
		},
		Findings: []model.Finding{
			{
				ID:       "find-001",
				Title:    "High CPU usage in main loop",
				Category: "cpu",
				Severity: "high",
				Score:    95,
			},
			{
				ID:       "find-002",
				Title:    "Memory allocation hotspot",
				Category: "alloc",
				Severity: "medium",
				Score:    75,
			},
		},
	}

	// Create test profile bundle
	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Timestamp:   time.Now(),
			DurationSec: 30,
			Service:     "test-service",
			Scenario:    "load-test",
			GitSha:      "abc123def456",
		},
		Target: model.Target{
			Type:    "http",
			BaseURL: "https://example.com",
		},
	}
	_ = bundle // Use bundle to avoid "declared and not used" error

	t.Run("Valid insights pass validation", func(t *testing.T) {
		validInsights := &model.InsightsBundle{
			SchemaVersion: "2.0",
			GeneratedAt:   time.Now(),
			Model:         "test-model",
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:         "Performance analysis complete",
				OverallSeverity: "high",
				Confidence:       90,
			},
			PerFinding: []model.FindingInsight{
				{
					FindingID:        "find-001",
					Narrative:        "Finding find-001 shows high CPU usage in the main loop function",
					Confidence:       95,
					CodeExamples:     []string{"optimize loop bounds"},
				},
			},
		}

		err := validateInsights(validInsights, findings)
		assert.NoError(t, err)
	})

	t.Run("Empty finding ID fails validation", func(t *testing.T) {
		invalidInsights := &model.InsightsBundle{
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:         "Analysis complete",
				OverallSeverity: "medium",
				Confidence:       80,
			},
			PerFinding: []model.FindingInsight{
				{
					FindingID: "", // Empty ID
					Narrative: "Some analysis",
					Confidence: 85,
				},
			},
		}

		err := validateInsights(invalidInsights, findings)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty finding_id")
	})

	t.Run("Non-existent finding ID fails validation", func(t *testing.T) {
		invalidInsights := &model.InsightsBundle{
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:         "Analysis complete",
				OverallSeverity: "medium",
				Confidence:       80,
			},
			PerFinding: []model.FindingInsight{
				{
					FindingID: "find-999", // Non-existent
					Narrative: "Some analysis",
					Confidence: 85,
				},
			},
		}

		err := validateInsights(invalidInsights, findings)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "non-existent finding_id")
	})

	t.Run("Missing evidence reference fails validation", func(t *testing.T) {
		invalidInsights := &model.InsightsBundle{
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:         "Analysis complete",
				OverallSeverity: "medium",
				Confidence:       80,
			},
			PerFinding: []model.FindingInsight{
				{
					FindingID: "find-001",
					Narrative: "High CPU usage detected", // No reference to find-001
					Confidence: 85,
				},
			},
		}

		err := validateInsights(invalidInsights, findings)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must contain evidence reference")
	})

	t.Run("Code example exceeding 200 chars fails validation", func(t *testing.T) {
		longCode := "This is a very long code example that exceeds the 200 character limit. " +
			"It contains multiple lines of code and detailed explanations that make it much longer than allowed. " +
			"The purpose is to test the validation that should catch this and return an error."

		invalidInsights := &model.InsightsBundle{
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:         "Analysis complete",
				OverallSeverity: "medium",
				Confidence:       80,
			},
			PerFinding: []model.FindingInsight{
				{
					FindingID:    "find-001",
					Narrative:    "Finding find-001 shows high CPU usage",
					Confidence:   85,
					CodeExamples: []string{longCode},
				},
			},
		}

		err := validateInsights(invalidInsights, findings)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds 200 character limit")
	})

	t.Run("Invalid confidence range fails validation", func(t *testing.T) {
		invalidInsights := &model.InsightsBundle{
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:         "Analysis complete",
				OverallSeverity: "medium",
				Confidence:       150, // Invalid
			},
			PerFinding: []model.FindingInsight{
				{
					FindingID:  "find-001",
					Narrative:  "Finding find-001 shows high CPU usage",
					Confidence: 85,
				},
			},
		}

		err := validateInsights(invalidInsights, findings)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between 0-100")
	})
}

func TestInsightsGuardrails(t *testing.T) {

	t.Run("Guardrails truncate long fields", func(t *testing.T) {
		insights := &model.InsightsBundle{
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:   "This is a very long overview that should be truncated by the guardrails function because it exceeds the maximum allowed length for this field",
				Confidence: 85,
			},
			PerFinding: []model.FindingInsight{
				{
					FindingID:  "test-finding",
					Narrative:  "This is an extremely long narrative that should definitely be truncated by the guardrails to ensure it doesn't exceed the maximum length",
					Confidence: 90,
					CodeExamples: []string{"This is a very long code example that exceeds the 200 character limit and should be truncated by the guardrails function to comply with the requirements"},
				},
			},
		}

		applyInsightsGuardrails(insights)

		// Check truncation
		assert.True(t, len(insights.ExecutiveSummary.Overview) <= 1000)
		assert.True(t, len(insights.PerFinding[0].Narrative) <= 2000)
		assert.True(t, len(insights.PerFinding[0].CodeExamples[0]) <= 200)
		assert.Contains(t, insights.PerFinding[0].CodeExamples[0], "...")
	})

	t.Run("Guardrails set default values", func(t *testing.T) {
		insights := &model.InsightsBundle{
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:         "Test overview",
				OverallSeverity: "high",
				Confidence:       85,
			},
		}

		applyInsightsGuardrails(insights)

		// Check defaults are set
		assert.False(t, insights.GeneratedAt.IsZero())
		assert.Equal(t, "2.0", insights.SchemaVersion)
		assert.Equal(t, "unknown", insights.Model)
	})
}

func TestInsightsRedaction(t *testing.T) {
	insights := &model.InsightsBundle{
		ExecutiveSummary: model.ExecutiveSummary{
			Overview:   "Service running on localhost:8080 with token abc123def456",
			Confidence: 85,
		},
		PerFinding: []model.FindingInsight{
			{
				FindingID:  "test-finding",
				Narrative:  "API key is secret1234567890 and password is hunter2",
				Confidence: 90,
				CodeExamples: []string{"http://localhost:3000/api?token=verysecrettoken123"},
			},
		},
	}

	applyRedactionToInsights(insights)

	// Check redactions
	assert.Contains(t, insights.ExecutiveSummary.Overview, "[REDACTED_HOSTNAME]")
	assert.Contains(t, insights.ExecutiveSummary.Overview, "[REDACTED_TOKEN]")
	assert.Contains(t, insights.PerFinding[0].Narrative, "[REDACTED]")
	assert.Contains(t, insights.PerFinding[0].CodeExamples[0], "[REDACTED_URL]")
}

func TestFindingReferenceDetection(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		findingID string
		want     bool
	}{
		{
			name:     "Direct reference",
			text:     "Analysis of finding find-001 shows issues",
			findingID: "find-001",
			want:     true,
		},
		{
			name:     "Finding prefix",
			text:     "Finding find-001 has high severity",
			findingID: "find-001",
			want:     true,
		},
		{
			name:     "Issue reference",
			text:     "Issue find-001 needs attention",
			findingID: "find-001",
			want:     true,
		},
		{
			name:     "No reference",
			text:     "High CPU usage detected in main loop",
			findingID: "find-001",
			want:     false,
		},
		{
			name:     "Different finding ID",
			text:     "Analysis of finding find-002",
			findingID: "find-001",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsFindingReference(tt.text, tt.findingID)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestInsightsGeneratorIntegration(t *testing.T) {
	// Create a mock provider that returns valid insights
	mockProvider := &MockProvider{
		insights: &model.InsightsBundle{
			SchemaVersion: "2.0",
			GeneratedAt:   time.Now(),
			ExecutiveSummary: model.ExecutiveSummary{
				Overview:         "Test analysis complete",
				OverallSeverity: "medium",
				Confidence:       85,
			},
			PerFinding: []model.FindingInsight{
				{
					FindingID:  "test-finding",
					Narrative:  "Finding test-finding shows performance issues",
					Confidence: 90,
				},
			},
		},
	}

	findings := &model.FindingsBundle{
		Findings: []model.Finding{
			{
				ID:       "test-finding",
				Title:    "Test finding",
				Category: "cpu",
			},
		},
	}

	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Service:  "test-service",
			Scenario: "test-scenario",
		},
	}

	generator := &InsightsGenerator{
		Provider:       mockProvider,
		Cache:          NewInsightsCache(CacheConfig{Enabled: false}),
		MaxPromptChars: 10000,
	}

	insights, err := generator.GenerateInsights(context.Background(), bundle, findings)
	require.NoError(t, err)
	assert.NotNil(t, insights)
	assert.Equal(t, "2.0", insights.SchemaVersion)
	assert.Equal(t, "Test analysis complete", insights.ExecutiveSummary.Overview)
}

// MockProvider for testing
type MockProvider struct {
	insights *model.InsightsBundle
	err      error
}

func (m *MockProvider) GenerateInsights(ctx context.Context, prompt string) (*model.InsightsBundle, error) {
	return m.insights, m.err
}

func (m *MockProvider) Name() string {
	return "mock"
}

func (m *MockProvider) Model() string {
	return "mock-model"
}