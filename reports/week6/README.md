# Week 6 Report – Trial Release and Handover Preparation

## Project
[SWP-Avito-project](https://github.com/gosha185/SWP-Avito-project) – Bonus System

**License:** [MIT](https://github.com/gosha185/SWP-Avito-project/blob/main/LICENSE)

---

## Sprint Overview

### Sprint Goal

**Deliver a stable trial release, complete customer-facing documentation (README, customer-handover, CONTRIBUTING, AGENTS), and conduct a transition-readiness meeting with the customer.**

### Sprint Dates

6–12 July 2026

### Sprint Milestone

[Sprint 4 – Trial Release](https://github.com/gosha185/SWP-Avito-project/milestone/4)

---

## Artifacts

### Backlogs
- [Product Backlog](https://github.com/users/gosha185/projects/1)
- [Sprint 4 Backlog](https://github.com/users/gosha185/projects/4)

### Sprint
- [Sprint 4 Milestone](https://github.com/gosha185/SWP-Avito-project/milestone/4)

---

## Documentation Links

- [README.md](https://github.com/gosha185/SWP-Avito-project/blob/main/README.md)
- [CONTRIBUTING.md](https://github.com/gosha185/SWP-Avito-project/blob/main/CONTRIBUTING.md)
- [AGENTS.md](https://github.com/gosha185/SWP-Avito-project/blob/main/AGENTS.md)
- [docs/customer-handover.md](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/customer-handover.md)
- [Hosted Documentation Site](https://gosha185.github.io/SWP-Avito-project/)
- [docs/roadmap.md](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/roadmap.md)
- [CHANGELOG.md](https://github.com/gosha185/SWP-Avito-project/blob/main/CHANGELOG.md)

---

## Quality, Testing and Architecture Documentation

- [Quality Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md)
- [Quality Requirement Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md)
- [Testing Documentation](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/testing.md)
- [User Acceptance Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-acceptance-tests.md)
- [Architecture README](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/README.md)
- [Static View (Component Diagram)](https://github.com/gosha185/SWP-Avito-project/tree/main/docs/architecture/static-view)
- [Dynamic View (Sequence Diagram)](https://github.com/gosha185/SWP-Avito-project/tree/main/docs/architecture/dynamic-view)
- [Deployment View](https://github.com/gosha185/SWP-Avito-project/tree/main/docs/architecture/deployment-view)
- [ADRs](https://github.com/gosha185/SWP-Avito-project/tree/main/docs/architecture/adr)
- [Process Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/Process_Requirements.md)
- [Definition of Done](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/definition-of-done.md)
- [Development Process](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/development-process.md)

---

## Deployment and Release

- [Deployment URL](http://10.93.26.189:8080/)
- [Week 6 Trial Release](https://github.com/gosha185/SWP-Avito-project/releases/tag/v0.5.0)
- [CHANGELOG.md](https://github.com/gosha185/SWP-Avito-project/blob/main/CHANGELOG.md)

---

## Customer Feedback Response

| Feedback point | Resulting PBI | Status | Response |
|---|---|---|---|
| HTTP 500 should be HTTP 400 for duplicate hold operations | [#177](https://github.com/gosha185/SWP-Avito-project/issues/177) | Done | Fixed in Sprint 4 |
| Idempotency retry should return 200 instead of 409 | — | Planned | Will be reviewed in Sprint 5 |
| UAT test status for user stories should be included in handover | — | Planned | Will be added to handover docs |

---

## UAT Results Summary

| UAT Scenario | Status | Customer Feedback |
|---|---|---|
| UAT-001: View balance | ✅ Passed | — |
| UAT-002: View holds | ✅ Passed | — |
| UAT-003: Award points | ✅ Passed | — |
| UAT-004: Place hold | ✅ Passed | — |
| UAT-005: Confirm hold | ✅ Passed | — |
| UAT-006: Cancel hold | ✅ Passed | — |
| UAT-007: Cancel already cancelled hold | ✅ Passed | Approved by customer in Sprint 4 |
| UAT-008: TTL worker — expired points burned | ✅ Passed | — |
| UAT-009: TTL worker — stale holds released | ✅ Passed | — |

---

## Customer-Facing Documentation Review

The customer reviewed the following documentation during the Week 6 meeting:

- `README.md` — approved as the main entry point
- `docs/customer-handover.md` — approved, handover status documented
- `CONTRIBUTING.md` — approved
- `AGENTS.md` — approved

**Customer feedback:** Documentation is comprehensive and ready for handover.

---

## Transition-Readiness Summary

**Handover level reached:** Ready for independent use

**Customer confirmation status:** Accepted with follow-up items

**What must still happen in Week 7:**
- Final handover meeting (Thursday)
- MVP v3 release
- Demo Day preparation

---

## Current Product Status

### Completed
- ✅ Trial release deployed
- ✅ README.md updated as main entry point
- ✅ docs/customer-handover.md created
- ✅ CONTRIBUTING.md created
- ✅ AGENTS.md created
- ✅ DoR/DoD checklists added to issue templates
- ✅ Week 6 Sprint Review conducted
- ✅ Customer feedback captured

### In Progress
- ⏳ Final handover meeting (Thursday)
- ⏳ MVP v3 release

### Next Steps
- Final handover meeting with customer
- Create SemVer release for MVP v3
- Prepare Week 7 reports and Demo Day presentation

---

## Screenshots

### Sprint Milestone
![Sprint Milestone](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/images/sprint-milestone.jpg)

### Sprint 4 Backlog
![Sprint 4 Backlog](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/images/sprint4-backlog.jpg)

### Latest CI Run
![Latest CI Run](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/images/ci-run.jpg)

### Week 6 Trial Release
![Week 6 Trial Release](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/images/release-week6.jpg)

### Reviewed PR
![Reviewed PR](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/images/reviewed-pr.jpg)

### Hosted Docs Site
![Hosted Docs](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/images/hosted-docs.jpg)

### Trial Deployment
![Trial Deployment](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/images/trial-deployment.jpg)

---

## Week 6 Files

- [Sprint Review Summary](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/sprint-review-summary.md)
- [Sprint Review Transcript](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/sprint-review-transcript.md)
- [Reflection](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/reflection.md)
- [Retrospective](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/retrospective.md)
- [LLM Report](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week6/llm-report.md)

---

## Contribution Traceability

| Team Member | Issues | PRs | Reviews | Testing/QA | Documentation | Transition |
|---|---|---|---|---|---|---|
| Leilia (Leilia34) | [#261](https://github.com/gosha185/SWP-Avito-project/issues/261), [#268](https://github.com/gosha185/SWP-Avito-project/issues/268), [#269](https://github.com/gosha185/SWP-Avito-project/issues/269), [#270](https://github.com/gosha185/SWP-Avito-project/issues/270), [#271](https://github.com/gosha185/SWP-Avito-project/issues/271), [#272](https://github.com/gosha185/SWP-Avito-project/issues/272), [#273](https://github.com/gosha185/SWP-Avito-project/issues/273), [#285](https://github.com/gosha185/SWP-Avito-project/issues/285), [#298](https://github.com/gosha185/SWP-Avito-project/issues/298) | ✅ | ✅ | — | ✅ | ✅ |
| Ekaterina (deadnothingness) | [#287](https://github.com/gosha185/SWP-Avito-project/issues/287), [#299](https://github.com/gosha185/SWP-Avito-project/issues/299) | ✅ | ✅ | — | ✅ | — |
| Stepan (Stepan4ick) | [#263](https://github.com/gosha185/SWP-Avito-project/issues/263), [#264](https://github.com/gosha185/SWP-Avito-project/issues/264), [#265](https://github.com/gosha185/SWP-Avito-project/issues/265), [#266](https://github.com/gosha185/SWP-Avito-project/issues/266), [#267](https://github.com/gosha185/SWP-Avito-project/issues/267), [#275](https://github.com/gosha185/SWP-Avito-project/issues/275), [#301](https://github.com/gosha185/SWP-Avito-project/issues/301) | ✅ | ✅ | — | ✅ | ✅ |
| Georgii (gosha185) | [#291](https://github.com/gosha185/SWP-Avito-project/issues/291), [#289](https://github.com/gosha185/SWP-Avito-project/issues/289) | ✅ | ✅ | ✅ | — | — |
| Ivan (Laplace-mt) |  [#246](https://github.com/gosha185/SWP-Avito-project/issues/246), [#281](https://github.com/gosha185/SWP-Avito-project/issues/281) | ✅ | ✅ | ✅ | — | — |
