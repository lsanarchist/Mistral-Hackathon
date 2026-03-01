// Copyright 2026 Mistral AI. All rights reserved.
// Use of this source code is governed by the Apache 2.0 license.

package llm

import (
	"context"
	"testing"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
)

// MockLLMProvider is a mock implementation for testing
type MockLLMProvider struct {
	response string
	err      error
}

func (m *MockLLMProvider) GenerateInsights(ctx context.Context, prompt string) (*model.InsightsBundle, error) {
	// For remediation testing, we don't need the full insights bundle
	var bundle model.InsightsBundle
	return &bundle, m.err
}

func (m *MockLLMProvider) Name() string {
	return "mock"
}

func (m *MockLLMProvider) Model() string {
	return "mock-model"
}

// Generate method for remediation testing
func (m *MockLLMProvider) Generate(ctx context.Context, prompt string) (string, error) {
	return m.response, m.err
}

func TestRemediationGenerator(t *testing.T) {
	// Create test findings
	findings := &model.FindingsBundle{
		Findings: []model.Finding{
			{
				ID:           "finding-001",
				Title:        "High CPU usage in JSON parsing",
				Category:     "cpu",
				Severity:     "high",
				Confidence:   0.95,
				ImpactSummary: "JSON parsing consumes 45% of CPU time",
				Evidence: []model.EvidenceItem{
					{
						Type:        "function",
						Description: "encoding/json.Unmarshal",
						Value:       "45.2",
						Weight:      0.8,
					},
				},
			},
		},
	}

	// Create mock LLM provider with valid JSON response
	mockProvider := &MockLLMProvider{
		response: `{
			"version": "1.0",
			"generated_by": "triageprof-remediation",
			"timestamp": "2026-03-01T00:00:00Z",
			"remediations": [
				{
					"finding_id": "finding-001",
					"title": "Optimize JSON parsing",
					"description": "Replace encoding/json with faster JSON library",
					"code_changes": [
						{
							"file_path": "main.go",
							"line_number": 42,
							"current_code": "json.Unmarshal(data, &v)",
							"suggested_code": "fastjson.Unmarshal(data, &v)",
							"change_type": "modify",
							"explanation": "fastjson is 3-5x faster than standard library"
						}
					],
					"confidence": 0.85,
					"impact_estimate": "high",
					"tags": ["performance", "json"],
					"evidence_refs": ["finding-001"]
				}
			],
			"summary": {
				"total_remediations": 1,
				"high_impact": 1,
				"medium_impact": 0,
				"low_impact": 0,
				"estimated_total_gain": "high (>30% improvement)",
				"confidence_score": 0.85
			}
		}`,
	}

	// Create remediation config
	config := model.RemediationConfig{
		Enabled:           true,
		MinConfidence:     0.7,
		MaxCodeChanges:    3,
		CodeChangeLimit:   200,
		Provider:          "mistral",
		Model:             "mistral-large-latest",
		Temperature:       0.3,
	}

	// Create remediation generator
	generator := NewRemediationGenerator(mockProvider, config, findings, nil)

	// Generate remediations
	remediations, err := generator.GenerateRemediations(context.Background())

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, remediations)
	assert.Equal(t, 1, len(remediations.Remediations))
	assert.Equal(t, "finding-001", remediations.Remediations[0].FindingID)
	assert.Equal(t, 1, remediations.Summary.TotalRemediations)
	assert.Equal(t, "high (>30% improvement)", remediations.Summary.EstimatedTotalGain)
}

func TestRemediationGeneratorDisabled(t *testing.T) {
	// Create remediation config with disabled remediation
	config := model.RemediationConfig{
		Enabled: false,
	}

	// Create remediation generator
	generator := NewRemediationGenerator(nil, config, nil, nil)

	// Generate remediations should fail when disabled
	remediations, err := generator.GenerateRemediations(context.Background())

	// Verify error
	assert.Error(t, err)
	assert.Nil(t, remediations)
	assert.Contains(t, err.Error(), "remediation is disabled")
}

func TestRemediationGeneratorNoFindings(t *testing.T) {
	// Create remediation config
	config := model.RemediationConfig{
		Enabled: true,
	}

	// Create remediation generator with no findings
	generator := NewRemediationGenerator(nil, config, &model.FindingsBundle{}, nil)

	// Generate remediations should fail when no findings
	remediations, err := generator.GenerateRemediations(context.Background())

	// Verify error
	assert.Error(t, err)
	assert.Nil(t, remediations)
	assert.Contains(t, err.Error(), "no findings available")
}