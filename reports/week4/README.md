# Week 4 Report – Quality and Testing Sprint

## Project
[SWP-Avito-project](https://github.com/gosha185/SWP-Avito-project) – Bonus System

**License:** [MIT](LICENSE)

---

## Sprint Overview

### Sprint Goal

**Stabilise the product and establish quality foundations: complete the remaining MVP v1 features, implement automated testing and CI, and ensure the application runs reliably and passes all tests.**

### Sprint Dates

22–28 June 2026

### Sprint Milestone

[Sprint 2 – Quality and Testing](https://github.com/gosha185/SWP-Avito-project/milestone/2)

---

## Artifacts

### Backlogs
- [Product Backlog](https://github.com/users/gosha185/projects/1)
- [Sprint Backlog](https://github.com/users/gosha185/projects/2)
- [Sprint 2 Backlog](https://github.com/users/gosha185/projects/3)

### Sprint
- [Sprint Milestone – Sprint 2](https://github.com/gosha185/SWP-Avito-project/milestone/2)
- [MVP v1 Scope](https://github.com/gosha185/SWP-Avito-project/issues?q=is%3Aopen+is%3Aissue+milestone%3A%22Sprint+1+–+MVP+v1%22)

---

## Quality Documentation

- [Quality Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md)
- [Quality Requirement Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md)
- [Testing Documentation](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/testing.md)
- [User Acceptance Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-acceptance-tests.md)

---

## Process Documentation

- [Definition of Done](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/definition-of-done.md)
- [Roadmap](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/roadmap.md)
- [Process Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/Process_Requirements.md)

---

## CI and Quality Gates

- [CI Pipeline](https://github.com/gosha185/SWP-Avito-project/actions/workflows/ci.yml)
- [Latest CI Run](https://github.com/gosha185/SWP-Avito-project/actions/workflows/ci.yml)
- [Branch Protection Rules](https://github.com/gosha185/SWP-Avito-project/settings/branches)

---

## Deployment and Release

- [Deployment URL](http://10.93.26.189:8080/)
- [SemVer Release v1.0.0](https://github.com/gosha185/SWP-Avito-project/releases/tag/v1.0.0)
- [CHANGELOG.md](https://github.com/gosha185/SWP-Avito-project/blob/main/CHANGELOG.md)

---

## Customer Feedback Response

| Feedback point | Resulting PBI or issue | Status | Response |
|---|---|---|---|
| Need ability to see points expiring soon | [#8](https://github.com/gosha185/SWP-Avito-project/issues/8) | In Progress | Implemented as technical PBI with configurable expiry window |
| Need automatic release of stale holds | [#32](https://github.com/gosha185/SWP-Avito-project/issues/32) | In Progress | TTL worker implemented for automatic hold release |
| Need better test coverage | [#143](https://github.com/gosha185/SWP-Avito-project/issues/143) | In Progress | Unit and integration tests added for critical modules |
| Need CI automation | [#144](https://github.com/gosha185/SWP-Avito-project/issues/144) | In Progress | CI pipeline configured with linting, tests, coverage, QRTs |

---

## UAT Results Summary

| UAT Scenario | Status | Customer Feedback | Resulting PBI |
|---|---|---|---|
| UAT-001: View current available balance | — | — | — |
| UAT-002: View points held for a concrete order | — | — | — |
| UAT-003: Accrue points with TTL | — | — | — |

*(Results to be filled after UAT session)*

---

## Quality Model

This Sprint introduces quality requirements based on **ISO/IEC 25010**:

| QR | Sub-characteristic | Scenario |
|---|---|---|
| QR-001 | Time behaviour | API response time ≤ 500ms |
| QR-002 | Availability | Service availability ≥ 99.5% |
| QR-003 | Testability | Critical module coverage ≥ 30% |

---

## Testing Status

| Critical Module | Coverage Target | Current Coverage | Status |
|---|---|---|---|
| Balance | ≥ 30% | — | — |
| Holds | ≥ 30% | — | — |
| Accrual | ≥ 30% | — | — |
| TTL Worker | ≥ 30% | — | — |
| Ledger | ≥ 30% | — | — |

---

## Current Product Status

### Completed
- ✅ US-06: View points held for a concrete order
- ✅ US-07: View current available point balance
- ✅ 18 technical tasks completed in Sprint 1
- ✅ Quality requirements defined (QR-001, QR-002, QR-003)
- ✅ UAT scenarios defined (UAT-001, UAT-002, UAT-003)
- ✅ Definition of Done updated for Assignment 4

### In Progress
- ⏳ Remaining MVP v1 tasks (#8, #9, #12, #14, #15, #16, #23, #32)
- ⏳ Unit and integration tests (#143)
- ⏳ CI pipeline configuration (#144)
- ⏳ Quality requirement tests (#133)

### Next Steps
- Complete remaining MVP v1 tasks
- Finalise CI pipeline with all quality gates
- Conduct UAT session with customer
- Conduct Sprint Review and Retrospective

---

## Contribution Traceability

| Team Member | Issues | PRs | Reviews | Testing/QA | Documentation |
|---|---|---|---|---|---|
| Leilia (Leilia34) | #132–#142 | ✅ | ✅ | — | ✅ |
| Ekaterina (deadnothingness) | #9, #11, #32–#36 | ✅ | ✅ | — | — |
| Stepan (Stepan4ick) | #12–#16, #28–#31, #144, #145 | ✅ | ✅ | — | — |
| Georgii (gosha185) | #24–#27, #146 | ✅ | ✅ | — | — |
| Ivan (Laplace-mt) | #10, #59, #101, #143 | ✅ | ✅ | ✅ | — |

---

## Links

- [SemVer Release v1.0.0](https://github.com/gosha185/SWP-Avito-project/releases/tag/v1.0.0)

---

## Week 4 Files

- [Customer Review Summary](customer-review-summary.md)
- [Reflection](reflection.md)
- [Retrospective](retrospective.md)
- [LLM Report](llm-report.md)
- [Presentation Slides](presentation.pdf) *(optional public copy)*

---

## Screenshots

### Sprint Milestone
![Sprint Milestone](images/sprint-milestone.png)

### Latest CI Run
![Latest CI Run](images/ci-run.png)

### Branch Protection Rules
![Branch Protection Rules](images/branch-protection.png)

### Coverage Report
![Coverage Report](images/coverage.png)

### Additional QA Check
![Additional QA Check](images/qa-check.png)

### SemVer Release
![SemVer Release](images/release.png)

### Reviewed PR
![Reviewed PR](images/reviewed-pr.png)

### Product Backlog
![Product Backlog](images/product-backlog.png)

### Sprint Backlog
![Sprint Backlog](images/sprint-backlog.png)
