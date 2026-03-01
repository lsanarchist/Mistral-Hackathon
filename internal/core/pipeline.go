package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
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
	alertsConfigFile string
	connectionQualityEnabled bool
	connectionQualityAlerts []webserver.ConnectionQualityAlert
	connectionQualityConfig webserver.ConnectionQualityConfig
	mlModelEnabled          bool
	advancedMLEnabled       bool
	phase4FeaturesEnabled  bool
	phase5FeaturesEnabled  bool
	phase6FeaturesEnabled  bool
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
func (p *Pipeline) WithLLM(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool) (*Pipeline, error) {
	var err error
	p.llmGenerator, err = llm.NewInsightsGenerator(apiKey, model, timeout, maxResponse, maxPromptChars, dryRun)
	return p, err
}

// WithPerformanceAlerts configures performance alerts from a JSON file
func (p *Pipeline) WithPerformanceAlerts(alertsFile string) *Pipeline {
	p.alertsConfigFile = alertsFile
	return p
}

// WithLLMWithProvider configures LLM insights generation with a specific provider
func (p *Pipeline) WithLLMWithProvider(config llm.ProviderConfig) (*Pipeline, error) {
	var err error
	p.llmGenerator, err = llm.NewInsightsGeneratorWithProvider(config)
	return p, err
}

// WithLLMWithCache configures LLM insights generation with caching
func (p *Pipeline) WithLLMWithCache(apiKey, model string, timeout, maxResponse, maxPromptChars int, dryRun bool, cacheConfig llm.CacheConfig) (*Pipeline, error) {
	config := llm.ProviderConfig{
		ProviderName: "mistral",
		APIKey:       apiKey,
		Model:        model,
		Timeout:      time.Duration(timeout) * time.Second,
		MaxResponse:  maxResponse,
		MaxPrompt:    maxPromptChars,
		DryRun:       dryRun,
	}
	
	generator, err := llm.NewInsightsGeneratorWithProvider(config)
	if err != nil {
		return p, err
	}
	
	// Note: Cache functionality would need to be implemented separately
	// For now, we'll use the basic generator
	p.llmGenerator = generator
	return p, nil
}

// WithWebSocketServer configures WebSocket server for real-time monitoring
func (p *Pipeline) WithWebSocketServer(port int, dataDir string, enableAuth bool, enableCompression bool, enableBatching bool, batchInterval time.Duration) {
	// Use the plugin directory from the pipeline's plugin manager
	pluginDir := "./plugins"
	if p.pluginManager != nil {
		pluginDir = p.pluginManager.PluginDir
	}
	
	// Load performance alerts if configured
	alertsConfig, err := webserver.LoadPerformanceAlertsFromFile(p.alertsConfigFile)
	if err != nil {
		log.Printf("Warning: Failed to load performance alerts: %v", err)
	}
	
	p.wsServer = webserver.NewWebSocketServer(port, dataDir, pluginDir, enableAuth, enableCompression, enableBatching, batchInterval, p.connectionQualityEnabled, alertsConfig, p.connectionQualityAlerts, p.connectionQualityConfig, p.mlModelEnabled, p.advancedMLEnabled, p.phase4FeaturesEnabled, p.phase5FeaturesEnabled, p.phase6FeaturesEnabled)
}

// WithWebSocketAutoRefresh configures auto-refresh interval for WebSocket server
func (p *Pipeline) WithWebSocketAutoRefresh(interval time.Duration) {
	if p.wsServer != nil {
		p.wsServer.StartAutoRefresh(interval)
	}
}

// WithWebSocketConnectionQuality enables WebSocket connection quality monitoring
func (p *Pipeline) WithWebSocketConnectionQuality(enabled bool) {
	p.connectionQualityEnabled = enabled
}

// WithWebSocketConnectionQualityAlerts configures connection quality alerts
func (p *Pipeline) WithWebSocketConnectionQualityAlerts(alerts []webserver.ConnectionQualityAlert) {
	p.connectionQualityAlerts = alerts
}

// WithWebSocketConnectionQualityConfig configures connection quality adaptation settings
func (p *Pipeline) WithWebSocketConnectionQualityConfig(config webserver.ConnectionQualityConfig) {
	p.connectionQualityConfig = config
}

// WithWebSocketMLModel enables ML-based anomaly detection for WebSocket connections
func (p *Pipeline) WithWebSocketMLModel(enabled bool) {
	p.mlModelEnabled = enabled
}

// WithWebSocketAdvancedML enables advanced ML features for WebSocket connections
func (p *Pipeline) WithWebSocketAdvancedML(enabled bool) {
	p.advancedMLEnabled = enabled
}

// WithWebSocketPhase4Features enables Phase 4 advanced features for WebSocket connections
func (p *Pipeline) WithWebSocketPhase4Features(enabled bool) {
	p.phase4FeaturesEnabled = enabled
}

// WithWebSocketPhase5Features enables Phase 5 advanced features for WebSocket connections
func (p *Pipeline) WithWebSocketPhase5Features(enabled bool) {
	p.phase5FeaturesEnabled = enabled
}

// WithWebSocketPhase6Features enables Phase 6 advanced features for WebSocket connections
func (p *Pipeline) WithWebSocketPhase6Features(enabled bool) {
	p.phase6FeaturesEnabled = enabled
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

// UpdateWebSocketData updates data in the WebSocket server and broadcasts to clients
func (p *Pipeline) UpdateWebSocketData(findings *model.FindingsBundle, insights *model.InsightsBundle) {
	if p.wsServer != nil {
		p.wsServer.UpdateData(findings, insights)
	}
}

// BroadcastWebSocketData sends current data to all WebSocket clients
func (p *Pipeline) BroadcastWebSocketData() {
	if p.wsServer != nil {
		p.wsServer.BroadcastData()
	}
}

// GetWebSocketClientCount returns the number of connected WebSocket clients
func (p *Pipeline) GetWebSocketClientCount() int {
	if p.wsServer == nil {
		return 0
	}
	return p.wsServer.GetClientCount()
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

// AnalyzeWithDeterministicRules performs deterministic analysis
func (p *Pipeline) AnalyzeWithDeterministicRules(ctx context.Context, bundlePath string, topN int, outPath string) (*model.FindingsBundle, error) {
	// Read bundle
	data, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, err
	}

	var profileBundle model.ProfileBundle
	if err := json.Unmarshal(data, &profileBundle); err != nil {
		return nil, err
	}

	// Analyze with deterministic rules
	findings, err := p.analyzer.AnalyzeWithDeterministicRules(profileBundle, topN)
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

	// Copy new web assets for professional report
	webAssets := []string{"report-template.html", "report.js", "style.css"}
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

	// Create report.html that loads the data with URL parameters
	findingsJSON := url.QueryEscape(string(data))
	var insightsJSON string
	if insights != nil {
		insightsData, _ := json.Marshal(insights)
		insightsJSON = url.QueryEscape(string(insightsData))
	}

	// Create the main report HTML file
	var reportHTML string
	if insights != nil {
		reportHTML = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>TriageProf Performance Report</title>
    <meta http-equiv="refresh" content="0; url=web/report-template.html?findings=%s&insights=%s">
</head>
<body>
    <p>Loading performance report...</p>
</body>
</html>`, findingsJSON, insightsJSON)
	} else {
		reportHTML = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>TriageProf Performance Report</title>
    <meta http-equiv="refresh" content="0; url=web/report-template.html?findings=%s">
</head>
<body>
    <p>Loading performance report...</p>
</body>
</html>`, findingsJSON)
	}

	// Create both report.html and index.html for compatibility
	reportPath := filepath.Join(outDir, "report.html")
	if err := os.WriteFile(reportPath, []byte(reportHTML), 0644); err != nil {
		return err
	}

	// Also create index.html for backward compatibility
	indexPath := filepath.Join(outDir, "index.html")
	return os.WriteFile(indexPath, []byte(reportHTML), 0644)
}

// GetCacheStats returns cache statistics if caching is enabled
func (p *Pipeline) GetCacheStats() (llm.CacheStats, bool) {
	if p.llmGenerator != nil && p.llmGenerator.Cache != nil {
		return p.llmGenerator.Cache.GetCacheStats(), true
	}
	return llm.CacheStats{}, false
}

// ClearLLMCache clears the LLM insights cache if caching is enabled
func (p *Pipeline) ClearLLMCache() error {
	if p.llmGenerator != nil && p.llmGenerator.Cache != nil {
		return p.llmGenerator.Cache.ClearCache()
	}
	return nil
}
