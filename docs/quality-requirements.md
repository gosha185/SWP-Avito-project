# Quality Requirements

This document defines measurable quality requirements for the Avito Bonus Points Service using ISO/IEC 25010 quality
characteristics.

---

## QR-001: API response time

**ISO/IEC 25010 sub-characteristic:** Time behaviour

**Scenario:** When an end user sends a request to any API endpoint under normal production-like load (up to 100
concurrent users), the system shall return a response within 500ms for 95% of requests.

**Why this matters:** Users need quick feedback when checking their balance, viewing holds, or performing accrual
operations. Slow responses degrade user experience and may cause timeouts in partner services (order processing, payment
systems).

**Measurable per endpoint:**

| Endpoint                                   | Expected response time (95th percentile) under 100 concurrent users |
|:-------------------------------------------|:--------------------------------------------------------------------|
| `GET /v1/users/{user_id}/balance`          | ≤ 500ms                                                             |
| `GET /v1/users/{user_id}/holds`            | ≤ 500ms                                                             |
| `GET /v1/users/{user_id}/holds/{order_id}` | ≤ 500ms                                                             |
| `POST /v1/accrue`                          | ≤ 500ms                                                             |
| `POST /v1/hold`                            | ≤ 500ms                                                             |
| `POST /v1/confirm`                         | ≤ 500ms                                                             |
| `POST /v1/cancel`                          | ≤ 500ms                                                             |

**Related ADRs:**

- [ADR-001](architecture/adr/ADR-001-leader-election.md) — Prevents duplicate background processing, reducing
  unnecessary database load.
- [ADR-002](architecture/adr/ADR-002-soft-expiration-cleanup.md) — Periodic cleanup keeps database tables smaller,
  improving query performance.
- [ADR-003](architecture/adr/ADR-003-transaction-isolation.md) — Efficient locking strategy minimizes contention during
  concurrent requests.

**Linked quality requirement tests:**

- [QRT-001](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md#qrt-001-api-response-time)

---

## QR-002: Service availability

**ISO/IEC 25010 sub-characteristic:** Availability

**Scenario:** When the service is deployed in production, the system shall remain available and respond to health checks
with HTTP 200 OK for at least 99.5% of the time over any 7-day period.

**Why this matters:** The bonus system must be reliable for users and partner services. Downtime causes lost
transactions, double processing, and user frustration.

**Related ADRs:**

- [ADR-001](architecture/adr/ADR-001-leader-election.md) — Automatic leader failover ensures background workers continue
  running if one instance fails.
- [ADR-003](architecture/adr/ADR-003-transaction-isolation.md) — ACID guarantees prevent data corruption that could
  require manual recovery and downtime.

**Linked quality requirement tests:**

- [QRT-002](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md#qrt-002-api-response-time)

---

## QR-003: Critical module testability

**ISO/IEC 25010 sub-characteristic:** Testability

**Scenario:** When a developer changes a critical product module under the standard CI environment, the module shall
have automated unit tests that achieve at least 30% line coverage for that module.

**Why this matters:** Critical product logic (balance operations, holds, accrual, TTL workers) must be directly
verifiable so defects can be detected before merge. This reduces regression risk and increases developer confidence when
making changes.

**Related ADRs:**

- None.

**Linked quality requirement tests:**

- [QRT-003](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md#qrt-003-api-response-time)

---

## QR-004: Data retention and availability

**ISO/IEC 25010 sub-characteristic:** Integrity

**Scenario:** When a user accesses their transaction history, the system shall return all transactions from the last 30
days without data loss or corruption.

**Why this matters:** Users need to trust that their bonus points are accurately tracked. Data loss or corruption
undermines user confidence and creates disputes.

**Measurable:**

| Requirement           | Expected                                    |
|:----------------------|:--------------------------------------------|
| Data retention period | Minimum 30 days                             |
| Data integrity        | No lost or corrupted transactions           |
| Query completeness    | All transactions in the period are returned |

**Related ADRs:**

- [ADR-002](architecture/adr/ADR-002-soft-expiration-cleanup.md) — Soft expiration preserves data during the retention
  period before cleanup.

**Linked quality requirement tests:**

- [QRT-004](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md#qrt-004-data-retention-and-integrity)

---

## QR-005: Fault tolerance and recovery

**ISO/IEC 25010 sub-characteristic:** Fault tolerance

**Scenario:** When a database connection fails during a critical operation (accrual, hold, confirm, cancel), the system
shall recover and complete the operation within 5 seconds without data loss.

**Why this matters:** Critical operations must be resilient to transient failures. Users should not lose points or
experience double processing due to temporary issues.

**Related ADRs:**

- [ADR-001](architecture/adr/ADR-001-leader-election.md) — Leader failover ensures background workers recover after
  instance failure.
- [ADR-003](architecture/adr/ADR-003-transaction-isolation.md) — ACID guarantees and transaction timeouts prevent
  partial updates and ensure recovery.

**Linked quality requirement tests:**

- [QRT-005](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirement-tests.md#qrt-005-api-response-time)
