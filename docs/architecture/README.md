# Architecture Documentation

## Purpose

This document is the main entry point for the project's architecture documentation. It describes the current software
architecture of the Bonus Service, links the maintained architecture views, and provides traceability to the project's
Architecture Decision Records (ADRs).

---

## Overview

The Bonus Service is a Go-based REST API responsible for managing user bonus points. The system follows a layered
architecture consisting of an HTTP API layer, business logic layer, storage layer, and PostgreSQL database.

The architecture is designed to support:

- **Horizontal scaling** — stateless API layer; multiple instances can run behind a load balancer.
- **Single-execution background workers** — leader-elected workers ensure maintenance tasks run once across all
  instances.
- **Data consistency** — ACID transactions and pessimistic locking for financial operations.
- **Auditability** — immutable ledger for all operations.

Background maintenance tasks (TTL expiry, batch expiry, and cleanup) are executed by leader-elected workers to support
horizontal scaling without duplicate processing.

---

## Architecture Views

### Static View

Directory: [static-view/](static-view/)

**Component Diagram:**

- Rendered SVG: [static-view/component-diagram.png](static-view/component-diagram.png)
- Source: [static-view/component-diagram.puml](static-view/component-diagram.puml)

The component diagram shows the main internal components of the system and their relationships:

- **External Service** — the caller (e.g., order service, payment service) that invokes the Bonus API via HTTP.
- **Router (HTTP Layer)** — handles incoming requests, routes them to the appropriate service methods, and applies
  validation and idempotency checks.
- **Service (Service Layer)** — implements business logic: accrual, hold, confirm, cancel, FEFO batch selection, and
  balance calculations.
- **Leader Service (Service Layer)** — manages leader election and controls which instance runs background workers.
- **BalanceRepo / BatchRepo / HoldRepo / HoldBatchRepo / LedgerRepo (Storage Layer)** — provide database operations for
  each domain entity. The storage layer encapsulates all SQL queries and transaction management.
- **PostgreSQL Database** — persistent data store with ACID transactions, row-level locking (`SELECT FOR UPDATE`), and
  indexes for FEFO queries.
- **Workers** — background maintenance tasks (TTL worker, batch expiry, batch cleanup, hold cleanup). Workers are only
  started on the elected leader instance.

- **Coupling and cohesion:**

- Clear separation of responsibilities between API, business logic, and persistence layers.
- Low coupling through well-defined package boundaries.
- High cohesion within each layer (each package has a single responsibility).
- The layered structure supports independent testing and maintainability.

**Quality support:**

- QR-001 (Performance) — indexes for FEFO queries; mass operations for workers.
- QR-003 (Testability) — separate layers are independently testable.
- QR-004 (Integrity) — ledger ensures auditability.

---

### Dynamic View

Directory: [dynamic-view/](dynamic-view/)

**Sequence Diagram (Hold → Confirm):**

- Rendered PNG: [dynamic-view/sequence-diagram.png](dynamic-view/sequence-diagram.png)
- Source: [dynamic-view/sequence-diagram.puml](dynamic-view/sequence-diagram.puml)

The sequence diagram illustrates the main interaction flows between components for user-initiated operations and
background maintenance tasks.

### User-initiated operations

The diagram shows how the Router passes requests to the Service, which then coordinates multiple repositories within a
single database transaction:

1. **Accrual** — the Service creates a new batch of points via `BatchRepo`, updates the user balance via `BalanceRepo`,
   and records the operation in the ledger via `LedgerRepo`.

2. **Hold** — the Service locks the user balance (`SELECT FOR UPDATE`), selects batches by FEFO order using `BatchRepo`,
   creates a hold record via `HoldRepo`, links the hold to the selected batches via `HoldBatchRepo`, updates the
   balance (`available--`, `held++`), and writes a `hold` entry to the ledger.

3. **Confirm** — the Service finds the active hold via `HoldRepo`, updates its status to `confirmed`, updates the
   balance (`held--`), and writes a `confirm` entry to the ledger.

4. **Cancel** — the Service finds the active hold, retrieves the associated hold-batch links via `HoldBatchRepo`,
   returns points to the original batches via `BatchRepo.IncreaseBatchRemaining()`, updates the balance (`available++`,
   `held--`), updates the hold status to `cancelled`, and writes a `cancel` entry to the ledger.

5. **Get Balance / Get History** — the Service reads data via `BalanceRepo` and `LedgerRepo` without modifying state.

### Background worker flows

The diagram also shows how the Leader Service determines whether the current instance is the leader. If it is, the
Workers are started and perform maintenance operations:

- **Subtract expired points from balance** — the Service uses `BatchRepo` to find expired batches and updates the
  balance via `BalanceRepo`.
- **Cancel expired holds** — the Service finds expired holds via `HoldRepo`, returns points to batches via
  `HoldBatchRepo` and `BatchRepo`, updates the balance, and writes to the ledger.
- **Delete old batches** — the Workers call `BatchRepo` to remove obsolete batches.
- **Delete old holds** — the Workers call `HoldRepo` to remove completed holds older than the retention period.

This flow demonstrates:

- How FEFO batch selection works.
- How transactions and row-level locking (`SELECT FOR UPDATE`) ensure consistency.
- How the ledger records every financial operation.
- How leader election prevents duplicate execution of background workers.

**Architecture decisions illustrated:**

- ADR-003 (Transaction Isolation) — `SELECT FOR UPDATE` and explicit `COMMIT`.
- ADR-002 (Soft Expiration) — batches have TTL; expired batches are processed by background workers.

---

### Deployment View

Directory: [deployment-view/](deployment-view/)

**Deployment Diagram:**

- Rendered PNG: [deployment-view/deployment-diagram.png](deployment-view/deployment-diagram.png)
- Source: [deployment-view/deployment-diagram.puml](deployment-view/deployment-diagram.puml)

The deployment diagram shows the physical deployment structure of the system:

- **End-User Devices** — clients (e.g., frontend applications, mobile apps) that users interact with directly.
- **API User Servers** — external systems (e.g., order service, payment service) that invoke the Bonus Service API on
  behalf of users or as part of business workflows.
- **API Servers** — one or more server instances running the Go API service. The service is stateless, allowing
  horizontal scaling.
    - **Usual Servers** — regular API instances that handle user requests.
    - **Leader Server** — the instance that has acquired the leader election lock and also runs background workers.
- **Workers** — background maintenance tasks (TTL, expiry, cleanup). Workers are co-located with the API but are only
  active on the leader server.
- **PostgreSQL Database** — the primary data store, running in a separate container. All API instances and workers
  connect to the same database.

**Communication paths:**

- End-user devices communicate with the Bonus Service via HTTP through the API User Servers.
- API User Servers invoke the Bonus Service API endpoints (`/v1/accrual`, `/v1/hold`, `/v1/balance`, etc.).
- All API instances connect to the PostgreSQL database for read and write operations.
- The leader server runs the background workers, which also connect to the PostgreSQL database.

**Why this deployment model was chosen:**

- Simple to start and operate (`docker compose up`).
- All components in one host — reduces operational complexity for MVP.
- The stateless API layer allows horizontal scaling without additional configuration.
- Leader election ensures that background workers run on only one instance, preventing duplicate processing.

**Deployment considerations:**

- The service is stateless, so multiple API instances can be deployed behind a load balancer.
- Leader election ensures that background workers run on only one instance.
- Future scaling may include deploying multiple API instances while preserving single-worker execution through leader
  election.

---

## Architecture Decision Records

The following Architecture Decision Records document major architectural decisions and their relationship to the
project's quality requirements.

| ADR                                               | Title                                                   | Status   | Related QR     |
|---------------------------------------------------|---------------------------------------------------------|----------|----------------|
| [ADR-001](adr/ADR-001-leader-election.md)         | PostgreSQL-based Leader Election for Background Workers | Accepted | QR-002, QR-005 |
| [ADR-002](adr/ADR-002-soft-expiration-cleanup.md) | Soft Expiration and Background Cleanup of Domain Data   | Accepted | QR-001, QR-004 |
| [ADR-003](adr/ADR-003-transaction-isolation.md)   | Transaction Isolation Strategy for Financial Operations | Accepted | QR-001, QR-005 |

All ADRs are linked to relevant quality requirements in [quality-requirements.md](../quality-requirements.md).

---

## Quality Support

The architecture supports the following quality requirements:

| QR                       | How the architecture supports it                                                                                                         |
|--------------------------|------------------------------------------------------------------------------------------------------------------------------------------|
| QR-001 (Performance)     | Indexes supporting FEFO queries; background cleanup reduces long-term database growth and helps maintain query performance.              |
| QR-002 (Availability)    | Leader election for automatic failover; health checks; graceful shutdown.                                                                |
| QR-003 (Testability)     | Layered architecture enables isolated unit testing; integration tests use Testcontainers where applicable.                               |
| QR-004 (Integrity)       | ACID transactions; immutable ledger; soft expiration with delayed cleanup.                                                               |
| QR-005 (Fault tolerance) | Transaction isolation prevents partial updates; leader election enables automatic recovery of background workers after instance failure. |

---

## Related Documentation

- [Quality Requirements](../quality-requirements.md)
- [Quality Requirement Tests](../quality-requirement-tests.md)
- [Testing Strategy](../testing.md)
- [Migrations](../migrations.md)
- [Storage Layer](../storage.md)