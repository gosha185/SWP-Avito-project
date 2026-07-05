# ADR-001: PostgreSQL-based Leader Election for Background Workers

## Status

Accepted

## Context

The bonus service executes several background maintenance workers, including TTL expiration, bonus batch expiration, and
cleanup tasks.

When multiple service instances are deployed for horizontal scaling, each instance would start its own background
workers. This could result in:

- duplicate processing of the same records;
- race conditions when updating shared data;
- unnecessary database load;
- duplicate ledger entries.

The system therefore requires a coordination mechanism that guarantees only one service instance executes background
workers at any given time.

## Decision

Use PostgreSQL as a lightweight distributed lock service.

A dedicated `leaders` table stores the current leader for each worker role together with a lease timestamp.

Each service instance periodically attempts to acquire or renew the lease using an atomic
`INSERT ... ON CONFLICT ... DO UPDATE` statement. The current leader periodically refreshes its lease, while other
instances continue attempting to acquire leadership.

If the leader becomes unavailable and its lease expires, another instance automatically becomes the leader and starts
all background workers.

## Consequences and Tradeoffs

### Advantages

- Prevents duplicate execution of background workers.
- Requires no external coordination service such as ZooKeeper or etcd.
- Reuses the existing PostgreSQL deployment.
- Provides automatic failover when the leader becomes unavailable.
- Supports horizontal scaling with multiple service instances.

### Disadvantages

- Introduces additional database traffic for lease renewal.
- Leader failover is limited by the configured lease TTL.
- PostgreSQL becomes the coordination backend for worker execution.

## Quality Requirements Addressed

- [QR-002](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md#qr-002-service-availability) –
  Ensures background workers continue operating after a leader failure.
- [QR-005](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md#qr-005-fault-tolerance-and-recovery) –
  Provides automatic recovery of background processing after service instance failure.

## Related ADRs

- [ADR-002](ADR-002-soft-expiration-cleanup.md) – Soft
  Expiration and Background Cleanup of Domain Data
- [ADR-003](ADR-003-transaction-isolation.md) –
  Transaction Isolation Strategy for Financial Operations

## Related Code

- `internal/storage/leader.go`
- `internal/service/leader.go`
- `migrations/000006_create_leaders_table.up.sql`