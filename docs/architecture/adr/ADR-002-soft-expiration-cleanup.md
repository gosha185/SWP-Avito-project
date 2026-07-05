# ADR-002: Soft Expiration and Background Cleanup of Domain Data

## Status

Accepted

## Context

Bonus batches and holds have a defined lifecycle. Holds expire after their configured TTL, while bonus batches expire
after their expiration date. Although expired records are no longer used in business operations, they may still be
required temporarily for auditing and reconciliation.

Keeping obsolete records indefinitely would increase database size and gradually degrade query performance. The system
therefore requires a mechanism that separates business-level expiration from physical data removal.

## Decision

Split the lifecycle of domain data into two phases:

1. **Soft expiration** – business operations mark data as inactive while preserving it for auditing. Expired holds are
   cancelled and expired bonus batches are depleted, with all changes recorded in the ledger.

2. **Background cleanup** – dedicated maintenance workers periodically remove obsolete records after the configured
   retention period.

This approach separates business logic from database maintenance and allows cleanup frequency and retention policies to
evolve independently.

## Consequences and Tradeoffs

### Advantages

- Separates business operations from database maintenance.
- Preserves expired data during the retention period for auditing.
- Prevents unbounded growth of database tables.
- Allows cleanup scheduling and retention policies to be configured independently.

### Disadvantages

- Requires additional background workers.
- Obsolete data remains in the database until cleanup.
- Cleanup operations must preserve referential integrity.

## Quality Requirements Addressed

- [QR-001](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md#qr-001-api-response-time) –
  Reduces table growth, helping maintain query performance over time.
- [QR-004](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md#qr-004-data-retention-and-availability) –
  Defines the data retention strategy before permanent deletion.

## Related ADRs

- [ADR-001](ADR-001-leader-election.md) – PostgreSQL-based Leader Election for Background Workers
- [ADR-003](ADR-003-transaction-isolation.md) – Transaction Isolation Strategy for Financial Operations

## Related Code

- `internal/service/expireAllBatches.go`
- `internal/service/cancelAllExpiredHolds.go`
- `internal/service/deleteExpiredBatches.go`
- `internal/service/cleanupOldHolds.go`