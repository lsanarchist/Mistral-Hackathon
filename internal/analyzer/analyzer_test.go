package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/pprof/profile"
	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineSeverity(t *testing.T) {
	tests := []struct {
		name     string
		frames   []model.StackFrame
		expected string
	}{
		{
			name: "critical",
			frames: []model.StackFrame{
				{Function: "test", Cum: 1500, Flat: 1500},
			},
			expected: "critical",
		},
		{
			name: "high",
			frames: []model.StackFrame{
				{Function: "test", Cum: 600, Flat: 600},
			},
			expected: "high",
		},
		{
			name: "medium",
			frames: []model.StackFrame{
				{Function: "test", Cum: 300, Flat: 300},
			},
			expected: "medium",
		},
		{
			name: "low",
			frames: []model.StackFrame{
				{Function: "test", Cum: 100, Flat: 100},
			},
			expected: "low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineSeverity(tt.frames)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		frames   []model.StackFrame
		expected int
	}{
		{
			name: "high score",
			frames: []model.StackFrame{
				{Function: "test", Cum: 1500, Flat: 1500},
			},
			expected: 90,
		},
		{
			name: "medium score",
			frames: []model.StackFrame{
				{Function: "test", Cum: 600, Flat: 600},
			},
			expected: 70,
		},
		{
			name: "low score",
			frames: []model.StackFrame{
				{Function: "test", Cum: 300, Flat: 300},
			},
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateScore(tt.frames)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyzeWithOptions(t *testing.T) {
	analyzer := NewAnalyzer()

	// Create a test bundle with a simple profile
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

	// Test basic analysis
	t.Run("basic analysis", func(t *testing.T) {
		findings, err := analyzer.AnalyzeWithOptions(bundle, 5, AnalyzeOptions{})
		require.NoError(t, err)
		require.NotNil(t, findings)
		assert.True(t, len(findings.Findings) > 0, "Should have at least one finding")
		if len(findings.Findings) > 0 {
			assert.Equal(t, "heap", findings.Findings[0].Category)
			assert.NotEmpty(t, findings.Findings[0].Top)
			assert.Nil(t, findings.Findings[0].Callgraph)
			assert.Nil(t, findings.Findings[0].Regression)
		}
	})

	// Test callgraph analysis
	t.Run("callgraph analysis", func(t *testing.T) {
		findings, err := analyzer.AnalyzeWithOptions(bundle, 5, AnalyzeOptions{
			EnableCallgraph: true,
			CallgraphDepth:  3,
		})
		require.NoError(t, err)
		require.NotNil(t, findings)
		assert.True(t, len(findings.Findings) > 0, "Should have at least one finding")
		if len(findings.Findings) > 0 {
			assert.NotEmpty(t, findings.Findings[0].Callgraph)
			assert.Nil(t, findings.Findings[0].Regression)
			// Verify callgraph has expected structure
			for _, node := range findings.Findings[0].Callgraph {
				assert.NotEmpty(t, node.Function)
				assert.GreaterOrEqual(t, node.Depth, 0)
				assert.Greater(t, node.Cum, 0.0)
			}
		}
	})

	// Test callgraph depth variation
	t.Run("callgraph depth variation", func(t *testing.T) {
		// Test depth 2
		findings2, err := analyzer.AnalyzeWithOptions(bundle, 5, AnalyzeOptions{
			EnableCallgraph: true,
			CallgraphDepth:  2,
		})
		require.NoError(t, err)
		
		// Test depth 4
		findings4, err := analyzer.AnalyzeWithOptions(bundle, 5, AnalyzeOptions{
			EnableCallgraph: true,
			CallgraphDepth:  4,
		})
		require.NoError(t, err)
		
		// Depth 4 should generally have more nodes than depth 2
		if len(findings2.Findings) > 0 && len(findings4.Findings) > 0 {
			nodes2 := countCallgraphNodes(findings2.Findings[0].Callgraph)
			nodes4 := countCallgraphNodes(findings4.Findings[0].Callgraph)
			assert.GreaterOrEqual(t, nodes4, nodes2, "Depth 4 should have at least as many nodes as depth 2")
		}
	})

	// Test regression analysis
	t.Run("regression analysis", func(t *testing.T) {
		// Create a baseline bundle
		baselineBundle := model.ProfileBundle{
			Metadata: model.Metadata{
				Timestamp:   time.Now(),
				DurationSec: 10,
				Service:     "test",
				Scenario:    "baseline",
				GitSha:      "baseline",
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

		// Create a current bundle
		currentBundle := model.ProfileBundle{
			Metadata: model.Metadata{
				Timestamp:   time.Now(),
				DurationSec: 10,
				Service:     "test",
				Scenario:    "current",
				GitSha:      "current",
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

		// Save baseline bundle to temp file
		baselineData, err := json.MarshalIndent(baselineBundle, "", "  ")
		require.NoError(t, err)
		baselinePath := filepath.Join(t.TempDir(), "baseline.json")
		require.NoError(t, os.WriteFile(baselinePath, baselineData, 0644))

		// Test regression analysis
		findings, err := analyzer.AnalyzeWithOptions(currentBundle, 5, AnalyzeOptions{
			EnableRegression:   true,
			BaselineBundlePath: baselinePath,
		})
		require.NoError(t, err)
		require.NotNil(t, findings)
		
		// Verify regression analysis results
		if len(findings.Findings) > 0 {
			finding := findings.Findings[0]
			assert.NotNil(t, finding.Regression, "Regression analysis should be present")
			assert.NotZero(t, finding.Regression.BaselineScore, "Baseline score should be calculated")
			assert.NotZero(t, finding.Regression.CurrentScore, "Current score should be calculated")
			assert.NotEmpty(t, finding.Regression.Severity, "Severity should be determined")
			assert.GreaterOrEqual(t, finding.Regression.Confidence, 50, "Confidence should be at least 50")
			assert.LessOrEqual(t, finding.Regression.Confidence, 100, "Confidence should be at most 100")
		}
	})
}

func TestCalculateProfileScore(t *testing.T) {
	// Create a test profile with concentrated hotspots
	prof := &profile.Profile{
		Sample: []*profile.Sample{
			{Value: []int64{1000}},
			{Value: []int64{500}},
			{Value: []int64{250}},
			{Value: []int64{100}},
			{Value: []int64{50}},
			{Value: []int64{25}},
			{Value: []int64{10}},
		},
	}

	score := calculateProfileScore(prof)
	assert.Greater(t, score, 50, "Expected high score for concentrated profile")
}

// countCallgraphNodes counts total nodes in callgraph (test helper)
func countCallgraphNodes(nodes []model.CallgraphNode) int {
	count := 0
	for _, node := range nodes {
		count += countCallgraphNode(&node)
	}
	return count
}

// countCallgraphNode recursively counts nodes (test helper)
func countCallgraphNode(node *model.CallgraphNode) int {
	count := 1
	for _, child := range node.Children {
		count += countCallgraphNode(&child)
	}
	return count
}

func TestMaxFunction(t *testing.T) {
	assert.Equal(t, 5, max(3, 5))
	assert.Equal(t, 5, max(5, 3))
	assert.Equal(t, 5, max(5, 5))
}

func TestAnalyzeAllocationPatterns(t *testing.T) {
	// Create a test profile with allocation data
	prof := &profile.Profile{
		Sample: []*profile.Sample{
			{Value: []int64{1000}, Location: []*profile.Location{
				{Line: []profile.Line{
					{Function: &profile.Function{Name: "allocFunction1", Filename: "file1.go"}, Line: 10},
				}},
			}},
			{Value: []int64{500}, Location: []*profile.Location{
				{Line: []profile.Line{
					{Function: &profile.Function{Name: "allocFunction2", Filename: "file2.go"}, Line: 20},
				}},
			}},
			{Value: []int64{250}, Location: []*profile.Location{
				{Line: []profile.Line{
					{Function: &profile.Function{Name: "allocFunction3", Filename: "file3.go"}, Line: 30},
				}},
			}},
			{Value: []int64{100}, Location: []*profile.Location{
				{Line: []profile.Line{
					{Function: &profile.Function{Name: "allocFunction4", Filename: "file4.go"}, Line: 40},
				}},
			}},
			{Value: []int64{50}, Location: []*profile.Location{
				{Line: []profile.Line{
					{Function: &profile.Function{Name: "allocFunction5", Filename: "file5.go"}, Line: 50},
				}},
			}},
		},
	}

	analysis := analyzeAllocationPatterns(prof)

	// Verify basic properties
	assert.Equal(t, 1900.0, analysis.TotalAllocations)
	assert.Greater(t, analysis.TopConcentration, 0.5)
	assert.NotEmpty(t, analysis.Severity)
	assert.Greater(t, analysis.Score, 0)
	assert.Equal(t, 5, len(analysis.Hotspots))

	// Verify hotspots are sorted by count
	for i := 0; i < len(analysis.Hotspots)-1; i++ {
		assert.GreaterOrEqual(t, analysis.Hotspots[i].Count, analysis.Hotspots[i+1].Count)
	}

	// Verify percentages sum correctly
	totalPercent := 0.0
	for _, hotspot := range analysis.Hotspots {
		totalPercent += hotspot.Percent
	}
	assert.Greater(t, totalPercent, 0.0)
}

func TestAllocationAnalysisIntegration(t *testing.T) {
	analyzer := NewAnalyzer()

	// Create a test bundle with allocation profile
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
				ProfileType: "allocs",
				Path:        "../../out/allocs.pb.gz",
				ContentType: "application/octet-stream",
			},
		},
	}

	// Test allocation analysis
	findings, err := analyzer.AnalyzeWithOptions(bundle, 5, AnalyzeOptions{})
	require.NoError(t, err)
	require.NotNil(t, findings)
	
	if len(findings.Findings) > 0 {
		finding := findings.Findings[0]
		assert.Equal(t, "allocs", finding.Category)
		assert.NotNil(t, finding.AllocationAnalysis, "Allocation analysis should be present for allocs profile")
		assert.NotEmpty(t, finding.AllocationAnalysis.Severity)
		assert.Greater(t, finding.AllocationAnalysis.Score, 0)
		assert.Greater(t, finding.AllocationAnalysis.TotalAllocations, 0.0)
	}
}
