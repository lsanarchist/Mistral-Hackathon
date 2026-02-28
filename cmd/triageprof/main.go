package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mistral-hackathon/triageprof/internal/core"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: triageprof <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  plugins list")
		fmt.Println("  collect --plugin <name> --target-url <url> --duration <sec> --out <path>")
		fmt.Println("  analyze --in <bundle.json> --out <findings.json> --top <N>")
		fmt.Println("  report --in <findings.json> --out <report.md>")
		fmt.Println("  run --plugin <name> --target-url <url> --duration <sec> --outdir <dir>")
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
	case "run":
		runRunCommand(pipeline)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func runPluginsCommand() {
	fmt.Println("Plugins command not yet implemented")
	os.Exit(1)
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

func runReportCommand(pipeline *core.Pipeline) {
	flagSet := flag.NewFlagSet("report", flag.ExitOnError)
	inPath := flagSet.String("in", "", "Input findings path")
	outPath := flagSet.String("out", "", "Output report path")
	flagSet.Parse(os.Args[2:])

	if *inPath == "" || *outPath == "" {
		fmt.Println("Required flags: --in, --out")
		os.Exit(1)
	}

	ctx := context.Background()
	err := pipeline.Report(ctx, *inPath, *outPath)
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
	flagSet.Parse(os.Args[2:])

	if *pluginName == "" || *targetURL == "" || *outDir == "" {
		fmt.Println("Required flags: --plugin, --target-url, --outdir")
		os.Exit(1)
	}

	ctx := context.Background()
	err := pipeline.Run(ctx, *pluginName, *targetURL, *duration, 20, *outDir)
	if err != nil {
		fmt.Printf("Run failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results saved to: %s/\n", *outDir)
}