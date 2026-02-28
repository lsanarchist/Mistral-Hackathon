# COMPASS — TriageProf Direction & Status (Living)

## North Star
Build a modular, language-agnostic profiling and bottleneck-analysis tool/library.
It integrates with existing profilers via plugins/adapters and produces evidence-backed bottleneck findings
(stacks/callgraphs/timelines/metrics), plus optional AI-assisted insights, and generates reports.

## Product Shape (today)
Pipeline: Collect → Analyze → Report.
Plugins are separate executables, discovered via manifests, and communicate with core using JSON-RPC 2.0 over stdio.

## Non-negotiable Architecture Rules
- Core is language-agnostic; language/profiler-specific logic must live in plugins/adapters.
- Plugin protocol/API must remain stable; changes must be backward-compatible or explicitly versioned.
- Deterministic analysis is the source of truth; AI/LLM features (if present) are optional and must not break the pipeline.
- One iteration = one PR-sized change; prefer minimal diffs and safe rollbacks.

## Current Focus (Now)
- <1–3 bullets: what we are actively improving right now>

## Next Milestones
- <1–5 bullets: concrete milestones>

## Feature Backlog (small + high signal)
- [ ] <user-visible capability>
- [ ] <user-visible capability>
- [ ] <user-visible capability>

## Plugins (status)
- go-pprof-http: ✅ (cpu/heap/mutex/block/goroutine via HTTP pprof)
- <future plugin>: ⏳

## Quick Verify
- Build: `make build`
- Tests: `make test` (or `go test ./...`)
- Demo: `make demo`

---

## Iteration Log (append-only)
(Always append new entries at the bottom; do not rewrite history.)

### Iter YYYYMMDD-HHMM — <timestamp / timezone>
**Type:** Maintenance | Feature | Plugin
**Objective:** <one line>

**Acceptance criteria (feature only)**
- [ ] ...
- [ ] ...

**Changes**
- ...
- ...

**Verification**
- Tests: `...`
- Lint/Format: `...`
- Build: `...`

**Risk/Notes**
- ...

**Git / Rollback**
- Branch: `agent/iter-YYYYMMDD-HHMM`
- Checkpoint tag: `agent-checkpoint-YYYYMMDD-HHMM`
- Commit: `<hash>`
- Rollback: revert commit OR reset to checkpoint