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
                Evidence: []model.EvidenceItem{
                    {
                        Type:        "profile",
                        Description: "Profile evidence",
                        Value:       "profile.pb.gz",
                        Weight:      1.0,
                    },
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
	generator, err := NewInsightsGenerator("test-key", "test-model", 10, 1000, 12000, true)
	require.NoError(t, err)
	require.NoError(t, err)
	
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
	generator, err := NewInsightsGenerator("", "test-model", 10, 1000, 12000, false)
	require.NoError(t, err)
	require.NoError(t, err)
	
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
	generator, err := NewInsightsGenerator("test-key", "test-model", 10, 1000, 100, false)
	require.NoError(t, err)
	require.NoError(t, err)
	
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
	generator, err := NewInsightsGenerator("test-key", "test-model", 10, 1000, 12000, false)
	require.NoError(t, err)
	require.NoError(t, err)
	
	// Verify configuration
	assert.NotNil(t, generator)
	assert.Equal(t, "test-model", generator.Provider.(*MistralProvider).modelName)
	assert.Equal(t, time.Duration(10)*time.Second, generator.Provider.(*MistralProvider).Timeout)
	assert.Equal(t, 1000, generator.Provider.(*MistralProvider).MaxResponse)
	assert.Equal(t, 12000, generator.MaxPromptChars)
	assert.False(t, generator.Provider.(*MistralProvider).DryRun)
}

func TestInsightsGenerator_WithLLM_DryRun(t *testing.T) {
	// Test dry-run mode
	generator, err := NewInsightsGenerator("test-key", "test-model", 10, 1000, 12000, true)
	require.NoError(t, err)
	require.NoError(t, err)
	
	assert.NotNil(t, generator)
	assert.True(t, generator.Provider.(*MistralProvider).DryRun)
}

func TestInsightsGenerator_WithLLM_NoAPIKey(t *testing.T) {
	// Test with empty API key
	generator, err := NewInsightsGenerator("", "test-model", 10, 1000, 12000, false)
	require.NoError(t, err)
	require.NoError(t, err)
	
	assert.NotNil(t, generator)
	assert.Equal(t, "", generator.Provider.(*MistralProvider).APIKey)
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
	generator, err := NewInsightsGenerator("test-key", "test-model", 10, 1000, 12000, false)
	require.NoError(t, err)
	
	assert.NotNil(t, generator)
	// generator.Client is not directly accessible, test through Provider interface
	assert.NotNil(t, generator.Provider)
	assert.False(t, generator.Provider.(*MistralProvider).DryRun)
}

func TestInsightsGenerator_WithRetries_DryRun(t *testing.T) {
	// Test insights generator with retries in dry-run mode
	generator, err := NewInsightsGenerator("test-key", "test-model", 10, 1000, 12000, true)
	require.NoError(t, err)
	
	assert.NotNil(t, generator)
	assert.NotNil(t, generator.Provider)
	assert.True(t, generator.Provider.(*MistralProvider).DryRun)
}

func TestEnhancedPromptBuilder(t *testing.T) {
	// Create test data with advanced features
	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Timestamp:   time.Now(),
			DurationSec: 30,
			Service:     "enhanced-test-service",
			Scenario:    "enhanced-test-scenario",
			GitSha:      "abc123def456ghi789",
		},
		Target: model.Target{
			Type:    "url",
			BaseURL: "http://localhost:6060/api",
		},
		Plugin: model.PluginRef{
			Name:    "enhanced-test-plugin",
			Version: "1.0.0",
		},
	}

	findings := &model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 85,
			TopIssueTags: []string{"cpu-optimization", "memory-leak", "concurrency"},
			Notes:        []string{"enhanced analysis required"},
		},
		Findings: []model.Finding{
			{
				Category:  "cpu",
				Title:     "CPU hotspot in critical path",
				Severity:  "high",
				Score:     92,
				Top: []model.StackFrame{
					{
						Function: "runtime.processRequest",
						File:     "/home/user/project/handler.go",
						Line:     42,
						Cum:      256.0,
						Flat:     128.0,
					},
				},
				Callgraph: []model.CallgraphNode{
					{
						Function: "runtime.processRequest",
						File:     "/home/user/project/handler.go",
						Line:     42,
						Depth:    0,
						Cum:      256.0,
						Flat:     128.0,
						Children: []model.CallgraphNode{
							{
								Function: "database.Query",
								File:     "/home/user/project/db.go",
								Line:     100,
								Depth:    1,
								Cum:      180.0,
								Flat:     90.0,
							},
						},
					},
				},
				AllocationAnalysis: &model.AllocationAnalysis{
					TotalAllocations: 1024.5,
					TopConcentration: 75.3,
					Severity:         "high",
					Score:            88,
					Hotspots: []model.AllocationHotspot{
						{
							Function: "memory.allocateBuffer",
							File:     "/home/user/project/memory.go",
							Line:     200,
							Count:    512.0,
							Percent:  50.0,
						},
					},
				},
				Regression: &model.RegressionAnalysis{
					BaselineScore: 75,
					CurrentScore:  92,
					Delta:         17,
					Percentage:    22.67,
					Severity:      "medium",
					Confidence:    85,
				},
                Evidence: []model.EvidenceItem{
                    {
                        Type:        "profile",
                        Description: "Profile evidence",
                        Value:       "profile.pb.gz",
                        Weight:      1.0,
                    },
                },
			},
		},
	}

	// Test enhanced prompt building
	builder := NewPromptBuilder(bundle, findings, 12000)
	prompt, err := builder.Build()
	require.NoError(t, err)
	require.NotEmpty(t, prompt)

	// Verify enhanced content
	assert.Contains(t, prompt, "=== TECHNICAL DEEP DIVE ===")
	assert.Contains(t, prompt, "Memory allocation patterns and optimization opportunities")
	assert.Contains(t, prompt, "CPU utilization breakdown by function and goroutine")
	assert.Contains(t, prompt, "Blocking operations and synchronization bottlenecks")
	assert.Contains(t, prompt, "Cache efficiency and data locality analysis")
	assert.Contains(t, prompt, "Algorithm complexity analysis and optimization suggestions")
	assert.Contains(t, prompt, "Concurrency patterns and parallelization opportunities")
	assert.Contains(t, prompt, "I/O patterns and optimization strategies")
	assert.Contains(t, prompt, "Garbage collection pressure and memory management")

	// Verify enhanced analysis requirements
	assert.Contains(t, prompt, "Performance impact: Quantitative estimate of improvement potential")
	assert.Contains(t, prompt, "Implementation complexity: Low/Medium/High with justification")
	assert.Contains(t, prompt, "ROI analysis: Cost-benefit assessment of proposed fixes")
	assert.Contains(t, prompt, "Include quantitative metrics and benchmarks where possible")
	assert.Contains(t, prompt, "Prioritize recommendations based on impact vs effort")

	// Verify callgraph analysis
	assert.Contains(t, prompt, "Callgraph Analysis:")
	assert.Contains(t, prompt, "runtime.processRequest (handler.go:42)")
	assert.Contains(t, prompt, "database.Query (db.go:100)")

	// Verify allocation analysis
	assert.Contains(t, prompt, "Allocation Analysis:")
	assert.Contains(t, prompt, "Total Allocations: 1024.50")
	assert.Contains(t, prompt, "Top Concentration: 75.30%")
	assert.Contains(t, prompt, "memory.allocateBuffer (memory.go:200)")

	// Verify regression analysis
	assert.Contains(t, prompt, "Regression Analysis:")
	assert.Contains(t, prompt, "Baseline Score: 75")
	assert.Contains(t, prompt, "Current Score: 92")
	assert.Contains(t, prompt, "Delta: 17 (22.67%)")
	assert.Contains(t, prompt, "Confidence: 85%")
}

func TestEnhancedInsightsSchema(t *testing.T) {
	// Test the new enhanced insights schema
	insights := &model.InsightsBundle{
		SchemaVersion:  "2.0",
		GeneratedAt:    time.Now(),
		Model:          "mistral-large-latest",
		RequestID:      "req-12345",
		ExecutiveSummary: model.ExecutiveSummary{
			Overview:        "Comprehensive performance analysis with enhanced insights",
			OverallSeverity: "high",
			KeyThemes:       []string{"CPU optimization", "Memory management", "Concurrency improvements"},
			Confidence:      90,
			PerformanceScore: 85,
			ImprovementPotential: 35, // 35% improvement potential
		},
		TopRisks: []model.RiskItem{
			{
				Description: "Critical CPU bottleneck in request processing",
				Severity:    "high",
				Impact:      "40% performance degradation",
				Likelihood:  "high",
				AffectedComponents: []string{"request handler", "database layer"},
				PotentialImpact: "30-50% latency reduction if fixed",
			},
		},
		TopActions: []model.ActionItem{
			{
				Description: "Optimize database query execution",
				Priority:       "high",
				EstimatedEffort: "medium",
				Categories:     []string{"database", "performance"},
				ImplementationComplexity: "medium",
				ExpectedImpact: "25-40% latency improvement",
				CodeExamples:    []string{"// Before: SELECT * FROM users WHERE active = true\n// After: SELECT id, name FROM users WHERE active = true LIMIT 1000"},
				ValidationMetrics: []string{"query execution time", "database load", "request latency"},
			},
		},
		PerformanceCategories: map[string]int{
			"cpu":      5,
			"memory":   3,
			"blocking": 2,
			"concurrency": 4,
		},
		PerFinding: []model.FindingInsight{
			{
				FindingID:        "cpu-hotspot-1",
				Narrative:        "Critical CPU bottleneck in request processing pipeline",
				LikelyRootCauses: []string{"Inefficient database queries", "Lack of connection pooling", "Suboptimal indexing"},
				Suggestions:      []string{"Implement query optimization", "Add database connection pooling", "Review and add appropriate indexes"},
				NextMeasurements:  []string{"Measure query execution time", "Monitor database connection usage", "Track index utilization"},
				Caveats:          []string{"Requires database schema changes", "May need application downtime for deployment"},
				Confidence:       88,
				PerformanceImpact: "30-50% latency reduction expected",
				ImplementationComplexity: "medium",
				CodeExamples:    []string{"// Add connection pooling\ndb, err := sql.Open(\"postgres\", connString)\ndb.SetMaxOpenConns(25)\ndb.SetMaxIdleConns(10)"},
				BeforeAfterMetrics: []string{"Before: 250ms avg query time", "After: 80ms avg query time", "Before: 100 connections/sec", "After: 50 connections/sec with pooling"},
			},
		},
		ROIAnalysis: []model.ROIItem{
			{
				ActionID:        "db-optimization",
				Description:     "Database query and connection optimization",
				EstimatedEffort: "2-3 days",
				ExpectedImpact:  "35% latency improvement, 20% resource reduction",
				CostBenefitRatio: 4.2,
				PriorityScore:   92,
			},
		},
		TechnicalDeepDive: model.TechnicalAnalysis{
			MemoryPatterns: []model.MemoryPattern{
				{
					PatternType:     "excessive allocations",
					Description:     "High allocation rate in JSON parsing",
					CurrentUsage:    "1024 allocations/sec",
					Optimization:    "Implement object pooling for JSON parsers",
					ExpectedSavings: "60% reduction in allocations",
					Implementation:  "Use sync.Pool for JSON parser instances",
				},
			},
			CPUUtilization: []model.CPUUtilization{
				{
					Component:      "request handler",
					CurrentUsage:   75.5,
					HotspotAnalysis: "70% time spent in database operations",
					Optimization:   "Implement caching layer for frequent queries",
					ExpectedGain:   40.0,
				},
			},
			BlockingOperations: []model.BlockingOperation{
				{
					OperationType:  "database query",
					Location:       "user authentication",
					CurrentLatency: "250ms average",
					RootCause:      "Missing indexes on user table",
					Solution:       "Add composite index on (email, active) columns",
					ExpectedGain:   "70% latency reduction",
				},
			},
		},
	}

	// Test JSON serialization
	data, err := json.MarshalIndent(insights, "", "  ")
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Test deserialization
	var deserialized model.InsightsBundle
	err = json.Unmarshal(data, &deserialized)
	require.NoError(t, err)

	// Verify enhanced fields
	assert.Equal(t, "2.0", deserialized.SchemaVersion)
	assert.Equal(t, "mistral-large-latest", deserialized.Model)
	assert.Equal(t, "req-12345", deserialized.RequestID)
	assert.Equal(t, 85, deserialized.ExecutiveSummary.PerformanceScore)
	assert.Equal(t, 35, deserialized.ExecutiveSummary.ImprovementPotential)
	assert.Equal(t, 1, len(deserialized.TopRisks))
	assert.Equal(t, "30-50% latency reduction if fixed", deserialized.TopRisks[0].PotentialImpact)
	assert.Equal(t, 2, len(deserialized.TopRisks[0].AffectedComponents))
	assert.Equal(t, 1, len(deserialized.TopActions))
	assert.Equal(t, "medium", deserialized.TopActions[0].ImplementationComplexity)
	assert.Equal(t, "25-40% latency improvement", deserialized.TopActions[0].ExpectedImpact)
	assert.Equal(t, 1, len(deserialized.TopActions[0].CodeExamples))
	assert.Equal(t, 1, len(deserialized.PerFinding))
	assert.Equal(t, "30-50% latency reduction expected", deserialized.PerFinding[0].PerformanceImpact)
	assert.Equal(t, "medium", deserialized.PerFinding[0].ImplementationComplexity)
	assert.Equal(t, 1, len(deserialized.ROIAnalysis))
	assert.Equal(t, 4.2, deserialized.ROIAnalysis[0].CostBenefitRatio)
	assert.Equal(t, 92, deserialized.ROIAnalysis[0].PriorityScore)
	assert.Equal(t, 1, len(deserialized.TechnicalDeepDive.MemoryPatterns))
	assert.Equal(t, 1, len(deserialized.TechnicalDeepDive.CPUUtilization))
	assert.Equal(t, 1, len(deserialized.TechnicalDeepDive.BlockingOperations))
}

func TestMistralClient_EnhancedConfiguration(t *testing.T) {
	// Test client creation with enhanced configuration
	client := NewMistralClientWithRetries("test-key", "mistral-large-latest", 30, 4096, 5, 2)
	
	assert.NotNil(t, client)
	assert.Equal(t, "test-key", client.APIKey)
	assert.Equal(t, "mistral-large-latest", client.Model)
	assert.Equal(t, 30, int(client.Timeout.Seconds()))
	assert.Equal(t, 4096, client.MaxResponse)
	assert.Equal(t, 5, client.MaxRetries)
	assert.Equal(t, 2, int(client.RetryDelay.Seconds()))
}
