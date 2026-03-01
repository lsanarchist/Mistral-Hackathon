# Suggested Next Steps — TriageProf Demo-Grade MVP (Go pprof + Mistral API)

> Goal: **one command → a “wow” report**.
> `triageprof demo --repo <url> [--ref <commit>] --out out/`  
> Outputs: `report.html` + `report.md` + `findings.json` (+ raw pprof artifacts).

---

## Definition of Done (what “MVP demo” means)

- Works on a **real Go open-source repo** that has benchmarks (or produces a clear “no benchmarks found” report).
- Produces a **clean HTML report** with:
  - Top 3–5 bottlenecks (ranked)
  - Evidence (pprof top funcs + stacks + file:line when available)
  - Actionable fix guidance (deterministic rules + optional LLM enrichment)
- **LLM is optional**:
  - If Mistral API is configured → adds “why / how to fix / trade-offs”
  - If not configured / API fails → report still builds, just without LLM sections
- Output is **reproducible** (pinned repo ref + stable schema + snapshot tests).

---

## Phase 0 — Repo Hygiene (make the project “demo clean”)

- **Narrow scope to Go for MVP**
  - Move non-Go examples to `examples/_archive/` (keep history, stop confusing the demo).
  - Keep only `examples/go/` and one pinned demo repo.
- **Delete/disable local-ML paths**
  - Keep LLM only as **remote API calls** (Mistral).
  - Ensure `--llm=off` is the default safe path.
- **Docs cleanup**
  - Replace the current `suggested_next_steps.md` with this document (no merge-conflict markers, no “COMPLETED” noise).

Acceptance criteria:
- `go test ./...` passes; `triageprof --help` is readable; repo looks focused.

---

## Phase 1 — Golden Path CLI: `triageprof demo`

Implement a single “golden” command that does everything:

- **Inputs**
  - `--repo` (git URL or local path)
  - `--ref` (commit/tag/branch; default: default branch)
  - `--bench` (regex; default: `.`)
  - `--pkg` (optional package selector)
  - `--duration` / `--count` (for stability)
  - `--out` (output dir)
- **Steps**
  1) Clone (or copy local)
  2) Detect benchmarks (fast scan: `go test ./... -list=.` or `-run=^$ -bench=.` with dry mode)
  3) Run benchmark(s) and generate:
     - CPU: `-cpuprofile`
     - Mem/alloc: `-memprofile` (+ `-benchmem`)
     - Optional: mutex/block profiles behind flags
  4) Save a `run.json` manifest (repo/ref, go env, flags, timestamps, versions)

Acceptance criteria:
- `triageprof demo --repo <repo> --out out/` produces `out/profiles/*` + `out/run.json` every time.

---

## Phase 2 — Deterministic Analyzer: `pprof → findings.json`

Create a stable schema and rules that never depend on LLM.

### Findings schema (v1)
Each finding should contain:
- `id`, `title`, `category` (cpu/alloc/heap/gc/mutex/block)
- `severity` (low/med/high/critical) + `confidence` (0..1)
- `impact_summary` (short)
- `evidence[]` (top funcs, flat/cum, stack excerpts, file:line where possible)
- `deterministic_hints[]` (rule-based fix ideas)
- `tags[]` (e.g., `allocation_churn`, `lock_contention`, `json_decode_hotpath`)

### Rules to implement first (highest demo value)
- CPU hotpath dominance (top N cumulative + dominance ratio)
- Allocation churn (mallocgc/memmove/bytes.growSlice patterns)
- JSON decode/encode hotspots (encoding/json heavy frames)
- String churn patterns (strings/bytes/regexp hotspots)
- GC pressure (runtime/GC frames noticeable in CPU profile)
- Mutex contention (single lock dominating)

Acceptance criteria:
- `triageprof analyze --in out/profiles --out out/` creates `out/findings.json` with 3–10 findings for a typical repo.

---

## Phase 3 — LLM Enrichment via Mistral API (no local ML)

Add an enrichment stage that takes **only the deterministic findings** and returns strictly structured output.

### Interface
- `--llm=mistral` (or `--llm=off`)
- Reads `MISTRAL_API_KEY`
- Adds `out/llm_enrichment.json` and merges into report sections.

### Guardrails (must-have)
- **Strict JSON-only output** (validate; if invalid → discard and continue)
- **No hallucinated claims**:
  - Require the model to cite `evidence_refs` that point to your evidence entries
  - If uncertain, it must put items into `unknowns[]`
- **Redaction**:
  - Strip secrets and overly-large code blobs
  - Send summaries + small evidence snippets, not whole repos
- **Caching**:
  - Hash of `findings.json` → cached enrichment to avoid repeat costs

Acceptance criteria:
- With a valid key: report includes an “LLM insights” section per top finding.
- Without a key: report still builds, clearly shows “LLM disabled”.

---

## Phase 4 — Report Generator Polish (this is what sells the demo)

Generate three outputs every run:

1) **`report.html`** (primary demo artifact)
   - Summary: repo, ref, runtime, top bottleneck cards
   - Findings list: sortable by severity/impact/confidence
   - Each finding: evidence tables + expandable stacks + (optional) LLM narrative
   - Footer: reproducibility info (command line, versions)

2) **`report.md`** (for GitHub issue/comment)
   - Short summary + top findings + bullet fixes
   - Links/paths to the evidence artifacts

3) **`findings.json`** (machine-readable contract)

Acceptance criteria:
- Opening `report.html` looks clean and understandable to a non-expert in <60 seconds.

---

## Phase 5 — Demo Kit (make it easy to show live)

- Add `./demo.sh` (or `make demo`) that runs the golden path on a **pinned** repo/ref.
- Commit `demo-output/` with one known-good output (so the UI can be shown even offline).
- Add README “Demo in 30 seconds”:
  - install
  - run demo
  - open report

Acceptance criteria:
- New person can reproduce the demo without asking questions.

---

## Phase 6 — Tests (so you can move fast safely)

- Unit tests:
  - pprof parsing wrappers
  - each rule emits expected finding fields
- Snapshot (golden) tests:
  - `findings.json` stable ordering
  - `report.md` stable text blocks
- E2E test:
  - run `triageprof demo` against a tiny fixture repo (or minimal internal test module)
  - assert files exist + basic invariants

Acceptance criteria:
- CI (or local `go test ./...`) catches regressions in schema/report.

---

## Phase 7 — Small “Nice” Extras (only if time remains)

- `triageprof serve --dir out/` to open a local report viewer (optional).
- “Compare two runs” (baseline vs current) for a single benchmark (simple diff in findings).
- Export a `github_issue.md` template file.

---

## Backlog (post-MVP, do not block the demo)

- Plugin SDK hardening + compatibility matrix
- Support more profilers/languages (keep as skeleton only for now)
- Advanced callgraph/dominator visualizations
- CI integration (GitHub Actions) and regression thresholds
