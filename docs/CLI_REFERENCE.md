# TriageProf CLI Reference

## Table of Contents

- [Command Structure](#command-structure)
- [Global Options](#global-options)
- [Commands](#commands)
  - [`demo`](#demo)
  - [`demo-kit`](#demo-kit)
  - [`plugins`](#plugins)
  - [`analyze`](#analyze)
  - [`report`](#report)
  - [`version`](#version)
- [Examples](#examples)

## Command Structure

```bash
triageprof <command> [options] [arguments]
```

## Global Options

| Option | Description | Default |
|--------|-------------|---------|
| `--debug` | Enable debug logging | false |
| `--config` | Config file path | `~/.triageprof/config.json` |
| `--llm` | Enable LLM insights | false |
| `--llm-provider` | LLM provider (mistral/openai) | mistral |
| `--cache-dir` | Cache directory | `~/.triageprof/cache` |

## Commands

### `demo`

Complete performance analysis workflow for Go repositories.

**Usage:**
```bash
triageprof demo --repo <repository> --out <output-directory> [options]
```

**Options:**

| Option | Description | Default |
|--------|-------------|---------|
| `--repo` | Repository path or Git URL | Required |
| `--out` | Output directory | Required |
| `--ref` | Git reference (branch/tag/commit) | `main` |
| `--duration` | Benchmark duration in seconds | 10 |
| `--concurrent` | Enable concurrent benchmark execution | false |
| `--max-workers` | Maximum concurrent workers | 2 |
| `--sampling-rate` | Profile sampling rate (0.1-1.0) | 1.0 |
| `--memory-optimization` | Enable memory optimization | false |
| `--large-codebase` | Optimize for large codebases | false |
| `--force` | Force analysis even if no benchmarks found | false |
| `--no-cleanup` | Keep temporary files | false |

**Examples:**

```bash
# Analyze local repository
triageprof demo --repo ./myapp --out analysis/

# Analyze Git repository
triageprof demo --repo https://github.com/user/repo.git --out analysis/

# With concurrent execution
triageprof demo --repo ./myapp --out analysis/ --concurrent --max-workers 4

# With sampling for large codebase
triageprof demo --repo ./myapp --out analysis/ --sampling-rate 0.5 --large-codebase
```

### `demo-kit`

Run built-in demo with sample benchmarks.

**Usage:**
```bash
triageprof demo-kit --out <output-directory> [options]
```

**Options:**

| Option | Description | Default |
|--------|-------------|---------|
| `--out` | Output directory | `out-demo/` |
| `--duration` | Benchmark duration in seconds | 5 |
| `--concurrent` | Enable concurrent execution | false |
| `--max-workers` | Maximum concurrent workers | 2 |
| `--sampling-rate` | Profile sampling rate | 1.0 |

**Examples:**

```bash
# Quick demo
triageprof demo-kit

# Custom output directory
triageprof demo-kit --out my-demo/

# Longer duration with concurrent execution
triageprof demo-kit --duration 15 --concurrent --max-workers 3
```

### `plugins`

Manage and test plugins.

**Usage:**
```bash
triageprof plugins <subcommand> [options]
```

**Subcommands:**

#### `list`

List available plugins.

```bash
triageprof plugins list
```

**Options:**

| Option | Description |
|--------|-------------|
| `--verbose` | Show detailed plugin information |
| `--json` | Output as JSON |

#### `show`

Show plugin details.

```bash
triageprof plugins show <plugin-name>
```

#### `test`

Test a plugin.

```bash
triageprof plugins test <plugin-name> [options]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--method` | Method to test (initialize, collectProfile, shutdown) |
| `--params` | JSON parameters for method |

**Examples:**

```bash
# List plugins
triageprof plugins list

# Show plugin details
triageprof plugins show go-pprof-http

# Test plugin initialization
triageprof plugins test go-pprof-http --method initialize

# Test profile collection
triageprof plugins test go-pprof-http --method collectProfile --params '{"profileType": "cpu"}'
```

### `analyze`

Analyze existing profiles.

**Usage:**
```bash
triageprof analyze --profiles <directory> --out <output-directory> [options]
```

**Options:**

| Option | Description | Default |
|--------|-------------|---------|
| `--profiles` | Directory containing profile files | Required |
| `--out` | Output directory | Required |
| `--findings` | Custom findings file | Auto-generated |
| `--insights` | Custom insights file | Auto-generated |

**Examples:**

```bash
# Analyze profiles
triageprof analyze --profiles ./profiles/ --out analysis/

# With custom findings
triageprof analyze --profiles ./profiles/ --out analysis/ --findings custom-findings.json
```

### `report`

Generate reports from findings.

**Usage:**
```bash
triageprof report --findings <file> --out <output-directory> [options]
```

**Options:**

| Option | Description | Default |
|--------|-------------|---------|
| `--findings` | Findings JSON file | Required |
| `--insights` | Insights JSON file | Optional |
| `--out` | Output directory | Required |
| `--format` | Report format (html, md, both) | both |

**Examples:**

```bash
# Generate HTML report
triageprof report --findings findings.json --out reports/ --format html

# Generate both formats
triageprof report --findings findings.json --insights insights.json --out reports/
```

### `version`

Show version information.

**Usage:**
```bash
triageprof version
```

**Options:**

| Option | Description |
|--------|-------------|
| `--verbose` | Show detailed version info |
| `--json` | Output as JSON |

**Examples:**

```bash
# Show version
triageprof version

# Verbose version info
triageprof version --verbose
```

## Examples

### Complete Workflow

```bash
# 1. Analyze a Go repository
triageprof demo --repo https://github.com/user/repo.git --out analysis/ --llm

# 2. Review findings
cat analysis/findings.json

# 3. Open HTML report
open analysis/report.html

# 4. Check run details
cat analysis/run.json
```

### Advanced Analysis

```bash
# Concurrent execution with sampling
triageprof demo \
  --repo ./myapp \
  --out analysis/ \
  --concurrent \
  --max-workers 4 \
  --sampling-rate 0.7 \
  --llm

# Analyze specific Git reference
triageprof demo \
  --repo https://github.com/user/repo.git \
  --ref v1.2.3 \
  --out analysis/
```

### Plugin Development

```bash
# Test plugin
triageprof plugins test my-plugin --method initialize

# Use custom plugin
triageprof demo --repo ./myapp --plugin my-plugin --out analysis/
```

### Report Generation

```bash
# Generate reports from existing data
triageprof report \
  --findings findings.json \
  --insights insights.json \
  --out reports/ \
  --format both

# HTML report only
triageprof report --findings findings.json --out reports/ --format html
```

## Configuration

### Config File

TriageProf supports a JSON configuration file:

```json
{
  "llm": {
    "enabled": true,
    "provider": "mistral",
    "apiKey": "your-api-key",
    "model": "mistral-small"
  },
  "performance": {
    "concurrent": true,
    "maxWorkers": 4,
    "samplingRate": 0.8
  },
  "output": {
    "format": "both",
    "directory": "analysis"
  }
}
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `TRIAGEPROF_DEBUG` | Enable debug logging |
| `MISTRAL_API_KEY` | Mistral API key |
| `OPENAI_API_KEY` | OpenAI API key |
| `TRIAGEPROF_CONFIG` | Config file path |
| `TRIAGEPROF_CACHE` | Cache directory |

## Output Structure

### Directory Structure

```
output-directory/
├── findings.json          # Structured findings
├── insights.json          # LLM insights (if enabled)
├── report.md              # Markdown report
├── report.html            # HTML report
├── bundle.json            # Complete data bundle
├── run.json               # Run manifest
├── profiles/              # Raw profile files
│   ├── cpu.pb.gz
│   ├── heap.pb.gz
│   ├── allocs.pb.gz
│   ├── block.pb.gz
│   └── mutex.pb.gz
└── web/                    # Web assets (for HTML report)
    ├── report.js
    ├── style.css
    └── chart.js
```

### Key Files

| File | Description |
|------|-------------|
| `findings.json` | Structured performance findings with evidence |
| `insights.json` | LLM-generated insights and recommendations |
| `report.md` | Human-readable Markdown report |
| `report.html` | Interactive HTML report with visualizations |
| `run.json` | Run metadata and configuration |
| `bundle.json` | Complete data bundle for archiving |

## Error Handling

### Error Codes

| Code | Meaning |
|------|---------|
| 1000 | Invalid input parameters |
| 1001 | Repository not found |
| 1002 | No benchmarks found |
| 1003 | Profile collection failed |
| 1004 | Analysis failed |
| 1005 | Report generation failed |
| 2000 | Plugin error |
| 2001 | Plugin not found |
| 2002 | Plugin initialization failed |
| 3000 | LLM error |
| 3001 | LLM API key missing |
| 3002 | LLM request failed |

### Error Context

Errors include structured context:

```json
{
  "error": {
    "code": 1002,
    "message": "No benchmarks found",
    "type": "validation",
    "details": "Repository does not contain any Go benchmark functions",
    "suggestions": [
      "Add benchmark functions to your test files",
      "Use --force flag to proceed without benchmarks"
    ],
    "isRecoverable": true
  }
}
```

## Best Practices

### Command Usage

1. **Start Simple**: Begin with basic analysis before adding complexity
2. **Use Concurrent Execution**: For faster benchmark runs on multi-core systems
3. **Enable Sampling**: For large codebases to reduce profile size
4. **Review Findings First**: Understand deterministic analysis before enabling LLM
5. **Check Run Manifest**: For troubleshooting and debugging information

### Performance Optimization

```bash
# For large repositories
triageprof demo \
  --repo ./large-app \
  --out analysis/ \
  --sampling-rate 0.5 \
  --large-codebase \
  --memory-optimization

# For CI/CD pipelines
triageprof demo \
  --repo . \
  --out analysis/ \
  --concurrent \
  --max-workers 2 \
  --duration 5
```

### Debugging

```bash
# Enable debug logging
export TRIAGEPROF_DEBUG=1
triageprof demo --repo ./myapp --out analysis/

# Check run details
cat analysis/run.json

# Verify profiles
ls -la analysis/profiles/
```
