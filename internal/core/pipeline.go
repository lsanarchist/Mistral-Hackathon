package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/analyzer"
	"github.com/mistral-hackathon/triageprof/internal/llm"
	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/mistral-hackathon/triageprof/internal/plugin"
	"github.com/mistral-hackathon/triageprof/internal/report"
)

type Pipeline struct {
	pluginManager *plugin.PluginManager
	analyzer      *analyzer.Analyzer
	reporter      *report.Reporter
	llmGenerator  *llm.InsightsGenerator
}

func NewPipeline(pluginDir string) *Pipeline {
	return &Pipeline{
		pluginManager: plugin.NewPluginManager(pluginDir),
		analyzer:      analyzer.NewAnalyzer(),
		reporter:      report.NewReporter(),
		llmGenerator:  nil, // LLM is optional and configured separately
	}
}

// WithLLM configures LLM insights generation
func (p *Pipeline) WithLLM(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool) *Pipeline {
	p.llmGenerator = llm.NewInsightsGenerator(apiKey, model, timeout, maxResponse, maxPromptChars, dryRun)
	return p
}

func (p *Pipeline) Collect(ctx context.Context, pluginName, targetURL string, durationSec, topN int, outDir string) (*model.ProfileBundle, error) {
	// Create output directory
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, err
	}

	// Resolve and validate plugin before launching
	manifest, _, err := p.pluginManager.ResolvePlugin(pluginName)
	if err != nil {
		return nil, err
	}

	// Validate target type compatibility
	if err := manifest.ValidateTarget("url"); err != nil {
		return nil, err
	}

	// Validate profile compatibility
	requestedProfiles := []string{"cpu", "heap", "mutex", "block", "goroutine", "allocs"}
	if err := manifest.ValidateProfiles(requestedProfiles); err != nil {
		return nil, err
	}

	// Launch plugin
	codec, err := p.pluginManager.LaunchPlugin(pluginName, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer codec.Close()

	// Get plugin info
	var pluginInfo model.PluginInfo
	if err := codec.Call("rpc.info", nil, &pluginInfo); err != nil {
		return nil, err
	}

	// Validate target (plugin-side validation)
	target := model.Target{Type: "url", BaseURL: targetURL}
	if err := codec.Call("rpc.validateTarget", target, nil); err != nil {
		return nil, err
	}

	// Prepare collect request
	req := model.CollectRequest{
		Target:     target,
		DurationSec: durationSec,
		Profiles:   []string{"cpu", "heap", "mutex", "block", "goroutine", "allocs"},
		OutDir:     outDir,
		Metadata: map[string]string{
			"service":  "demo",
			"scenario": "default",
			"gitSha":   "",
		},
	}

	// Collect artifacts
	var bundle model.ArtifactBundle
	if err := codec.Call("rpc.collect", req, &bundle); err != nil {
		return nil, err
	}

	// Create profile bundle
	profileBundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Timestamp:   time.Now(),
			DurationSec: durationSec,
			Service:     "demo",
			Scenario:    "default",
			GitSha:      "",
		},
		Target: target,
		Plugin: model.PluginRef{
			Name:    pluginInfo.Name,
			Version: pluginInfo.Version,
		},
		Artifacts: bundle.Artifacts,
	}

	// Save bundle
	bundleData, err := json.MarshalIndent(profileBundle, "", "  ")
	if err != nil {
		return nil, err
	}

	bundlePath := filepath.Join(outDir, "bundle.json")
	if err := os.WriteFile(bundlePath, bundleData, 0644); err != nil {
		return nil, err
	}

	return profileBundle, nil
}

// CoreAnalyzeOptions configure analysis behavior
type CoreAnalyzeOptions struct {
	EnableCallgraph bool
	EnableRegression bool
	BaselineBundlePath string
}

func (p *Pipeline) Analyze(ctx context.Context, bundlePath string, topN int, outPath string) (*model.FindingsBundle, error) {
	return p.AnalyzeWithOptions(ctx, bundlePath, topN, outPath, CoreAnalyzeOptions{})
}

func (p *Pipeline) AnalyzeWithOptions(ctx context.Context, bundlePath string, topN int, outPath string, options CoreAnalyzeOptions) (*model.FindingsBundle, error) {
	// Read bundle
	data, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, err
	}

	var profileBundle model.ProfileBundle
	if err := json.Unmarshal(data, &profileBundle); err != nil {
		return nil, err
	}

	// Analyze with options
	analyzerOptions := analyzer.AnalyzeOptions{
		EnableCallgraph: options.EnableCallgraph,
		EnableRegression: options.EnableRegression,
		BaselineBundlePath: options.BaselineBundlePath,
	}
	findings, err := p.analyzer.AnalyzeWithOptions(profileBundle, topN, analyzerOptions)
	if err != nil {
		return nil, err
	}

	// Save findings
	findingsData, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(outPath, findingsData, 0644); err != nil {
		return nil, err
	}

	return findings, nil
}

func (p *Pipeline) Report(ctx context.Context, findingsPath, outPath string) error {
	// Read findings
	data, err := os.ReadFile(findingsPath)
	if err != nil {
		return err
	}

	var findings model.FindingsBundle
	if err := json.Unmarshal(data, &findings); err != nil {
		return err
	}

	// Generate report
	reportData, err := p.reporter.Generate(findings)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, []byte(reportData), 0644)
}

func (p *Pipeline) Run(ctx context.Context, pluginName, targetURL string, durationSec, topN int, outDir string) error {
	// Collect
	_, err := p.Collect(ctx, pluginName, targetURL, durationSec, topN, outDir)
	if err != nil {
		return err
	}

	// Analyze
	bundlePath := filepath.Join(outDir, "bundle.json")
	findingsPath := filepath.Join(outDir, "findings.json")
	_, err = p.Analyze(ctx, bundlePath, topN, findingsPath)
	if err != nil {
		return err
	}

	// Generate LLM insights (optional)
	var insights *model.InsightsBundle
	if p.llmGenerator != nil {
		insights, err = p.GenerateInsights(ctx, bundlePath, findingsPath)
		if err != nil {
			// LLM failure is non-fatal
			fmt.Printf("Warning: LLM insights generation failed: %v\n", err)
		}
	}

	// Report
	reportPath := filepath.Join(outDir, "report.md")
	return p.ReportWithInsights(ctx, findingsPath, insights, reportPath)
}

// GenerateInsights creates LLM insights from bundle and findings
func (p *Pipeline) GenerateInsights(ctx context.Context, bundlePath, findingsPath string) (*model.InsightsBundle, error) {
	if p.llmGenerator == nil {
		return nil, fmt.Errorf("LLM generator not configured")
	}

	// Load bundle
	bundleData, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read bundle: %w", err)
	}
	var bundle model.ProfileBundle
	if err := json.Unmarshal(bundleData, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse bundle: %w", err)
	}

	// Load findings
	findingsData, err := os.ReadFile(findingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read findings: %w", err)
	}
	var findings model.FindingsBundle
	if err := json.Unmarshal(findingsData, &findings); err != nil {
		return nil, fmt.Errorf("failed to parse findings: %w", err)
	}

	// Generate insights
	return p.llmGenerator.GenerateInsights(ctx, &bundle, &findings)
}

// ReportWithInsights generates a report with optional LLM insights
func (p *Pipeline) ReportWithInsights(ctx context.Context, findingsPath string, insights *model.InsightsBundle, outPath string) error {
	// Read findings
	data, err := os.ReadFile(findingsPath)
	if err != nil {
		return err
	}

	var findings model.FindingsBundle
	if err := json.Unmarshal(data, &findings); err != nil {
		return err
	}

	// Generate report with insights
	reportData, err := p.reporter.GenerateWithInsights(findings, insights)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, []byte(reportData), 0644)
}

// ReportJSONWithInsights generates a JSON report with optional LLM insights
func (p *Pipeline) ReportJSONWithInsights(ctx context.Context, findingsPath string, insights *model.InsightsBundle, outPath string, prettyPrint bool) error {
	// Read findings
	data, err := os.ReadFile(findingsPath)
	if err != nil {
		return err
	}

	var findings model.FindingsBundle
	if err := json.Unmarshal(data, &findings); err != nil {
		return err
	}

	// Generate JSON report with insights
	reportData, err := p.reporter.GenerateJSON(findings, insights, model.JSONReportOptions{
		IncludeInsights: insights != nil,
		PrettyPrint:     prettyPrint,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, reportData, 0644)
}