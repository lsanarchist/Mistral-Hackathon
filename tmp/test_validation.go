package main

import (
	"fmt"
	"github.com/mistral-hackathon/triageprof/internal/llm"
	"github.com/mistral-hackathon/triageprof/internal/model"
)

func main() {
	// Test finding reference detection
	text := "Analysis of finding find-001 shows issues"
	findingID := "find-001"
	result := llm.ContainsFindingReference(text, findingID)
	fmt.Printf("Finding reference test: %t\n", result)
	
	// Test validation
	findings := &model.FindingsBundle{
		Findings: []model.Finding{
			{ID: "find-001", Title: "Test", Category: "cpu"},
		},
	}
	
	insights := &model.InsightsBundle{
		ExecutiveSummary: model.ExecutiveSummary{
			Overview:         "Test analysis",
			OverallSeverity: "medium",
			Confidence:       85,
		},
		PerFinding: []model.FindingInsight{
			{
				FindingID:  "find-001",
				Narrative:  "Finding find-001 shows issues",
				Confidence: 90,
			},
		},
	}
	
	err := llm.ValidateInsights(insights, findings)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("Validation passed!")
	}
	
	// Test invalid finding reference
	invalidInsights := &model.InsightsBundle{
		ExecutiveSummary: model.ExecutiveSummary{
			Overview:         "Test analysis",
			OverallSeverity: "medium",
			Confidence:       85,
		},
		PerFinding: []model.FindingInsight{
			{
				FindingID:  "find-001",
				Narrative:  "High CPU usage detected", // No reference to find-001
				Confidence: 90,
			},
		},
	}
	
	err = llm.ValidateInsights(invalidInsights, findings)
	if err != nil {
		fmt.Printf("Expected validation error: %v\n", err)
	} else {
		fmt.Println("Unexpected: validation should have failed!")
	}
}