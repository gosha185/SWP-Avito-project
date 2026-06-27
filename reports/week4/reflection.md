# Week 4 Reflection

## Learning points

1. **Quality requirements are hard to define** – Turning vague "be fast" into measurable "500ms for 95% of requests" required careful thought and team discussion.

2. **Documentation first, then tests** – We started defining quality requirements before writing tests. This gave us clear targets and made it easier to design tests.

3. **UAT with the customer is valuable** – The customer executed our UAT scenarios and gave immediate feedback. This is much more effective than guessing what they want.

4. **Daily scrum is needed** - We discovered that daily scrums enhance mutual understanding of workload and foster greater team cohesion.
---

## Validated assumptions

| Assumption | Result | Evidence |
|------------|--------|----------|
| Quality requirements should be based on ISO/IEC 25010 | **Confirmed** – It provided clear structure for QR-001, QR-002, QR-003 | `docs/quality-requirements.md` |
| Coverage ≥ 30% is achievable for critical modules | **To be validated** – Tests are still in progress | `docs/testing.md` (pending) |
| UAT catches issues early | **Confirmed** – Customer feedback from UAT session | `docs/user-acceptance-tests.md` |
| CI pipeline is essential for quality | **Confirmed** – Team agreed to prioritise it | `#144` |

---

## Friction and gaps

1. **QR/QRT approvals delayed** – `quality-requirements.md` and `quality-requirement-tests.md` were waiting for review for couple days. This blocked `testing.md` and slowed progress.

2. **Late UAT session** – UAT was scheduled later in the sprint, leaving less time to address feedback.

---

## Planned response

1. **Finalise QR/QRT** – Get approvals as soon as possible, then complete `testing.md`.

2. **Capture UAT feedback** – Convert customer feedback into actionable PBIs for Sprint 3.

3. **Prepare for Sprint 3** – Define Sprint 3 scope based on remaining MVP v1 tasks and UAT feedback.

4. **Create SemVer release** – Tag and release the Sprint 2 increment.
