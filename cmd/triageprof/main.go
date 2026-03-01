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
		fmt.Println("  analyze --in <bundle.json> --out <findings.json> --top <N> [--callgraph --callgraph-depth <depth>] [--regression --baseline <path>]")
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
	flagSet.Parse(os.Args[2:])

	if *inPath == "" || *outPath == "" {
		fmt.Println("Required flags: --in, --out")
		os.Exit(1)
	}

	if *regression && *baseline == "" {
		fmt.Println("Regression analysis requires --baseline flag")
		os.Exit(1)
	}

	ctx := context.Background()
	options := core.CoreAnalyzeOptions{
		EnableCallgraph:    *callgraph,
		CallgraphDepth:     *callgraphDepth,
		EnableRegression:   *regression,
		BaselineBundlePath: *baseline,
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
	if *regression {
		fmt.Println("✓ Regression analysis completed")
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
				llmModel = flagSet.String("llm-model", "gpt-3.5-turbo", "OpenAI model name")
			} else {
				llmModel = flagSet.String("llm-model", "devstral-small-latest", "Mistral model name")
			}
		}
		
		// Create provider config
		config := llm.ProviderConfig{
			ProviderName: *llmProvider,
			Model:        *llmModel,
			APIKey:       apiKey,
			Timeout:      time.Duration(*llmTimeout) * time.Second,
			MaxResponse:  4096,
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

	// Create performance configuration
	perfConfig := &model.PerformanceOptimizationConfig{
		EnableConcurrentBenchmarks: *concurrentBenchmarks,
		MaxConcurrentWorkers:       *maxWorkers,
		EnableProfileSampling:      *samplingRate < 1.0,
		SamplingRate:               *samplingRate,
		EnableMemoryOptimization:   *memoryOptimization,
		LargeCodebaseMode:          *largeCodebase,
	}

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

	// Create performance configuration
	perfConfig := &model.PerformanceOptimizationConfig{
		EnableConcurrentBenchmarks: *concurrentBenchmarks,
		MaxConcurrentWorkers:       *maxWorkers,
		EnableProfileSampling:      *samplingRate < 1.0,
		SamplingRate:               *samplingRate,
		EnableMemoryOptimization:   *memoryOptimization,
		LargeCodebaseMode:          *largeCodebase,
	}

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
		
		fmt.Println("\n🎉 Demo Kit Features Showcased:")
		fmt.Println("  ✓ Go benchmark detection")
		fmt.Println("  ✓ CPU, Heap, Allocs, Block, and Mutex profiling")
		fmt.Println("  ✓ Deterministic performance analysis")
		fmt.Println("  ✓ Structured findings with evidence")
		fmt.Println("  ✓ Markdown report generation")
		fmt.Println("  ✓ Profile artifact collection")
		
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


