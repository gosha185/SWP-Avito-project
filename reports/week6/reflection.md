# Week 6 Reflection

## Learning points

1. **Internal conflict affects the whole team** — When two team members disagree, it creates friction that slows everyone down. As a team lead, I learned that addressing conflict early is critical. For this Sprint, I used a text-based communication channel through myself to keep work moving.

2. **Documentation is critical for handover** — Preparing `README.md`, `customer-handover.md`, `CONTRIBUTING.md`, and `AGENTS.md` was more work than expected, but it made the Sprint Review much smoother. The customer appreciated the clarity.

3. **Daily Scrums need structure** — Without consistent attendance, alignment suffers. The team decided to introduce 15-minute reminders before each sync. This improved attendance slightly, but we still need to work on discipline.

4. **Conflict resolution takes time** — It is not always possible to resolve disagreements immediately. Sometimes the best option is to create a temporary workaround (text-based communication) and address the root cause later.

5. **Trial release validates progress** — Deploying a trial release and showing it to the customer gave us confidence that the product is ready for handover. Even with internal issues, the team delivered.

## Validated assumptions

| Assumption | Result | Evidence |
|------------|--------|----------|
| The product is ready for handover | **Confirmed** — customer approved the trial release | Sprint Review feedback |
| Documentation is sufficient for handover | **Confirmed** — customer reviewed and approved | Sprint Review feedback |
| Daily Scrum is the best sync method | **Partially rejected** — attendance was inconsistent; text-based sync was used as fallback | Daily Scrum logs |
| DoR/DoD templates improve quality | **To be validated** — added this week, will evaluate next Sprint | Issue templates |
| Team conflict can be managed via text | **Partially confirmed** — work continued, but conflict remains unresolved | Team communication |

## Friction and gaps

1. **Team conflict** — Two team members disagree on optimisation priorities. The conflict is unresolved and may affect Sprint 5 if not addressed.

2. **Daily Scrum attendance** — Attendance was inconsistent. One day, no one attended. A 15-minute reminder was added, but discipline is still an issue.

3. **Team member illness** — One member was sick, causing minor delays. The work was completed, but slightly later than planned.

4. **Idempotency handling feedback** — The customer suggested improving idempotency retry handling (return 200 instead of 409). This is a useful improvement for the final handover.

## Planned response

1. **Resolve conflict before Sprint 5** — The team lead will mediate a discussion between the two members to align on optimisation priorities.

2. **Improve Daily Scrum discipline** — Continue using 15-minute reminders. If attendance remains low, the team will switch to written async updates.

3. **Address customer feedback** — Review idempotency handling (HTTP 200 on retry with same key) and include UAT test status for user stories in the handover documentation.

4. **Complete handover** — Schedule and conduct the final handover meeting (Thursday), create SemVer release for MVP v3, and prepare for Demo Day.
