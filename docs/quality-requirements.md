# Quality Requirements

This document defines measurable quality requirements for the Avito Bonus Points Service using ISO/IEC 25010 quality characteristics.

---

## QR-001: API response time

**ISO/IEC 25010 sub-characteristic:** Time behaviour

**Scenario:** When an end user sends a request to any API endpoint under normal production-like load (up to 100 concurrent users), the system shall return a response within 500ms for 95% of requests.

**Why this matters:** Users need quick feedback when checking their balance, viewing holds, or performing accrual operations. Slow responses degrade user experience and may cause timeouts in partner services (order processing, payment systems).

**Measurable per endpoint:**
| Endpoint | Expected response time (95th percentile) | 
|----------|-------------------------------------------|
| `GET /v1/users/{user_id}/balance` | ≤ 500ms | 
| `GET /v1/users/{user_id}/holds` | ≤ 500ms | 
| `GET /v1/users/{user_id}/holds/{order_id}` | ≤ 500ms | 
| `POST /v1/accrue` | ≤ 500ms | 
| `POST /v1/hold` | ≤ 500ms | 
| `POST /v1/confirm` | ≤ 500ms | 
| `POST /v1/cancel` | ≤ 500ms |

**Linked quality requirement tests:** [QRT-001](quality-requirement-tests.md#qrt-001-api-response-time)

---

## QR-002: Service availability

**ISO/IEC 25010 sub-characteristic:** Availability

**Scenario:** When the service is deployed in production, the system shall remain available and respond to health checks with HTTP 200 OK for at least 99.5% of the time over any 7-day period.

**Why this matters:** The bonus system must be reliable for users and partner services. Downtime causes lost transactions, double processing, and user frustration.

**Linked quality requirement tests:** [QRT-002](quality-requirement-tests.md#qrt-002-health-check-availability)

---

## QR-003: Critical module testability

**ISO/IEC 25010 sub-characteristic:** Testability

**Scenario:** When a developer changes a critical product module under the standard CI environment, the module shall have automated unit tests that achieve at least 30% line coverage for that module.

**Why this matters:** Critical product logic (balance operations, holds, accrual, TTL workers) must be directly verifiable so defects can be detected before merge. This reduces regression risk and increases developer confidence when making changes.

**Linked quality requirement tests:** [QRT-003](quality-requirement-tests.md#qrt-003-critical-module-unit-coverage)

---

## QR-004: Data retention and availability

**ISO/IEC 25010 sub-characteristic:** Integrity

**Scenario:** When a user accesses their transaction history, the system shall return all transactions from the last 30 days without data loss or corruption.

**Why this matters:** Users need to trust that their bonus points are accurately tracked. Data loss or corruption undermines user confidence and creates disputes.

**Measurable:**
| Requirement | Expected |
|-------------|----------|
| Data retention period | Minimum 30 days |
| Data integrity | No lost or corrupted transactions |
| Query completeness | All transactions in the period are returned |

**Linked quality requirement tests:** [QRT-004](quality-requirement-tests.md#qrt-004-data-retention-and-integrity)

---

## QR-005: Fault tolerance and recovery

**ISO/IEC 25010 sub-characteristic:** Fault tolerance

**Scenario:** When a database connection fails during a critical operation (accrual, hold, confirm, cancel), the system shall recover and complete the operation within 5 seconds without data loss.

**Why this matters:** Critical operations must be resilient to transient failures. Users should not lose points or experience double processing due to temporary issues.

**Linked quality requirement tests:** [QRT-005](quality-requirement-tests.md#qrt-005-fault-tolerance-and-recovery)
