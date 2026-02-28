package core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
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
	requestedProfiles := []string{"cpu", "heap", "mutex", "block", "goroutine"}
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
		Profiles:   []string{"cpu", "heap", "mutex", "block", "goroutine"},
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

func (p *Pipeline) Analyze(ctx context.Context, bundlePath string, topN int, outPath string) (*model.FindingsBundle, error) {
	// Read bundle
	data, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, err
	}

	var profileBundle model.ProfileBundle
	if err := json.Unmarshal(data, &profileBundle); err != nil {
		return nil, err
	}

	// Analyze
	findings, err := p.analyzer.Analyze(profileBundle, topN)
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

	// Report
	reportPath := filepath.Join(outDir, "report.md")
	return p.Report(ctx, findingsPath, reportPath)
}