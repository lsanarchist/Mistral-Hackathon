package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/analyzer"
	"github.com/mistral-hackathon/triageprof/internal/llm"
	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/mistral-hackathon/triageprof/internal/plugin"
	"github.com/mistral-hackathon/triageprof/internal/report"
	"github.com/mistral-hackathon/triageprof/internal/webserver"
)

// containsString checks if a string exists in a slice
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type Pipeline struct {
	pluginManager *plugin.PluginManager
	analyzer      *analyzer.Analyzer
	reporter      *report.Reporter
	llmGenerator  *llm.InsightsGenerator
	wsServer      *webserver.WebSocketServer
}

func NewPipeline(pluginDir string) *Pipeline {
	return &Pipeline{
		pluginManager: plugin.NewPluginManager(pluginDir),
		analyzer:      analyzer.NewAnalyzer(),
		reporter:      report.NewReporter(),
		llmGenerator:  nil, // LLM disabled by default
	}
}

// WithLLM configures LLM insights generation
func (p *Pipeline) WithLLM(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool) *Pipeline {
	p.llmGenerator = llm.NewInsightsGenerator(apiKey, model, timeout, maxResponse, maxPromptChars, dryRun)
	return p
}

// WithWebSocketServer configures WebSocket server for real-time monitoring
func (p *Pipeline) WithWebSocketServer(port int, dataDir string) {
	p.wsServer = webserver.NewWebSocketServer(port, dataDir)
}

// StartWebSocketServer starts the WebSocket server
func (p *Pipeline) StartWebSocketServer() error {
	if p.wsServer == nil {
		return fmt.Errorf("WebSocket server not configured")
	}
	return p.wsServer.Start()
}

// StopWebSocketServer stops the WebSocket server
func (p *Pipeline) StopWebSocketServer() error {
	if p.wsServer == nil {
		return nil
	}
	return p.wsServer.Stop()
}

// LoadWebSocketData loads data into the WebSocket server
func (p *Pipeline) LoadWebSocketData(findingsPath, insightsPath string) error {
	if p.wsServer == nil {
		return fmt.Errorf("WebSocket server not configured")
	}
	return p.wsServer.LoadData(findingsPath, insightsPath)
}

// BroadcastWebSocketData sends current data to all WebSocket clients
func (p *Pipeline) BroadcastWebSocketData() {
	if p.wsServer != nil {
		p.wsServer.BroadcastData()
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
		// For command-based targets, we need to determine the type from the plugin's capabilities
		// Since the plugin manifest tells us what targets it supports, we'll use that
		// to determine the appropriate target type
		if containsString(manifest.Capabilities.Targets, "python") {
			targetType = "python"
		} else if containsString(manifest.Capabilities.Targets, "node") {
			targetType = "node"
		}
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
	} else if pluginName == "ruby-stackprof" {
		// Ruby stackprof plugin has different profile capabilities
		requestedProfiles = []string{"cpu", "memory", "object_allocation"}
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
	} else if pluginName == "ruby-stackprof" {
		// Ruby stackprof plugin has different profile capabilities
		profiles = []string{"cpu", "memory", "object_allocation"}
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
	return p.ReportWithInsights(ctx, findingsPath, outPath, nil)
}

func (p *Pipeline) ReportWithInsights(ctx context.Context, findingsPath, outPath string, insights *model.InsightsBundle) error {
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
	var reportData string
	var reportErr error
	if insights != nil {
		reportData, reportErr = p.reporter.GenerateWithInsights(findings, insights)
	} else {
		reportData, reportErr = p.reporter.Generate(findings)
	}
	if reportErr != nil {
		return reportErr
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

	// Generate LLM insights (optional)
	var insights *model.InsightsBundle
	if p.llmGenerator != nil {
		// Read bundle for LLM
		bundleData, err := os.ReadFile(bundlePath)
		if err != nil {
			return err
		}
		var profileBundle model.ProfileBundle
		if err := json.Unmarshal(bundleData, &profileBundle); err != nil {
			return err
		}

		// Read findings for LLM
		findingsData, err := os.ReadFile(findingsPath)
		if err != nil {
			return err
		}
		var findingsBundle model.FindingsBundle
		if err := json.Unmarshal(findingsData, &findingsBundle); err != nil {
			return err
		}

		// Generate insights
		insights, err = p.llmGenerator.GenerateInsights(ctx, &profileBundle, &findingsBundle)
		if err != nil {
			return err
		}

		// Save insights
		if insights != nil {
			insightsData, err := json.MarshalIndent(insights, "", "  ")
			if err != nil {
				return err
			}
			insightsPath := filepath.Join(outDir, "insights.json")
			if err := os.WriteFile(insightsPath, insightsData, 0644); err != nil {
				return err
			}
		}
	}

	// Report (with optional insights)
	reportPath := filepath.Join(outDir, "report.md")
	if insights != nil {
		if err := p.ReportWithInsights(ctx, findingsPath, reportPath, insights); err != nil {
			return err
		}
		// Generate web report when LLM insights are available
		if err := p.GenerateWebReport(ctx, findingsPath, outDir, insights); err != nil {
			return err
		}
	} else {
		if err := p.Report(ctx, findingsPath, reportPath); err != nil {
			return err
		}
	}
	return nil
}

// GenerateWebReport creates a web-based HTML report
func (p *Pipeline) GenerateWebReport(ctx context.Context, findingsPath, outDir string, insights *model.InsightsBundle) error {
	// Read findings
	data, err := os.ReadFile(findingsPath)
	if err != nil {
		return err
	}

	var findings model.FindingsBundle
	if err := json.Unmarshal(data, &findings); err != nil {
		return err
	}

	// Create web directory
	webDir := filepath.Join(outDir, "web")
	if err := os.MkdirAll(webDir, 0755); err != nil {
		return err
	}

	// Copy web assets
	webAssets := []string{"index.html", "style.css", "app.js"}
	for _, asset := range webAssets {
		srcPath := filepath.Join("web", asset)
		dstPath := filepath.Join(webDir, asset)
		
		// Read source file
		assetData, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}
		
		// Write to destination
		if err := os.WriteFile(dstPath, assetData, 0644); err != nil {
			return err
		}
	}

	// Create data directory for JSON files
	dataDir := filepath.Join(webDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// Copy findings.json to web/data/
	findingsDst := filepath.Join(dataDir, "findings.json")
	if err := os.WriteFile(findingsDst, data, 0644); err != nil {
		return err
	}

	// Copy insights.json if available
	if insights != nil {
		insightsData, err := json.MarshalIndent(insights, "", "  ")
		if err != nil {
			return err
		}
		insightsDst := filepath.Join(dataDir, "insights.json")
		if err := os.WriteFile(insightsDst, insightsData, 0644); err != nil {
			return err
		}
	}

	// Create index.html that loads the data automatically
	indexHTML := `<!DOCTYPE html>
<html>
<head>
    <title>TriageProf Web Report</title>
    <meta http-equiv="refresh" content="0; url=web/index.html">
</head>
<body>
    <p>Redirecting to web report...</p>
</body>
</html>`

	indexPath := filepath.Join(outDir, "index.html")
	return os.WriteFile(indexPath, []byte(indexHTML), 0644)
}
