# Week 5 Reflection

## Learning points

1. **Architecture documentation is harder than it looks** – Creating component, sequence, and deployment diagrams required a solid understanding of the system.

2. **ADRs help preserve decisions** – Writing Architecture Decision Records forced us to explain why we made certain choices.
3. 
4. **Daily Scrums improve alignment** – We held Daily Scrums every day this week. Even though not everyone attended every day (3–4 people on average), they helped keep the team synchronised and reduced blockers.

5. **Customer feedback on architecture is valuable** – The customer's suggestion to split sequence diagrams into smaller, scenario-specific diagrams (with `alt` blocks) was a useful insight. Smaller, focused diagrams are more readable than one large diagram.

6. **Manual testing is not enough** – During the Sprint Review, we demonstrated workers using manual text-based checks. The customer suggested using unit tests instead. Automated tests are more reliable and repeatable.

7. **Technical issues disrupt meetings** – Microphone problems during the Sprint Review disrupted communication. We need to ensure all team members test their audio/video before meetings.

## Validated assumptions

| Assumption | Result | Evidence |
|------------|--------|----------|
| MVP v2 would require architecture documentation | **Confirmed** – We created static, dynamic, and deployment views | `docs/architecture/` |
| ADRs would help explain decisions | **Confirmed** – 3 ADRs were created and reviewed | `docs/architecture/adr/` |
| Daily Scrums would improve team alignment | **Partially confirmed** – Helped, but attendance was inconsistent | Sprint Review feedback |
| Sequence diagrams should be large and comprehensive | **Rejected** – Customer advised splitting into smaller, scenario-specific diagrams | Customer feedback |
| Manual testing of workers is sufficient | **Rejected** – Customer advised using unit tests instead | Customer feedback |

## Friction and gaps

1. **Microphone issues** – During the Sprint Review, one member had microphone problems, which disrupted communication. Need to ensure technical setup is stable for future meetings.

2. **Daily Scrum attendance** – Not all team members attended every Daily Scrum. This reduced the effectiveness of the synchronisation.

3. **PR review process** – Pull Requests were approved but still blocked from merging. This suggests a configuration issue with GitHub branch protection rules that needs investigation.

4. **Worker demonstration** – Workers were demonstrated using manual text-based checks instead of automated tests. This made the demo slower and less reliable.

5. **HTTP 500 issue** – The confirm/cancel flow still returns HTTP 500 because the team forgot to merge the PR with the fix.

## Planned response

1. **Fix PR review process** – Investigate GitHub branch protection rules to ensure approved PRs can be merged.

2. **Add DoD/DoR checklists to Issue templates** – As suggested by the customer, add Definition of Done and Definition of Ready checklists to GitHub Issue templates.

3. **Split sequence diagrams** – Refactor the sequence diagram into smaller, scenario-specific diagrams (balance flow, hold flow with alt blocks for confirm/cancel/timeout).

4. **Improve meeting setup** – Ensure all team members test their audio/video before meetings.
