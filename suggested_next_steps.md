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
  - Moved non-Go examples to `examples/_archive/` (kept history)
  - Kept only `examples/demo-server/` and Go-focused demo
- **Delete/disable local-ML paths** ✅
  - No local ML implementations found - already remote API only
- **LLM is optional and safe** ✅
  - `--llm=off` is the default (llmGenerator nil by default)
  - Remote API calls only (Mistral/OpenAI)
- **Docs cleanup** ✅
  - Simplified README.md to focus on Go MVP
  - Removed merge-conflict markers and noise

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

## Phase 3 — LLM Enrichment via Mistral API

Add optional enrichment that takes deterministic findings and returns structured insights.

### Guardrails (Critical)
- **Strict JSON validation** - Discard invalid responses
- **Evidence citations** - Require `evidence_refs` in responses
- **Redaction** - Strip secrets, limit code snippets to 200 chars
- **Caching** - Hash findings.json → cache insights

### Implementation
```go
func (g *InsightsGenerator) GenerateInsights(ctx context.Context, findings *model.FindingsBundle) (*model.InsightsBundle, error) {
    // Build prompt with redacted findings
    // Call Mistral API with structured prompt
    // Validate JSON response strictly
    // Return insights or disabled bundle
}
```

---

## Quick Verification Checklist

- [x] Phase 0: Repo hygiene complete
- [ ] Phase 1: `triageprof demo` command works
- [ ] Phase 2: Deterministic analyzer produces findings.json
- [ ] Phase 3: LLM enrichment works with API key
- [ ] Phase 4: HTML report looks professional
- [ ] Phase 5: Demo kit with pinned repo works

---

## Immediate Next Action

**Start Phase 1**: Implement `triageprof demo` command with repo cloning, benchmark detection, and profile collection.

Expected time: 4-8 hours for basic implementation + tests.
