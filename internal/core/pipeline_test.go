package core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateWebReport(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test findings data
	findings := model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
			TopIssueTags:  []string{"cpu", "memory"},
			Notes:         []string{"Test analysis"},
		},
		Findings: []model.Finding{
			{
				ID:            "test-finding-1",
				Title:         "High CPU Usage",
				Category:      "cpu",
				Severity:      "high",
				Confidence:    0.95,
				ImpactSummary: "CPU usage is high in main functions",
				Evidence: []model.EvidenceItem{
					{
						Type:        "profile",
						Description: "CPU profile shows high usage",
						Value:       "cpu.pb.gz",
						Weight:      0.8,
					},
				},
				DeterministicHints: []string{"Optimize hot loops", "Reduce allocations"},
				Tags:              []string{"performance", "cpu"},
			},
		},
	}

	// Create test insights data
	insights := &model.InsightsBundle{
		ExecutiveSummary: model.ExecutiveSummary{
			Overview:         "CPU bottlenecks detected",
			OverallSeverity:  "high",
			Confidence:       90,
			KeyThemes:        []string{"optimization", "efficiency"},
		},
		TopRisks: []model.RiskItem{
			{
				Description: "High CPU usage in main functions",
				Severity:    "high",
				Impact:      "significant",
				Likelihood:  "likely",
			},
		},
		TopActions: []model.ActionItem{
			{
				Description:      "Optimize CPU-intensive functions",
				Priority:         "high",
				EstimatedEffort:  "medium",
				Categories:       []string{"cpu", "performance"},
			},
		},
		PerFinding: []model.FindingInsight{
			{
				FindingID:        "test-finding-1",
				Narrative:        "The CPU usage is concentrated in a few key functions",
				LikelyRootCauses: []string{"Inefficient algorithms", "Excessive allocations"},
				Suggestions:      []string{"Profile specific functions", "Optimize data structures"},
				NextMeasurements: []string{"Run targeted benchmarks", "Analyze memory usage"},
				Caveats:          []string{"Results may vary by workload"},
				Confidence:       85,
			},
		},
	}

	// Create findings.json file
	findingsData, err := json.MarshalIndent(findings, "", "  ")
	require.NoError(t, err)
	
	findingsPath := filepath.Join(tempDir, "findings.json")
	err = os.WriteFile(findingsPath, findingsData, 0644)
	require.NoError(t, err)

	// Create web directory and copy necessary assets
	webDir := filepath.Join(tempDir, "web")
	err = os.MkdirAll(webDir, 0755)
	require.NoError(t, err)
	
	// Copy minimal required assets for testing
	// Create a minimal report-template.html
	reportTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>Test Report</title>
</head>
<body>
    <h1>Performance Report</h1>
    <div id="content"></div>
    <script>
        // Minimal script for testing
        const urlParams = new URLSearchParams(window.location.search);
        console.log("Findings param:", urlParams.get('findings'));
        console.log("Insights param:", urlParams.get('insights'));
    </script>
</body>
</html>`
	
	err = os.WriteFile(filepath.Join(webDir, "report-template.html"), []byte(reportTemplate), 0644)
	require.NoError(t, err)
	
	// Create minimal report.js
	reportJS := `console.log("Report JS loaded");`
	err = os.WriteFile(filepath.Join(webDir, "report.js"), []byte(reportJS), 0644)
	require.NoError(t, err)
	
	// Create minimal style.css
	styleCSS := `body { font-family: Arial, sans-serif; }`
	err = os.WriteFile(filepath.Join(webDir, "style.css"), []byte(styleCSS), 0644)
	require.NoError(t, err)

	// Get original working directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	
	// Change to project root directory where web assets are located
	err = os.Chdir("/home/doomguy/Documents/hackaton/Mistral-Hackathon")
	require.NoError(t, err)
	
	// Restore original directory after test
	defer os.Chdir(originalDir)

	// Create pipeline
	pipeline := &Pipeline{}

	// Test report generation
	err = pipeline.GenerateWebReport(context.Background(), findingsPath, tempDir, insights)
	require.NoError(t, err)

	// Verify files were created
	expectedFiles := []string{
		"report.html",
		"index.html",
		"web/report-template.html",
		"web/report.js",
		"web/style.css",
		"web/data/findings.json",
		"web/data/insights.json",
	}

	for _, expectedFile := range expectedFiles {
		fullPath := filepath.Join(tempDir, expectedFile)
		assert.FileExists(t, fullPath, "Expected file %s to exist", expectedFile)
	}

	// Verify report.html content contains expected elements
	reportHTML, err := os.ReadFile(filepath.Join(tempDir, "report.html"))
	require.NoError(t, err)
	
	reportContent := string(reportHTML)
	assert.Contains(t, reportContent, "TriageProf Performance Report")
	assert.Contains(t, reportContent, "report-template.html")
	assert.Contains(t, reportContent, "findings=")
	assert.Contains(t, reportContent, "insights=")

	// Verify findings.json was copied correctly
	copiedFindings, err := os.ReadFile(filepath.Join(tempDir, "web", "data", "findings.json"))
	require.NoError(t, err)
	
	var copiedFindingsData model.FindingsBundle
	err = json.Unmarshal(copiedFindings, &copiedFindingsData)
	require.NoError(t, err)
	assert.Equal(t, findings.Summary.OverallScore, copiedFindingsData.Summary.OverallScore)
	assert.Len(t, copiedFindingsData.Findings, len(findings.Findings))

	// Verify insights.json was copied correctly
	copiedInsights, err := os.ReadFile(filepath.Join(tempDir, "web", "data", "insights.json"))
	require.NoError(t, err)
	
	var copiedInsightsData model.InsightsBundle
	err = json.Unmarshal(copiedInsights, &copiedInsightsData)
	require.NoError(t, err)
	assert.Equal(t, insights.ExecutiveSummary.Overview, copiedInsightsData.ExecutiveSummary.Overview)
}

func TestGenerateWebReportWithoutInsights(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test findings data
	findings := model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 60,
			TopIssueTags:  []string{"memory"},
		},
		Findings: []model.Finding{
			{
				ID:       "test-finding-no-insights",
				Title:    "Memory Leak",
				Category: "memory",
				Severity: "medium",
			},
		},
	}

	// Create findings.json file
	findingsData, err := json.MarshalIndent(findings, "", "  ")
	require.NoError(t, err)
	
	findingsPath := filepath.Join(tempDir, "findings.json")
	err = os.WriteFile(findingsPath, findingsData, 0644)
	require.NoError(t, err)

	// Create web directory and copy necessary assets
	webDir := filepath.Join(tempDir, "web")
	err = os.MkdirAll(webDir, 0755)
	require.NoError(t, err)
	
	// Copy minimal required assets for testing
	// Create a minimal report-template.html
	reportTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>Test Report</title>
</head>
<body>
    <h1>Performance Report</h1>
    <div id="content"></div>
    <script>
        // Minimal script for testing
        const urlParams = new URLSearchParams(window.location.search);
        console.log("Findings param:", urlParams.get('findings'));
    </script>
</body>
</html>`
	
	err = os.WriteFile(filepath.Join(webDir, "report-template.html"), []byte(reportTemplate), 0644)
	require.NoError(t, err)
	
	// Create minimal report.js
	reportJS := `console.log("Report JS loaded");`
	err = os.WriteFile(filepath.Join(webDir, "report.js"), []byte(reportJS), 0644)
	require.NoError(t, err)
	
	// Create minimal style.css
	styleCSS := `body { font-family: Arial, sans-serif; }`
	err = os.WriteFile(filepath.Join(webDir, "style.css"), []byte(styleCSS), 0644)
	require.NoError(t, err)

	// Get original working directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	
	// Change to project root directory where web assets are located
	err = os.Chdir("/home/doomguy/Documents/hackaton/Mistral-Hackathon")
	require.NoError(t, err)
	
	// Restore original directory after test
	defer os.Chdir(originalDir)

	// Create pipeline
	pipeline := &Pipeline{}

	// Test report generation without insights
	err = pipeline.GenerateWebReport(context.Background(), findingsPath, tempDir, nil)
	require.NoError(t, err)

	// Verify files were created (should not have insights.json)
	assert.FileExists(t, filepath.Join(tempDir, "report.html"))
	assert.FileExists(t, filepath.Join(tempDir, "index.html"))
	assert.FileExists(t, filepath.Join(tempDir, "web", "data", "findings.json"))
	assert.NoFileExists(t, filepath.Join(tempDir, "web", "data", "insights.json"))

	// Verify report.html content
	reportHTML, err := os.ReadFile(filepath.Join(tempDir, "report.html"))
	require.NoError(t, err)
	
	reportContent := string(reportHTML)
	assert.Contains(t, reportContent, "TriageProf Performance Report")
	assert.Contains(t, reportContent, "report-template.html")
	assert.Contains(t, reportContent, "findings=")
	// Should not contain insights parameter when insights are nil
	assert.NotContains(t, reportContent, "insights=")
}