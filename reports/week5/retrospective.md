# Sprint Retrospective – Week 5

## What went well

1. **Architecture documentation completed** – We successfully created static, dynamic, and deployment views. This was a new challenge for the team, and we delivered it on time.

2. **MVP v2 delivered** – We completed the remaining product features: TTL workers, expiring points, view all held points, and fixed HTTP 500 → HTTP 400 error handling. The customer approved the increment.

3. **Customer feedback incorporated** – The customer provided valuable feedback on architecture, testing, and process. We captured all feedback and created action points for future improvements.

## What did not go well

1. **PR review process blocked** – Pull Requests were approved but still blocked from merging. This created unnecessary delays and confusion. The issue needs investigation.

2. **Daily Scrum attendance** – Not all team members attended every Daily Scrum (3–4 people on average). This reduced the effectiveness of the synchronisation and created misalignment.

3. **HTTP 500 issue persists** – The confirm/cancel flow returned HTTP 500 because of forgetting to merge the fix.

## What changed from the previous Sprint

| Previous Sprint | This Sprint |
|-----------------|-------------|
| No architecture documentation | Static, dynamic, deployment views added |
| No ADRs | 3 ADRs created and reviewed |
| Manual testing of workers | Customer suggested unit tests — planned for next sprint |
| No development process docs | `docs/development-process.md` created |
| Daily Scrums in 2-3 days | Every day, but not all members |

## Process improvements for the next Sprint

1. **Investigate and fix PR review blocking** – The team will investigate why approved PRs are still blocked from merging and resolve the configuration issue.

2. **Improve Daily Scrum attendance** – The team will commit to attending Daily Scrums every day. The team lead will send reminders before each session.

3. **Add DoD/DoR checklists to Issue templates** – As suggested by the customer, add Definition of Done and Definition of Ready checklists to GitHub Issue templates.

4. **Use unit tests for worker testing** – Replace manual text-based checks with automated unit tests for workers.

5. **Split sequence diagrams** – Refactor the sequence diagram into smaller, scenario-specific diagrams as suggested by the customer.
