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
- LLM calls are always backed by structured profiling data;
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


## Current Focus Areas

check out suggested_next_steps.md