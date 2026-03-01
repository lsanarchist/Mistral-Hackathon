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
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: triageprof <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  plugins list")
		fmt.Println("  collect --plugin <name> --target-url <url> --duration <sec> --out <path>")
		fmt.Println("  analyze --in <bundle.json> --out <findings.json> --top <N> [--callgraph --callgraph-depth <depth>] [--regression --baseline <path>]")
		fmt.Println("  report --in <findings.json> --out <report.md|json> --output markdown|json")
		fmt.Println("  llm --bundle <bundle.json> --findings <findings.json> --out <insights.json> [--provider <provider>] [--model <model>] [--timeout <sec>] [--dry-run]")
		fmt.Println("  run --plugin <name> --target-url <url> --duration <sec> --outdir <dir>")
		fmt.Println("  run --plugin <name> --target-type python --target-command <cmd> --duration <sec> --outdir <dir>")
		fmt.Println("  run --plugin <name> --target-type node --target-command <cmd> --duration <sec> --outdir <dir>")
		fmt.Println("  web --in <findings.json> --outdir <dir> [--insights <insights.json>]")
		fmt.Println("  websocket --findings <findings.json> [--insights <insights.json>] [--port <port>] [--data-dir <dir>] [--compression]")
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
	flagSet.Parse(os.Args[2:])

	if *findingsPath == "" {
		fmt.Println("Required flag: --findings")
		os.Exit(1)
	}

	// Configure WebSocket server
	pipeline.WithWebSocketServer(*port, *dataDir, false, *compression)

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


