package analyzer

import (
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
		})
		require.NoError(t, err)
		require.NotNil(t, findings)
		assert.True(t, len(findings.Findings) > 0, "Should have at least one finding")
		if len(findings.Findings) > 0 {
			assert.NotEmpty(t, findings.Findings[0].Callgraph)
			assert.Nil(t, findings.Findings[0].Regression)
		}
	})

	// Test regression analysis - disabled for now due to path issues in tests
	// t.Run("regression analysis", func(t *testing.T) {
	//  ...
	// })
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



func TestMaxFunction(t *testing.T) {
	assert.Equal(t, 5, max(3, 5))
	assert.Equal(t, 5, max(5, 3))
	assert.Equal(t, 5, max(5, 5))
}