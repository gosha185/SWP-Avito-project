# Product Roadmap

## Sprint 1. MVP v1

**Dates:** 15–21 June 2026  
**Sprint Goal:** Deliver core bonus operations: users can view their available balance and see points held for orders. The system supports idempotent accrual, two-phase redemption with FEFO, and automatic release of stale holds.

**Milestone:** [Sprint 1. MVP v1](https://github.com/gosha185/SWP-Avito-project/milestone/1)

**Planned items:**
- US-06: View points held for a concrete order
- US-07: View current available point balance
- 23 technical and infrastructure PBIs (accrual, holds, TTL, locking, PostgreSQL, deployment)

---

## Sprint 2. MVP v2

**Dates:** 22–28 June 2026  
**Sprint Goal:** Enhance system flexibility and observability for production readiness.

**Milestone:**  [Sprint 2. MVP v2](https://github.com/gosha185/SWP-Avito-project/milestone/2)

**Planned items:**
- **Flexible TTL per accrual** – allow specifying TTL days in the accrual request
- **Resource management:**
  - Connection pool tuning (MaxOpenConns/MaxIdleConns based on load)
  - Graceful shutdown (SIGTERM handling with in-flight request completion)
  - Prometheus metrics (goroutines, memory, DB pool stats, request latency)
- **OpenAPI spec + deployment documentation** – finalise API documentation
- **Optional Redis caching** – for frequently accessed balances (requires invalidation strategy)

---

## Sprint 3. MVP v3. Database Scaling & Final Polish (Week 3)

**Dates:** 29 June – 5 July 2026  
**Sprint Goal:** Prepare database layer for production-scale usage.

**Milestone:**  [Sprint 3. MVP v3](https://github.com/gosha185/SWP-Avito-project/milestone/3)

**Planned items:**
- **Read replicas** – route GET /balance and history queries to replicas
- **Sharding** – distribute data by `user_id` (Citus or manual partitioning)
- **Connection pool monitoring** – track `WaitCount` and `WaitDuration` to detect bottlenecks
- Final testing, documentation, and deployment polish
- Remaining Could Have features (if time permits)
