package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		fmt.Println("  analyze --in <bundle.json> --out <findings.json> --top <N>")
		fmt.Println("  report --in <findings.json> --out <report.md>")
		fmt.Println("  llm --bundle <path> --findings <path> --out <insights.json>")
		fmt.Println("  run --plugin <name> --target-url <url> --duration <sec> --outdir <dir>")
		fmt.Println("\nLLM Options for 'run' command:")
		fmt.Println("  --llm (enable LLM insights)")
		fmt.Println("  --llm-model <model> (default: devstral-small-latest)")
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
	targetURL := flagSet.String("target-url", "", "Target URL")
	duration := flagSet.Int("duration", 15, "Duration in seconds")
	outPath := flagSet.String("out", "", "Output bundle path")
	flagSet.Parse(os.Args[2:])

	if *pluginName == "" || *targetURL == "" || *outPath == "" {
		fmt.Println("Required flags: --plugin, --target-url, --out")
		os.Exit(1)
	}

	ctx := context.Background()
	_, err := pipeline.Collect(ctx, *pluginName, *targetURL, *duration, 20, filepath.Dir(*outPath))
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
	flagSet.Parse(os.Args[2:])

	if *inPath == "" || *outPath == "" {
		fmt.Println("Required flags: --in, --out")
		os.Exit(1)
	}

	ctx := context.Background()
	_, err := pipeline.Analyze(ctx, *inPath, *topN, *outPath)
	if err != nil {
		fmt.Printf("Analyze failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Findings saved to: %s\n", *outPath)
}

func runLLMCommand() {
	flagSet := flag.NewFlagSet("llm", flag.ExitOnError)
	bundlePath := flagSet.String("bundle", "", "Input bundle path")
	findingsPath := flagSet.String("findings", "", "Input findings path")
	outPath := flagSet.String("out", "", "Output insights path")
	model := flagSet.String("model", "devstral-small-latest", "Mistral model name")
	timeout := flagSet.Int("timeout", 20, "API timeout in seconds")
	maxResponse := flagSet.Int("max-response", 4096, "Max response tokens")
	maxPromptChars := flagSet.Int("max-prompt-chars", 12000, "Max prompt characters")
	dryRun := flagSet.Bool("dry-run", false, "Dry run - save prompt without API call")
	flagSet.Parse(os.Args[2:])

	if *bundlePath == "" || *findingsPath == "" || *outPath == "" {
		fmt.Println("Required flags: --bundle, --findings, --out")
		os.Exit(1)
	}

	ctx := context.Background()
	apiKey := os.Getenv("MISTRAL_API_KEY")

	err := llm.GenerateInsightsFromFiles(ctx, *bundlePath, *findingsPath, *outPath, 
		apiKey, *model, *timeout, *maxResponse, *maxPromptChars, *dryRun)
	if err != nil {
		fmt.Printf("LLM insights generation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Insights saved to: %s\n", *outPath)
}

func runReportCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("report", flag.ExitOnError)
	inPath := flagSet.String("in", "", "Input findings path")
	outPath := flagSet.String("out", "", "Output report path")
	insightsPath := flagSet.String("insights", "", "Optional insights path")
	flagSet.Parse(os.Args[2:])

	if *inPath == "" || *outPath == "" {
		fmt.Println("Required flags: --in, --out")
		os.Exit(1)
	}

	ctx := context.Background()
	
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

	err := pipeline.ReportWithInsights(ctx, *inPath, insights, *outPath)
	if err != nil {
		fmt.Printf("Report failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Report saved to: %s\n", *outPath)
}

func runRunCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("run", flag.ExitOnError)
	pluginName := flagSet.String("plugin", "", "Plugin name")
	targetURL := flagSet.String("target-url", "", "Target URL")
	duration := flagSet.Int("duration", 15, "Duration in seconds")
	outDir := flagSet.String("outdir", "", "Output directory")
	llmEnabled := flagSet.Bool("llm", false, "Enable LLM insights")
	llmModel := flagSet.String("llm-model", "devstral-small-latest", "Mistral model name")
	llmTimeout := flagSet.Int("llm-timeout", 20, "LLM API timeout in seconds")
	llmMaxChars := flagSet.Int("llm-max-chars", 12000, "Max prompt characters")
	llmDryRun := flagSet.Bool("llm-dry-run", false, "Dry run - save prompt without API call")
	flagSet.Parse(os.Args[2:])

	if *pluginName == "" || *targetURL == "" || *outDir == "" {
		fmt.Println("Required flags: --plugin, --target-url, --outdir")
		os.Exit(1)
	}

	ctx := context.Background()

	// Configure LLM if enabled
	if *llmEnabled {
		apiKey := os.Getenv("MISTRAL_API_KEY")
		pipeline.WithLLM(apiKey, *llmModel, *llmTimeout, 4096, *llmMaxChars, *llmDryRun)
	}

	err := pipeline.Run(ctx, *pluginName, *targetURL, *duration, 20, *outDir)
	if err != nil {
		fmt.Printf("Run failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results saved to: %s/\n", *outDir)
}