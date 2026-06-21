# Sprint Retrospective – Week 3

## What went well

1. **Clear backlog refinement** – We successfully refined the Product Backlog into 27 well-defined PBIs with clear acceptance criteria, story points, and MoSCoW priorities. This made Sprint Planning smooth and gave everyone a clear understanding of what to do.

2. **Effective use of GitHub Projects** – We set up both Product Backlog and Sprint Backlog with proper columns (To Do, In Progress, Review, Done). This gave full visibility into progress.

3. **Customer feedback incorporated** – The customer's advice on estimation, DoD/DoR, and incremental delivery was valuable. We now have a clearer understanding of how to manage research tasks and decompose large PBIs.

## What did not go well

1. **Late collaborator permissions** – Not all team members had write access to the repository initially. This delayed issue assignment and project setup. The team spent unnecessary time waiting for permissions.

2. **Story Point estimation uncertainty** – Some tasks (#8 with SP=13, #9 with SP=8) may be too large. The customer confirmed that tasks above 5 story points are harder to estimate accurately and recommended keeping estimates ≤ 5.

3. **PR and branch naming chaos** – We initially created branches with incorrect naming (`Leilia34-patch-*`) instead of the required `<issue-number>-short-description` format. This caused confusion and required rework. In the future, we will enforce naming conventions from the start.

## Action points

1. **Split large tasks** – #8 (SP=13) and #9 (SP=8) will be split into smaller subtasks for Sprint 2. Goal: keep all tasks ≤ 5 story points where possible.

2. **Set up public deployment** – Deploy MVP v1 to Render or ngrok for external TA access if needed.

3. **Enforce branch naming** – All team members must use `<issue-number>-short-description` format for branches from the start of the sprint.

4. **Fix collaborator permissions** – The repository owner will grant write access to all team members to avoid future delays.

5. **Write down problems as they occur** – Maintain a shared document to capture issues, blockers, and process friction during the sprint.

6. **New tests for Sprint 2** – Create test Issues for:
   - TTL worker for expired point batches (#9)
   - TTL worker for expired held points (#32)
   - View all current held points (#23)
   - Edge cases for FEFO spending order
   - Idempotency with external_key under load
