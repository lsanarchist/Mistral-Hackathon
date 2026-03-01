# TriageProf - Go Profiling Triage Tool

A focused Go-based tool for collecting, analyzing, and reporting performance profiles with optional AI-powered insights.

## Quick Start

### Build

```bash
make build
```

### Run Demo

```bash
make demo
```

This runs a comprehensive demo that:
- Collects profiles from a Go demo server
- Analyzes performance bottlenecks
- Generates professional reports
- Optionally adds AI insights (if API key configured)

## Features

- **Go pprof Support**: Built-in plugin for Go HTTP pprof endpoints
- **Deterministic Analysis**: Rule-based analysis without LLM dependencies
- **Optional AI Insights**: Mistral/OpenAI integration for enhanced analysis
- **Professional Reports**: Markdown and HTML reports with actionable insights

## LLM Configuration

LLM is **disabled by default** for safety. To enable:

```bash
# Set your Mistral API key
export MISTRAL_API_KEY="your-api-key"

# Run with LLM insights
triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --llm
```

## Development

```bash
# Run tests
go test ./...

# Build and install
make build
sudo make install
```

## Documentation

- [COMPASS.md](COMPASS.md) - Project direction and architecture
- [AGENTS.md](AGENTS.md) - Development guidelines
- [change.log](change.log) - Recent changes and roadmap
