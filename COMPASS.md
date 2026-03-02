# COMPASS — TriageProf (Demo MVP)

## North Star
## North Star
One command → Go pprof triage report with **mandatory Mistral LLM enrichment** grounded in deterministic evidence.


## Demo Promise
With Go installed and API key set:
- `make demo` succeeds and produces:
  - `findings.json` (deterministic)
  - `llm_enrichment.json` (mandatory)
  - `report.md` + web report folder

If LLM cannot run (missing key/provider), `make demo` exits non-zero with clear instructions.

## Product Shape
Collect → Analyze (deterministic) → Enrich (LLM mandatory) → Report.

## Rules
- Deterministic outputs are truth; LLM adds “why/how”, never invents numbers.
- Enrichment references evidence IDs from findings.
- Cache enrichment by input hash for speed.

## Scope
Go-only golden path. Archive or quarantine non-Go plugins and “enterprise/phase4/phase5/remediation” features away from default demo.

## Now
Make `make demo` + `demo-kit` the single golden path and keep CI green.