Your mistral key is at apikey.swaga
Final Goal
Build a modular, language-agnostic profiling and bottleneck-analysis library that integrates with existing profilers (e.g., Go pprof, Linux perf/eBPF, Java async-profiler/JFR, etc.) and uses a combination of heuristics and AI-assisted analysis to identify performance bottlenecks, explain likely root causes with evidence (stacks, call graphs, timelines, metrics), and produce actionable reports.

The system must be extensible via a stable plugin/adapters interface so external contributors can add new languages/profilers (e.g., Java, C/C++) without changing the core. The core owns: a unified profiling schema, normalization pipeline, analyzers, reporting, and plugin lifecycle. Language/profiler-specific logic must live in plugins.


> NOTE: This is a deep reference document.
> Agents must read AGENTS.md and COMPASS.md first.
> Update this file when major architecture/protocol/schema decisions change.
# TriageProf - Plugin-based Profiling Triage Tool

A Go-based tool for collecting, analyzing, and reporting performance profiles from various sources using a plugin architecture.

## Features

- **Plugin Architecture**: Support for multiple profiler plugins via JSON-RPC
- **Go pprof Support**: Built-in plugin for Go HTTP pprof endpoints
- **Analysis Pipeline**: Collect → Analyze → Report workflow
- **Deterministic Analysis**: Rule-based analysis without LLM dependencies
- **Markdown Reports**: Professional performance reports

## Quick Start

### Build

```bash
make build
```

This builds:
- `bin/triageprof` - Main CLI tool
- `plugins/bin/go-pprof-http` - Go pprof plugin

### Run Demo

```bash
make demo
```

This will:
1. Start a demo server with performance issues
2. Generate load to create profiles
3. Run the full triage pipeline
4. Generate a report in `out/report.md`

### Manual Usage

```bash
# Start demo server
go run examples/demo-server/main.go

# Generate load (in another terminal)
./examples/load.sh

# Run full pipeline
bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 10 --outdir results/

# Individual steps
bin/triageprof collect --plugin go-pprof-http --target-url http://localhost:6060 --duration 10 --out results/bundle.json
bin/triageprof analyze --in results/bundle.json --out results/findings.json --top 20
bin/triageprof report --in results/findings.json --out results/report.md
```

## Architecture

### Core Components

- **CLI**: `cmd/triageprof/main.go` - Command-line interface
- **Core**: `internal/core/` - Pipeline orchestration
- **Model**: `internal/model/` - Data structures and schemas
- **Plugin**: `internal/plugin/` - Plugin management and JSON-RPC
- **Analyzer**: `internal/analyzer/` - Profile analysis
- **Reporter**: `internal/report/` - Markdown report generation

### Plugin System

Plugins are separate executables that communicate via JSON-RPC 2.0 over stdin/stdout.

**Plugin Protocol Methods:**
- `rpc.info` - Get plugin metadata
- `rpc.validateTarget` - Validate target configuration
- `rpc.collect` - Collect profiles and return artifact bundle

### Data Flow

```
Target → Plugin → ArtifactBundle → Analyzer → FindingsBundle → Reporter → Markdown Report
```

## Project Structure

```
.
├── bin/                  # Built binaries
├── cmd/triageprof/       # Main CLI
├── internal/             # Core packages
│   ├── analyzer/         # Profile analysis
│   ├── core/             # Pipeline orchestration
│   ├── model/            # Data models
│   ├── plugin/           # Plugin management
│   └── report/           # Report generation
├── plugins/              # Plugin system
│   ├── bin/              # Built plugins
│   ├── manifests/        # Plugin manifests
│   └── src/              # Plugin source code
│       └── go-pprof-http/ # Go pprof plugin
├── examples/             # Demo and examples
│   ├── demo-server/      # Demo server with issues
│   └── load.sh           # Load generation script
├── testdata/             # Test fixtures
└── Makefile              # Build automation
```

## Plugin Development

To create a new plugin:

1. **Create plugin directory**: `plugins/src/<plugin-name>/`
2. **Implement JSON-RPC interface**: Handle `rpc.info`, `rpc.validateTarget`, `rpc.collect`
3. **Add manifest**: `plugins/manifests/<plugin-name>.json`
4. **Build**: `go build -o plugins/bin/<plugin-name> ./plugins/src/<plugin-name>`

## Configuration

### Environment Variables

- `TRIAGEPROF_PLUGINS`: Plugin directory (default: `./plugins`)

### Plugin Manifest Format

```json
{
  "name": "plugin-name",
  "version": "0.1.0",
  "sdkVersion": "1.0",
  "capabilities": {
    "targets": ["url"],
    "profiles": ["cpu", "heap", "mutex", "block", "goroutine"]
  },
  "description": "Plugin description",
  "author": "Author name"
}
```

## Testing

```bash
make test
```

## Clean

```bash
make clean
```

## License

MIT

## Roadmap

- [x] Core pipeline implementation
- [x] Go pprof plugin
- [x] Basic analyzer
- [x] Markdown reporter
- [x] Demo server
- [ ] Additional plugins (Python, Java, etc.)
- [ ] Advanced analysis rules
- [ ] Web UI for reports
- [ ] CI/CD integration

## Contributing

Contributions welcome! Please open issues and pull requests.

## Support

For issues and questions, please open a GitHub issue.
