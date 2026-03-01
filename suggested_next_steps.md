# Suggested Next Steps — TriageProf Demo-Grade MVP (Go pprof + Mistral API)

> Goal: **one command → a "wow" report**.
> `triageprof demo --repo <url> [--ref <commit>] --out out/`  
> Outputs: `report.html` + `report.md` + `findings.json` (+ raw pprof artifacts).

---

## Definition of Done (what "MVP demo" means)

- Works on a **real Go open-source repo** that has benchmarks (or produces a clear "no benchmarks found" report).
- Produces a **clean HTML report** with:
  - Top 3–5 bottlenecks (ranked)
  - Evidence (pprof top funcs + stacks + file:line when available)
  - Actionable fix guidance (deterministic rules + optional LLM enrichment)
- **LLM is optional**:
  - If Mistral API is configured → adds "why / how to fix / trade-offs"
  - If not configured / API fails → report still builds, just without LLM sections
- Output is **reproducible** (pinned repo ref + stable schema + snapshot tests).

---

## ✅ Phase 0 — Repo Hygiene (COMPLETED)

- **Narrow scope to Go for MVP** ✅
  - Moved non-Go examples (node, python, ruby) to `examples/_archive/`
  - Kept only `examples/demo-server/` and Go-focused demo
- **Delete/disable local-ML paths** ✅
  - No local ML implementations found - already remote API only
- **LLM is optional and safe** ✅
  - `--llm=off` is the default (llmGenerator nil by default)
  - Remote API calls only (Mistral/OpenAI)
- **Docs cleanup** ✅
  - Simplified README.md to focus on Go MVP
  - Removed merge-conflict markers and noise

## ✅ Phase 1 — Golden Path CLI: `triageprof demo` (COMPLETED)

**NEXT PRIORITY**: Implement the single "golden" command that does everything.

### Implementation Plan

1. **Add `demo` command** to `cmd/triageprof/main.go` ✅
2. **Implement repo cloning** with git support ✅
3. **Add benchmark detection** using `go test` patterns ✅
4. **Create profile collection** with proper pprof flags ✅
5. **Generate run manifest** with metadata ✅

### Acceptance Criteria
- ✅ `triageprof demo --repo <repo> --out out/` produces `out/profiles/*` + `out/run.json`
- ✅ Works with both local paths and git URLs
- ✅ Handles benchmark detection gracefully
- ✅ Generates complete reports and findings

### Estimated Files to Modify
- ✅ `cmd/triageprof/main.go` - Add demo command handler
- ✅ `internal/core/demo.go` - Complete demo implementation
- ✅ `internal/core/pipeline.go` - Add Demo() method integration
- ✅ `internal/model/types.go` - Add RunManifest struct
- ✅ `internal/core/pipeline_test.go` - Fix test paths

## ✅ Phase 2 — Deterministic Analyzer: `pprof → findings.json` (COMPLETED)

Create stable schema and rules that never depend on LLM.

### Findings Schema (v1)
```go
type Finding struct {
    ID               string
    Title            string
    Category         string // cpu, alloc, heap, gc, mutex, block
    Severity         string // low, medium, high, critical
    Confidence       float64 // 0.0-1.0
    ImpactSummary    string
    Evidence         []EvidenceItem
    DeterministicHints []string
    Tags             []string
}
```

### Rules Implemented
1. ✅ **CPU hotpath dominance** - Top N functions consuming >70% cumulative time
2. ✅ **Allocation churn** - High mallocgc/memmove/bytes.growSlice patterns
3. ✅ **JSON hotspots** - encoding/json decode/encode in top functions
4. ✅ **String churn** - strings.Builder, bytes.Buffer, regexp in hot paths
5. ✅ **GC pressure** - runtime.gcBgMarkWorker, runtime.gcAssistAlloc
6. ✅ **Mutex contention** - sync.(*Mutex).Lock with high contention
7. ✅ **Heap allocation concentration** - Top N functions consuming >60% of heap allocations
8. ✅ **Block contention** - runtime.chan, select, Wait, Sleep patterns

### Implementation Details
- **New deterministic analyzer** in `internal/analyzer/deterministic.go`
- **Updated model types** with new Finding schema and backward compatibility
- **Enhanced report generation** to handle both new and legacy evidence formats
- **Updated LLM client** to work with new evidence structure
- **Comprehensive test coverage** for all deterministic rules

### Files Modified
- ✅ `internal/model/types.go` - Added new Finding schema with EvidenceItem
- ✅ `internal/analyzer/deterministic.go` - New deterministic analyzer implementation
- ✅ `internal/analyzer/analyzer.go` - Added deterministic analysis method
- ✅ `internal/core/pipeline.go` - Added deterministic analysis pipeline method
- ✅ `internal/core/demo.go` - Updated to use deterministic analysis
- ✅ `internal/report/report.go` - Enhanced to handle both evidence formats
- ✅ `internal/llm/client.go` - Updated for new evidence structure
- ✅ `internal/analyzer/analyzer_test.go` - Added deterministic analyzer tests
- ✅ `internal/core/pipeline_test.go` - Updated test evidence format

---

## 🚀 Phase 1 — Golden Path CLI: `triageprof demo`

**NEXT PRIORITY**: Implement the single "golden" command that does everything.

### Implementation Plan

1. **Add `demo` command** to `cmd/triageprof/main.go`
2. **Implement repo cloning** with git support
3. **Add benchmark detection** using `go test` patterns
4. **Create profile collection** with proper pprof flags
5. **Generate run manifest** with metadata

### Acceptance Criteria
- `triageprof demo --repo <repo> --out out/` produces `out/profiles/*` + `out/run.json`
- Works with both local paths and git URLs
- Handles benchmark detection gracefully

### Estimated Files to Modify
- `cmd/triageprof/main.go` - Add demo command handler
- `internal/core/pipeline.go` - Add Demo() method
- `internal/model/types.go` - Add RunManifest struct

---

## Phase 2 — Deterministic Analyzer: `pprof → findings.json`

Create stable schema and rules that never depend on LLM.

### Findings Schema (v1)
```go
type Finding struct {
    ID               string
    Title            string
    Category         string // cpu, alloc, heap, gc, mutex, block
    Severity         string // low, medium, high, critical
    Confidence       float64 // 0.0-1.0
    ImpactSummary    string
    Evidence         []EvidenceItem
    DeterministicHints []string
    Tags             []string
}
```

### Rules to Implement
1. **CPU hotpath dominance** - Top N functions consuming >70% cumulative time
2. **Allocation churn** - High mallocgc/memmove/bytes.growSlice patterns
3. **JSON hotspots** - encoding/json decode/encode in top functions
4. **String churn** - strings.Builder, bytes.Buffer, regexp in hot paths
5. **GC pressure** - runtime.gcBgMarkWorker, runtime.gcAssistAlloc
6. **Mutex contention** - sync.(*Mutex).Lock with high contention

---

## ✅ Phase 3 — LLM Enrichment via Mistral API ✅ COMPLETED

Add optional enrichment that takes deterministic findings and returns structured insights.

### Guardrails Implemented (Critical)
1. ✅ **Strict JSON validation** - Discard invalid responses with comprehensive schema validation
2. ✅ **Evidence citations** - Require evidence_refs in responses with finding ID validation
3. ✅ **Redaction** - Strip secrets, limit code snippets to 200 chars using regex patterns
4. ✅ **Caching** - Hash findings.json → cache insights with size and age limits
5. ✅ **Confidence validation** - Ensure confidence scores are in valid 0-100 range
6. ✅ **Finding reference validation** - Ensure all insights reference valid finding IDs
7. ✅ **Code example length limits** - Enforce 200 character maximum for all code examples

### Implementation Details
- **Enhanced InsightsGenerator** with comprehensive validation and guardrails
- **Strict JSON parsing** with fallback to legacy text format for backward compatibility
- **Automatic caching** enabled by default with configurable limits
- **Redaction patterns** for tokens, secrets, URLs, and hostnames
- **Field truncation** to prevent excessive output sizes
- **Schema version 2.0** for structured insights

### Key Features
- **Caching**: Insights are automatically cached based on findings hash
- **Validation**: Strict validation of all LLM-generated content
- **Redaction**: Automatic removal of sensitive information
- **Truncation**: Enforcement of field length limits
- **Backward Compatibility**: Support for both JSON and legacy text responses
- **Error Handling**: Graceful degradation when validation fails

### Files Modified
- ✅ `internal/llm/insights.go` - Enhanced with validation, guardrails, caching, and redaction
- ✅ `internal/llm/mistral.go` - Updated with strict JSON parsing and validation
- ✅ `internal/llm/client.go` - Enhanced prompt building with evidence citation requirements
- ✅ `internal/llm/cache.go` - Caching enabled by default with proper configuration
- ✅ `internal/llm/insights_test.go` - Comprehensive test coverage for validation and guardrails

### Validation Rules Implemented
1. **Schema Validation**: Executive summary, per-finding insights, top risks, top actions
2. **Confidence Range**: 0-100 for all confidence scores
3. **Finding ID Validation**: All finding IDs must exist in findings bundle
4. **Evidence Citations**: All narratives must reference their finding IDs
5. **Code Example Limits**: Maximum 200 characters for all code examples
6. **Severity Validation**: Valid severity levels (low, medium, high, critical)
7. **Priority Validation**: Valid priority levels (low, medium, high, critical)

---

## Quick Verification Checklist

- [x] Phase 0: Repo hygiene complete
- [x] Phase 1: `triageprof demo` command works
- [x] Phase 2: Deterministic analyzer produces findings.json
- [x] Phase 3: LLM enrichment works with API key
- [x] Phase 4: HTML report looks professional
- [ ] Phase 5: Demo kit with pinned repo works

---

## Immediate Next Action

**Start Phase 5**: Implement demo kit with pinned repository that works end-to-end.

Expected time: 4-8 hours for demo kit implementation + testing.
