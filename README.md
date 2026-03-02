# TriageProf - AI-Powered Performance Profiling

[![Go](https://img.shields.io/badge/Go-1.24%2B-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Mistral AI](https://img.shields.io/badge/Mistral-AI%20Powered-orange)](https://mistral.ai)

**TriageProf** is a modular profiling and bottleneck-analysis tool for Go applications. It collects pprof profiles from a live service, runs deterministic analysis, then enriches findings with **Mistral AI** to explain *why* things are slow and *exactly* what to fix — producing a beautiful, self-contained HTML report.

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/lsanarchist/Mistral-Hackathon.git
cd Mistral-Hackathon

# Build the binary and plugins (single step)
make build

# List available plugins
./bin/triageprof plugins
```

### Run Against Your Own Service

```bash
# Profile a live Go service and get an AI-powered report
export MISTRAL_API_KEY="your-key-here"

./bin/triageprof run \
  --plugin go-pprof-http \
  --target-url http://localhost:6060 \
  --duration 15 \
  --outdir ./analysis \
  --llm

# After the run, triageprof will offer to open the report in your browser automatically
```

## ✨ Features

### Core Capabilities

- **🔍 Comprehensive Profiling**: CPU, heap, allocation, block, and mutex profiling
- **🤖 Deterministic Analysis**: 8+ rule-based patterns, no LLM required for findings
- **🧠 Mistral AI Enrichment**: `mistral-large-latest` explains root causes, suggests fixes, estimates effort
- **📊 Self-contained HTML Report**: Interactive dark-theme report — charts, hotspot bars, per-finding AI cards
- **🌐 Auto Browser Serve**: After each run, offers to serve the report and open it in your browser
- **🔌 Plugin Architecture**: Extensible via JSON-RPC plugin SDK

### Analysis Rules

1. **CPU Hotpath Dominance**: Detects functions consuming CPU time
2. **Allocation Churn**: Identifies high mallocgc/memmove patterns
3. **JSON Hotspots**: Finds encoding/json bottlenecks
4. **String Churn**: Analyzes strings.Builder/bytes.Buffer usage
5. **GC Pressure**: Measures runtime.gcBgMarkWorker impact
6. **Mutex Contention**: Detects sync.(*Mutex).Lock issues
7. **Heap Allocation**: Identifies memory allocation hotspots
8. **Block Contention**: Analyzes runtime.chan/select patterns

### Report Contents

- **Score gauge** — overall health at a glance
- **Severity breakdown** — donut chart + filter by critical / high / medium / low
- **AI Executive Summary** — Mistral's verdict with confidence score and key themes
- **Top Risks** — what will blow up under load
- **Recommendations** — prioritized actions with effort, complexity, code examples, validation metrics
- **Per-finding AI cards** — root causes ↔ suggestions, next measurements ↔ caveats, before/after metrics

## 🔧 Configuration

### LLM Options

```bash
export MISTRAL_API_KEY="your-api-key"

./bin/triageprof run \
  --plugin go-pprof-http \
  --target-url http://localhost:6060 \
  --outdir ./analysis \
  --llm \
  --llm-model mistral-large-latest \
  --llm-timeout 90
```

| Flag | Default | Description |
|---|---|---|
| `--llm` | off | Enable Mistral AI enrichment |
| `--llm-model` | `mistral-large-latest` | Model to use |
| `--llm-provider` | `mistral` | `mistral` or `openai` |
| `--llm-timeout` | `20` | API timeout in seconds |

## 🎯 Commands

```bash
# Full pipeline: collect → analyze → AI enrich → HTML report
./bin/triageprof run --plugin go-pprof-http --target-url <url> --outdir <dir> [--llm]

# Just collect profiles into a bundle
./bin/triageprof collect --plugin go-pprof-http --target-url <url> --out bundle.json

# Analyze an existing bundle
./bin/triageprof analyze --in bundle.json --out findings.json

# Generate LLM insights from findings
./bin/triageprof llm --bundle bundle.json --findings findings.json --out insights.json

# List plugins
./bin/triageprof plugins
```

## 📦 Output Structure

```
analysis/
├── report.html            # Self-contained interactive HTML report (~300KB)
├── index.html             # Alias for report.html (same content)
├── findings.json          # Structured performance findings
├── insights.json          # Mistral AI analysis (if --llm enabled)
├── report.md              # Markdown report
└── bundle.json            # Raw collected profiles bundle
```

## 🔌 Plugin System

Profilers are separate executables communicating via JSON-RPC over stdio. The `go-pprof-http` plugin is included and works out of the box for any Go service with `net/http/pprof` enabled.

```bash
# List available plugins
./bin/triageprof plugins list

# Your service just needs this import:
import _ "net/http/pprof"
```

See **[API Documentation](docs/API_DOCUMENTATION.md)** for the plugin development guide.

## 📚 Documentation

- **[User Guide](docs/USER_GUIDE.md)**: Complete usage guide with examples
- **[CLI Reference](docs/CLI_REFERENCE.md)**: Detailed command reference
- **[API Documentation](docs/API_DOCUMENTATION.md)**: Plugin SDK and JSON-RPC API
- **[Contributing Guide](docs/CONTRIBUTING.md)**: Development and contribution guidelines

## 🤝 Contributing

- **Issues**: Report bugs and request features
- **Pull Requests**: Contributions welcome!

See **[CONTRIBUTING.md](docs/CONTRIBUTING.md)** for guidelines.

## 📜 License

TriageProf is licensed under the **[MIT License](LICENSE)**.

---

**TriageProf** — built for the [Mistral AI Hackathon](https://mistral.ai) 🚀
