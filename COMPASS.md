# COMPASS — TriageProf Direction & Status

## IMPORTANT RULE: THIS FILE MUST NOT BE MODIFIED DIRECTLY

**All changes must be appended to `change.log` instead.**

This document establishes the foundational principles and direction for TriageProf.
To propose changes, add entries to the change.log file following the format below.

---

## North Star

Build a production-grade, modular profiling and bottleneck-analysis tool — the `osv-scanner` of performance.
It integrates with existing profilers via a well-defined plugin SDK, produces evidence-backed bottleneck findings
(stacks, callgraphs, timelines, metrics), and generates structured, machine-readable reports.

**Killer feature: AI/LLM-powered analysis** — deterministic profiling data feeds an LLM that explains *why* a bottleneck exists, suggests concrete fixes with code examples, and ranks issues by impact.

**Philosophy: go deep, not wide. Make each layer excellent before adding the next.**

## Product Shape (today)

Pipeline: Collect → Analyze → Report.
Plugins are separate executables, discovered via manifests, and communicate with core using JSON-RPC 2.0 over stdio.

## Non-negotiable Architecture Rules

- Core is language-agnostic; language/profiler-specific logic must live in plugins.
- Plugin protocol/API must remain stable; breaking changes require explicit versioning.
- Deterministic profiling data is always collected first and is the source of truth.
- **AI/LLM analysis is a first-class feature**, not optional glue — it must be designed as a proper pipeline stage with a clean interface, not bolted on. apikey for dev is at apikey.swaga
- LLM calls are always backed by structured profiling data; the LLM never guesses — it reasons over real evidence.
- LLM stage must be skippable via `--no-ai` flag so the tool works without an API key.
- **Depth over breadth**: do not add new plugins unless the existing ones are excellent and the SDK is solid.
- Every feature must have tests. No untested code merges.
- Big PRs are fine when a feature warrants it. Correctness and completeness over minimal diffs.

## Change Management Process

### How to Propose Changes

1. **Do NOT modify this file directly**
2. Create or append to `change.log` using the format below
3. Changes will be reviewed and incorporated by maintainers

### Change Log Format

```
## Change Log

### YYYY-MM-DD
- **Objective**: Brief description of the change
- **Rationale**: Why this change is needed
- **Impact**: What will be affected
- **Changes**:
  - Specific change 1
  - Specific change 2
- **Testing**: How this will be validated
```

### Example Change Entry

```
## Change Log

### 2026-02-28
- **Objective**: Add plugin capability validation
- **Rationale**: Prevent incompatible plugins from being loaded
- **Impact**: Plugin loading process, error handling
- **Changes**:
  - Add manifest-based plugin discovery
  - Implement capability validation before plugin launch
  - Add SDK version compatibility checking
  - Enhance error messages for plugin issues
- **Testing**: Unit tests for manifest parsing and validation
```

## Current Focus Areas

1. **Plugin System Maturity**
   - Manifest-based discovery and validation 
   - Capability checking and compatibility
   - Plugin lifecycle management

2. **LLM Integration** !!!!!
   - Safe prompt generation with redaction 
   - Structured insights generation
   - Report enhancement with LLM sections

3. **Core Robustness**
   - Error handling and recovery
   - Test coverage and validation
   - Performance optimization

4. **Cool web page of results**

5. **Make more clear where each plugin-module lives**

## Decision Log

### 2026-02-28: Plugin Manifest Approach
**Decision**: Use JSON manifests in `plugins/manifests/` for plugin discovery
**Rationale**: Simple, declarative, language-agnostic, easy to validate
**Alternatives Considered**:
- YAML manifests (rejected: more complex parsing)
- Code-based registration (rejected: less flexible)
- Database storage (rejected: overkill for this use case)

### 2026-02-28: LLM Safety Design
**Decision**: Redact sensitive data before sending to LLM
**Rationale**: Security and privacy by default
**Implementation**:
- Remove hostnames, tokens, absolute paths
- Truncate long strings
- Send only derived data, never raw profiles

## Iteration Log

### 2026-02-28 14:00: Initial Implementation
- **Objective**: Implement manifest-based plugin discovery
- **Changes**:
  - Created `internal/plugin/manifest.go` with Manifest struct
  - Added `LoadManifest()` with strict JSON parsing
  - Implemented `DiscoverManifests()` for directory walking
  - Added `ResolvePlugin()` for validation
  - Enhanced PluginManager with manifest support
  - Updated core pipeline to validate plugins before launch
- **Testing**: Unit tests for parsing, discovery, and validation
- **Result**: Plugins now discovered from manifests with capability checking

### 2026-02-28 15:30: LLM Insights Feature
- **Objective**: Add optional LLM augmentation
- **Changes**:
  - Created `internal/model/insights.go` with InsightsBundle
  - Implemented `internal/llm/client.go` for Mistral API
  - Added `internal/llm/prompt.go` with redaction
  - Created `internal/llm/insights.go` for orchestration
  - Enhanced pipeline with optional LLM step
  - Updated reporter to include LLM insights
  - Added CLI commands for LLM control
- **Testing**: Unit tests for client, prompt building, and integration
- **Result**: Optional LLM insights with safety features and proper error handling

### 2026-02-28 16:00: LLM Integration Completion
- **Objective**: Complete LLM integration with core pipeline
- **Changes**:
  - Updated `internal/core/pipeline.go` with LLM generator support
  - Added `WithLLM()` method to configure LLM insights
  - Enhanced `RunWithTarget()` to generate and save insights
  - Updated `ReportWithInsights()` for enhanced report generation
  - Modified CLI to support LLM flags in run command
  - Added LLM command for standalone insights generation
- **Testing**: Integration tests for full pipeline with LLM
- **Result**: End-to-end LLM integration with optional insights generation

---

## How to Contribute

1. **Read this document** to understand the vision and rules
2. **Propose changes** by adding to `change.log`
3. **Discuss** the proposed changes with the team
4. **Implement** the approved changes
5. **Update change.log** with actual results

## Maintenance

- **Do NOT edit this file directly**
- All historical changes preserved in iteration log
- Change.log entries become part of the permanent record
- This ensures traceability and accountability

**Remember**: The compass points the way, but the journey is recorded in the logs.

---

## 🎯 Killer Feature Focus: AI-Powered Performance Triage

### The Vision

**TriageProf's killer feature is transforming raw profiling data into actionable insights using AI.**

While traditional profilers show *what* is slow, TriageProf explains *why* it's slow and *how* to fix it.

### Demo Goal: "Wow" Factor

**Objective**: Create a demo that makes developers say "Wow, I need this!"

**Success Criteria**:
- Input: Complex Go application with real performance issues
- Process: Automatic collection → analysis → LLM insights
- Output: Professional report with:
  - ✅ Executive summary with severity assessment
  - ✅ Top 3 risks with impact analysis
  - ✅ Actionable suggestions with code examples
  - ✅ Confidence scores and caveats

### Demo Flow

```
1. Start demo server with intentional issues
   - CPU hotspots
   - Memory allocation patterns
   - Mutex contention
   - Blocking operations

2. Run TriageProf with LLM enabled
   bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 30 --outdir demo --llm

3. Show the "Wow" moments:
   - ✨ Automatic plugin discovery and validation
   - 🔍 Detailed profile collection (CPU, heap, mutex, block, goroutine)
   - 🤖 LLM-generated insights with root cause analysis
   - 📊 Professional markdown report with executive summary

4. Highlight key differentiators:
   - "Traditional profilers show you the hotspots"
   - "TriageProf explains why they're hot and how to fix them"
   - "From data to insights in one command"

### Cool Demo Features to Showcase

1. **Plugin Discovery Magic**
   - Run `triageprof plugins` to show available plugins
   - Emphasize manifest-based discovery

2. **End-to-End Workflow**
   - Single command collects, analyzes, and reports
   - Optional LLM augmentation for deeper insights

3. **LLM Insights**
   - Executive summary with severity assessment
   - Per-finding narrative explanations
   - Suggested fixes and next measurements

4. **Safety Features**
   - Sensitive data redaction
   - Graceful degradation without API key
   - Test verification before committing

5. **Professional Output**
   - Structured markdown reports
   - LLM insights clearly marked
   - Confidence indicators

### Demo Script Example

```bash
# Start the demo server
go run examples/demo-server/main.go &

# Generate some load
./examples/load.sh &

# Run the full pipeline with LLM
export MISTRAL_API_KEY="your-key"
bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 30 --outdir demo --llm

# Show the results
cat demo/report.md

# Highlight the LLM insights
bin/triageprof report --in demo/findings.json --insights demo/insights.json --out demo/enhanced-report.md
```

### Demo Success Metrics

**Technical**:
- ✅ All profiles collected successfully
- ✅ Analysis completes without errors
- ✅ LLM generates meaningful insights
- ✅ Report includes all expected sections

**User Experience**:
- ✅ "Wow" reaction from viewers
- ✅ Clear understanding of the value proposition
- ✅ Desire to use the tool on their own projects
- ✅ Appreciation of the AI augmentation

**Visual Appeal**:
- ✅ Professional markdown formatting
- ✅ Clear separation of deterministic vs. LLM content
- ✅ Helpful visual elements (tables, bullet points)
- ✅ Confidence indicators for transparency

### Cool Factor Checklist

- [x] **Automatic**: One command does everything
- [x] **Smart**: AI explains the why behind performance issues
- [x] **Safe**: Redacts sensitive data automatically
- [x] **Professional**: Generates reports suitable for stakeholders
- [x] **Extensible**: Plugin architecture for any profiler
- [x] **Optional**: Works great without LLM too

### Demo Tips

1. **Start with a problem**: Show a slow application
2. **Run the analysis**: Single command magic
3. **Show the insights**: Focus on LLM explanations
4. **Compare approaches**: "Traditional vs. TriageProf"
5. **Highlight safety**: Emphasize data redaction
6. **Show flexibility**: Demonstrate plugin system

**Goal**: Make developers excited about performance analysis again!

---

## 🚀 Delivery Timeline

### Phase 1: Core Complete ✅
- Plugin discovery and validation
- Basic analysis pipeline
- Markdown report generation

### Phase 2: LLM Integration ✅
- Mistral API client
- Secure prompt generation
- Insights integration
- Enhanced reports

### Phase 3: Demo Polish (Current Focus)
- Perfect the demo script
- Create compelling sample application
- Refine visual output
- Prepare presentation materials

### Phase 4: Launch Ready
- Final testing and validation
- Documentation polish
- Screencast recording
- Website/update materials

**Target**: Have a "wow"-worthy demo ready for showcase!