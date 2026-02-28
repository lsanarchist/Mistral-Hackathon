package core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/analyzer"
	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/mistral-hackathon/triageprof/internal/plugin"
	"github.com/mistral-hackathon/triageprof/internal/report"
)

type Pipeline struct {
	pluginManager *plugin.PluginManager
	analyzer      *analyzer.Analyzer
	reporter      *report.Reporter
}

func NewPipeline(pluginDir string) *Pipeline {
	return &Pipeline{
		pluginManager: plugin.NewPluginManager(pluginDir),
		analyzer:      analyzer.NewAnalyzer(),
		reporter:      report.NewReporter(),
	}
}



func (p *Pipeline) Collect(ctx context.Context, pluginName, targetURL string, durationSec, topN int, outDir string) (*model.ProfileBundle, error) {
	return p.CollectWithTarget(ctx, pluginName, targetURL, "", durationSec, topN, outDir)
}

func (p *Pipeline) CollectWithTarget(ctx context.Context, pluginName, targetURL, targetCommand string, durationSec, topN int, outDir string) (*model.ProfileBundle, error) {
	// Create output directory
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, err
	}

	// Resolve and validate plugin before launching
	manifest, _, err := p.pluginManager.ResolvePlugin(pluginName)
	if err != nil {
		return nil, err
	}

	// Determine target type and validate compatibility
	targetType := "url"
	if targetCommand != "" {
		targetType = "python"
	}
	
	if err := manifest.ValidateTarget(targetType); err != nil {
		return nil, err
	}

	// Determine profiles based on target type
	requestedProfiles := []string{"cpu", "heap", "mutex", "block", "goroutine", "allocs"}
	if targetType == "python" {
		requestedProfiles = []string{"cpu", "heap", "allocs", "memory-leak"}
	} else if targetType == "node" {
		requestedProfiles = []string{"cpu", "heap", "allocs"}
	}
	
	// Validate profile compatibility
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

	// Create target based on type
	target := model.Target{Type: targetType}
	if targetType == "url" {
		target.BaseURL = targetURL
	} else if targetType == "python" || targetType == "node" {
		// For Python and Node.js targets, parse the command string into a list
		// Simple shell-like parsing (basic space splitting, no complex shell features)
		cmdParts := strings.Fields(targetCommand)
		target.Command = cmdParts
	}

	// Validate target (plugin-side validation)
	if err := codec.Call("rpc.validateTarget", target, nil); err != nil {
		return nil, err
	}

	// Prepare collect request with appropriate profiles
	profiles := []string{"cpu", "heap", "mutex", "block", "goroutine", "allocs"}
	if targetType == "python" {
		profiles = []string{"cpu", "heap", "allocs", "memory-leak"}
	} else if targetType == "node" {
		profiles = []string{"cpu", "heap", "allocs"}
	}
	
	req := model.CollectRequest{
		Target:      target,
		DurationSec: durationSec,
		Profiles:    profiles,
		OutDir:      outDir,
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
	EnableCallgraph    bool
	CallgraphDepth     int
	EnableRegression   bool
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
		EnableCallgraph:    options.EnableCallgraph,
		CallgraphDepth:     options.CallgraphDepth,
		EnableRegression:   options.EnableRegression,
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
	return p.RunWithTarget(ctx, pluginName, targetURL, "", durationSec, topN, outDir)
}

func (p *Pipeline) RunWithTarget(ctx context.Context, pluginName, targetURL, targetCommand string, durationSec, topN int, outDir string) error {
	// Collect
	_, err := p.CollectWithTarget(ctx, pluginName, targetURL, targetCommand, durationSec, topN, outDir)
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

	// Report
	reportPath := filepath.Join(outDir, "report.md")
	return p.Report(ctx, findingsPath, reportPath)
}
