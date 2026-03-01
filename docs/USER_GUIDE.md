# TriageProf User Guide

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Commands](#core-commands)
- [Demo Workflow](#demo-workflow)
- [Performance Analysis](#performance-analysis)
- [LLM Integration](#llm-integration)
- [Report Interpretation](#report-interpretation)
- [Troubleshooting](#troubleshooting)
- [Advanced Usage](#advanced-usage)

## Introduction

TriageProf is a production-grade profiling and bottleneck-analysis tool for Go applications. It integrates with existing profilers via a plugin SDK, produces evidence-backed bottleneck findings, and generates structured, machine-readable reports with optional AI-powered insights.

## Installation

### Prerequisites

- Go 1.20+
- Git
- Make
- Optional: Mistral/OpenAI API key for AI insights

### Build from Source

```bash
# Clone the repository
git clone https://github.com/your-repo/triageprof.git
cd triageprof

# Build the main binary
make build

# Build plugins
make plugins

# Install (optional)
sudo make install
```

### Verify Installation

```bash
./bin/triageprof version
./bin/triageprof plugins list
```

## Quick Start

### Run the Demo

```bash
# Quick demo with built-in demo server
make demo

# This will:
# 1. Start a demo Go server
# 2. Collect CPU, heap, and allocation profiles
# 3. Analyze performance bottlenecks
# 4. Generate reports in out-demo/
```

### Analyze Your Own Repository

```bash
# Analyze a local Go repository
./bin/triageprof demo --repo /path/to/your/repo --out my-analysis/

# Analyze a Git repository
./bin/triageprof demo --repo https://github.com/your/repo.git --out my-analysis/
```

## Core Commands

### `demo` - Complete Analysis Workflow

```bash
triageprof demo --repo <repository> --out <output-directory>
```

**Options:**
- `--repo`: Path to local repository or Git URL
- `--out`: Output directory for reports and profiles
- `--ref`: Git reference (branch, tag, or commit)
- `--llm`: Enable LLM insights (requires API key)
- `--concurrent`: Enable concurrent benchmark execution
- `--max-workers`: Maximum concurrent workers (default: 2)
- `--sampling-rate`: Profile sampling rate (0.1-1.0)

### `demo-kit` - Built-in Demo

```bash
triageprof demo-kit --out <output-directory> --duration <seconds>
```

**Options:**
- `--out`: Output directory (default: out-demo/)
- `--duration`: Benchmark duration in seconds (default: 5)
- `--llm`: Enable LLM insights
- `--concurrent`: Enable concurrent execution

### `plugins` - Plugin Management

```bash
# List available plugins
triageprof plugins list

# Show plugin details
triageprof plugins show <plugin-name>
```

## Demo Workflow

### Step-by-Step Analysis

1. **Clone and Prepare Repository**
   ```bash
   triageprof demo --repo https://github.com/your/repo.git --out analysis/
   ```

2. **Benchmark Detection**
   - Automatically detects Go benchmark functions
   - Supports `Benchmark*` functions in `_test.go` files

3. **Profile Collection**
   - CPU profiles
   - Heap allocation profiles
   - Memory allocation profiles
   - Block contention profiles
   - Mutex contention profiles

4. **Deterministic Analysis**
   - Identifies CPU hotpaths
   - Detects allocation churn
   - Finds JSON hotspots
   - Analyzes string operations
   - Measures GC pressure
   - Identifies mutex contention

5. **Report Generation**
   - `findings.json` - Structured findings
   - `report.md` - Markdown report
   - `report.html` - Interactive HTML report
   - Raw profile files in `profiles/`

## Performance Analysis

### Understanding Findings

Each finding includes:
- **ID**: Unique identifier
- **Title**: Descriptive title
- **Category**: cpu, alloc, heap, gc, mutex, block
- **Severity**: low, medium, high, critical
- **Confidence**: 0.0-1.0 score
- **Evidence**: Supporting data and metrics
- **Deterministic Hints**: Optimization suggestions

### Common Performance Issues

#### CPU Hotpaths
- Functions consuming >70% of CPU time
- Look for algorithmic complexity issues
- Consider caching or memoization

#### Allocation Churn
- High `mallocgc` or `memmove` operations
- Use object pooling or sync.Pool
- Reduce temporary allocations

#### JSON Hotspots
- `encoding/json` in top functions
- Consider jsoniter or easyjson
- Use streaming JSON parsers

#### GC Pressure
- High `runtime.gcBgMarkWorker` time
- Reduce heap allocations
- Use stack allocations where possible

## LLM Integration

### Configuration

```bash
# Set Mistral API key
export MISTRAL_API_KEY="your-api-key"

# Or set OpenAI API key
export OPENAI_API_KEY="your-api-key"
```

### Using LLM Insights

```bash
# Enable LLM for demo
triageprof demo --repo ./myapp --out analysis/ --llm

# Enable LLM for demo-kit
triageprof demo-kit --out analysis/ --llm
```

### LLM Features

- **Executive Summary**: High-level overview
- **Per-Finding Insights**: Detailed analysis for each finding
- **Top Risks**: Prioritized issues
- **Top Actions**: Recommended fixes
- **Code Examples**: Concrete implementation suggestions

## Report Interpretation

### HTML Report Structure

1. **Overview Section**: Summary statistics and charts
2. **Findings Table**: Filterable list of all findings
3. **Severity Charts**: Visual distribution by severity
4. **Category Breakdown**: Analysis by performance category
5. **Detailed Findings**: Expandable cards with evidence
6. **AI Insights**: LLM-generated recommendations (if enabled)

### Markdown Report

```markdown
# Performance Analysis Report

## Summary
- Total Findings: X
- Critical: Y, High: Z, Medium: A, Low: B
- Categories: CPU, Alloc, Heap, GC, Mutex, Block

## Top Findings

### [HIGH] CPU Hotpath in FunctionX
- Severity: high
- Confidence: 0.95
- Evidence: FunctionX consumes 85% of CPU time
- Impact: Significant performance bottleneck
- Recommendation: Optimize algorithm or add caching
```

## Troubleshooting

### Common Issues

#### "No benchmarks found"
- **Cause**: Repository doesn't contain Go benchmark functions
- **Solution**: Add benchmark functions or use `--force` flag

#### "Git not found"
- **Cause**: Git binary not in PATH
- **Solution**: Install Git or use local repository path

#### "Go not found"
- **Cause**: Go toolchain not installed
- **Solution**: Install Go 1.20+ and add to PATH

#### "Profile collection failed"
- **Cause**: Benchmark execution error
- **Solution**: Check benchmark code and dependencies

### Debugging

```bash
# Enable verbose logging
export TRIAGEPROF_DEBUG=1

# Check run manifest for errors
cat analysis/run.json

# Verify profiles were collected
ls -la analysis/profiles/
```

## Advanced Usage

### Performance Optimization

```bash
# Concurrent benchmark execution
triageprof demo --repo ./myapp --out analysis/ --concurrent --max-workers 4

# Profile sampling for large codebases
triageprof demo --repo ./myapp --out analysis/ --sampling-rate 0.5

# Memory optimization
triageprof demo --repo ./myapp --out analysis/ --memory-optimization
```

### Custom Analysis

```bash
# Analyze specific profiles
triageprof analyze --profiles ./profiles/ --out analysis/

# Generate reports from existing findings
triageprof report --findings ./findings.json --out analysis/
```

### Plugin Development

```bash
# Create new plugin
triageprof plugin init my-plugin

# Build plugin
cd plugins/src/my-plugin
make build

# Test plugin
triageprof plugin test my-plugin
```

## Best Practices

1. **Start with Deterministic Analysis**: Run without LLM first
2. **Focus on High-Severity Findings**: Prioritize critical issues
3. **Use Sampling for Large Projects**: Reduce profile size
4. **Enable Concurrent Execution**: Faster benchmark runs
5. **Review Evidence Carefully**: Understand the root cause
6. **Test Optimizations**: Verify fixes with before/after analysis
