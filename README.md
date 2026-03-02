# TriageProf - AI-Powered Go Performance Profiling

[![Go](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org)
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
./bin/triageprof plugins list
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
| `--llm-timeout` | `20s` | API timeout in seconds |

## 🎯 Commands

```bash
# Full pipeline: collect → analyze → AI enrich → HTML report
./bin/triageprof run --plugin go-pprof-http --target-url <url> --outdir <dir> [--llm]

# Just collect profiles into a bundle
./bin/triageprof collect --plugin go-pprof-http --target-url <url> --out bundle.json

# Analyze an existing bundle
./bin/triageprof analyze --in bundle.json --out findings.json

# Generate LLM insights from findings
./bin/triageprof llm --findings findings.json --out insights.json

# List plugins
./bin/triageprof plugins list
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

---

## Quick Start

```bash
# 1. Build
make build

# 2. Run the built-in demo (no API key needed)
make demo

# 3. Or run with Mistral AI insights
export MISTRAL_API_KEY="your-key-here"
./bin/triageprof run \
  --plugin go-pprof-http \
  --target-url http://localhost:6060 \
  --duration 15 \
  --outdir ./analysis \
  --llm

# Open the report
open analysis/report.html
```

---

## How it works

```
Live Go service (pprof)
        │
        ▼
  [ Collect profiles ]   CPU · heap · allocs · mutex · block
        │
        ▼
  [ Deterministic analysis ]   8+ rule-based patterns, scored findings
        │
        ▼
  [ Mistral AI enrichment ]   mistral-large-latest
        │  → root causes, fix suggestions, code examples
        │  → effort estimates, complexity, validation metrics
        │  → before/after impact predictions
        ▼
  [ Self-contained HTML report ]   charts · per-finding AI cards · recommendations
```

---

## Output

A single `report.html` file (~300KB, no server needed):

- **Score gauge** — overall health at a glance
- **Severity breakdown** — chart + filter by critical / high / medium / low
- **AI Executive Summary** — Mistral's overall verdict with confidence score
- **Top Risks** — what will blow up under load
- **Recommendations** — prioritized actions with effort, complexity, code examples, and how-to-validate
- **Per-finding cards** — expandable with hotspot bars, root causes ↔ suggestions, next measurements ↔ caveats, before/after metrics

---

## Requirements

- Go 1.21+
- A Go service exposing `net/http/pprof` (e.g. `import _ "net/http/pprof"`)
- `MISTRAL_API_KEY` for AI enrichment (optional, but recommended)

---

## Commands

```bash
# Full run (collect → analyze → AI enrich → HTML report)
./bin/triageprof run --plugin go-pprof-http --target-url <url> --outdir <dir> [--llm]

# Just collect profiles
./bin/triageprof collect --plugin go-pprof-http --target-url <url> --out bundle.json

# Just analyze an existing bundle
./bin/triageprof analyze --in bundle.json --out findings.json

# Built-in demo (starts demo server automatically)
make demo
```

### LLM flags

| Flag | Default | Description |
|---|---|---|
| `--llm` | off | Enable Mistral AI enrichment |
| `--llm-model` | `mistral-large-latest` | Model to use |
| `--llm-timeout` | `20s` | API timeout |
| `--llm-provider` | `mistral` | `mistral` or `openai` |

---

## Plugin System

Profilers are separate executables communicating over JSON-RPC. The `go-pprof-http` plugin is included and works out of the box for any Go service with pprof enabled.

```bash
./bin/triageprof plugins   # list available plugins
```

---

## Built for the Mistral Hackathon

This project demonstrates grounded AI analysis — Mistral only adds *why/how*, never invents numbers. All findings are backed by real pprof data. The LLM references specific function names and hotspot percentages from the deterministic analysis.

---

## License

MIT
