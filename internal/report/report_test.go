package report

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJSON(t *testing.T) {
	reporter := NewReporter()

	// Create test findings
	findings := model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
			TopIssueTags: []string{"performance", "memory"},
			Notes:        []string{"test analysis"},
		},
		Findings: []model.Finding{
			{
				Category: "cpu",
				Title:    "Top CPU hotspots",
				Severity: "medium",
				Score:    80,
				Top: []model.StackFrame{
					{
						Function: "runtime.allocm",
						File:     "proc.go",
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

	// Test basic JSON generation
	reportData, err := reporter.GenerateJSON(findings, nil, model.JSONReportOptions{
		IncludeInsights: false,
		PrettyPrint:     false,
	})
	assert.NoError(t, err)

	// Verify it's valid JSON
	var jsonReport model.JSONReport
	err = json.Unmarshal(reportData, &jsonReport)
	assert.NoError(t, err)

	// Verify structure
	assert.Equal(t, "1.0", jsonReport.SchemaVersion)
	assert.Equal(t, 1, len(jsonReport.Findings))
	assert.Equal(t, "cpu", jsonReport.Findings[0].Category)
	assert.Equal(t, "Top CPU hotspots", jsonReport.Findings[0].Title)
	assert.Equal(t, "medium", jsonReport.Findings[0].Severity)
	assert.Equal(t, 80, jsonReport.Findings[0].Score)
	assert.Equal(t, 1, len(jsonReport.Findings[0].TopHotspots))

	// Test pretty print
	prettyData, err := reporter.GenerateJSON(findings, nil, model.JSONReportOptions{
		IncludeInsights: false,
		PrettyPrint:     true,
	})
	assert.NoError(t, err)
	assert.Contains(t, string(prettyData), "\n")
	assert.Contains(t, string(prettyData), "  ")
}

func TestDetermineSeverity(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{90, "critical"},
		{70, "high"},
		{50, "medium"},
		{30, "low"},
		{10, "info"},
		{0, "info"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("score_%d", tt.score), func(t *testing.T) {
			result := determineSeverity(tt.score)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJSONReportWithAllocationAnalysis(t *testing.T) {
	reporter := NewReporter()

	// Create test findings with allocation analysis
	findings := model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
			TopIssueTags: []string{"performance", "memory"},
		},
		Findings: []model.Finding{
			{
				Category: "allocs",
				Title:    "Top allocation hotspots",
				Severity: "high",
				Score:    85,
				AllocationAnalysis: &model.AllocationAnalysis{
					TotalAllocations: 10000.0,
					TopConcentration: 0.75,
					Severity:         "high",
					Score:            85,
					Hotspots: []model.AllocationHotspot{
						{
							Function: "runtime.allocm",
							File:     "proc.go",
							Line:     2276,
							Count:    5000.0,
							Percent:  50.0,
						},
						{
							Function: "main.allocateMemory",
							File:     "main.go",
							Line:     42,
							Count:    2500.0,
							Percent:  25.0,
						},
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

	// Test JSON generation with allocation analysis
	reportData, err := reporter.GenerateJSON(findings, nil, model.JSONReportOptions{
		IncludeInsights: false,
		PrettyPrint:     false,
	})
	assert.NoError(t, err)

	// Verify it's valid JSON
	var jsonReport model.JSONReport
	err = json.Unmarshal(reportData, &jsonReport)
	assert.NoError(t, err)

	// Verify allocation analysis is included
	assert.Equal(t, 1, len(jsonReport.Findings))
	assert.Equal(t, "allocs", jsonReport.Findings[0].Category)
	assert.NotNil(t, jsonReport.Findings[0].AllocationAnalysis)
	assert.Equal(t, 10000.0, jsonReport.Findings[0].AllocationAnalysis.TotalAllocations)
	assert.Equal(t, 0.75, jsonReport.Findings[0].AllocationAnalysis.TopConcentration)
	assert.Equal(t, 2, len(jsonReport.Findings[0].AllocationAnalysis.Hotspots))
	assert.Equal(t, "runtime.allocm", jsonReport.Findings[0].AllocationAnalysis.Hotspots[0].Function)
}

func TestCallgraphStatistics(t *testing.T) {
	reporter := NewReporter()

	// Create test callgraph
	callgraph := []model.CallgraphNode{
		{
			Function: "root",
			Depth:    0,
			Cum:      100.0,
			Flat:     50.0,
			Children: []model.CallgraphNode{
				{
					Function: "child1",
					Depth:    1,
					Cum:      60.0,
					Flat:     30.0,
					Children: []model.CallgraphNode{
						{
							Function: "grandchild1",
							Depth:    2,
							Cum:      40.0,
							Flat:     20.0,
						},
					},
				},
				{
					Function: "child2",
					Depth:    1,
					Cum:      40.0,
					Flat:     20.0,
				},
			},
		},
	}

	// Test node counting
	totalNodes := countCallgraphNodes(callgraph)
	assert.Equal(t, 4, totalNodes, "Should count all nodes including children")

	// Test max depth finding
	maxDepth := findMaxCallgraphDepth(callgraph)
	assert.Equal(t, 2, maxDepth, "Should find maximum depth")

	// Test JSON generation with callgraph
	findings := model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
			TopIssueTags: []string{"performance"},
		},
		Findings: []model.Finding{
			{
				Category:  "cpu",
				Title:     "Top CPU hotspots",
				Severity:  "medium",
				Score:     80,
				Callgraph: callgraph,
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

	// Generate JSON report
	reportData, err := reporter.GenerateJSON(findings, nil, model.JSONReportOptions{
		IncludeInsights: false,
		PrettyPrint:     false,
	})
	assert.NoError(t, err)

	// Verify callgraph is included in JSON
	var jsonReport model.JSONReport
	err = json.Unmarshal(reportData, &jsonReport)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jsonReport.Findings))
	assert.Equal(t, 1, len(jsonReport.Findings[0].Callgraph))
	assert.Equal(t, "root", jsonReport.Findings[0].Callgraph[0].Function)
}
