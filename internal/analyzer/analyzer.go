package analyzer

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/pprof/profile"
	"github.com/mistral-hackathon/triageprof/internal/model"
)

type Analyzer struct {
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(bundle model.ProfileBundle, topN int) (*model.FindingsBundle, error) {
	findings := []model.Finding{}

	// Analyze each artifact
	for _, artifact := range bundle.Artifacts {
		if artifact.Kind != "pprof" {
			continue
		}

		// Read profile
		data, err := os.ReadFile(artifact.Path)
		if err != nil {
			continue
		}

		prof, err := profile.ParseData(data)
		if err != nil {
			continue
		}

		// Extract top functions
		topFuncs := extractTopFunctions(prof, topN)

		// Create finding
		finding := model.Finding{
			Category:  artifact.ProfileType,
			Title:     fmt.Sprintf("Top %s hotspots", artifact.ProfileType),
			Severity:  "medium",
			Score:     calculateScore(topFuncs),
			Top:       topFuncs,
			Evidence: model.Evidence{
				ArtifactPath: artifact.Path,
				ProfileType:  artifact.ProfileType,
				ExtractedAt:  time.Now(),
			},
		}

		findings = append(findings, finding)
	}

	// Create summary
	summary := model.Summary{
		TopIssueTags: []string{"performance"},
		OverallScore: 75,
		Notes:       []string{"Analysis completed successfully"},
	}

	return &model.FindingsBundle{
		Summary:  summary,
		Findings: findings,
	}, nil
}

func extractTopFunctions(prof *profile.Profile, topN int) []model.StackFrame {
	samples := []*profile.Sample{}
	for _, sample := range prof.Sample {
		samples = append(samples, sample)
	}

	// Sort by cumulative
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Value[0] > samples[j].Value[0]
	})

	frames := []model.StackFrame{}
	for i, sample := range samples {
		if i >= topN {
			break
		}

		for _, location := range sample.Location {
			for _, line := range location.Line {
				frame := model.StackFrame{
					Function: line.Function.Name,
					File:     line.Function.Filename,
					Line:     int(line.Line),
					Cum:      float64(sample.Value[0]),
					Flat:     float64(sample.Value[0]),
				}
				frames = append(frames, frame)
			}
		}
	}

	return frames
}

func calculateScore(frames []model.StackFrame) int {
	total := 0.0
	for _, frame := range frames {
		total += frame.Cum
	}

	if total > 1000 {
		return 90
	} else if total > 500 {
		return 70
	}
	return 50
}