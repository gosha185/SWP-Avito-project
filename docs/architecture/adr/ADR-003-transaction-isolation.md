# ADR-003: Transaction Isolation Strategy for Financial Operations

## Status

Accepted

## Context

The bonus service performs financial operations that modify multiple related entities, including balances, bonus
batches, holds, and ledger entries. These operations must remain atomic and consistent, even when multiple requests for
the same user are processed concurrently or when background workers execute alongside user requests.

The system therefore requires a transaction strategy that guarantees data consistency while maintaining acceptable
performance.

## Decision

Use PostgreSQL's default **READ COMMITTED** isolation level together with explicit row-level locking (
`SELECT ... FOR UPDATE`) for all operations that modify financial data.

The transaction strategy is based on the following principles:

- all financial operations execute within a single database transaction;
- balance rows are locked before modification using pessimistic locking;
- database rows are always locked in a consistent order (`balances` → `batches` → `holds` → `ledger`) to reduce the risk
  of deadlocks;
- background workers should prefer set-based SQL operations over per-record processing whenever possible to minimize
  transaction duration and lock contention;
- database operations execute with bounded timeouts to prevent indefinite waits.

## Consequences and Tradeoffs

### Advantages

- Guarantees ACID-compliant financial operations.
- Prevents lost updates and duplicate spending.
- Uses PostgreSQL's native concurrency mechanisms.
- Reduces lock contention by encouraging set-based operations.

### Disadvantages

- Concurrent operations affecting the same user become serialized.
- Incorrect lock ordering may still result in deadlocks.
- High contention can reduce throughput.

## Alternatives Considered

- **Optimistic locking** - rejected because it requires retry logic and additional complexity.
- **SERIALIZABLE isolation** - rejected because it introduces unnecessary transaction conflicts.
- **External distributed locking** - rejected because PostgreSQL provides sufficient consistency guarantees for this
  service.

## Quality Requirements Addressed

- [QR-001](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md#qr-001-api-response-time) - Reduces table growth, helping maintain query performance over time.
- [QR-005](https://github.com/gosha185/SWP-Avito-project/blob/main/docs/quality-requirements.md#qr-005-fault-tolerance-and-recovery) - Fault tolerance and recovery after service instance failure.

## Related ADRs

- [ADR-001](ADR-001-leader-election.md) - PostgreSQL-based Leader Election for Background Workers
- [ADR-002](ADR-002-soft-expiration-cleanup.md) - Soft Expiration and Background Cleanup of Domain Data

## Related Code

- `internal/storage/`
- `internal/service/`