package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/analyzer"
	"github.com/mistral-hackathon/triageprof/internal/llm"
	"github.com/mistral-hackathon/triageprof/internal/auth"
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
	performanceGateConfig model.PerformanceGateConfig
	auditLogger    *AuditLogger
	rbacManager    *auth.RBACManager
	enterpriseConfig model.EnterpriseConfig
}

func NewPipeline(pluginDir string) *Pipeline {
	return &Pipeline{
		performanceGateConfig: model.DefaultPerformanceGateConfig(),
		pluginManager: plugin.NewPluginManager(pluginDir),
		analyzer:      analyzer.NewAnalyzer(),
		reporter:      report.NewReporter(),
		llmGenerator:  nil, // LLM disabled by default
		auditLogger:    NewAuditLogger(AuditLoggerConfig{Enabled: false}),
		rbacManager:    auth.NewRBACManager(),
		enterpriseConfig: model.EnterpriseConfig{},
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

// WithPerformanceGates configures performance gate settings
func (p *Pipeline) WithPerformanceGates(config model.PerformanceGateConfig) *Pipeline {
	p.performanceGateConfig = config
	return p
}

// WithEnterpriseConfig configures enterprise features
func (p *Pipeline) WithEnterpriseConfig(config model.EnterpriseConfig) *Pipeline {
	p.enterpriseConfig = config
	
	// Initialize audit logger if enterprise features are enabled
	if config.Enabled && config.AuditLogging {
		p.auditLogger = NewAuditLogger(AuditLoggerConfig{
			Enabled:    true,
			LogDir:     ".", // Default to current directory, will be updated when output dir is known
			MaxEntries: 1000,
		})
	} else {
		p.auditLogger = NewAuditLogger(AuditLoggerConfig{Enabled: false})
	}
	
	// Initialize RBAC if enabled
	if config.Enabled && config.RBACEnabled {
		p.rbacManager = auth.NewRBACManager()
	} else {
		p.rbacManager = auth.NewRBACManager() // Still initialize but with default roles
	}
	
	return p
}

// UpdateAuditLogDirectory updates the audit log directory after output directory is known
func (p *Pipeline) UpdateAuditLogDirectory(outDir string) {
	if p.auditLogger != nil && p.auditLogger.IsEnabled() {
		// Create a new audit logger with the correct output directory
		oldLogger := p.auditLogger
		p.auditLogger = NewAuditLogger(AuditLoggerConfig{
			Enabled:    true,
			LogDir:     outDir,
			MaxEntries: 1000,
		})
		// Copy any existing logs from the old logger
		if oldLogger != nil {
			for _, entry := range oldLogger.GetLogs() {
				p.auditLogger.LogAction(entry.UserID, entry.Action, entry.Resource, entry.Details, entry.Status)
			}
		}
	}
}

// LogAuditAction logs an audit action for enterprise auditing
func (p *Pipeline) LogAuditAction(userID, action, resource, details, status string) {
	if p.auditLogger != nil && p.auditLogger.IsEnabled() {
		p.auditLogger.LogAction(userID, action, resource, details, status)
	}
}

// GetAuditSummary returns audit log summary
func (p *Pipeline) GetAuditSummary() map[string]interface{} {
	if p.auditLogger != nil {
		return p.auditLogger.GetAuditSummary()
	}
	return map[string]interface{}{"enabled": false}
}

// GetRBACManager returns the RBAC manager
func (p *Pipeline) GetRBACManager() *auth.RBACManager {
	return p.rbacManager
}

// GetEnterpriseConfig returns the enterprise configuration
func (p *Pipeline) GetEnterpriseConfig() model.EnterpriseConfig {
	return p.enterpriseConfig
}

// CheckPerformanceGates checks findings against configured performance gates
func (p *Pipeline) CheckPerformanceGates(findings []model.Finding) (model.PerformanceGateResult, error) {
	if !p.performanceGateConfig.Enabled {
		return model.PerformanceGateResult{
			Passed: true,
			Message: "Performance gates disabled",
		}, nil
	}
	
	// Count findings by severity
	severityCounts := make(map[string]int)
	for _, finding := range findings {
		severityCounts[finding.Severity]++
	}
	
	// Check thresholds
	var warnings []string
	var errors []string
	
	if criticalCount, ok := severityCounts["critical"]; ok && criticalCount > p.performanceGateConfig.CriticalFindingsThreshold {
		if p.performanceGateConfig.FailOnCriticalThreshold {
			errors = append(errors, fmt.Sprintf("Critical findings threshold exceeded: %d > %d", criticalCount, p.performanceGateConfig.CriticalFindingsThreshold))
		} else {
			warnings = append(warnings, fmt.Sprintf("Critical findings threshold exceeded: %d > %d", criticalCount, p.performanceGateConfig.CriticalFindingsThreshold))
		}
	}
	
	if highCount, ok := severityCounts["high"]; ok && highCount > p.performanceGateConfig.HighFindingsThreshold {
		if p.performanceGateConfig.FailOnHighThreshold {
			errors = append(errors, fmt.Sprintf("High findings threshold exceeded: %d > %d", highCount, p.performanceGateConfig.HighFindingsThreshold))
		} else {
			warnings = append(warnings, fmt.Sprintf("High findings threshold exceeded: %d > %d", highCount, p.performanceGateConfig.HighFindingsThreshold))
		}
	}
	
	if mediumCount, ok := severityCounts["medium"]; ok && mediumCount > p.performanceGateConfig.MediumFindingsThreshold {
		if p.performanceGateConfig.WarnOnMediumThreshold {
			warnings = append(warnings, fmt.Sprintf("Medium findings threshold exceeded: %d > %d", mediumCount, p.performanceGateConfig.MediumFindingsThreshold))
		}
	}
	
	passed := len(errors) == 0
	message := "Performance gates passed"
	if len(errors) > 0 {
		message = fmt.Sprintf("Performance gates failed: %s", strings.Join(errors, "; "))
	} else if len(warnings) > 0 {
		message = fmt.Sprintf("Performance gates passed with warnings: %s", strings.Join(warnings, "; "))
	}
	
	return model.PerformanceGateResult{
		Passed:    passed,
		Message:   message,
		Warnings:  warnings,
		Errors:    errors,
		SeverityCounts: severityCounts,
	}, nil
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
	return p.AnalyzeWithDeterministicRulesAndOptions(ctx, bundlePath, topN, outPath, nil)
}

// AnalyzeWithDeterministicRulesAndOptions performs deterministic analysis with performance options
func (p *Pipeline) AnalyzeWithDeterministicRulesAndOptions(ctx context.Context, bundlePath string, topN int, outPath string, perfConfig *model.PerformanceOptimizationConfig) (*model.FindingsBundle, error) {
	// Read bundle
	data, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, err
	}

	var profileBundle model.ProfileBundle
	if err := json.Unmarshal(data, &profileBundle); err != nil {
		return nil, err
	}

	// Analyze with deterministic rules and performance options
	findings, err := p.analyzer.AnalyzeWithDeterministicRulesAndOptions(profileBundle, topN, perfConfig)
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

// AnalyzeWithBaselineComparison performs comparative analysis against a baseline
func (p *Pipeline) AnalyzeWithBaselineComparison(ctx context.Context, bundlePath string, topN int, outPath string, baselineComparison model.BaselineComparison) (*model.FindingsBundle, error) {
	// Read bundle
	data, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, err
	}

	var profileBundle model.ProfileBundle
	if err := json.Unmarshal(data, &profileBundle); err != nil {
		return nil, err
	}

	// Perform baseline comparison analysis
	findings, err := p.analyzer.AnalyzeWithBaselineComparison(profileBundle, topN, baselineComparison)
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

// AnalyzePerformanceTrends analyzes performance trends across multiple runs
func (p *Pipeline) AnalyzePerformanceTrends(ctx context.Context, currentBundlePath, baselineBundlePath, outPath string) ([]model.PerformanceTrend, error) {
	// Read current bundle
	currentData, err := os.ReadFile(currentBundlePath)
	if err != nil {
		return nil, err
	}

	var currentBundle model.ProfileBundle
	if err := json.Unmarshal(currentData, &currentBundle); err != nil {
		return nil, err
	}

	// Read baseline bundle
	baselineData, err := os.ReadFile(baselineBundlePath)
	if err != nil {
		return nil, err
	}

	var baselineBundle model.ProfileBundle
	if err := json.Unmarshal(baselineData, &baselineBundle); err != nil {
		return nil, err
	}

	// Analyze performance trends
	trends, err := p.analyzer.AnalyzePerformanceTrends(currentBundle, baselineBundle)
	if err != nil {
		return nil, err
	}

	// Save trends
	trendsData, err := json.MarshalIndent(trends, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(outPath, trendsData, 0644); err != nil {
		return nil, err
	}

	return trends, nil
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

	// Generate self-contained HTML report (all data embedded inline)
	htmlReporter := report.NewHTMLReporter()
	html, err := htmlReporter.Generate(findings, insights)
	if err != nil {
		return fmt.Errorf("generate HTML report: %w", err)
	}

	// Write report.html (self-contained, no redirect)
	reportPath := filepath.Join(outDir, "report.html")
	if err := os.WriteFile(reportPath, []byte(html), 0644); err != nil {
		return err
	}

	// Also write index.html for backward compatibility
	indexPath := filepath.Join(outDir, "index.html")
	if err := os.WriteFile(indexPath, []byte(html), 0644); err != nil {
		return err
	}

	// Keep data directory for raw JSON access
	webDir := filepath.Join(outDir, "web")
	if err := os.MkdirAll(webDir, 0755); err != nil {
		return err
	}
	dataDir := filepath.Join(webDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dataDir, "findings.json"), data, 0644); err != nil {
		return err
	}
	if insights != nil {
		insightsData, err := json.MarshalIndent(insights, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dataDir, "insights.json"), insightsData, 0644); err != nil {
			return err
		}
	}
	return nil
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

// GenerateRemediations creates automated code fix suggestions from findings and insights
func (p *Pipeline) GenerateRemediations(ctx context.Context, findingsPath, insightsPath, outPath string, config model.RemediationConfig) (*model.RemediationBundle, error) {
	// Read findings
	findingsData, err := os.ReadFile(findingsPath)
	if err != nil {
		return nil, err
	}

	var findingsBundle model.FindingsBundle
	if err := json.Unmarshal(findingsData, &findingsBundle); err != nil {
		return nil, err
	}

	// Read insights if available
	var insightsBundle *model.InsightsBundle
	if insightsPath != "" {
		if _, err := os.Stat(insightsPath); err == nil {
			// File exists, try to read it
			insightsData, err := os.ReadFile(insightsPath)
			if err != nil {
				log.Printf("Warning: failed to read insights file: %v", err)
			} else {
				if err := json.Unmarshal(insightsData, &insightsBundle); err != nil {
					log.Printf("Warning: failed to parse insights file: %v", err)
				}
			}
		}
		// If file doesn't exist or there were errors, insightsBundle remains nil
	}

	// Generate remediations
	if p.llmGenerator == nil {
		return nil, fmt.Errorf("LLM generator is not configured")
	}

	remediations, err := p.llmGenerator.GenerateRemediations(ctx, &findingsBundle, insightsBundle, config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate remediations: %w", err)
	}

	// Save remediations
	if remediations != nil {
		remediationsData, err := json.MarshalIndent(remediations, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(outPath, remediationsData, 0644); err != nil {
			return nil, err
		}
	}

	return remediations, nil
}
