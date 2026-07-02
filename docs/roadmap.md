# Product Roadmap

## Sprint 1 – MVP v1 (Completed)

**Dates:** 15–21 June 2026  
**Sprint Goal:** Deliver core bonus operations: users can view their available balance and see points held for orders. The system supports idempotent accrual, two-phase redemption with FEFO, and automatic release of stale holds.

**Milestone:** [Sprint 1 – MVP v1](https://github.com/gosha185/SWP-Avito-project/milestone/1)

**Delivered:**
- US-06: View points held for a concrete order
- US-07: View current available point balance
- 18 technical and infrastructure PBIs completed
- 7 tasks moved to Sprint 2

---

## Sprint 2 – Quality and Testing (Completed)

**Dates:** 22–28 June 2026  
**Sprint Goal:** Stabilise the product and establish quality foundations: complete the remaining MVP v1 features, implement automated testing and CI, and ensure the application runs reliably and passes all tests.

**Milestone:** [Sprint 2 – Quality and Testing](https://github.com/gosha185/SWP-Avito-project/milestone/2)

**Completed:**
- Quality requirements (QR-001 to QR-005)
- Quality requirement tests (QRT-001 to QRT-005)
- Unit and integration tests (coverage 53.2%)
- CI pipeline with linting, tests, coverage, QRTs, additional QA
- UAT scenarios (7 scenarios, 6 passed, 1 failed)
- Updated Definition of Done
- Updated Roadmap
- Week 4 reports and documentation

**Focus areas:**
- Quality requirements definition
- Automated testing (unit, integration, concurrency)
- CI pipeline configuration with quality gates
- UAT with customer
- Documentation and reporting

---

## Sprint 3 – MVP v2 (Current)

**Dates:** 29 June – 5 July 2026  
**Sprint Goal:** Deliver MVP v2: complete remaining features, document architecture (static, dynamic, deployment views), create 3 ADRs, and document development process with git workflow diagram.

**Milestone:** [Sprint 3 – MVP v2](https://github.com/gosha185/SWP-Avito-project/milestone/3)

**Planned items:**

| # | Title | Type | SP | Assignee |
|---|-------|------|----|----------|
| #8 | Technical: Make possible to see points expiring in giving time | Technical | 13 | Georgii, Ekaterina, Stepan |
| #9 | Technical: Make the TTL worker manage expired points batches | Technical | 8 | Ekaterina |
| #23 | Technical: Make possible to see all current held points | Technical | 1 | Georgii, Ekaterina, Stepan |
| #32 | Technical: Make the TTL worker manage expired held points | Technical | 8 | Ekaterina |
| #177 | Bug: Cancel hold returns HTTP 500 (should be HTTP 400) | Bug | 3 | Stepan |
| #— | Architecture: Static view (component diagram) | Architecture | 5 | Ekaterina |
| #— | Architecture: Dynamic view (sequence diagram) | Architecture | 5 | Georgii |
| #— | Architecture: Deployment view | Architecture | 5 | Stepan |
| #— | ADRs: Create 3 Architecture Decision Records | Architecture | 5 | Georgii |
| #198 | Docs: Development process documentation with git workflow | Documentation | 3 | Leilia |
| #— | Docs: Quality requirements link to ADRs | Documentation | 2 | Leilia |
| #— | Docs: Update testing.md for MVP v2 | Documentation | 2 | Ivan |
| #— | Docs: Add 2 new UAT scenarios for MVP v2 | Documentation | 3 | Ivan |

**Total Sprint Size:** ~70 Story Points

**Focus areas:**
- Complete remaining MVP v1 features (TTL workers, expiring points, held points)
- Document architecture (static, dynamic, deployment views)
- Create 3 ADRs (Architecture Decision Records)
- Document development process with Mermaid gitGraph diagram
- Extend testing and CI for MVP v2
- Deploy and release MVP v2
- Conduct UAT and Sprint Review with customer

---

## Sprint 4 – Post-MVP Improvements (Planned)

**Dates:** 6–12 July 2026  
**Sprint Goal:** Enhance system flexibility and observability for production readiness.

**Milestone:** [Sprint 4 – Post-MVP Improvements](https://github.com/gosha185/SWP-Avito-project/milestone/4)

**Planned items:**
- Flexible TTL per accrual – allow specifying TTL days in the accrual request
- Resource management:
  - Connection pool tuning (MaxOpenConns/MaxIdleConns based on load)
  - Graceful shutdown (SIGTERM handling with in-flight request completion)
  - Prometheus metrics (goroutines, memory, DB pool stats, request latency)
- OpenAPI spec + deployment documentation
- Optional Redis caching for frequently accessed balances (invalidation strategy required)
