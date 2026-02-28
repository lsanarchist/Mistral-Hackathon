## Feature Policy (MANDATORY for feature iterations)
A "feature iteration" must:
1) Define the feature in 2–5 bullet acceptance criteria.
2) Write code.
3) Add tests covering the new behavior (or a deterministic verification harness if tests are hard).
4) Update  COMPASS.md + project_context.md to reflect:
   - what the feature does
   - how it integrates with core schema / pipeline
   - plugin impact (if any)
5) Preserve backwards compatibility unless explicitly allowed to break it.
   - If breaking changes are unavoidable: version the plugin API or add compatibility notes.

## Stop Conditions (Avoid feature creep)
- If no feature fits within the change budget and verification gate, do NOT implement a half-feature.
- Instead, update COMPASS.md with a "Feature Backlog" entry (3–5 ideas) and stop.

## Mandatory Reading Order (every iteration)
1) Read COMPASS.md first (North Star, constraints, current focus, recent changes).
2) Consult project_context.md only if deeper architectural detail is needed.