# Week 4 Report – Quality and Testing Sprint

## Project

[SWP-Avito-project](https://github.com/gosha185/SWP-Avito-project) – Bonus System

**License:** [MIT](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/LICENSE)

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

- [MVP v1 Scope](https://github.com/gosha185/SWP-Avito-project/issues?q=is%3Aclosed%20is%3Aissue%20milestone%3A%22Sprint%201.%20MVP%20v1%22)

---

## Quality Documentation

- [Quality Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md)
- [Quality Requirement Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md)
- [Testing Documentation] WAIT CI
- [User Acceptance Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-acceptance-tests.md)

---

## Process Documentation

- [Definition of Done](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/definition-of-done.md)
- [Roadmap](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/roadmap.md)
- [Process Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/Process_Requirements.md)

---

## CI and Quality Gates

- [CI Pipeline](https://github.com/gosha185/SWP-Avito-project/actions) — link will be updated after CI is configured
- [Latest CI Run](https://github.com/gosha185/SWP-Avito-project/actions) — link will be updated after CI is configured
- [Branch Protection](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/branch-protection.jpg)
---

## Deployment and Release

- [Deployment URL](http://10.93.26.189:8080/)
- [SemVer Release v1.0.0](https://github.com/gosha185/SWP-Avito-project/releases/tag/v1.0.0)
- [CHANGELOG.md](https://github.com/gosha185/SWP-Avito-project/blob/main/CHANGELOG.md)
- [Public Demo Video](https://drive.google.com/file/d/1nJIoO6pt99Obi517iMFWUHnlumuFBYgq/view?usp=sharing)
---

## Customer Feedback Response

| Feedback point | Resulting PBI or issue | Status | Response |
|---|---|---|---|
| Cancel hold should return HTTP 400, not HTTP 500 | [#177](https://github.com/gosha185/SWP-Avito-project/issues/177) | To Do | Fix error handling to return appropriate HTTP status codes |
| FEFO (first-expiring first-out) should be verified | — | Planned | Will be tested in the next sprint |
| Sprint Goal should be defined before the sprint | — | Applied | Sprint Goals will be defined in advance for future sprints |

---

## UAT Results Summary

| UAT Scenario | Status | Customer Feedback | Resulting PBI |
|---|---|---|---|
| UAT-001: View current available balance | ✅ Passed | "That looks correct." | — |
| UAT-002: View points held for a concrete order | ✅ Passed | "All 700 points have been successfully reserved." | — |
| UAT-003: Award points with TTL | ✅ Passed | "Works correctly." | — |
| UAT-004: Place points on hold | ✅ Passed | "All 700 points have been successfully reserved." | — |
| UAT-005: Confirm a hold | ✅ Passed | "That's exactly how it should work." | — |
| UAT-006: Cancel a hold | ✅ Passed | "We held 200 points, cancelled the hold, and those points became available again." | — |
| UAT-007: Cancel an already cancelled hold (error handling) | ❌ Failed | "This definitely needs improvement. I'd expect HTTP 400." | [#177](https://github.com/gosha185/SWP-Avito-project/issues/177) |

---

## Quality Model

This Sprint introduces quality requirements based on **ISO/IEC 25010**:

| QR | Sub-characteristic | Scenario |
|---|---|---|
| QR-001 | Time behaviour | API response time ≤ 500ms for 95% of requests |
| QR-002 | Availability | Service availability ≥ 99.5% |
| QR-003 | Testability | Critical module coverage ≥ 30% |
| QR-004 | Integrity | Data retention and availability (30 days) |
| QR-005 | Fault tolerance | Recovery within 5 seconds without data loss |

---

## Testing Status
| Critical Module | Coverage Target | Current Coverage | Status |
|---|---|---|---|
| Balance | ≥ 30% | 87.5% | ✅ |
| Batch | ≥ 30% | 81.8% | ✅ |
| Total Coverage | — | 53.2% | ✅ |

## Current Product Status

### Completed

- ✅ US-06: View points held for a concrete order
- ✅ US-07: View current available point balance
- ✅ 18 technical tasks completed in Sprint 1
- ✅ Quality requirements defined (QR-001 to QR-005)
- ✅ UAT scenarios defined and executed (UAT-001 to UAT-007)
- ✅ Definition of Done updated for Assignment 4
- ✅ Daily Scrums introduced
- ✅ Sprint Board reorganised into 4 columns


### Next Steps

- Complete remaining MVP v1 tasks
- Finalise CI pipeline with all quality gates
- Fix error handling (HTTP 500 → HTTP 400)
- Prepare Sprint 3 scope based on customer feedback

---

## Contribution Traceability

| Team Member | Issues | PRs | Reviews | Testing/QA | Documentation |
|---|---|---|---|---|---|
| Leilia (Leilia34) | #132–#142 | ✅ | ✅ | — | ✅ |
| Ekaterina (deadnothingness) | #9, #11, #32–#36 | ✅ | ✅ | — | — |
| Stepan (Stepan4ick) | #12–#16, #28–#31, #144, #145 | ✅ | ✅ | — | — |
| Georgii (gosha185) | #24–#27, #146 | ✅ | ✅ | — | — |
| Ivan (Laplace-mt) | #10, #59, #101, #143, #177 | ✅ | ✅ | ✅ | — |

---

## Week 4 Files

- [Customer Review Transcript](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week4/customer-review-transcript.md)
- [Customer Review Summary](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week4/customer-review-summary.md)
- [Reflection](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week4/reflection.md)
- [Retrospective](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week4/retrospective.md)
- [LLM Report](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week4/llm-report.md)

---

## Screenshots

### Sprint Milestone
![Sprint Milestone](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/sprint-milestone.jpg)

### Latest CI Run
![Latest CI Run](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/latest-ci-run.jpg)

### Branch Protection
![Branch Protection](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/branch-protection.jpg)

### Coverage Report
![Coverage Report](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/coverage.jpg)

### Additional QA Check
![Additional QA Check](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/qa-check.jpg)

### SemVer Release
![SemVer Release] update

### Reviewed PR
![Reviewed PR](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/reviewed-pr.jpg)

### Product Backlog
![Product Backlog](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/product-backlog.jpg)

### Sprint Backlog
![Sprint Backlog](https://github.com/gosha185/SWP-Avito-project/blob/138-week4-report-index/reports/week4/images/sprint-backlog.jpg)
