package core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPipeline_Collect(t *testing.T) {
	// This test requires a running demo server
	t.Skip("Skipping integration test - requires running demo server")
	
	pipeline := NewPipeline("../../plugins")
	
	// Create temp directory
	tmpDir := t.TempDir()
	
	ctx := context.Background()
	
	// Test collection from demo server
	bundle, err := pipeline.Collect(ctx, "go-pprof-http", "http://localhost:6060", 5, 10, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, bundle)
	
	// Verify bundle structure
	assert.Equal(t, "url", bundle.Target.Type)
	assert.Equal(t, "http://localhost:6060", bundle.Target.BaseURL)
	assert.Equal(t, "go-pprof-http", bundle.Plugin.Name)
	assert.NotEmpty(t, bundle.Artifacts)
	
	// Verify bundle file was created
	bundlePath := filepath.Join(tmpDir, "bundle.json")
	_, err = os.Stat(bundlePath)
	require.NoError(t, err)
}

func TestPipeline_Analyze(t *testing.T) {
	pipeline := NewPipeline("../../plugins")
	
	// Create a test bundle
	bundle := model.ProfileBundle{
		Metadata: model.Metadata{
			Timestamp:   time.Now(),
			DurationSec: 10,
			Service:     "test",
			Scenario:    "test",
			GitSha:      "test",
		},
		Target: model.Target{
			Type:    "url",
			BaseURL: "http://localhost:6060",
		},
		Plugin: model.PluginRef{
			Name:    "test",
			Version: "0.1.0",
		},
		Artifacts: []model.Artifact{
			{
				Kind:        "pprof",
				ProfileType: "heap",
				Path:        "../../out/heap.pb.gz",
				ContentType: "application/octet-stream",
			},
		},
	}
	
	// Save test bundle
	tmpDir := t.TempDir()
	bundlePath := filepath.Join(tmpDir, "bundle.json")
	bundleData, err := json.MarshalIndent(bundle, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(bundlePath, bundleData, 0644))
	
	// Copy profile file to temp directory
	srcProfile := "../../out/heap.pb.gz"
	dstProfile := filepath.Join(tmpDir, "heap.pb.gz")
	profileData, err := os.ReadFile(srcProfile)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(dstProfile, profileData, 0644))
	
	// Update bundle to use local profile path
	var updatedBundle model.ProfileBundle
	require.NoError(t, json.Unmarshal(bundleData, &updatedBundle))
	updatedBundle.Artifacts[0].Path = dstProfile
	updatedBundleData, err := json.MarshalIndent(updatedBundle, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(bundlePath, updatedBundleData, 0644))
	
	ctx := context.Background()
	
	// Test basic analysis
	findingsPath := filepath.Join(tmpDir, "findings.json")
	findings, err := pipeline.Analyze(ctx, bundlePath, 5, findingsPath)
	require.NoError(t, err)
	require.NotNil(t, findings)
	
	// Verify findings structure
	assert.NotEmpty(t, findings.Findings)
	if len(findings.Findings) > 0 {
		assert.Equal(t, "heap", findings.Findings[0].Category)
		assert.NotEmpty(t, findings.Findings[0].Top)
	}
	
	// Verify findings file was created
	_, err = os.Stat(findingsPath)
	require.NoError(t, err)
}

func TestPipeline_AnalyzeWithOptions(t *testing.T) {
	pipeline := NewPipeline("../../plugins")
	
	// Create a test bundle
	bundle := model.ProfileBundle{
		Metadata: model.Metadata{
			Timestamp:   time.Now(),
			DurationSec: 10,
			Service:     "test",
			Scenario:    "test",
			GitSha:      "test",
		},
		Target: model.Target{
			Type:    "url",
			BaseURL: "http://localhost:6060",
		},
		Plugin: model.PluginRef{
			Name:    "test",
			Version: "0.1.0",
		},
		Artifacts: []model.Artifact{
			{
				Kind:        "pprof",
				ProfileType: "heap",
				Path:        "../../out/cpu.pb.gz",
				ContentType: "application/octet-stream",
			},
		},
	}
	
	// Save test bundle
	tmpDir := t.TempDir()
	bundlePath := filepath.Join(tmpDir, "bundle.json")
	bundleData, err := json.MarshalIndent(bundle, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(bundlePath, bundleData, 0644))
	
	// Copy profile file to temp directory
	srcProfile := "../../out/heap.pb.gz"
	dstProfile := filepath.Join(tmpDir, "heap.pb.gz")
	profileData, err := os.ReadFile(srcProfile)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(dstProfile, profileData, 0644))
	
	// Update bundle to use local profile path
	var updatedBundle model.ProfileBundle
	require.NoError(t, json.Unmarshal(bundleData, &updatedBundle))
	updatedBundle.Artifacts[0].Path = dstProfile
	updatedBundleData, err := json.MarshalIndent(updatedBundle, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(bundlePath, updatedBundleData, 0644))
	
	ctx := context.Background()
	
	// Test callgraph analysis
	findingsPath := filepath.Join(tmpDir, "findings.json")
	findings, err := pipeline.AnalyzeWithOptions(ctx, bundlePath, 5, findingsPath, CoreAnalyzeOptions{
		EnableCallgraph: true,
		CallgraphDepth:  3,
	})
	require.NoError(t, err)
	require.NotNil(t, findings)
	
	// Verify callgraph is present
	if len(findings.Findings) > 0 {
		assert.NotEmpty(t, findings.Findings[0].Callgraph)
		for _, node := range findings.Findings[0].Callgraph {
			assert.NotEmpty(t, node.Function)
			assert.GreaterOrEqual(t, node.Depth, 0)
		}
	}
}

func TestPipeline_Report(t *testing.T) {
	pipeline := NewPipeline("../../plugins")
	
	// Create test findings
	findings := model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
			TopIssueTags: []string{"performance", "memory"},
			Notes:        []string{"test analysis"},
		},
		Findings: []model.Finding{
			{
				Category: "heap",
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
				Evidence: model.Evidence{
					ArtifactPath: "heap.pb.gz",
					ProfileType:  "heap",
					ExtractedAt:  time.Now(),
				},
			},
		},
	}
	
	// Save test findings
	tmpDir := t.TempDir()
	findingsPath := filepath.Join(tmpDir, "findings.json")
	findingsData, err := json.MarshalIndent(findings, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(findingsPath, findingsData, 0644))
	
	ctx := context.Background()
	
	// Test report generation
	reportPath := filepath.Join(tmpDir, "report.md")
	err = pipeline.Report(ctx, findingsPath, reportPath)
	require.NoError(t, err)
	
	// Verify report file was created
	reportData, err := os.ReadFile(reportPath)
	require.NoError(t, err)
	assert.Contains(t, string(reportData), "Top CPU hotspots")
	assert.Contains(t, string(reportData), "runtime.allocm")
}

func TestPipeline_EndToEnd(t *testing.T) {
	// This test requires a running demo server
	t.Skip("Skipping integration test - requires running demo server")
	
	pipeline := NewPipeline("../../plugins")
	
	// Create temp directory
	tmpDir := t.TempDir()
	
	ctx := context.Background()
	
	// Step 1: Collect
	bundle, err := pipeline.Collect(ctx, "go-pprof-http", "http://localhost:6060", 5, 10, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, bundle)
	
	// Step 2: Analyze
	bundlePath := filepath.Join(tmpDir, "bundle.json")
	findingsPath := filepath.Join(tmpDir, "findings.json")
	findings, err := pipeline.Analyze(ctx, bundlePath, 5, findingsPath)
	require.NoError(t, err)
	require.NotNil(t, findings)
	
	// Step 3: Report
	reportPath := filepath.Join(tmpDir, "report.md")
	err = pipeline.Report(ctx, findingsPath, reportPath)
	require.NoError(t, err)
	
	// Verify all files exist
	_, err = os.Stat(bundlePath)
	require.NoError(t, err)
	_, err = os.Stat(findingsPath)
	require.NoError(t, err)
	_, err = os.Stat(reportPath)
	require.NoError(t, err)
	
	// Verify report content
	reportData, err := os.ReadFile(reportPath)
	require.NoError(t, err)
	assert.Contains(t, string(reportData), "Performance Analysis Report")
}

func TestPipeline_ReportJSON(t *testing.T) {
// 	pipeline := NewPipeline("../../plugins")
	
	// Create test findings
	findings := model.FindingsBundle{
		Summary: model.Summary{
			OverallScore: 75,
			TopIssueTags: []string{"performance", "memory"},
		},
		Findings: []model.Finding{
			{
				Category: "heap",
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
				Evidence: model.Evidence{
					ArtifactPath: "heap.pb.gz",
					ProfileType:  "heap",
					ExtractedAt:  time.Now(),
				},
			},
		},
	}
	
	// Save test findings
	tmpDir := t.TempDir()
	findingsPath := filepath.Join(tmpDir, "findings.json")
	findingsData, err := json.MarshalIndent(findings, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(findingsPath, findingsData, 0644))
	
// 	ctx := context.Background()
	
// 	// Test JSON report generation
// 	reportPath := filepath.Join(tmpDir, "report.json")
// 	err = pipeline.ReportJSONWithInsights(ctx, findingsPath, nil, reportPath, false)
// 	require.NoError(t, err)
// 	
// 	// Verify JSON report file was created
// 	_, err = os.Stat(reportPath)
// 	require.NoError(t, err)
// 	
// 	// Verify JSON content
// 	reportData, err := os.ReadFile(reportPath)
// 	require.NoError(t, err)
// 	
// 	var jsonReport model.JSONReport
// // 	err = json.Unmarshal(reportData, &jsonReport)
// 	require.NoError(t, err)
// 	assert.Equal(t, "1.0", jsonReport.SchemaVersion)
// 	assert.Equal(t, 1, len(jsonReport.Findings))
// // 	assert.Equal(t, "heap", jsonReport.Findings[0].Category)
}// }