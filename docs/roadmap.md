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

## Sprint 2 – Quality and Testing (Current)

**Dates:** 22–28 June 2026  
**Sprint Goal:** Stabilise the product and establish quality foundations: complete the remaining MVP v1 features, implement automated testing and CI, and ensure the application runs reliably and passes all tests.

**Milestone:** [Sprint 2 – Quality and Testing](https://github.com/gosha185/SWP-Avito-project/milestone/2)

**Planned items:**

| # | Title | Type | SP |
|---|-------|------|----|
| #8 | Technical: Make possible to see points expiring in giving time | Technical | 13 |
| #9 | Technical: Make the TTL worker manage expired points batches | Technical | 8 |
| #23 | Technical: Make possible to see all current held points | Technical | 1 |
| #32 | Technical: Make the TTL worker manage expired held points | Technical | 8 |
| #12 | Public deployment of MVP v1 | Infrastructure | 2 |
| #14 | Update Postman collection | Infrastructure | 5 |
| #15 | SemVer release v1.0.0 | Infrastructure | 2 |
| #16 | CHANGELOG.md update | Infrastructure | 1 |
| #132 | Docs: Add quality requirements | Documentation | 5 |
| #133 | Docs: Add quality requirement tests | Documentation | 3 |
| #134 | Docs: Add testing documentation | Documentation | 3 |
| #135 | Docs: Add user acceptance tests | Documentation | 3 |
| #136 | Docs: Update Definition of Done | Documentation | 2 |
| #137 | Docs: Update roadmap | Documentation | 1 |
| #138 | Docs: Add Week 4 report index | Documentation | 5 |
| #139 | Docs: Add customer review summary | Documentation | 3 |
| #140 | Docs: Add Week 4 reflection | Documentation | 2 |
| #141 | Docs: Add Week 4 retrospective | Documentation | 2 |
| #142 | Docs: Add LLM usage report | Documentation | 1 |
| #143 | Test: Add unit and integration tests | Testing | 15 |
| #144 | Infra: Configure CI pipeline | Infrastructure | 8 |
| #145 | refactor: fix handlers | Technical | 5 |
| #146 | refactor: Remove unused files | Technical | 3 |

**Total Sprint Size:** 101 Story Points

**Focus areas:**
- Complete remaining MVP v1 features
- Implement automated testing (unit, integration, concurrency)
- Configure CI pipeline with quality gates
- Define quality requirements and UAT scenarios
- Documentation and reporting

---

## Sprint 3 – Post-MVP Improvements (Planned)

**Dates:** 29   June – 5 July 2026  
**Sprint Goal:** Enhance system flexibility and observability for production readiness.

**Milestone:** [Sprint 3 – Post-MVP Improvements](https://github.com/gosha185/SWP-Avito-project/milestone/3)

**Planned items:**
- Flexible TTL per accrual – allow specifying TTL days in the accrual request
- Resource management:
  - Connection pool tuning (MaxOpenConns/MaxIdleConns based on load)
  - Graceful shutdown (SIGTERM handling with in-flight request completion)
  - Prometheus metrics (goroutines, memory, DB pool stats, request latency)
- OpenAPI spec + deployment documentation
- Optional Redis caching for frequently accessed balances (invalidation strategy required)

---

## Sprint 4 – Database Scaling & Final Polish (Planned)

**Dates:** 6–12 July 2026  
**Sprint Goal:** Prepare database layer for production-scale usage.

**Milestone:** -

**Planned items:**
- Read replicas – route GET /balance and history queries to replicas
- Sharding – distribute data by user_id (Citus or manual partitioning)
- Connection pool monitoring – track WaitCount and WaitDuration to detect bottlenecks
- Final testing, documentation, and deployment polish
- Remaining Could Have features (if time permits)
