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

### Iter 20240301-1730 — UTC
**Type:** Maintenance
**Objective:** Verify comprehensive error handling for plugin compatibility

**Acceptance criteria (maintenance)**
- [x] Plugin discovery system with manifest-based validation is working
- [x] All error scenarios produce clear, actionable error messages
- [x] Backward compatibility maintained - existing workflows unchanged
- [x] All tests passing
- [x] Demo workflow produces expected output

**Verification**
- Tests: `go test ./...` - all passing
- Error scenarios tested via CLI:
  - `bin/triageprof collect --plugin non-existent` → "plugin not found. Available plugins: go-pprof-http"
  - Missing binary → "manifest found but binary missing at ..."
  - SDK mismatch → "plugin requires sdkVersion 2.0, but core supports 1.0"
  - Unsupported target → "target type not supported. Supported targets: ..."
  - Unsupported profile → "profiles not supported. Supported profiles: ..."
- Build: `make build` - successful
- Demo: `make demo` - produces report.md successfully

**Risk/Notes**
- No breaking changes - all existing functionality preserved
- Error handling is backward compatible
- User experience improved with clear, actionable error messages
- All validation happens before plugin execution
- Feature was already implemented in previous iteration, this commit verifies it works correctly

**Git / Rollback**
- Branch: `main`
- Checkpoint tag: N/A (direct commit to main)
- Commit: `1a61b9d`
- Rollback: `git revert 1a61b9d`

### Iter 20240301-1200 — UTC
**Type:** Maintenance
**Objective:** Add comprehensive error handling for plugin compatibility

**Acceptance criteria (maintenance)**
- [x] Clear error when plugin manifest is missing
- [x] Clear error when plugin binary is missing
- [x] Clear error when SDK version is incompatible
- [x] Clear error when target type is unsupported
- [x] Clear error when profile type is unsupported

**Changes**
- `internal/plugin/manifest.go`: Enhanced error messages for all validation scenarios
- `internal/core/pipeline.go`: Integrated validation in Collect() method
- `cmd/triageprof/main.go`: CLI properly displays validation errors to users
- `internal/plugin/manifest_test.go`: Comprehensive tests for all error scenarios

**Verification**
- Tests: `go test ./internal/plugin/...` - all passing
- Integration: Tested all error scenarios via CLI
  - `bin/triageprof collect --plugin non-existent` → "plugin not found. Available plugins: ..."
  - Missing binary → "manifest found but binary missing at ..."
  - SDK mismatch → "plugin requires sdkVersion X, but core supports Y"
  - Unsupported target → "target type not supported. Supported targets: ..."
  - Unsupported profile → "profiles not supported. Supported profiles: ..."
- Build: `make build` - successful
- Demo: `make demo` - produces report.md successfully

**Risk/Notes**
- No breaking changes - all existing functionality preserved
- Error handling is backward compatible
- User experience improved with clear, actionable error messages
- All validation happens before plugin execution

**Git / Rollback**
- Branch: `main`
- Checkpoint tag: N/A (direct commit to main)
- Commit: `<to be filled after commit>`
- Rollback: `git revert <commit-hash>`