# TriageProf

**One command → AI-powered Go performance report.**

TriageProf collects pprof profiles from a live Go service, runs deterministic bottleneck analysis, then enriches findings with [Mistral AI](https://mistral.ai) to produce a beautiful, self-contained HTML report explaining *why* things are slow and *exactly* what to fix.

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
