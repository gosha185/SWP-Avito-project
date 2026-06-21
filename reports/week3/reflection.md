# Week 3 Reflection

## Learning points

1. **Estimation is an art** – The customer explained the difference between implementation tasks (easier to estimate) and research tasks (more variance). We learned to keep estimates ≤ 5 story points and to use time‑boxed research (spikes) for unknown work.

2. **DoR vs DoD** – We now understand the distinction between Definition of Ready (when a task can be started) and Definition of Done (when it's truly complete). We plan to use GitHub Issue templates to automate these checklists.

3. **Incremental delivery** – The customer emphasised delivering fully completed increments rather than partial progress on many problems. This will guide our future sprint planning.

4. **Retrospective preparation** – We learned to write down problems as they occur and to separate retro discussions into metrics and process/feelings.

5. **Public deployment** – We realised that MVP v1 must be publicly accessible, which requires setting up deployment (Render/ngrok) – a lesson for future projects.

## Validated assumptions

| Assumption | Result | Evidence |
|------------|--------|----------|
| Showing available balance is critical for MVP | **Confirmed** – US-07 is Must Have | Customer feedback (Week 2 & Week 3) |
| Idempotency is the right solution for duplicate prevention | **Confirmed** – customer recommended this approach | Customer feedback |
| A Kanban‑style board helps visualise progress | **Confirmed** – customer advised setting it up | Customer feedback |
| Internal VM would be sufficient for TA evaluation | **Rejected** – requires public URL or university network access | Assignment requirement |
| The team can deliver 25 tasks in one sprint | **Partially confirmed** – 17 completed, 8 moved to next sprint | Sprint Backlog |

## Needs clarification

- **Public deployment solution** – We still need to decide between Render, Railway, or ngrok. For now, the VM is accessible within the university network.

## Planned response

1. **Complete remaining MVP v1 tasks** – #8, #9, #12, #14, #15, #16, #23, #32 will be completed in Sprint 2.
2. **Set up public deployment** – Use Render or ngrok for MVP v1 if needed for external TA access.
3. **Create SemVer release v1.0.0** – Tag the final MVP v1 code (Stepan will do this after all merges).
4. **Enforce DoD** – All PRs reviewed, tests pass, CHANGELOG updated.
5. **Implement retrospective improvements** – Use customer advice for future sprints.
6. **New tests for Sprint 2** – The following test tasks are planned:
   - `Test: TTL worker for expired point batches` (#9)
   - `Test: TTL worker for expired held points` (#32)
   - `Test: view all current held points` (#23)
   - `Test: edge cases for FEFO spending order`
   - `Test: idempotency with external_key under load`
