# TriageProf - Production-Grade Performance Profiling

[![Go](https://img.shields.io/badge/Go-1.20%2B-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/triageprof/triageprof)

**TriageProf** is a production-grade, modular profiling and bottleneck-analysis tool for Go applications. It integrates with existing profilers via a well-defined plugin SDK, produces evidence-backed bottleneck findings, and generates structured, machine-readable reports with optional AI-powered insights.

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/triageprof/triageprof.git
cd triageprof

# Build the main binary and plugins
make build
make plugins

# Verify installation
./bin/triageprof version
./bin/triageprof plugins list
```

### Run the Demo

```bash
# Quick demo with built-in demo server
make demo

# This will:
# 1. Start a demo Go server with realistic benchmarks
# 2. Collect CPU, heap, and allocation profiles
# 3. Analyze performance bottlenecks deterministically
# 4. Generate professional HTML and Markdown reports
# 5. Optionally add AI insights (if API key configured)
```

### Analyze Your Own Repository

```bash
# Analyze a local Go repository
./bin/triageprof demo --repo /path/to/your/repo --out my-analysis/

# Analyze a Git repository
./bin/triageprof demo --repo https://github.com/your/repo.git --out my-analysis/

# With LLM insights (requires API key)
./bin/triageprof demo --repo ./myapp --out analysis/ --llm
```

## ✨ Features

### Core Capabilities

- **🔍 Comprehensive Profiling**: CPU, heap, allocation, block, and mutex profiling
- **🤖 Deterministic Analysis**: 8+ rule-based analysis patterns without LLM dependencies
- **🧠 Optional AI Insights**: Mistral/OpenAI integration for enhanced analysis
- **📊 Professional Reports**: Interactive HTML and structured Markdown reports
- **⚡ Performance Optimized**: Concurrent execution, sampling, and memory optimization
- **🔌 Plugin Architecture**: Extensible via JSON-RPC plugin SDK

### Analysis Rules

1. **CPU Hotpath Dominance**: Detects functions consuming >70% CPU time
2. **Allocation Churn**: Identifies high mallocgc/memmove patterns
3. **JSON Hotspots**: Finds encoding/json bottlenecks
4. **String Churn**: Analyzes strings.Builder/bytes.Buffer usage
5. **GC Pressure**: Measures runtime.gcBgMarkWorker impact
6. **Mutex Contention**: Detects sync.(*Mutex).Lock issues
7. **Heap Allocation**: Identifies memory allocation hotspots
8. **Block Contention**: Analyzes runtime.chan/select patterns

### Report Formats

- **HTML Report**: Interactive web interface with charts and filtering
- **Markdown Report**: Structured text report for documentation
- **JSON Findings**: Machine-readable findings with evidence
- **Raw Profiles**: Original pprof files for manual analysis

## 📚 Documentation

### User Guide

- **[User Guide](docs/USER_GUIDE.md)**: Complete usage guide with examples
- **[CLI Reference](docs/CLI_REFERENCE.md)**: Detailed command reference
- **[API Documentation](docs/API_DOCUMENTATION.md)**: Plugin SDK and JSON-RPC API
- **[Contributing Guide](docs/CONTRIBUTING.md)**: Development and contribution guidelines

### Core Concepts

- **[COMPASS.md](COMPASS.md)**: Project direction, architecture, and roadmap
- **[AGENTS.md](AGENTS.md)**: Development guidelines and feature policy
- **[change.log](change.log)**: Recent changes and implementation details

## 🔧 Configuration

### LLM Configuration

LLM is **disabled by default** for safety and privacy. To enable:

```bash
# Set your Mistral API key
export MISTRAL_API_KEY="your-api-key"

# Or set OpenAI API key
export OPENAI_API_KEY="your-api-key"

# Run with LLM insights
triageprof demo --repo ./myapp --out analysis/ --llm
```

### Performance Optimization

```bash
# Concurrent benchmark execution
triageprof demo --repo ./myapp --out analysis/ --concurrent --max-workers 4

# Profile sampling for large codebases
triageprof demo --repo ./myapp --out analysis/ --sampling-rate 0.5

# Memory optimization
triageprof demo --repo ./myapp --out analysis/ --memory-optimization
```

## 🎯 Use Cases

### Performance Triage
```bash
# Quick analysis of performance issues
triageprof demo --repo ./myapp --out analysis/

# Review findings
cat analysis/findings.json

# Open interactive report
open analysis/report.html
```

### CI/CD Integration
```bash
# Add to your CI pipeline
triageprof demo --repo . --out analysis/ --concurrent --duration 5

# Fail on critical findings
if [ $(jq '.findings | length' analysis/findings.json) -gt 0 ]; then
  echo "Performance issues detected!"
  exit 1
fi
```

### Plugin Development
```bash
# Create a new plugin
triageprof plugin init my-plugin

# Test your plugin
triageprof plugins test my-plugin

# Use your plugin
triageprof demo --repo ./myapp --plugin my-plugin --out analysis/
```

## 📦 Output Structure

```
output-directory/
├── findings.json          # Structured performance findings
├── insights.json          # LLM insights (if enabled)
├── report.md              # Markdown report
├── report.html            # Interactive HTML report
├── bundle.json            # Complete data bundle
├── run.json               # Run metadata and configuration
├── profiles/              # Raw profile files
│   ├── cpu.pb.gz
│   ├── heap.pb.gz
│   ├── allocs.pb.gz
│   ├── block.pb.gz
│   └── mutex.pb.gz
└── web/                    # Web assets for HTML report
    ├── report.js
    ├── style.css
    └── chart.js
```

## 🔌 Plugin System

TriageProf uses a plugin architecture where profilers are separate executables that communicate via JSON-RPC over stdio.

### Built-in Plugins

- **go-pprof-http**: Go HTTP pprof endpoint profiler
- **node-inspector**: Node.js profiling (archived)
- **python-cprofile**: Python cProfile support (archived)
- **ruby-stackprof**: Ruby stackprof integration (archived)

### Creating Plugins

See **[API Documentation](docs/API_DOCUMENTATION.md)** for complete plugin development guide.

## 🤝 Community

- **Issues**: Report bugs and request features
- **Discussions**: Ask questions and share ideas
- **Contributions**: Pull requests welcome!
- **Documentation**: Help improve our docs

See **[CONTRIBUTING.md](docs/CONTRIBUTING.md)** for contribution guidelines.

## 📜 License

TriageProf is licensed under the **[MIT License](LICENSE)**.

## 🎯 Roadmap

Check **[suggested_next_steps.md](suggested_next_steps.md)** for upcoming features and priorities.

## 💡 Quick Examples

```bash
# Basic analysis
./bin/triageprof demo --repo ./myapp --out analysis/

# With LLM insights
export MISTRAL_API_KEY="your-key"
./bin/triageprof demo --repo ./myapp --out analysis/ --llm

# Concurrent execution
./bin/triageprof demo --repo ./myapp --out analysis/ --concurrent --max-workers 4

# Built-in demo kit
./bin/triageprof demo-kit --out my-demo/ --duration 10

# List available plugins
./bin/triageprof plugins list

# Test a plugin
./bin/triageprof plugins test go-pprof-http --method initialize
```

## 🚀 Getting Help

- **Documentation**: Check our comprehensive guides in `docs/`
- **Issues**: Search existing issues or create new ones
- **Discussions**: Ask questions in GitHub Discussions
- **Community**: Join our growing community of contributors

---

**TriageProf** - The `osv-scanner` of performance analysis! 🚀
