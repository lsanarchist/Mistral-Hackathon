package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"os/signal"
	"syscall"

	"github.com/mistral-hackathon/triageprof/internal/core"
	"github.com/mistral-hackathon/triageprof/internal/llm"
	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/mistral-hackathon/triageprof/internal/plugin"
	"github.com/mistral-hackathon/triageprof/internal/webserver"
	"log"
)

// boolToStatus converts a boolean to ENABLED/DISABLED string
func boolToStatus(value bool) string {
	if value {
		return "ENABLED"
	}
	return "DISABLED"
}

// loadConnectionQualityAlerts loads connection quality alerts from a JSON file
func loadConnectionQualityAlerts(filePath string) ([]webserver.ConnectionQualityAlert, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read connection quality alerts file: %w", err)
	}

	var alerts []webserver.ConnectionQualityAlert
	if err := json.Unmarshal(data, &alerts); err != nil {
		return nil, fmt.Errorf("failed to parse connection quality alerts: %w", err)
	}

	return alerts, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: triageprof <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  plugins list")
		fmt.Println("  collect --plugin <name> --target-url <url> --duration <sec> --out <path>")
		fmt.Println("  analyze --in <bundle.json> --out <findings.json> --top <N> [--callgraph --callgraph-depth <depth>] [--regression --baseline <path> --baseline-ref <ref> --current-ref <ref>] [--trend-analysis --baseline <path> --trend-output <path>]")
		fmt.Println("  report --in <findings.json> --out <report.md|json> --output markdown|json")
		fmt.Println("  llm --bundle <bundle.json> --findings <findings.json> --out <insights.json> [--provider <provider>] [--model <model>] [--timeout <sec>] [--dry-run]")
		fmt.Println("  run --plugin <name> --target-url <url> --duration <sec> --outdir <dir> [--websocket-port <port>] [--websocket-auth] [--websocket-compression] [--websocket-batching] [--websocket-batch-interval <ms>]")
		fmt.Println("  run --plugin <name> --target-type python --target-command <cmd> --duration <sec> --outdir <dir> [--websocket-port <port>] [--websocket-auth] [--websocket-compression] [--websocket-batching] [--websocket-batch-interval <ms>]")
		fmt.Println("  run --plugin <name> --target-type node --target-command <cmd> --duration <sec> --outdir <dir> [--websocket-port <port>] [--websocket-auth] [--websocket-compression] [--websocket-batching] [--websocket-batch-interval <ms>]")
		fmt.Println("  web --in <findings.json> --outdir <dir> [--insights <insights.json>]")
		fmt.Println("  websocket --findings <findings.json> [--insights <insights.json>] [--port <port>] [--data-dir <dir>] [--compression] [--batching] [--batch-interval <ms>] [--phase5]")
		fmt.Println("  demo --repo <url> --out <dir> [--ref <branch/commit>] [--duration <sec>] [--concurrent] [--max-workers <N>] [--sampling-rate <rate>] [--memory-optimization] [--large-codebase]")
		fmt.Println("  demo-kit --out <dir> [--duration <sec>] [--concurrent] [--max-workers <N>] [--sampling-rate <rate>] [--memory-optimization] [--large-codebase] (uses built-in demo repository)")
		fmt.Println("\nLLM Options for 'run' command:")
		fmt.Println("  --llm (enable LLM insights)")
		fmt.Println("  --llm-provider <provider> (mistral, openai - default: mistral)")
		fmt.Println("  --llm-model <model> (default: provider-specific)")
		fmt.Println("  --llm-timeout <seconds> (default: 20)")
		fmt.Println("  --llm-max-chars <chars> (default: 12000)")
		fmt.Println("  --llm-dry-run (print prompt without API call)")
		os.Exit(1)
	}

	cmd := os.Args[1]

	// Determine plugin directory
	pluginDir := "./plugins"
	if envDir := os.Getenv("TRIAGEPROF_PLUGINS"); envDir != "" {
		pluginDir = envDir
	}

	pipeline := core.NewPipeline(pluginDir)

	switch cmd {
	case "plugins":
		runPluginsCommand()
	case "collect":
		runCollectCommand(pipeline)
	case "analyze":
		runAnalyzeCommand(pipeline)
	case "report":
		runReportCommand(pipeline)
		case "llm":
			runLLMCommand()
	case "run":
		runRunCommand(pipeline)
	case "web":
		runWebCommand(pipeline)
	case "websocket":
		runWebSocketCommand(pipeline)
	case "demo":
		runDemoCommand(pipeline)
	case "demo-kit":
		runDemoKitCommand(pipeline)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func runPluginsCommand() {
	// Determine plugin directory
	pluginDir := "./plugins"
	if envDir := os.Getenv("TRIAGEPROF_PLUGINS"); envDir != "" {
		pluginDir = envDir
	}

	pm := plugin.NewPluginManager(pluginDir)
	manifests, err := pm.ListPlugins()
	if err != nil {
		fmt.Printf("Error listing plugins: %v\n", err)
		os.Exit(1)
	}

	if len(manifests) == 0 {
		fmt.Println("No plugins found.")
		return
	}

	fmt.Println("Available plugins:")
	for _, m := range manifests {
		fmt.Printf("  %s v%s (sdk %s)\n", m.Name, m.Version, m.SDKVersion)
		if m.Description != "" {
			fmt.Printf("    %s\n", m.Description)
		}
		fmt.Printf("    targets: %s\n", strings.Join(m.Capabilities.Targets, ", "))
		fmt.Printf("    profiles: %s\n", strings.Join(m.Capabilities.Profiles, ", "))
		if m.Author != "" {
			fmt.Printf("    author: %s\n", m.Author)
		}
		fmt.Println()
	}
}

func runCollectCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("collect", flag.ExitOnError)
	pluginName := flagSet.String("plugin", "", "Plugin name")
	targetURL := flagSet.String("target-url", "", "Target URL (for URL-based targets)")
	targetType := flagSet.String("target-type", "", "Target type (url/python)")
	targetCommand := flagSet.String("target-command", "", "Target command (for Python targets)")
	duration := flagSet.Int("duration", 15, "Duration in seconds")
	outPath := flagSet.String("out", "", "Output bundle path")
	flagSet.Parse(os.Args[2:])

	if *pluginName == "" || *outPath == "" {
		fmt.Println("Required flags: --plugin, --out")
		fmt.Println("For URL targets: --target-url")
		fmt.Println("For Python targets: --target-type python --target-command")
		fmt.Println("For Node.js targets: --target-type node --target-command")
		os.Exit(1)
	}

	// Validate target parameters
	if *targetType == "python" || *targetType == "node" {
		if *targetCommand == "" {
			fmt.Printf("%s target requires --target-command\n", strings.Title(*targetType))
			os.Exit(1)
		}
	} else if *targetURL == "" {
		// Default to URL target if target-type not specified
		*targetType = "url"
		if *targetURL == "" {
			fmt.Println("URL target requires --target-url")
			os.Exit(1)
		}
	}

	ctx := context.Background()
	var err error
	if *targetType == "python" || *targetType == "node" {
		_, err = pipeline.CollectWithTarget(ctx, *pluginName, "", *targetCommand, *duration, 20, filepath.Dir(*outPath))
	} else {
		_, err = pipeline.Collect(ctx, *pluginName, *targetURL, *duration, 20, filepath.Dir(*outPath))
	}
	if err != nil {
		fmt.Printf("Collect failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Bundle saved to: %s\n", *outPath)
}

func runAnalyzeCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("analyze", flag.ExitOnError)
	inPath := flagSet.String("in", "", "Input bundle path")
	outPath := flagSet.String("out", "", "Output findings path")
	topN := flagSet.Int("top", 20, "Top N findings")
	callgraph := flagSet.Bool("callgraph", false, "Enable callgraph analysis")
	callgraphDepth := flagSet.Int("callgraph-depth", 3, "Callgraph maximum depth (default: 3)")
	regression := flagSet.Bool("regression", false, "Enable regression analysis")
	baseline := flagSet.String("baseline", "", "Baseline bundle path for regression analysis")
	baselineRef := flagSet.String("baseline-ref", "", "Baseline reference (commit/branch/tag)")
	currentRef := flagSet.String("current-ref", "", "Current reference (commit/branch/tag)")
	trendAnalysis := flagSet.Bool("trend-analysis", false, "Enable performance trend analysis")
	trendOutput := flagSet.String("trend-output", "", "Output path for trend analysis results")
	flagSet.Parse(os.Args[2:])

	if *inPath == "" || *outPath == "" {
		fmt.Println("Required flags: --in, --out")
		os.Exit(1)
	}

	if *regression && *baseline == "" {
		fmt.Println("Regression analysis requires --baseline flag")
		os.Exit(1)
	}

	if *trendAnalysis && (*baseline == "" || *trendOutput == "") {
		fmt.Println("Trend analysis requires --baseline and --trend-output flags")
		os.Exit(1)
	}

	ctx := context.Background()
	
	// Handle trend analysis first if requested
	if *trendAnalysis {
		trends, err := pipeline.AnalyzePerformanceTrends(ctx, *inPath, *baseline, *trendOutput)
		if err != nil {
			fmt.Printf("Trend analysis failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✓ Performance trend analysis completed\n")
		fmt.Printf("📊 Found %d performance trends\n", len(trends))
		fmt.Printf("📈 Trend analysis saved to: %s\n", *trendOutput)
		
		// Print summary of trends
		criticalCount := 0
		highCount := 0
		improvedCount := 0
		
		for _, trend := range trends {
			switch trend.Severity {
			case "critical":
				criticalCount++
			case "high":
				highCount++
			case "improved":
				improvedCount++
			}
		}
		
		if criticalCount > 0 {
			fmt.Printf("🔴 %d critical regressions detected\n", criticalCount)
		}
		if highCount > 0 {
			fmt.Printf("🟠 %d high severity regressions detected\n", highCount)
		}
		if improvedCount > 0 {
			fmt.Printf("📈 %d performance improvements detected\n", improvedCount)
		}
	}

	// Handle regular analysis with baseline comparison
	if *regression {
		// Create baseline comparison configuration
		baselineComparison := model.BaselineComparison{
			BaselinePath: *baseline,
			CurrentPath:  *inPath,
			BaselineRef:   *baselineRef,
			CurrentRef:    *currentRef,
			Threshold:     10.0, // default 10% threshold
		}
		
		findings, err := pipeline.AnalyzeWithBaselineComparison(ctx, *inPath, *topN, *outPath, baselineComparison)
		if err != nil {
			fmt.Printf("Baseline comparison failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Findings saved to: %s\n", *outPath)
		fmt.Println("✓ Baseline comparison analysis completed")
		
		// Count regression findings
		regressionCount := 0
		for _, finding := range findings.Findings {
			if finding.Regression != nil && finding.Regression.Severity != "none" && finding.Regression.Severity != "low" && finding.Regression.Severity != "improved" {
				regressionCount++
			}
		}
		
		if regressionCount > 0 {
			fmt.Printf("🔴 %d findings with potential regressions detected\n", regressionCount)
		} else {
			fmt.Println("📈 No significant regressions detected")
		}
		
		if *callgraph {
			fmt.Printf("✓ Callgraph analysis completed (depth %d)\n", *callgraphDepth)
		}
	} else {
		// Regular analysis without regression
		options := core.CoreAnalyzeOptions{
			EnableCallgraph:    *callgraph,
			CallgraphDepth:     *callgraphDepth,
			EnableRegression:   false,
			BaselineBundlePath: "",
		}

		_, err := pipeline.AnalyzeWithOptions(ctx, *inPath, *topN, *outPath, options)
		if err != nil {
			fmt.Printf("Analyze failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Findings saved to: %s\n", *outPath)
		if *callgraph {
			fmt.Printf("✓ Callgraph analysis completed (depth %d)\n", *callgraphDepth)
		}
	}
}



func runReportCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("report", flag.ExitOnError)
	inPath := flagSet.String("in", "", "Input findings path")
	outPath := flagSet.String("out", "", "Output report path")
	// insightsPath := flagSet.String("insights", "", "Optional insights path")
	outputFormat := flagSet.String("output", "markdown", "Output format: markdown or json")
	// prettyPrint := flagSet.Bool("pretty", false, "Pretty print JSON output")
	flagSet.Parse(os.Args[2:])

	if *inPath == "" || *outPath == "" {
		fmt.Println("Required flags: --in, --out")
		os.Exit(1)
	}

	ctx := context.Background()

	// Load insights if provided
	// var insights *model.InsightsBundle
	// if *insightsPath != "" {
	// 	data, err := os.ReadFile(*insightsPath)
	// 	if err != nil {
	// 		fmt.Printf("Warning: failed to read insights: %v\n", err)
	// 	} else {
	// 		var ib model.InsightsBundle
	// 		if err := json.Unmarshal(data, &ib); err != nil {
	// 			fmt.Printf("Warning: failed to parse insights: %v\n", err)
	// 		} else {
	// 			insights = &ib
	// 		}
	// 	}
	// }

	// Generate report based on format
	switch *outputFormat {
	case "json":
		// err := pipeline.ReportJSONWithInsights(ctx, *inPath, insights, *outPath, *prettyPrint)
		// if err != nil {
		// 	fmt.Printf("JSON report failed: %v\n", err)
		// 	os.Exit(1)
		// }
		fmt.Printf("JSON report functionality temporarily disabled\n")
		os.Exit(1)
	case "markdown":
		err := pipeline.Report(ctx, *inPath, *outPath)
		if err != nil {
			fmt.Printf("Markdown report failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Markdown report saved to: %s\n", *outPath)
	default:
		fmt.Printf("Unknown output format: %s. Using markdown.\n", *outputFormat)
		err := pipeline.Report(ctx, *inPath, *outPath)
		if err != nil {
			fmt.Printf("Markdown report failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Markdown report saved to: %s\n", *outPath)
	}
}

func runRunCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("run", flag.ExitOnError)
	pluginName := flagSet.String("plugin", "", "Plugin name")
	targetURL := flagSet.String("target-url", "", "Target URL (for URL-based targets)")
	targetType := flagSet.String("target-type", "", "Target type (url/python)")
	targetCommand := flagSet.String("target-command", "", "Target command (for Python targets)")
	duration := flagSet.Int("duration", 15, "Duration in seconds")
	outDir := flagSet.String("outdir", "", "Output directory")
	llmEnabled := flagSet.Bool("llm", false, "Enable LLM insights")
	llmProvider := flagSet.String("llm-provider", "mistral", "LLM provider (mistral, openai)")
	llmModel := flagSet.String("llm-model", "", "Model name (provider-specific default if empty)")
	llmTimeout := flagSet.Int("llm-timeout", 20, "LLM API timeout in seconds")
	llmMaxChars := flagSet.Int("llm-max-chars", 12000, "Max prompt characters")
	llmDryRun := flagSet.Bool("llm-dry-run", false, "Dry run - save prompt without API call")
	websocketPort := flagSet.Int("websocket-port", 0, "WebSocket server port (0 to disable)")
	websocketAuth := flagSet.Bool("websocket-auth", false, "Enable WebSocket authentication")
	websocketCompression := flagSet.Bool("websocket-compression", false, "Enable WebSocket message compression")
	websocketBatching := flagSet.Bool("websocket-batching", false, "Enable WebSocket message batching")
	websocketBatchInterval := flagSet.Int("websocket-batch-interval", 100, "WebSocket batch interval in milliseconds")
	websocketConnectionQuality := flagSet.Bool("websocket-connection-quality", false, "Enable WebSocket connection quality monitoring")
	websocketQualityAlerts := flagSet.String("websocket-quality-alerts", "", "WebSocket connection quality alert configuration file (JSON)")
	websocketAdaptiveUpdates := flagSet.Bool("websocket-adaptive-updates", true, "Enable adaptive updates based on connection quality")
	websocketBandwidthThrottling := flagSet.Bool("websocket-bandwidth-throttling", true, "Enable bandwidth throttling based on connection quality")
	websocketMLModel := flagSet.Bool("websocket-ml-model", false, "Enable ML-based anomaly detection for WebSocket connections")
	websocketAdvancedML := flagSet.Bool("websocket-advanced-ml", false, "Enable advanced ML features (root cause analysis, predictions, correlations)")
	websocketPhase4Features := flagSet.Bool("websocket-phase4-features", false, "Enable Phase 4 advanced ML features (deep learning, time series forecasting, automated root cause analysis)")
	websocketPhase5Features := flagSet.Bool("websocket-phase5-features", false, "Enable Phase 5 advanced ML features (anomaly correlation detection, predictive maintenance, enhanced root cause analysis)")
	performanceAlerts := flagSet.String("performance-alerts", "", "Performance alert configuration file (JSON)")
	flagSet.Parse(os.Args[2:])

	if *pluginName == "" || *outDir == "" {
		fmt.Println("Required flags: --plugin, --outdir")
		fmt.Println("For URL targets: --target-url")
		fmt.Println("For Python targets: --target-type python --target-command")
		os.Exit(1)
	}

	// Validate target parameters
	if *targetType == "python" || *targetType == "node" {
		if *targetCommand == "" {
			fmt.Printf("%s target requires --target-command\n", strings.Title(*targetType))
			os.Exit(1)
		}
	} else if *targetURL == "" {
		// Default to URL target if target-type not specified
		*targetType = "url"
		if *targetURL == "" {
			fmt.Println("URL target requires --target-url")
			os.Exit(1)
		}
	}

	ctx := context.Background()

	// Configure LLM if enabled
	if *llmEnabled {
		// Get API key based on provider
		apiKey := ""
		apiKeyEnv := "MISTRAL_API_KEY"
		if *llmProvider == "openai" {
			apiKeyEnv = "OPENAI_API_KEY"
		}
		apiKey = os.Getenv(apiKeyEnv)
		
		// Set default model if not specified
		if *llmModel == "" {
			if *llmProvider == "openai" {
				defaultModel := "gpt-3.5-turbo"
				llmModel = &defaultModel
			} else {
				defaultModel := "mistral-large-latest"
				llmModel = &defaultModel
			}
		}
		
		// Create provider config
		config := llm.ProviderConfig{
			ProviderName: *llmProvider,
			Model:        *llmModel,
			APIKey:       apiKey,
			Timeout:      time.Duration(*llmTimeout) * time.Second,
			MaxResponse:  8192,
			MaxPrompt:    *llmMaxChars,
			DryRun:       *llmDryRun,
		}
		
		_, err := pipeline.WithLLMWithProvider(config)
		if err != nil {
			fmt.Printf("Failed to configure LLM: %v\n", err)
			os.Exit(1)
		}
	}

	// Configure WebSocket server if port is specified
	if *websocketPort > 0 {
		batchInterval := time.Duration(*websocketBatchInterval) * time.Millisecond
		pipeline.WithWebSocketServer(*websocketPort, *outDir, *websocketAuth, *websocketCompression, *websocketBatching, batchInterval)
		pipeline.WithWebSocketConnectionQuality(*websocketConnectionQuality)
		
		// Configure connection quality enhancements
		if *websocketConnectionQuality {
			// Load quality alerts if specified
			if *websocketQualityAlerts != "" {
				qualityAlerts, err := loadConnectionQualityAlerts(*websocketQualityAlerts)
				if err != nil {
					log.Printf("Warning: Failed to load connection quality alerts: %v", err)
				} else {
					pipeline.WithWebSocketConnectionQualityAlerts(qualityAlerts)
				}
			}
			
			// Configure quality-based adaptations
			qualityConfig := webserver.ConnectionQualityConfig{
				AdaptiveUpdatesEnabled: *websocketAdaptiveUpdates,
				BandwidthThrottlingEnabled: *websocketBandwidthThrottling,
				UpdateIntervals: webserver.UpdateIntervals{
					Excellent: 1 * time.Second,
					Good:     2 * time.Second,
					Fair:     5 * time.Second,
					Poor:     10 * time.Second,
				},
				ThrottlingThresholds: webserver.ThrottlingThresholds{
					Excellent: 1000000, // 1MB/s
					Good:     500000,  // 500KB/s
					Fair:     200000,  // 200KB/s
					Poor:     50000,   // 50KB/s
				},
			}
			pipeline.WithWebSocketConnectionQualityConfig(qualityConfig)
		}
		
		// Configure ML-based anomaly detection
		pipeline.WithWebSocketMLModel(*websocketMLModel)
		pipeline.WithWebSocketAdvancedML(*websocketAdvancedML)
		pipeline.WithWebSocketPhase4Features(*websocketPhase4Features)
		pipeline.WithWebSocketPhase5Features(*websocketPhase5Features)
		
		if *performanceAlerts != "" {
			pipeline.WithPerformanceAlerts(*performanceAlerts)
		}
		fmt.Printf("WebSocket server enabled on port %d\n", *websocketPort)
		fmt.Printf("  Connection Quality: %s\n", boolToStatus(*websocketConnectionQuality))
		fmt.Printf("  ML Anomaly Detection: %s\n", boolToStatus(*websocketMLModel))
		fmt.Printf("  Advanced ML Features: %s\n", boolToStatus(*websocketAdvancedML))
		fmt.Printf("  Phase 4 Features: %s\n", boolToStatus(*websocketPhase4Features))
		fmt.Printf("  Phase 5 Features: %s\n", boolToStatus(*websocketPhase5Features))
		if *websocketCompression {
			fmt.Println("WebSocket compression: ENABLED")
		} else {
			fmt.Println("WebSocket compression: DISABLED")
		}
		if *websocketAuth {
			fmt.Println("WebSocket authentication: ENABLED")
		} else {
			fmt.Println("WebSocket authentication: DISABLED")
		}
		if *websocketBatching {
			fmt.Printf("WebSocket batching: ENABLED (interval: %dms)\n", *websocketBatchInterval)
		} else {
			fmt.Println("WebSocket batching: DISABLED")
		}
		if *websocketConnectionQuality {
			fmt.Println("WebSocket connection quality monitoring: ENABLED")
			fmt.Printf("WebSocket adaptive updates: %s\n", boolToStatus(*websocketAdaptiveUpdates))
			fmt.Printf("WebSocket bandwidth throttling: %s\n", boolToStatus(*websocketBandwidthThrottling))
			if *websocketQualityAlerts != "" {
				fmt.Printf("WebSocket quality alerts: LOADED from %s\n", *websocketQualityAlerts)
			}
		} else {
			fmt.Println("WebSocket connection quality monitoring: DISABLED")
		}
	}

	var err error
	if *targetType == "python" || *targetType == "node" {
		err = pipeline.RunWithTarget(ctx, *pluginName, "", *targetCommand, *duration, 20, *outDir)
	} else {
		err = pipeline.Run(ctx, *pluginName, *targetURL, *duration, 20, *outDir)
	}

	if err != nil {
		fmt.Printf("Run failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results saved to: %s/\n", *outDir)
}

func runLLMCommand() {
	flagSet := flag.NewFlagSet("llm", flag.ExitOnError)
	bundlePath := flagSet.String("bundle", "", "Input bundle path")
	findingsPath := flagSet.String("findings", "", "Input findings path")
	outPath := flagSet.String("out", "", "Output insights path")
	provider := flagSet.String("provider", "mistral", "LLM provider (mistral, openai)")
	llmModel := flagSet.String("model", "", "Model name (provider-specific default if empty)")
	timeout := flagSet.Int("timeout", 20, "API timeout in seconds")
	maxResponse := flagSet.Int("max-response", 4096, "Max response tokens")
	maxPromptChars := flagSet.Int("max-prompt-chars", 12000, "Max prompt characters")
	dryRun := flagSet.Bool("dry-run", false, "Dry run - save prompt without API call")
	flagSet.Parse(os.Args[2:])

	if *bundlePath == "" || *findingsPath == "" || *outPath == "" {
		fmt.Println("Required flags: --bundle, --findings, --out")
		os.Exit(1)
	}

	// Read bundle
	bundleData, err := os.ReadFile(*bundlePath)
	if err != nil {
		fmt.Printf("Failed to read bundle: %v\n", err)
		os.Exit(1)
	}

	var profileBundle model.ProfileBundle
	if err := json.Unmarshal(bundleData, &profileBundle); err != nil {
		fmt.Printf("Failed to parse bundle: %v\n", err)
		os.Exit(1)
	}

	// Read findings
	findingsData, err := os.ReadFile(*findingsPath)
	if err != nil {
		fmt.Printf("Failed to read findings: %v\n", err)
		os.Exit(1)
	}

	var findingsBundle model.FindingsBundle
	if err := json.Unmarshal(findingsData, &findingsBundle); err != nil {
		fmt.Printf("Failed to parse findings: %v\n", err)
		os.Exit(1)
	}

	// Get API key from environment based on provider
	apiKey := ""
	apiKeyEnv := "MISTRAL_API_KEY"
	if *provider == "openai" {
		apiKeyEnv = "OPENAI_API_KEY"
	}
	
	apiKey = os.Getenv(apiKeyEnv)
	if apiKey == "" && !*dryRun {
		fmt.Printf("%s environment variable not set\n", apiKeyEnv)
		os.Exit(1)
	}

	// Set default model if not specified
	if *llmModel == "" {
		if *provider == "openai" {
			llmModel = flagSet.String("model", "gpt-3.5-turbo", "OpenAI model name")
		} else {
			llmModel = flagSet.String("model", "devstral-small-latest", "Mistral model name")
		}
	}

	// Create provider config
	config := llm.ProviderConfig{
		ProviderName: *provider,
		Model:        *llmModel,
		APIKey:       apiKey,
		Timeout:      time.Duration(*timeout) * time.Second,
		MaxResponse:  *maxResponse,
		MaxPrompt:    *maxPromptChars,
		DryRun:       *dryRun,
	}

	// Create insights generator with specific provider
	generator, err := llm.NewInsightsGeneratorWithProvider(config)
	if err != nil {
		fmt.Printf("Failed to create insights generator: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	insights, err := generator.GenerateInsights(ctx, &profileBundle, &findingsBundle)
	if err != nil {
		fmt.Printf("LLM insights generation failed: %v\n", err)
		os.Exit(1)
	}

	// Save insights
	insightsData, err := json.MarshalIndent(insights, "", "  ")
	if err != nil {
		fmt.Printf("Failed to serialize insights: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outPath, insightsData, 0644); err != nil {
		fmt.Printf("Failed to write insights: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Insights saved to: %s\n", *outPath)
	if *dryRun {
		fmt.Println("✓ Dry-run mode: prompt saved to llm_prompt.json")
	}
}

func runWebCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("web", flag.ExitOnError)
	inPath := flagSet.String("in", "", "Input findings path")
	outDir := flagSet.String("outdir", "", "Output directory")
	insightsPath := flagSet.String("insights", "", "Optional insights path")
	flagSet.Parse(os.Args[2:])

	if *inPath == "" || *outDir == "" {
		fmt.Println("Required flags: --in, --outdir")
		fmt.Println("Optional flags: --insights")
		os.Exit(1)
	}

	// Load insights if provided
	var insights *model.InsightsBundle
	if *insightsPath != "" {
		data, err := os.ReadFile(*insightsPath)
		if err != nil {
			fmt.Printf("Warning: failed to read insights: %v\n", err)
		} else {
			var ib model.InsightsBundle
			if err := json.Unmarshal(data, &ib); err != nil {
				fmt.Printf("Warning: failed to parse insights: %v\n", err)
			} else {
				insights = &ib
			}
		}
	}

	ctx := context.Background()
	err := pipeline.GenerateWebReport(ctx, *inPath, *outDir, insights)
	if err != nil {
		fmt.Printf("Web report generation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Web report generated in: %s/\n", *outDir)
	fmt.Println("Open index.html in your browser to view the interactive report")
}

func runWebSocketCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("websocket", flag.ExitOnError)
	findingsPath := flagSet.String("findings", "", "Path to findings.json file")
	insightsPath := flagSet.String("insights", "", "Optional path to insights.json file")
	port := flagSet.Int("port", 8080, "WebSocket server port")
	dataDir := flagSet.String("data-dir", "./data", "Directory for data files")
	autoRefresh := flagSet.Int("auto-refresh", 0, "Auto-refresh interval in seconds (0 to disable)")
	compression := flagSet.Bool("compression", false, "Enable WebSocket message compression")
	batching := flagSet.Bool("batching", false, "Enable WebSocket message batching")
	batchInterval := flagSet.Int("batch-interval", 100, "WebSocket batch interval in milliseconds")
	connectionQuality := flagSet.Bool("connection-quality", false, "Enable WebSocket connection quality monitoring")
	qualityAlerts := flagSet.String("quality-alerts", "", "WebSocket connection quality alert configuration file (JSON)")
	adaptiveUpdates := flagSet.Bool("adaptive-updates", true, "Enable adaptive updates based on connection quality")
	bandwidthThrottling := flagSet.Bool("bandwidth-throttling", true, "Enable bandwidth throttling based on connection quality")
	phase4Features := flagSet.Bool("phase4-features", false, "Enable Phase 4 advanced ML features (deep learning, time series forecasting, automated root cause analysis)")
	phase5Features := flagSet.Bool("phase5-features", false, "Enable Phase 5 advanced ML features (anomaly correlation detection, predictive maintenance, enhanced root cause analysis)")
	flagSet.Parse(os.Args[2:])

	if *findingsPath == "" {
		fmt.Println("Required flag: --findings")
		os.Exit(1)
	}

	// Configure WebSocket server
	batchIntervalDuration := time.Duration(*batchInterval) * time.Millisecond
	pipeline.WithWebSocketServer(*port, *dataDir, false, *compression, *batching, batchIntervalDuration)
	pipeline.WithWebSocketConnectionQuality(*connectionQuality)
	pipeline.WithWebSocketPhase4Features(*phase4Features)
	pipeline.WithWebSocketPhase5Features(*phase5Features)
	
	// Configure connection quality enhancements
	if *connectionQuality {
		// Load quality alerts if specified
		if *qualityAlerts != "" {
			qualityAlertsList, err := loadConnectionQualityAlerts(*qualityAlerts)
			if err != nil {
				log.Printf("Warning: Failed to load connection quality alerts: %v", err)
			} else {
				pipeline.WithWebSocketConnectionQualityAlerts(qualityAlertsList)
			}
		}
		
		// Configure quality-based adaptations
		qualityConfig := webserver.ConnectionQualityConfig{
			AdaptiveUpdatesEnabled: *adaptiveUpdates,
			BandwidthThrottlingEnabled: *bandwidthThrottling,
			UpdateIntervals: webserver.UpdateIntervals{
				Excellent: 1 * time.Second,
				Good:     2 * time.Second,
				Fair:     5 * time.Second,
				Poor:     10 * time.Second,
			},
			ThrottlingThresholds: webserver.ThrottlingThresholds{
				Excellent: 1000000, // 1MB/s
				Good:     500000,  // 500KB/s
				Fair:     200000,  // 200KB/s
				Poor:     50000,   // 50KB/s
			},
		}
		pipeline.WithWebSocketConnectionQualityConfig(qualityConfig)
	}

	// Configure auto-refresh if enabled
	if *autoRefresh > 0 {
		pipeline.WithWebSocketAutoRefresh(time.Duration(*autoRefresh) * time.Second)
		fmt.Printf("Auto-refresh enabled: %d seconds\n", *autoRefresh)
	}

	// Load data
	err := pipeline.LoadWebSocketData(*findingsPath, *insightsPath)
	if err != nil {
		fmt.Printf("Failed to load data: %v\n", err)
		os.Exit(1)
	}

	// Start WebSocket server
	fmt.Printf("Starting WebSocket server on port %d...\n", *port)
	fmt.Printf("WebSocket endpoint: ws://localhost:%d/ws\n", *port)
	fmt.Printf("Health check: http://localhost:%d/health\n", *port)
	fmt.Printf("Authentication: DISABLED (basic WebSocket connections)\n")
	if compression != nil && *compression {
		fmt.Printf("Compression: ENABLED (WebSocket messages will be compressed)\n")
	} else {
		fmt.Printf("Compression: DISABLED (WebSocket messages will be sent uncompressed)\n")
	}
	if batching != nil && *batching {
		fmt.Printf("Batching: ENABLED (interval: %dms)\n", *batchInterval)
	} else {
		fmt.Printf("Batching: DISABLED (messages sent immediately)\n")
	}
	if connectionQuality != nil && *connectionQuality {
		fmt.Printf("Connection quality: ENABLED (tracking latency and packet loss)\n")
		fmt.Printf("Adaptive updates: %s\n", boolToStatus(*adaptiveUpdates))
		fmt.Printf("Bandwidth throttling: %s\n", boolToStatus(*bandwidthThrottling))
		if *qualityAlerts != "" {
			fmt.Printf("Quality alerts: LOADED from %s\n", *qualityAlerts)
		}
		if phase4Features != nil && *phase4Features {
			fmt.Printf("Phase 4 features: ENABLED (deep learning, time series forecasting, automated root cause analysis)\n")
		} else {
			fmt.Printf("Phase 4 features: DISABLED\n")
		}
	} else {
		fmt.Printf("Connection quality: DISABLED\n")
	}
	fmt.Println("Press Ctrl+C to stop the server")

	// Start server in goroutine
	go func() {
		if err := pipeline.StartWebSocketServer(); err != nil {
			fmt.Printf("WebSocket server failed: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Cleanup
	fmt.Println("\nShutting down WebSocket server...")
	pipeline.StopWebSocketServer()
	fmt.Println("WebSocket server stopped")
}

func runDemoCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("demo", flag.ExitOnError)
	repoURL := flagSet.String("repo", "", "Repository URL")
	outDir := flagSet.String("out", "", "Output directory")
	ref := flagSet.String("ref", "", "Branch, tag, or commit (optional)")
	duration := flagSet.Int("duration", 15, "Benchmark duration in seconds")
	concurrentBenchmarks := flagSet.Bool("concurrent", false, "Enable concurrent benchmark execution")
	maxWorkers := flagSet.Int("max-workers", 2, "Maximum concurrent workers for benchmark execution")
	samplingRate := flagSet.Float64("sampling-rate", 1.0, "Profile sampling rate (0.1-1.0)")
	memoryOptimization := flagSet.Bool("memory-optimization", false, "Enable memory optimization for large profiles")
	largeCodebase := flagSet.Bool("large-codebase", false, "Optimize for large codebases")
	remediationEnabled := flagSet.Bool("remediation", false, "Enable automated code remediation suggestions")
	remediationConfidence := flagSet.Float64("remediation-confidence", 0.7, "Minimum confidence threshold for remediation suggestions (0.0-1.0)")
	remediationMaxChanges := flagSet.Int("remediation-max-changes", 3, "Maximum number of code changes per finding")
	remediationCodeLimit := flagSet.Int("remediation-code-limit", 200, "Maximum characters per code change")
	
	// CI/CD Performance Gates
	criticalThreshold := flagSet.Int("critical-threshold", 5, "Critical findings threshold for CI/CD gates")
	highThreshold := flagSet.Int("high-threshold", 10, "High findings threshold for CI/CD gates")
	mediumThreshold := flagSet.Int("medium-threshold", 20, "Medium findings threshold for CI/CD gates")
	failOnCritical := flagSet.Bool("fail-on-critical", true, "Fail build on critical threshold exceedance")
	failOnHigh := flagSet.Bool("fail-on-high", false, "Fail build on high threshold exceedance")
	warnOnMedium := flagSet.Bool("warn-on-medium", true, "Warn on medium threshold exceedance")
	
	// Enterprise Features
	enterpriseEnabled := flagSet.Bool("enterprise", false, "Enable enterprise features")
	teamName := flagSet.String("team", "", "Team name for enterprise collaboration")
	userName := flagSet.String("user", "", "User name for enterprise auditing")
	auditLogging := flagSet.Bool("audit-logging", false, "Enable audit logging for enterprise compliance")
	rbacEnabled := flagSet.Bool("rbac", false, "Enable role-based access control")
	maxUsers := flagSet.Int("max-users", 10, "Maximum number of users for enterprise license")
	maxTeams := flagSet.Int("max-teams", 5, "Maximum number of teams for enterprise license")
	
	flagSet.Parse(os.Args[2:])

	if *repoURL == "" || *outDir == "" {
		fmt.Println("Required flags: --repo, --out")
		fmt.Println("Optional flags: --ref, --duration")
		os.Exit(1)
	}

	ctx := context.Background()

	fmt.Printf("🚀 Starting demo for repository: %s\n", *repoURL)
	if *ref != "" {
		fmt.Printf("📌 Using reference: %s\n", *ref)
	}
	fmt.Printf("⏱  Benchmark duration: %d seconds\n", *duration)
	fmt.Printf("📁 Output directory: %s\n", *outDir)
	
	// Display enterprise settings
	if *enterpriseEnabled {
		fmt.Printf("\n🏢 Enterprise Features:\n")
		fmt.Printf("   ✅ Enterprise mode: ENABLED\n")
		if *teamName != "" {
			fmt.Printf("   👥 Team: %s\n", *teamName)
		}
		if *userName != "" {
			fmt.Printf("   👤 User: %s\n", *userName)
		}
		if *auditLogging {
			fmt.Printf("   📝 Audit logging: ENABLED\n")
		} else {
			fmt.Printf("   ❌ Audit logging: DISABLED\n")
		}
		if *rbacEnabled {
			fmt.Printf("   🔐 RBAC: ENABLED\n")
		} else {
			fmt.Printf("   ❌ RBAC: DISABLED\n")
		}
		fmt.Printf("   👥 Max users: %d\n", *maxUsers)
		fmt.Printf("   👥 Max teams: %d\n", *maxTeams)
	} else {
		fmt.Printf("\n🏢 Enterprise Features: DISABLED\n")
	}
	
	// Display performance optimization settings
	fmt.Printf("\n🔧 Performance Optimization Settings:\n")
	if *concurrentBenchmarks {
		fmt.Printf("   ✅ Concurrent benchmarks: ENABLED (max workers: %d)\n", *maxWorkers)
	} else {
		fmt.Printf("   ❌ Concurrent benchmarks: DISABLED\n")
	}
	if *samplingRate < 1.0 {
		fmt.Printf("   ✅ Profile sampling: ENABLED (rate: %.1f)\n", *samplingRate)
	} else {
		fmt.Printf("   ❌ Profile sampling: DISABLED\n")
	}
	if *memoryOptimization {
		fmt.Printf("   ✅ Memory optimization: ENABLED\n")
	} else {
		fmt.Printf("   ❌ Memory optimization: DISABLED\n")
	}
	if *largeCodebase {
		fmt.Printf("   ✅ Large codebase mode: ENABLED\n")
	} else {
		fmt.Printf("   ❌ Large codebase mode: DISABLED\n")
	}
	fmt.Printf("\n")

	// Display remediation settings
	fmt.Printf("🔧 Remediation Settings:\n")
	if *remediationEnabled {
		fmt.Printf("   ✅ Automated remediation: ENABLED\n")
		fmt.Printf("   📊 Minimum confidence: %.1f\n", *remediationConfidence)
		fmt.Printf("   🔄 Max changes per finding: %d\n", *remediationMaxChanges)
		fmt.Printf("   📝 Max code characters: %d\n", *remediationCodeLimit)
	} else {
		fmt.Printf("   ❌ Automated remediation: DISABLED\n")
	}
	fmt.Printf("\n")

	// Check for API key and enable LLM automatically (mandatory per COMPASS.md)
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if apiKey != "" {
		fmt.Printf("🤖 LLM Insights: ENABLED (using %s provider)\n", 
			func() string {
				if os.Getenv("MISTRAL_API_KEY") != "" {
					return "Mistral"
				} else {
					return "OpenAI"
				}
			}())
		
		// Configure LLM with automatic provider detection
		var err error
		if os.Getenv("MISTRAL_API_KEY") != "" {
			pipeline, err = pipeline.WithLLM(apiKey, "mistral-large-latest", 30, 8000, 4000, false)
		} else {
			pipeline, err = pipeline.WithLLM(apiKey, "gpt-4o", 30, 8000, 4000, false)
		}
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to configure LLM: %v\n", err)
			fmt.Println("💡 Continuing without LLM insights (deterministic analysis only)")
		}
	} else {
		fmt.Println("🤖 LLM Insights: DISABLED (no API key found)")
		fmt.Println("📋 To enable LLM insights, set MISTRAL_API_KEY or OPENAI_API_KEY environment variable")
		fmt.Println("   Example: export MISTRAL_API_KEY='your-api-key-here'")
	}
	fmt.Printf("\n")

	// Create performance configuration
	perfConfig := &model.PerformanceOptimizationConfig{
		EnableConcurrentBenchmarks: *concurrentBenchmarks,
		MaxConcurrentWorkers:       *maxWorkers,
		EnableProfileSampling:      *samplingRate < 1.0,
		SamplingRate:               *samplingRate,
		EnableMemoryOptimization:   *memoryOptimization,
		LargeCodebaseMode:          *largeCodebase,
	}

	// Create remediation configuration
	remediationConfig := &model.RemediationConfig{
		Enabled:           *remediationEnabled,
		MinConfidence:     *remediationConfidence,
		MaxCodeChanges:    *remediationMaxChanges,
		CodeChangeLimit:   *remediationCodeLimit,
		Provider:          "mistral",
		Model:             "mistral-large-latest",
		Temperature:       0.3,
	}

	// Create performance gate configuration
	performanceGateConfig := model.PerformanceGateConfig{
		Enabled:                    true,
		CriticalFindingsThreshold:  *criticalThreshold,
		HighFindingsThreshold:      *highThreshold,
		MediumFindingsThreshold:    *mediumThreshold,
		FailOnCriticalThreshold:    *failOnCritical,
		FailOnHighThreshold:        *failOnHigh,
		WarnOnMediumThreshold:      *warnOnMedium,
	}

	// Create enterprise configuration
	enterpriseConfig := model.EnterpriseConfig{
		Enabled:          *enterpriseEnabled,
		TeamName:         *teamName,
		UserName:         *userName,
		AuditLogging:     *auditLogging,
		RBACEnabled:      *rbacEnabled,
		MaxUsers:         *maxUsers,
		MaxTeams:         *maxTeams,
	}

	// Configure enterprise features in pipeline
	pipeline.WithEnterpriseConfig(enterpriseConfig)

	// Configure performance gates in pipeline
	pipeline.WithPerformanceGates(performanceGateConfig)

	// Run the demo workflow with performance configuration
	manifest, err := pipeline.DemoWithPerformance(ctx, *repoURL, *ref, *outDir, *duration, perfConfig)
	if err != nil {
		fmt.Printf("❌ Demo failed: %v\n", err)

		// Provide enhanced error information if available
		if manifest != nil && manifest.ErrorContext != nil {
			fmt.Printf("\n🔍 Error Details:\n")
			fmt.Printf("   Type: [%s:%s]\n", manifest.ErrorContext.ErrorType, manifest.ErrorContext.ErrorCode)
			fmt.Printf("   Message: %s\n", manifest.ErrorContext.Message)
			if manifest.ErrorContext.Details != "" {
				fmt.Printf("   Details: %s\n", manifest.ErrorContext.Details)
			}
			if manifest.ErrorContext.Suggestion != "" {
				fmt.Printf("   💡 Suggestion: %s\n", manifest.ErrorContext.Suggestion)
			}
			if manifest.ErrorContext.IsRecoverable {
				fmt.Printf("   🔄 This error is recoverable\n")
			}
		}
		
		os.Exit(1)
	}

	// Add enterprise config to manifest
	manifest.EnterpriseConfig = enterpriseConfig

	// Add performance gate config to manifest
	manifest.PerformanceGateConfig = performanceGateConfig

	// Log audit action if enterprise features are enabled
	if enterpriseConfig.Enabled && enterpriseConfig.AuditLogging {
		userID := *userName
		if userID == "" {
			userID = "system"
		}
		pipeline.LogAuditAction(userID, "demo_completed", *repoURL, 
			fmt.Sprintf("Benchmarks: %d, Profiles: %d, Duration: %ds", len(manifest.Benchmarks), len(manifest.Profiles), *duration), "success")
		
		// Display audit summary
		auditSummary := pipeline.GetAuditSummary()
		if auditSummary["enabled"].(bool) {
			fmt.Printf("\n📝 Audit Log Summary:\n")
			fmt.Printf("   Total entries: %d\n", auditSummary["total_entries"])
			fmt.Printf("   Status: ENABLED\n")
		}
	}

	// Add enterprise config to manifest
	manifest.EnterpriseConfig = enterpriseConfig

	// Log audit action if enterprise features are enabled
	if enterpriseConfig.Enabled && enterpriseConfig.AuditLogging {
		userID := *userName
		if userID == "" {
			userID = "system"
		}
		pipeline.LogAuditAction(userID, "demo_kit_completed", "built-in-demo-repo", 
			fmt.Sprintf("Benchmarks: %d, Profiles: %d, Duration: %ds", len(manifest.Benchmarks), len(manifest.Profiles), *duration), "success")
		
		// Display audit summary
		auditSummary := pipeline.GetAuditSummary()
		if auditSummary["enabled"].(bool) {
			fmt.Printf("\n📝 Audit Log Summary:\n")
			fmt.Printf("   Total entries: %d\n", auditSummary["total_entries"])
			fmt.Printf("   Status: ENABLED\n")
		}
	}

	// Save run manifest
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Printf("⚠️  Warning: failed to serialize run manifest: %v\n", err)
	} else {
		manifestPath := filepath.Join(*outDir, "run.json")
		if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
			fmt.Printf("⚠️  Warning: failed to write run manifest: %v\n", err)
		} else {
			fmt.Printf("📋 Run manifest saved to: %s\n", manifestPath)
		}
	}

	if manifest.Success {
		fmt.Println("\n✅ Demo completed successfully!")
		fmt.Printf("📊 Found %d benchmarks\n", len(manifest.Benchmarks))
		fmt.Printf("📈 Generated %d profiles\n", len(manifest.Profiles))
		fmt.Println("\n📄 Generated files:")
		fmt.Printf("  📋 %s\n", filepath.Join(*outDir, "run.json"))
		fmt.Printf("  📦 %s\n", filepath.Join(*outDir, "bundle.json"))
		fmt.Printf("  🔍 %s\n", filepath.Join(*outDir, "findings.json"))
		fmt.Printf("  📊 %s\n", filepath.Join(*outDir, "report.md"))
		for _, profile := range manifest.Profiles {
			fmt.Printf("  📈 %s\n", profile)
		}

		// Check performance gates
		fmt.Println("\n🚦 Checking performance gates...")
		findingsPath := filepath.Join(*outDir, "findings.json")
		findingsData, err := os.ReadFile(findingsPath)
		if err != nil {
			fmt.Printf("⚠️  Warning: failed to read findings for performance gate check: %v\n", err)
		} else {
			var findingsBundle struct {
				Findings []model.Finding `json:"findings"`
			}
			if err := json.Unmarshal(findingsData, &findingsBundle); err != nil {
				fmt.Printf("⚠️  Warning: failed to parse findings: %v\n", err)
			} else {
				gateResult, err := pipeline.CheckPerformanceGates(findingsBundle.Findings)
				if err != nil {
					fmt.Printf("⚠️  Warning: performance gate check failed: %v\n", err)
				} else {
					if gateResult.Passed {
						fmt.Printf("✅ Performance gates PASSED\n")
						if len(gateResult.Warnings) > 0 {
							fmt.Printf("⚠️  Warnings:\n")
							for _, warning := range gateResult.Warnings {
								fmt.Printf("   - %s\n", warning)
							}
						}
					} else {
						fmt.Printf("❌ Performance gates FAILED\n")
						fmt.Printf("🔍 Reason: %s\n", gateResult.Message)
						if len(gateResult.Errors) > 0 {
							fmt.Printf("💥 Errors:\n")
							for _, error := range gateResult.Errors {
								fmt.Printf("   - %s\n", error)
							}
						}
						if performanceGateConfig.FailOnCriticalThreshold || performanceGateConfig.FailOnHighThreshold {
							fmt.Printf("\n💥 Build failed due to performance gate violations\n")
							os.Exit(1)
							return
						}
					}
					
					// Display severity distribution
					fmt.Printf("\n📊 Finding Severity Distribution:\n")
					for severity, count := range gateResult.SeverityCounts {
						fmt.Printf("   %s: %d\n", severity, count)
					}
				}
			}
		}

		// Generate remediations if enabled
		if *remediationEnabled {
			fmt.Printf("\n🔧 Generating automated remediations...\n")
			
			findingsPath := filepath.Join(*outDir, "findings.json")
			insightsPath := filepath.Join(*outDir, "insights.json")
			remediationsPath := filepath.Join(*outDir, "remediations.json")
			
			// Generate remediations
			remediations, err := pipeline.GenerateRemediations(ctx, findingsPath, insightsPath, remediationsPath, *remediationConfig)
			if err != nil {
				fmt.Printf("⚠️  Warning: failed to generate remediations: %v\n", err)
			} else if remediations != nil {
				fmt.Printf("✅ Generated %d remediation suggestions\n", len(remediations.Remediations))
				fmt.Printf("   Estimated total gain: %s\n", remediations.Summary.EstimatedTotalGain)
				fmt.Printf("   Average confidence: %.1f%%\n", remediations.Summary.ConfidenceScore*100)
				fmt.Printf("  🛠️  %s\n", remediationsPath)
			}
		}
	} else {
		fmt.Printf("\n❌ Demo failed: %s\n", manifest.Error)
		
		// Provide enhanced error information if available
		if manifest.ErrorContext != nil {
			fmt.Printf("\n🔍 Error Details:\n")
			fmt.Printf("   Type: [%s:%s]\n", manifest.ErrorContext.ErrorType, manifest.ErrorContext.ErrorCode)
			fmt.Printf("   Message: %s\n", manifest.ErrorContext.Message)
			if manifest.ErrorContext.Details != "" {
				fmt.Printf("   Details: %s\n", manifest.ErrorContext.Details)
			}
			if manifest.ErrorContext.Suggestion != "" {
				fmt.Printf("   💡 Suggestion: %s\n", manifest.ErrorContext.Suggestion)
			}
			if manifest.ErrorContext.IsRecoverable {
				fmt.Printf("   🔄 This error is recoverable\n")
			}
		}
		
		os.Exit(1)
	}
}

func runDemoKitCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("demo-kit", flag.ExitOnError)
	outDir := flagSet.String("out", "", "Output directory")
	duration := flagSet.Int("duration", 15, "Benchmark duration in seconds")
	concurrentBenchmarks := flagSet.Bool("concurrent", false, "Enable concurrent benchmark execution")
	maxWorkers := flagSet.Int("max-workers", 2, "Maximum concurrent workers for benchmark execution")
	samplingRate := flagSet.Float64("sampling-rate", 1.0, "Profile sampling rate (0.1-1.0)")
	memoryOptimization := flagSet.Bool("memory-optimization", false, "Enable memory optimization for large profiles")
	largeCodebase := flagSet.Bool("large-codebase", false, "Optimize for large codebases")
	remediationEnabled := flagSet.Bool("remediation", false, "Enable automated code remediation suggestions")
	remediationConfidence := flagSet.Float64("remediation-confidence", 0.7, "Minimum confidence threshold for remediation suggestions (0.0-1.0)")
	remediationMaxChanges := flagSet.Int("remediation-max-changes", 3, "Maximum number of code changes per finding")
	remediationCodeLimit := flagSet.Int("remediation-code-limit", 200, "Maximum characters per code change")
	
	// CI/CD Performance Gates
	criticalThreshold := flagSet.Int("critical-threshold", 5, "Critical findings threshold for CI/CD gates")
	highThreshold := flagSet.Int("high-threshold", 10, "High findings threshold for CI/CD gates")
	mediumThreshold := flagSet.Int("medium-threshold", 20, "Medium findings threshold for CI/CD gates")
	failOnCritical := flagSet.Bool("fail-on-critical", true, "Fail build on critical threshold exceedance")
	failOnHigh := flagSet.Bool("fail-on-high", false, "Fail build on high threshold exceedance")
	warnOnMedium := flagSet.Bool("warn-on-medium", true, "Warn on medium threshold exceedance")
	
	// Enterprise Features
	enterpriseEnabled := flagSet.Bool("enterprise", false, "Enable enterprise features")
	teamName := flagSet.String("team", "", "Team name for enterprise collaboration")
	userName := flagSet.String("user", "", "User name for enterprise auditing")
	auditLogging := flagSet.Bool("audit-logging", false, "Enable audit logging for enterprise compliance")
	rbacEnabled := flagSet.Bool("rbac", false, "Enable role-based access control")
	maxUsers := flagSet.Int("max-users", 10, "Maximum number of users for enterprise license")
	maxTeams := flagSet.Int("max-teams", 5, "Maximum number of teams for enterprise license")
	
	flagSet.Parse(os.Args[2:])

	if *outDir == "" {
		fmt.Println("Required flags: --out")
		fmt.Println("Optional flags: --duration")
		os.Exit(1)
	}

	// Use the built-in demo repository
	demoRepoPath := filepath.Join("examples", "demo-repo-simple")
	
	ctx := context.Background()

	fmt.Println("🚀 Starting TriageProf Demo Kit")
	fmt.Printf("📌 Using built-in demo repository: %s\n", demoRepoPath)
	fmt.Printf("⏱  Benchmark duration: %d seconds\n", *duration)
	fmt.Printf("📁 Output directory: %s\n", *outDir)
	
	fmt.Println("🛫 Pre-flight validation...")
	
	// Validate environment dependencies
	if errContext, ok := core.ValidateDemoEnvironment(ctx); !ok {
		fmt.Printf("❌ Environment validation failed: %s\n", errContext.Message)
		fmt.Printf("💡 Suggestion: %s\n", errContext.Suggestion)
		if errContext.IsRecoverable {
			fmt.Println("🔄 This error is recoverable")
		}
		os.Exit(1)
	}
	
	// Validate duration
	if *duration < 1 {
		errContext := model.NewErrorContext(
			model.ErrorTypeValidation,
			model.ErrorCodeInvalidInput,
			"Invalid benchmark duration",
			fmt.Sprintf("Duration must be at least 1 second, got: %d", *duration),
			"Use --duration flag with a value >= 1",
			true,
		)
		fmt.Printf("❌ Invalid duration: %s\n", errContext.Message)
		fmt.Printf("💡 Suggestion: %s\n", errContext.Suggestion)
		os.Exit(1)
	}
	
	// Validate output directory
	if *outDir == "" {
		errContext := model.NewErrorContext(
			model.ErrorTypeValidation,
			model.ErrorCodeInvalidInput,
			"Output directory not specified",
			"Output directory path is empty",
			"Use --out flag to specify output directory",
			true,
		)
		fmt.Printf("❌ Invalid output directory: %s\n", errContext.Message)
		fmt.Printf("💡 Suggestion: %s\n", errContext.Suggestion)
		os.Exit(1)
	}
	
	fmt.Println("✅ Pre-flight validation completed")
	fmt.Println("🚀 Running demo workflow...")
	
	// Display performance optimization settings
	fmt.Printf("\n🔧 Performance Optimization Settings:\n")
	if *concurrentBenchmarks {
		fmt.Printf("   ✅ Concurrent benchmarks: ENABLED (max workers: %d)\n", *maxWorkers)
	} else {
		fmt.Printf("   ❌ Concurrent benchmarks: DISABLED\n")
	}
	if *samplingRate < 1.0 {
		fmt.Printf("   ✅ Profile sampling: ENABLED (rate: %.1f)\n", *samplingRate)
	} else {
		fmt.Printf("   ❌ Profile sampling: DISABLED\n")
	}
	if *memoryOptimization {
		fmt.Printf("   ✅ Memory optimization: ENABLED\n")
	} else {
		fmt.Printf("   ❌ Memory optimization: DISABLED\n")
	}
	if *largeCodebase {
		fmt.Printf("   ✅ Large codebase mode: ENABLED\n")
	} else {
		fmt.Printf("   ❌ Large codebase mode: DISABLED\n")
	}
	fmt.Printf("\n")

	// Check for API key and enable LLM automatically (mandatory per COMPASS.md)
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if apiKey != "" {
		fmt.Printf("🤖 LLM Insights: ENABLED (using %s provider)\n", 
			func() string {
				if os.Getenv("MISTRAL_API_KEY") != "" {
					return "Mistral"
				} else {
					return "OpenAI"
				}
			}())
		
		// Configure LLM with automatic provider detection
		var err error
		if os.Getenv("MISTRAL_API_KEY") != "" {
			pipeline, err = pipeline.WithLLM(apiKey, "mistral-large-latest", 30, 8000, 4000, false)
		} else {
			pipeline, err = pipeline.WithLLM(apiKey, "gpt-4o", 30, 8000, 4000, false)
		}
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to configure LLM: %v\n", err)
			fmt.Println("💡 Continuing without LLM insights (deterministic analysis only)")
		}
	} else {
		fmt.Println("🤖 LLM Insights: DISABLED (no API key found)")
		fmt.Println("📋 To enable LLM insights, set MISTRAL_API_KEY or OPENAI_API_KEY environment variable")
		fmt.Println("   Example: export MISTRAL_API_KEY='your-api-key-here'")
	}
	fmt.Printf("\n")

	// Display enterprise settings
	if *enterpriseEnabled {
		fmt.Printf("\n🏢 Enterprise Features:\n")
		fmt.Printf("   ✅ Enterprise mode: ENABLED\n")
		if *teamName != "" {
			fmt.Printf("   👥 Team: %s\n", *teamName)
		}
		if *userName != "" {
			fmt.Printf("   👤 User: %s\n", *userName)
		}
		if *auditLogging {
			fmt.Printf("   📝 Audit logging: ENABLED\n")
		} else {
			fmt.Printf("   ❌ Audit logging: DISABLED\n")
		}
		if *rbacEnabled {
			fmt.Printf("   🔐 RBAC: ENABLED\n")
		} else {
			fmt.Printf("   ❌ RBAC: DISABLED\n")
		}
		fmt.Printf("   👥 Max users: %d\n", *maxUsers)
		fmt.Printf("   👥 Max teams: %d\n", *maxTeams)
	} else {
		fmt.Printf("\n🏢 Enterprise Features: DISABLED\n")
	}

	// Create performance configuration
	perfConfig := &model.PerformanceOptimizationConfig{
		EnableConcurrentBenchmarks: *concurrentBenchmarks,
		MaxConcurrentWorkers:       *maxWorkers,
		EnableProfileSampling:      *samplingRate < 1.0,
		SamplingRate:               *samplingRate,
		EnableMemoryOptimization:   *memoryOptimization,
		LargeCodebaseMode:          *largeCodebase,
	}

	// Create remediation configuration
	remediationConfig := &model.RemediationConfig{
		Enabled:           *remediationEnabled,
		MinConfidence:     *remediationConfidence,
		MaxCodeChanges:    *remediationMaxChanges,
		CodeChangeLimit:   *remediationCodeLimit,
		Provider:          "mistral",
		Model:             "mistral-large-latest",
		Temperature:       0.3,
	}

	// Create performance gate configuration
	performanceGateConfig := model.PerformanceGateConfig{
		Enabled:                    true,
		CriticalFindingsThreshold:  *criticalThreshold,
		HighFindingsThreshold:      *highThreshold,
		MediumFindingsThreshold:    *mediumThreshold,
		FailOnCriticalThreshold:    *failOnCritical,
		FailOnHighThreshold:        *failOnHigh,
		WarnOnMediumThreshold:      *warnOnMedium,
	}

	// Create enterprise configuration
	enterpriseConfig := model.EnterpriseConfig{
		Enabled:          *enterpriseEnabled,
		TeamName:         *teamName,
		UserName:         *userName,
		AuditLogging:     *auditLogging,
		RBACEnabled:      *rbacEnabled,
		MaxUsers:         *maxUsers,
		MaxTeams:         *maxTeams,
	}

	// Configure enterprise features in pipeline
	pipeline.WithEnterpriseConfig(enterpriseConfig)

	// Configure performance gates in pipeline
	pipeline.WithPerformanceGates(performanceGateConfig)

	// Run the demo workflow with the local demo repository
	manifest, err := pipeline.DemoWithPerformance(ctx, demoRepoPath, "", *outDir, *duration, perfConfig)
	if err != nil {
		fmt.Printf("❌ Demo kit failed: %v\n", err)
		
		// Provide enhanced error information if available
		if manifest != nil && manifest.ErrorContext != nil {
			fmt.Printf("\n🔍 Error Details:\n")
			fmt.Printf("   Type: [%s:%s]\n", manifest.ErrorContext.ErrorType, manifest.ErrorContext.ErrorCode)
			fmt.Printf("   Message: %s\n", manifest.ErrorContext.Message)
			if manifest.ErrorContext.Details != "" {
				fmt.Printf("   Details: %s\n", manifest.ErrorContext.Details)
			}
			if manifest.ErrorContext.Suggestion != "" {
				fmt.Printf("   💡 Suggestion: %s\n", manifest.ErrorContext.Suggestion)
			}
			if manifest.ErrorContext.IsRecoverable {
				fmt.Printf("   🔄 This error is recoverable\n")
			}
		}
		
		os.Exit(1)
	}

	// Save run manifest
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Printf("⚠️  Warning: failed to serialize run manifest: %v\n", err)
	} else {
		manifestPath := filepath.Join(*outDir, "run.json")
		if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
			fmt.Printf("⚠️  Warning: failed to write run manifest: %v\n", err)
		} else {
			fmt.Printf("📋 Run manifest saved to: %s\n", manifestPath)
		}
	}
	
	fmt.Println("✅ Demo workflow completed")

	// Log audit action if enterprise features are enabled
	if enterpriseConfig.Enabled && enterpriseConfig.AuditLogging {
		userID := *userName
		if userID == "" {
			userID = "system"
		}
		pipeline.LogAuditAction(userID, "demo_kit_completed", "built-in-demo-repo", 
			fmt.Sprintf("Benchmarks: %d, Profiles: %d, Duration: %ds", len(manifest.Benchmarks), len(manifest.Profiles), *duration), "success")
		
		// Display audit summary
		auditSummary := pipeline.GetAuditSummary()
		if auditSummary["enabled"].(bool) {
			fmt.Printf("\n📝 Audit Log Summary:\n")
			fmt.Printf("   Total entries: %d\n", auditSummary["total_entries"])
			fmt.Printf("   Status: ENABLED\n")
		}
	}

	if manifest.Success {
		fmt.Println("\n✅ Demo Kit completed successfully!")
		fmt.Printf("📊 Found %d benchmarks\n", len(manifest.Benchmarks))
		fmt.Printf("📈 Generated %d profiles\n", len(manifest.Profiles))
		fmt.Println("\n📄 Generated files:")
		fmt.Printf("  📋 %s\n", filepath.Join(*outDir, "run.json"))
		fmt.Printf("  📦 %s\n", filepath.Join(*outDir, "bundle.json"))
		fmt.Printf("  🔍 %s\n", filepath.Join(*outDir, "findings.json"))
		fmt.Printf("  📊 %s\n", filepath.Join(*outDir, "report.md"))
		for _, profile := range manifest.Profiles {
			fmt.Printf("  📈 %s\n", profile)
		}

		// Check performance gates
		fmt.Println("\n🚦 Checking performance gates...")
		findingsPath := filepath.Join(*outDir, "findings.json")
		findingsData, err := os.ReadFile(findingsPath)
		if err != nil {
			fmt.Printf("⚠️  Warning: failed to read findings for performance gate check: %v\n", err)
		} else {
			var findingsBundle struct {
				Findings []model.Finding `json:"findings"`
			}
			if err := json.Unmarshal(findingsData, &findingsBundle); err != nil {
				fmt.Printf("⚠️  Warning: failed to parse findings: %v\n", err)
			} else {
				gateResult, err := pipeline.CheckPerformanceGates(findingsBundle.Findings)
				if err != nil {
					fmt.Printf("⚠️  Warning: performance gate check failed: %v\n", err)
				} else {
					if gateResult.Passed {
						fmt.Printf("✅ Performance gates PASSED\n")
						if len(gateResult.Warnings) > 0 {
							fmt.Printf("⚠️  Warnings:\n")
							for _, warning := range gateResult.Warnings {
								fmt.Printf("   - %s\n", warning)
							}
						}
					} else {
						fmt.Printf("❌ Performance gates FAILED\n")
						fmt.Printf("🔍 Reason: %s\n", gateResult.Message)
						if len(gateResult.Errors) > 0 {
							fmt.Printf("💥 Errors:\n")
							for _, error := range gateResult.Errors {
								fmt.Printf("   - %s\n", error)
							}
						}
						if performanceGateConfig.FailOnCriticalThreshold || performanceGateConfig.FailOnHighThreshold {
							fmt.Printf("\n💥 Build failed due to performance gate violations\n")
							os.Exit(1)
							return
						}
					}
					
					// Display severity distribution
					fmt.Printf("\n📊 Finding Severity Distribution:\n")
					for severity, count := range gateResult.SeverityCounts {
						fmt.Printf("   %s: %d\n", severity, count)
					}
				}
			}
		}

		// Generate remediations if enabled
		if *remediationEnabled {
			fmt.Printf("\n🔧 Generating automated remediations...\n")
			
			findingsPath := filepath.Join(*outDir, "findings.json")
			insightsPath := filepath.Join(*outDir, "insights.json")
			remediationsPath := filepath.Join(*outDir, "remediations.json")
			
			// Generate remediations
			remediations, err := pipeline.GenerateRemediations(ctx, findingsPath, insightsPath, remediationsPath, *remediationConfig)
			if err != nil {
				fmt.Printf("⚠️  Warning: failed to generate remediations: %v\n", err)
			} else if remediations != nil {
				fmt.Printf("✅ Generated %d remediation suggestions\n", len(remediations.Remediations))
				fmt.Printf("   Estimated total gain: %s\n", remediations.Summary.EstimatedTotalGain)
				fmt.Printf("   Average confidence: %.1f%%\n", remediations.Summary.ConfidenceScore*100)
				fmt.Printf("  🛠️  %s\n", remediationsPath)
			}
		}
		
		// Post-demo verification
		fmt.Println("\n🔍 Post-demo verification...")
		
		// Verify expected output files exist
		fmt.Println("📋 Verifying output files...")
		expectedFiles := []string{"bundle.json", "findings.json", "report.md"}
		missingFiles := 0
		
		for _, file := range expectedFiles {
			filePath := filepath.Join(*outDir, file)
			if _, err := os.Stat(filePath); err != nil {
				fmt.Printf("❌ Missing expected file: %s\n", file)
				missingFiles++
			} else {
				fmt.Printf("✅ Found: %s\n", file)
			}
		}
		
		if missingFiles > 0 {
			fmt.Printf("⚠️  Warning: %d expected files missing\n", missingFiles)
		} else {
			fmt.Println("✅ All expected output files present")
		}
		
		// Verify profiles were generated
		if len(manifest.Profiles) == 0 {
			fmt.Println("⚠️  Warning: No profiles were generated")
		} else {
			fmt.Printf("✅ Generated %d profiles\n", len(manifest.Profiles))
		}
		
		// Verify benchmarks were found
		if len(manifest.Benchmarks) == 0 {
			fmt.Println("⚠️  Warning: No benchmarks were found")
		} else {
			fmt.Printf("✅ Found %d benchmarks\n", len(manifest.Benchmarks))
		}
		
		fmt.Println("✅ Demo verification completed!")
		
		fmt.Println("\n🎉 Demo Kit Features Showcased:")
		fmt.Println("  ✓ Go benchmark detection")
		fmt.Println("  ✓ CPU, Heap, Allocs, Block, and Mutex profiling")
		fmt.Println("  ✓ Deterministic performance analysis")
		fmt.Println("  ✓ Structured findings with evidence")
		fmt.Println("  ✓ Markdown report generation")
		fmt.Println("  ✓ Profile artifact collection")
		if *remediationEnabled {
			fmt.Println("  ✓ Automated code remediation suggestions")
		}
		
		fmt.Println("\n🚀 Next Steps:")
		fmt.Println("  • Try with your own Go repository: triageprof demo --repo <your-repo-url> --out <dir>")
		fmt.Println("  • Enable LLM insights: Set MISTRAL_API_KEY environment variable")
		fmt.Println("  • Generate HTML reports: triageprof web --in findings.json --outdir <dir>")
		fmt.Println("  • Explore WebSocket dashboard: triageprof websocket --findings findings.json")
	} else {
		fmt.Printf("\n❌ Demo Kit failed: %s\n", manifest.Error)
		
		// Provide enhanced error information if available
		if manifest.ErrorContext != nil {
			fmt.Printf("\n🔍 Error Details:\n")
			fmt.Printf("   Type: [%s:%s]\n", manifest.ErrorContext.ErrorType, manifest.ErrorContext.ErrorCode)
			fmt.Printf("   Message: %s\n", manifest.ErrorContext.Message)
			if manifest.ErrorContext.Details != "" {
				fmt.Printf("   Details: %s\n", manifest.ErrorContext.Details)
			}
			if manifest.ErrorContext.Suggestion != "" {
				fmt.Printf("   💡 Suggestion: %s\n", manifest.ErrorContext.Suggestion)
			}
			if manifest.ErrorContext.IsRecoverable {
				fmt.Printf("   🔄 This error is recoverable\n")
			}
		}
		
		os.Exit(1)
	}
}


