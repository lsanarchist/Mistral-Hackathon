package analyzer

import (
	"encoding/json"
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

// AnalyzeOptions configure analysis behavior
type AnalyzeOptions struct {
	EnableCallgraph    bool
	CallgraphDepth     int
	EnableRegression   bool
	BaselineBundlePath string
}

func (a *Analyzer) Analyze(bundle model.ProfileBundle, topN int) (*model.FindingsBundle, error) {
	return a.AnalyzeWithOptions(bundle, topN, AnalyzeOptions{})
}

func (a *Analyzer) AnalyzeWithOptions(bundle model.ProfileBundle, topN int, options AnalyzeOptions) (*model.FindingsBundle, error) {
	findings := []model.Finding{}

	// Load baseline for regression analysis if enabled
	var baselineBundle *model.ProfileBundle
	if options.EnableRegression && options.BaselineBundlePath != "" {
		baselineData, err := os.ReadFile(options.BaselineBundlePath)
		if err == nil {
			var tempBundle model.ProfileBundle
			if err := json.Unmarshal(baselineData, &tempBundle); err == nil {
				baselineBundle = &tempBundle
			}
		}
	}

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

		// Build callgraph if enabled
		var callgraph []model.CallgraphNode
		if options.EnableCallgraph {
			maxDepth := options.CallgraphDepth
			if maxDepth <= 0 {
				maxDepth = 3 // default depth
			}
			callgraph = buildCallgraph(prof, topN, maxDepth)
		}

		// Perform regression analysis if baseline available
		var regression *model.RegressionAnalysis
		if baselineBundle != nil && options.EnableRegression {
			regression = analyzeRegression(baselineBundle, &bundle, artifact.ProfileType)
		}

		// Create finding
		finding := model.Finding{
			Category:   artifact.ProfileType,
			Title:      fmt.Sprintf("Top %s hotspots", artifact.ProfileType),
			Severity:   determineSeverity(topFuncs),
			Score:      calculateScore(topFuncs),
			Top:        topFuncs,
			Callgraph:  callgraph,
			Regression: regression,
			Evidence: model.Evidence{
				ArtifactPath: artifact.Path,
				ProfileType:  artifact.ProfileType,
				ExtractedAt:  time.Now(),
			},
		}

		// Add allocation-specific analysis for allocs profiles
		if artifact.ProfileType == "allocs" {
			allocationAnalysis := analyzeAllocationPatterns(prof)
			finding.Severity = allocationAnalysis.Severity
			finding.Score = allocationAnalysis.Score
			finding.AllocationAnalysis = &allocationAnalysis
		}

		findings = append(findings, finding)
	}

	// Create summary
	summary := model.Summary{
		TopIssueTags: []string{"performance"},
		OverallScore: 75,
		Notes:        []string{"Analysis completed successfully"},
	}

	// Add analysis notes
	if options.EnableCallgraph {
		summary.Notes = append(summary.Notes, "Callgraph analysis enabled (depth 3)")
	}
	if options.EnableRegression && baselineBundle != nil {
		summary.Notes = append(summary.Notes, "Regression analysis performed against baseline")
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

func determineSeverity(frames []model.StackFrame) string {
	total := 0.0
	for _, frame := range frames {
		total += frame.Cum
	}

	if total > 1000 {
		return "critical"
	} else if total > 500 {
		return "high"
	} else if total > 200 {
		return "medium"
	}
	return "low"
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

// buildCallgraph constructs a callgraph tree from the profile data
func buildCallgraph(prof *profile.Profile, topN, maxDepth int) []model.CallgraphNode {
	// Group samples by root function
	sampleMap := make(map[string]*profile.Sample)
	for _, sample := range prof.Sample {
		if len(sample.Location) > 0 {
			rootFunc := sample.Location[0].Line[0].Function.Name
			if existing, exists := sampleMap[rootFunc]; !exists || sample.Value[0] > existing.Value[0] {
				sampleMap[rootFunc] = sample
			}
		}
	}

	// Convert to slice and sort by cumulative value
	samples := []*profile.Sample{}
	for _, sample := range sampleMap {
		samples = append(samples, sample)
	}
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Value[0] > samples[j].Value[0]
	})

	// Build callgraph nodes for top N functions
	nodes := []model.CallgraphNode{}
	for i, sample := range samples {
		if i >= topN {
			break
		}

		rootNode := buildCallgraphNode(sample, 0, maxDepth)
		if rootNode.Function != "" {
			nodes = append(nodes, rootNode)
		}
	}

	return nodes
}

// buildCallgraphNode recursively builds a callgraph node
func buildCallgraphNode(sample *profile.Sample, depth, maxDepth int) model.CallgraphNode {
	if depth >= maxDepth || len(sample.Location) == 0 {
		return model.CallgraphNode{}
	}

	location := sample.Location[0]
	if len(location.Line) == 0 {
		return model.CallgraphNode{}
	}

	line := location.Line[0]
	totalValue := float64(sample.Value[0])
	
	node := model.CallgraphNode{
		Function: line.Function.Name,
		File:     line.Function.Filename,
		Line:     int(line.Line),
		Depth:    depth,
		Cum:      totalValue,
		Flat:     totalValue,
	}

	// Build children from call stack
	if depth < maxDepth-1 && len(sample.Location) > 1 {
		for _, childLoc := range sample.Location[1:] {
			if len(childLoc.Line) > 0 {
				childSample := &profile.Sample{
					Value:    []int64{sample.Value[0]},
					Location: []*profile.Location{childLoc},
				}
				childNode := buildCallgraphNode(childSample, depth+1, maxDepth)
				if childNode.Function != "" {
					node.Children = append(node.Children, childNode)
				}
			}
		}
	}

	return node
}

// analyzeRegression compares current profile with baseline
func analyzeRegression(baseline, current *model.ProfileBundle, profileType string) *model.RegressionAnalysis {
	// Find matching artifacts
	var baselineArtifact, currentArtifact *model.Artifact
	for _, artifact := range baseline.Artifacts {
		if artifact.ProfileType == profileType {
			baselineArtifact = &artifact
			break
		}
	}
	for _, artifact := range current.Artifacts {
		if artifact.ProfileType == profileType {
			currentArtifact = &artifact
			break
		}
	}

	if baselineArtifact == nil || currentArtifact == nil {
		return nil
	}

	// Read and parse profiles
	baselineData, err := os.ReadFile(baselineArtifact.Path)
	if err != nil {
		return nil
	}
	currentData, err := os.ReadFile(currentArtifact.Path)
	if err != nil {
		return nil
	}

	baselineProf, err := profile.ParseData(baselineData)
	if err != nil {
		return nil
	}
	currentProf, err := profile.ParseData(currentData)
	if err != nil {
		return nil
	}

	// Calculate scores
	baselineScore := calculateProfileScore(baselineProf)
	currentScore := calculateProfileScore(currentProf)

	// Calculate regression metrics
	delta := currentScore - baselineScore
	percentage := 0.0
	if baselineScore > 0 {
		percentage = float64(delta) / float64(baselineScore) * 100
	}

	// Determine severity
	severity := "none"
	if delta > 20 {
		severity = "critical"
	} else if delta > 10 {
		severity = "high"
	} else if delta > 5 {
		severity = "medium"
	} else if delta < -10 {
		severity = "improved"
	}

	// Confidence based on sample count
	baselineSamples := len(baselineProf.Sample)
	currentSamples := len(currentProf.Sample)
	confidence := 50
	if baselineSamples > 100 && currentSamples > 100 {
		confidence = 80
	} else if baselineSamples > 50 && currentSamples > 50 {
		confidence = 60
	}

	return &model.RegressionAnalysis{
		BaselineScore: baselineScore,
		CurrentScore:  currentScore,
		Delta:         delta,
		Percentage:    percentage,
		Severity:      severity,
		Confidence:    confidence,
	}
}

// calculateProfileScore calculates a score for a profile based on hotspot concentration
func calculateProfileScore(prof *profile.Profile) int {
	if len(prof.Sample) == 0 {
		return 0
	}

	// Sort samples by cumulative value
	samples := prof.Sample
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Value[0] > samples[j].Value[0]
	})

	// Calculate top 5% concentration
	top5Percent := max(1, len(samples)/20)
	topTotal := 0.0
	for i := 0; i < top5Percent && i < len(samples); i++ {
		topTotal += float64(samples[i].Value[0])
	}

	total := 0.0
	for _, sample := range samples {
		total += float64(sample.Value[0])
	}

	if total == 0 {
		return 0
	}

	concentration := topTotal / total

	// Score based on concentration (higher = worse)
	if concentration > 0.8 {
		return 95
	} else if concentration > 0.6 {
		return 80
	} else if concentration > 0.4 {
		return 60
	} else if concentration > 0.2 {
		return 40
	}
	return 20
}

// analyzeAllocationPatterns performs allocation-specific analysis
func analyzeAllocationPatterns(prof *profile.Profile) model.AllocationAnalysis {
	if len(prof.Sample) == 0 {
		return model.AllocationAnalysis{}
	}

	// Sort samples by allocation count
	samples := prof.Sample
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Value[0] > samples[j].Value[0]
	})

	totalAllocations := 0.0
	for _, sample := range samples {
		totalAllocations += float64(sample.Value[0])
	}

	if totalAllocations == 0 {
		return model.AllocationAnalysis{}
	}

	// Calculate top 10% concentration
	top10Percent := max(1, len(samples)/10)
	topAllocations := 0.0
	for i := 0; i < top10Percent && i < len(samples); i++ {
		topAllocations += float64(samples[i].Value[0])
	}

	concentration := topAllocations / totalAllocations

	// Determine allocation severity
	severity := "low"
	if concentration > 0.7 {
		severity = "critical"
	} else if concentration > 0.5 {
		severity = "high"
	} else if concentration > 0.3 {
		severity = "medium"
	}

	// Calculate score based on concentration
	score := 0
	if concentration > 0.7 {
		score = 90
	} else if concentration > 0.5 {
		score = 70
	} else if concentration > 0.3 {
		score = 50
	} else {
		score = 30
	}

	// Identify top allocation hotspots
	hotspots := []model.AllocationHotspot{}
	for i, sample := range samples {
		if i >= 5 { // Top 5 hotspots
			break
		}
		
		if len(sample.Location) > 0 && len(sample.Location[0].Line) > 0 {
			line := sample.Location[0].Line[0]
			hotspots = append(hotspots, model.AllocationHotspot{
				Function: line.Function.Name,
				File:     line.Function.Filename,
				Line:     int(line.Line),
				Count:    float64(sample.Value[0]),
				Percent:  float64(sample.Value[0]) / totalAllocations * 100,
			})
		}
	}

	return model.AllocationAnalysis{
		TotalAllocations: totalAllocations,
		TopConcentration: concentration,
		Severity:         severity,
		Score:            score,
		Hotspots:         hotspots,
	}
}

// calculateAllocationScore calculates a specialized score for allocation profiles
func calculateAllocationScore(prof *profile.Profile) int {
	analysis := analyzeAllocationPatterns(prof)
	return analysis.Score
}

// determineAllocationSeverity determines severity specifically for allocation profiles
func determineAllocationSeverity(prof *profile.Profile) string {
	analysis := analyzeAllocationPatterns(prof)
	return analysis.Severity
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
