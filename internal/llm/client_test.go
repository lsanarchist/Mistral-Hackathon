package llm

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromptBuilder_Build(t *testing.T) {
	// Create test data
	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Timestamp:   time.Now(),
			DurationSec: 10,
			Service:     "test-service",
			Scenario:    "test-scenario",
			GitSha:      "abc123def456",
		},
		Target: model.Target{
			Type:    "url",
			BaseURL: "http://localhost:6060",
		},
		Plugin: model.PluginRef{
			Name:    "test-plugin",
			Version: "0.1.0",
		},
	}

	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
			TopIssueTags: []string{"performance", "memory"},
			Notes:        []string{"test analysis"},
		},
		Findings: []model.Finding{
			{
				Category:  "cpu",
				Title:     "Top CPU hotspots",
				Severity:  "medium",
				Score:     80,
				Top: []model.StackFrame{
					{
						Function: "runtime.allocm",
						File:     "/home/user/project/proc.go",
						Line:     2276,
						Cum:      256.0,
						Flat:     256.0,
					},
				},
				Evidence: model.Evidence{
					ArtifactPath: "cpu.pb.gz",
					ProfileType:  "cpu",
					ExtractedAt:  time.Now(),
				},
			},
		},
	}

	// Test prompt building
	builder := NewPromptBuilder(bundle, findings, 12000)
	prompt, err := builder.Build()
	require.NoError(t, err)
	require.NotEmpty(t, prompt)

	// Verify redaction
	assert.Contains(t, prompt, "=== PROFILE METADATA ===")
	assert.Contains(t, prompt, "Service: test-service")
	assert.Contains(t, prompt, "Scenario: test-scenario")
	assert.Contains(t, prompt, "Git SHA: abc123d") // Should be truncated
	assert.Contains(t, prompt, "Target: http://[REDACTED_HOSTNAME]") // Should be redacted

	// Verify findings summary
	assert.Contains(t, prompt, "=== FINDINGS SUMMARY ===")
	assert.Contains(t, prompt, "Overall Score: 75/100 (medium)")
	assert.Contains(t, prompt, "Top Issues: performance, memory")
	assert.Contains(t, prompt, "Finding: Top CPU hotspots")
	assert.Contains(t, prompt, "Category: cpu")
	assert.Contains(t, prompt, "Severity: medium")
	assert.Contains(t, prompt, "Score: 80")

	// Verify enhanced analysis context
	assert.Contains(t, prompt, "=== ANALYSIS CONTEXT ===")
	assert.Contains(t, prompt, "You are an expert performance engineer analyzing profiling data.")
	assert.Contains(t, prompt, "Provide deep technical analysis with actionable insights.")
	assert.Contains(t, prompt, "=== ANALYSIS REQUIREMENTS ===")
	assert.Contains(t, prompt, "Narrative explanation: Clear technical explanation of the root cause")
	assert.Contains(t, prompt, "Likely root causes: 2-4 specific technical reasons with evidence")
	assert.Contains(t, prompt, "Concrete suggestions: Actionable recommendations with code examples")
	assert.Contains(t, prompt, "=== EXECUTIVE SUMMARY REQUIREMENTS ===")
	assert.Contains(t, prompt, "Executive summary: Concise overview with overall severity assessment")
	assert.Contains(t, prompt, "Top 3 risks: Most critical issues with impact analysis")
	assert.Contains(t, prompt, "Top 3 action items: Prioritized recommendations with effort estimates")
	assert.Contains(t, prompt, "Key themes: Patterns and common issues across findings")
	assert.Contains(t, prompt, "Performance categories: Distribution of issues by type")
	assert.Contains(t, prompt, "=== OUTPUT FORMAT REQUIREMENTS ===")
	assert.Contains(t, prompt, "Use JSON format with the exact schema provided")
	assert.Contains(t, prompt, "Be specific and technical in explanations")
	assert.Contains(t, prompt, "Provide code examples where applicable")

	// Verify function name redaction
	assert.Contains(t, prompt, "runtime.allocm (proc.go:2276)")
	assert.NotContains(t, prompt, "/home/user/project/") // Path should be redacted
}

func TestPromptBuilder_Redaction(t *testing.T) {
	builder := &PromptBuilder{
		MaxSize: 10000,
	}

	// Test URL redaction
	redactedURL := builder.redactURL("http://localhost:6060/debug/pprof/heap?token=secret123")
	assert.Equal(t, "http://[REDACTED_HOSTNAME]", redactedURL)

	// Test path redaction
	redactedPath := builder.redactPath("/home/user/project/main.go")
	assert.Equal(t, "main.go", redactedPath)

	// Test sensitive info redaction
	redactedInfo := builder.redactSensitiveInfo("token=abc123def456 secret=password123")
	assert.Contains(t, redactedInfo, "token=[REDACTED]")
	assert.Contains(t, redactedInfo, "secret=[REDACTED]")

	// Test long token redaction
	redactedToken := builder.redactSensitiveInfo("abc123def456ghi789jkl012mno345pqr678stu901")
	assert.Contains(t, redactedToken, "[REDACTED_TOKEN]")

	// Test function name redaction
	redactedFunc := builder.redactFunctionName("processRequestWithToken_abc123def456ghi789jkl012mno345pqr678stu901")
	assert.Contains(t, redactedFunc, "[REDACTED_TOKEN]")
}

func TestPromptBuilder_SizeLimit(t *testing.T) {
	// Create large bundle and findings
	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Service:  "test",
			Scenario: "test",
		},
		Target: model.Target{
			Type:    "url",
			BaseURL: "http://localhost:6060",
		},
	}

	// Create many findings to exceed size limit
	var findings model.FindingsBundle
	for i := 0; i < 100; i++ {
		findings.Findings = append(findings.Findings, model.Finding{
			Category:  "cpu",
			Title:     "Test finding",
			Severity:  "medium",
			Score:     80,
			Top: []model.StackFrame{
				{
					Function: "testFunction",
					File:     "test.go",
					Line:     1,
					Cum:      100.0,
					Flat:     100.0,
				},
			},
		})
	}

	// Test with small size limit
	builder := NewPromptBuilder(bundle, &findings, 100) // Very small limit
	_, err := builder.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum size")
}

func TestMistralClient_GenerateInsights_NoAPIKey(t *testing.T) {
	client := NewMistralClient("", "test-model", 10, 1000)
	
	insights, err := client.GenerateInsights(context.Background(), "test prompt")
	require.NoError(t, err)
	require.NotNil(t, insights)
	assert.Equal(t, "MISTRAL_API_KEY environment variable not set", insights.DisabledReason)
}

func TestMistralClient_GenerateInsights_PromptTooLarge(t *testing.T) {
	client := NewMistralClient("test-key", "test-model", 10, 1000)
	
	// Create very large prompt
	largePrompt := "x"
	for i := 0; i < 15000; i++ {
		largePrompt += "x"
	}
	
	insights, err := client.GenerateInsights(context.Background(), largePrompt)
	require.NoError(t, err)
	require.NotNil(t, insights)
	assert.Contains(t, insights.DisabledReason, "prompt too large")
}

func TestInsightsGenerator_GenerateInsights_DryRun(t *testing.T) {
	// Create test data
	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Service:  "test",
			Scenario: "test",
		},
		Target: model.Target{
			Type:    "url",
			BaseURL: "http://localhost:6060",
		},
	}

	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
		},
		Findings: []model.Finding{
			{
				Category:  "cpu",
				Title:     "Test finding",
				Severity:  "medium",
				Score:     80,
			},
		},
	}

	// Test dry-run mode
	generator := NewInsightsGenerator("test-key", "test-model", 10, 1000, 12000, true)
	
	insights, err := generator.GenerateInsights(context.Background(), bundle, findings)
	require.NoError(t, err)
	require.NotNil(t, insights)
	assert.Equal(t, "dry-run mode enabled - no API call made", insights.DisabledReason)
	
	// Verify prompt file was created
	promptData, err := os.ReadFile("llm_prompt.json")
	require.NoError(t, err)
	assert.NotEmpty(t, promptData)
	
	// Clean up
	os.Remove("llm_prompt.json")
}

func TestInsightsGenerator_GenerateInsights_NoAPIKey(t *testing.T) {
	// Create test data
	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Service:  "test",
			Scenario: "test",
		},
		Target: model.Target{
			Type:    "url",
			BaseURL: "http://localhost:6060",
		},
	}

	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
		},
		Findings: []model.Finding{
			{
				Category:  "cpu",
				Title:     "Test finding",
				Severity:  "medium",
				Score:     80,
			},
		},
	}

	// Test with no API key
	generator := NewInsightsGenerator("", "test-model", 10, 1000, 12000, false)
	
	insights, err := generator.GenerateInsights(context.Background(), bundle, findings)
	require.NoError(t, err)
	require.NotNil(t, insights)
	assert.Equal(t, "MISTRAL_API_KEY environment variable not set", insights.DisabledReason)
}

func TestInsightsGenerator_GenerateInsights_PromptTooLarge(t *testing.T) {
	// Create test data with many findings to exceed size limit
	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Service:  "test",
			Scenario: "test",
		},
		Target: model.Target{
			Type:    "url",
			BaseURL: "http://localhost:6060",
		},
	}

	var findings model.FindingsBundle
	for i := 0; i < 1000; i++ {
		findings.Findings = append(findings.Findings, model.Finding{
			Category:  "cpu",
			Title:     "Test finding",
			Severity:  "medium",
			Score:     80,
			Top: []model.StackFrame{
				{
					Function: "testFunction",
					File:     "test.go",
					Line:     1,
					Cum:      100.0,
					Flat:     100.0,
				},
			},
		})
	}

	// Test with small size limit
	generator := NewInsightsGenerator("test-key", "test-model", 10, 1000, 100, false)
	
	insights, err := generator.GenerateInsights(context.Background(), bundle, &findings)
	require.NoError(t, err)
	require.NotNil(t, insights)
	assert.Contains(t, insights.DisabledReason, "failed to build prompt")
}

func TestInsightsBundle_Serialization(t *testing.T) {
	// Test JSON serialization of insights bundle
	insights := &model.InsightsBundle{
		SchemaVersion:  "1.0",
		GeneratedAt:    time.Now(),
		Model:          "test-model",
		DisabledReason: "test reason",
		ExecutiveSummary: model.ExecutiveSummary{
			Overview:        "test overview",
			OverallSeverity: "medium",
			KeyThemes:       []string{"theme1", "theme2"},
			Confidence:      85,
		},
		TopRisks: []model.RiskItem{
			{
				Description: "test risk",
				Severity:    "high",
				Impact:      "performance",
				Likelihood:  "high",
			},
		},
		TopActions: []model.ActionItem{
			{
				Description:    "test action",
				Priority:       "high",
				EstimatedEffort: "medium",
				Categories:     []string{"code", "optimization"},
			},
		},
		PerformanceCategories: map[string]int{
			"cpu":      3,
			"memory":   2,
			"blocking": 1,
		},
		PerFinding: []model.FindingInsight{
			{
				FindingID:        "cpu",
				Narrative:        "test narrative",
				LikelyRootCauses: []string{"cause1", "cause2"},
				Suggestions:      []string{"suggestion1", "suggestion2"},
				NextMeasurements: []string{"measurement1"},
				Caveats:          []string{"caveat1"},
				Confidence:       80,
			},
		},
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(insights, "", "  ")
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Deserialize back
	var deserialized model.InsightsBundle
	err = json.Unmarshal(data, &deserialized)
	require.NoError(t, err)

	// Verify fields
	assert.Equal(t, "1.0", deserialized.SchemaVersion)
	assert.Equal(t, "test-model", deserialized.Model)
	assert.Equal(t, "test reason", deserialized.DisabledReason)
	assert.Equal(t, "test overview", deserialized.ExecutiveSummary.Overview)
	assert.Equal(t, 85, deserialized.ExecutiveSummary.Confidence)
	assert.Equal(t, 1, len(deserialized.TopRisks))
	assert.Equal(t, 1, len(deserialized.TopActions))
	assert.Equal(t, 3, len(deserialized.PerformanceCategories))
	assert.Equal(t, 3, deserialized.PerformanceCategories["cpu"])
	assert.Equal(t, 2, deserialized.PerformanceCategories["memory"])
	assert.Equal(t, 1, deserialized.PerformanceCategories["blocking"])
	assert.Equal(t, 1, len(deserialized.PerFinding))
}

func TestInsightsGenerator_WithLLM(t *testing.T) {
	// Test that insights generator can be configured
	generator := NewInsightsGenerator("test-key", "test-model", 10, 1000, 12000, false)
	
	// Verify configuration
	assert.NotNil(t, generator)
	assert.Equal(t, "test-model", generator.Client.Model)
	assert.Equal(t, time.Duration(10)*time.Second, generator.Client.Timeout)
	assert.Equal(t, 1000, generator.Client.MaxResponse)
	assert.Equal(t, 12000, generator.MaxPromptChars)
	assert.False(t, generator.DryRun)
}

func TestInsightsGenerator_WithLLM_DryRun(t *testing.T) {
	// Test dry-run mode
	generator := NewInsightsGenerator("test-key", "test-model", 10, 1000, 12000, true)
	
	assert.NotNil(t, generator)
	assert.True(t, generator.DryRun)
}

func TestInsightsGenerator_WithLLM_NoAPIKey(t *testing.T) {
	// Test with empty API key
	generator := NewInsightsGenerator("", "test-model", 10, 1000, 12000, false)
	
	assert.NotNil(t, generator)
	assert.Equal(t, "", generator.Client.APIKey)
}

func TestMistralClient_WithRetries(t *testing.T) {
	// Test client creation with retry configuration
	client := NewMistralClientWithRetries("test-key", "test-model", 10, 1000, 5, 2)
	
	assert.NotNil(t, client)
	assert.Equal(t, "test-key", client.APIKey)
	assert.Equal(t, "test-model", client.Model)
	assert.Equal(t, 10, int(client.Timeout.Seconds()))
	assert.Equal(t, 1000, client.MaxResponse)
	assert.Equal(t, 5, client.MaxRetries)
	assert.Equal(t, 2, int(client.RetryDelay.Seconds()))
}

func TestInsightsGenerator_WithRetries(t *testing.T) {
	// Test insights generator with retry configuration
	generator := NewInsightsGeneratorWithRetries("test-key", "test-model", 10, 1000, 12000, 3, 1, false)
	
	assert.NotNil(t, generator)
	assert.Equal(t, "test-key", generator.Client.APIKey)
	assert.Equal(t, 3, generator.Client.MaxRetries)
	assert.Equal(t, 1, int(generator.Client.RetryDelay.Seconds()))
	assert.False(t, generator.DryRun)
}

func TestInsightsGenerator_WithRetries_DryRun(t *testing.T) {
	// Test insights generator with retries in dry-run mode
	generator := NewInsightsGeneratorWithRetries("test-key", "test-model", 10, 1000, 12000, 3, 1, true)
	
	assert.NotNil(t, generator)
	assert.True(t, generator.DryRun)
	assert.Equal(t, 3, generator.Client.MaxRetries)
}
