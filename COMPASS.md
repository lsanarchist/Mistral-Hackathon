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
- **AI/LLM analysis is a first-class feature**, not optional glue — it must be designed as a proper pipeline stage with a clean interface, not bolted on.
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
   - Manifest-based discovery and validation ✅
   - Capability checking and compatibility
   - Plugin lifecycle management

2. **LLM Integration**
   - Safe prompt generation with redaction ✅
   - Structured insights generation
   - Report enhancement with LLM sections

3. **Core Robustness**
   - Error handling and recovery
   - Test coverage and validation
   - Performance optimization

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