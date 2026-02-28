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
- Plugin discovery system with manifest-based validation
- Enhanced error handling for plugin compatibility

## Next Milestones
- Additional plugins for other languages (Python, Java, Node.js)
- Advanced analysis rules and heuristics
- Historical comparison features
- Web UI for report visualization

## Feature Backlog (small + high signal)
- [ ] Plugin marketplace/repository
- [ ] CI/CD pipeline integration
- [ ] Cloud deployment options

## Plugins (status)
- go-pprof-http: ✅ (cpu/heap/mutex/block/goroutine via HTTP pprof)
- Plugin discovery system: ✅ (manifest-based validation and capability checking)
- <future plugin>: ⏳

## Quick Verify
- Build: `make build`
- Tests: `make test` (or `go test ./...`)
- Demo: `make demo`

---

## Iteration Log (append-only)
(Always append new entries at the bottom; do not rewrite history.)

### Iter 20240228-1530 — UTC
**Type:** Feature
**Objective:** Add plugin discovery system with manifest-based validation

**Acceptance criteria (feature only)**
- [x] `triageprof plugins` command lists all available plugins with capabilities
- [x] Plugin resolution validates manifest and binary existence before launch
- [x] Target and profile capability validation prevents incompatible plugin use
- [x] Clear error messages for missing plugins, SDK mismatches, and capability issues
- [x] Backward compatibility maintained - existing workflows unchanged

**Changes**
- `internal/plugin/manifest.go`: Added Manifest model with strict parsing, discovery, and validation
- `internal/plugin/jsonrpc.go`: Enhanced PluginManager with ListPlugins() and ResolvePlugin() methods
- `cmd/triageprof/main.go`: Added `plugins` command to list available plugins
- `internal/plugin/manifest_test.go`: Comprehensive unit tests for manifest functionality

**Verification**
- Tests: `go test ./internal/plugin/...` - all passing
- Lint/Format: `gofmt -d .` - no formatting issues
- Build: `make build` - successful
- Integration: `bin/triageprof plugins` - lists go-pprof-http with capabilities

**Risk/Notes**
- No breaking changes to existing functionality
- Plugin discovery is backward compatible
- All existing workflows continue to work unchanged

**Git / Rollback**
- Branch: `main`
- Checkpoint tag: N/A (direct commit to main)
- Commit: `<to be filled after commit>`
- Rollback: `git revert <commit-hash>`