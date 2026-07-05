# Week 5 Report – MVP v2

## Project
[SWP-Avito-project](https://github.com/gosha185/SWP-Avito-project) – Bonus System

**License:** [MIT](LICENSE)

---

## Sprint Overview

### Sprint Goal

**Deliver MVP v2, document architecture (static, dynamic, deployment views), create 3 ADRs, and document development process with git workflow diagram.**

### Sprint Dates

29 June – 5 July 2026

### Sprint Milestone

[Sprint 3 – MVP v2](https://github.com/gosha185/SWP-Avito-project/milestone/3)

---

## Artifacts

### Backlogs
- [Product Backlog](https://github.com/users/gosha185/projects/1)
- [Sprint 1 Backlog](https://github.com/users/gosha185/projects/2)
- [Sprint 2 Backlog](https://github.com/users/gosha185/projects/3)
- [Sprint 3 Backlog](https://github.com/users/gosha185/projects/4)

### Sprint
- [MVP v2 Scope](https://github.com/gosha185/SWP-Avito-project/issues?q=is%3Aopen+is%3Aissue+milestone%3A%22Sprint+3+–+MVP+v2%22)

---

## Quality Documentation

- [Quality Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md)
- [Quality Requirement Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md)
- [Testing Documentation](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/testing.md)
- [User Acceptance Tests](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/user-acceptance-tests.md)

---

## Architecture Documentation

- [Architecture README](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/README.md)
- [Static View (Component Diagram)](https://github.com/gosha185/SWP-Avito-project/tree/main/docs/architecture/static-view)
- [Dynamic View (Sequence Diagram)](https://github.com/gosha185/SWP-Avito-project/tree/main/docs/architecture/dynamic-view)
- [Deployment View](https://github.com/gosha185/SWP-Avito-project/tree/main/docs/architecture/deployment-view)
- [ADRs](https://github.com/gosha185/SWP-Avito-project/tree/main/docs/architecture/adr)

---

## Architecture Summary

The system follows a layered architecture:

- **HTTP Layer** – handles requests and responses
- **Service/Business Logic Layer** – implements core algorithms
- **Storage Layer** – manages database operations
- **PostgreSQL** – persistent storage
- **TTL Workers** – background processes for expired points and holds

Three architecture views were documented:
- **Static View:** Component diagram showing system structure and external dependencies
- **Dynamic View:** Sequence diagram for hold/confirm/cancel flows
- **Deployment View:** Deployment diagram showing Go API service, PostgreSQL, and Docker Compose

### Architecture Decision Records (ADRs)

Three ADRs were created to document key architecture decisions:

| ADR | Decision | Quality Requirements Addressed |
|-----|----------|-------------------------------|
| [ADR-001](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/adr/ADR-001-leader-election.md) | PostgreSQL-based Leader Election for Background Workers | QR-002, QR-005 |
| [ADR-002](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/adr/ADR-002-soft-expiration-cleanup.md) | Soft Expiration and Background Cleanup of Domain Data | QR-001, QR-004 |
| [ADR-003](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/architecture/adr/ADR-003-transaction-isolation.md) | Transaction Isolation Strategy for Financial Operations | QR-001, QR-005 |

All ADRs are linked from `docs/quality-requirements.md`.

---

## Process Documentation

- [Development Process](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/development-process.md)
- [Definition of Done](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/definition-of-done.md)
- [Roadmap](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/roadmap.md)
- [Process Requirements](https://github.com/gosha185/SWP-Avito-project/blob/main/Process_Requirements.md)

---

## CI and Quality Gates

- [CI Pipeline](https://github.com/gosha185/SWP-Avito-project/actions)
- [Latest CI Run](https://github.com/gosha185/SWP-Avito-project/actions)
- [Branch Protection Rules](https://github.com/gosha185/SWP-Avito-project/settings/branches)

---

## Deployment and Release

- [Deployment URL](http://10.93.26.189:8080/)
- [SemVer Release MVP v2](https://github.com/gosha185/SWP-Avito-project/releases#release-v2.0.0) 
- [CHANGELOG.md](https://github.com/gosha185/SWP-Avito-project/blob/main/CHANGELOG.md)
- [Public Demo Video](https://drive.google.com/drive/folders/1ErYiGQRrbtJMk0coST_A4IIJJGk0-AWv)

---

## Testing and CI Status

| Gate | Status |
|------|--------|
| Linting | ✅ Passing |
| Build | ✅ Passing |
| Unit tests | ✅ Passing |
| Integration tests | ✅ Passing |
| Coverage | 53.2% |
| QRTs | ✅ Passing |
| Additional QA | ✅ Passing |

---

## Customer Feedback Response

| Feedback point | Resulting PBI or issue | Status | Response |
|---|---|---|---|
| Cancel hold should return HTTP 400, not HTTP 500 | [#177](https://github.com/gosha185/SWP-Avito-project/issues/177) | Done | Fixed in Sprint 3 — now returns HTTP 400 |
| FEFO should be verified | — | Planned | Will be tested in Sprint 4 |
| Architecture should be documented | — | Done | Architecture docs and 3 ADRs created |
| Sequence diagrams should be scenario-specific | — | Planned | Will be refactored in Sprint 4 |
| Add DoD/DoR checklists to Issue templates | — | Planned | Will be added to GitHub Issue templates |
| Use unit tests instead of manual testing for workers | — | Planned | Will be implemented in Sprint 4 |

---

## UAT Results Summary

| UAT Scenario | Status | Customer Feedback | Resulting PBI |
|---|---|---|---|
| UAT-001: View balance | ✅ Passed | — | — |
| UAT-002: View holds | ✅ Passed | — | — |
| UAT-003: Award points | ✅ Passed | — | — |
| UAT-004: Place hold | ✅ Passed | — | — |
| UAT-005: Confirm hold | ✅ Passed | — | — |
| UAT-006: Cancel hold | ✅ Passed | — | — |
| UAT-007: Cancel already cancelled hold | ✅ Passed | Fixed in Sprint 3 — now returns HTTP 400 | [#177](https://github.com/gosha185/SWP-Avito-project/issues/177) |
| UAT-008: TTL worker — expired points burned | ✅ Passed | "Workers demonstrated working correctly." | — |
| UAT-009: TTL worker — stale holds released | ✅ Passed | "Workers demonstrated working correctly." | — |

---

## Current Product Status

### Completed
- ✅ MVP v1 (Sprint 1)
- ✅ Quality and Testing (Sprint 2)
- ✅ TTL workers for expired points and holds (#9, #32)
- ✅ View points expiring in given time (#8)
- ✅ View all current held points (#23)
- ✅ Architecture documentation (static, dynamic, deployment views)
- ✅ 3 ADRs created (ADR-001, ADR-002, ADR-003)
- ✅ Development process documentation
- ✅ UAT scenarios (9 scenarios, all passed)
- ✅ Sprint Review with customer (4 July 2026)

### In Progress
- testing for MVP v2

### Next Steps
- Complete SemVer release for MVP v2
- Refactor sequence diagrams into scenario-specific views
- Add DoD/DoR checklists to Issue templates
- Replace manual worker testing with unit tests
- Sprint 4: Post-MVP improvements (flexible TTL, resource management)

---

## Week 5 Files

- [Sprint Review Summary](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week5/sprint-review-summary.md)
- [Sprint Review Transcript](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week5/sprint-review-transcript.md)
- [Reflection](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week5/reflection.md)
- [Retrospective](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week5/retrospective.md)
- [LLM Report](https://github.com/gosha185/SWP-Avito-project/blob/main/reports/week5/llm-report.md)
- [Hosted Docs Site](https://gosha185.github.io/SWP-Avito-project/)
---

## Screenshots

### Sprint Milestone
![Sprint Milestone]()

### Board/Project Workflow View
![Board View]()

### Latest CI Run
![Latest CI Run]()

### SemVer Release
![SemVer Release]()

### Reviewed PR
![Reviewed PR]()

### Hosted Docs Site
![Hosted Docs]()
