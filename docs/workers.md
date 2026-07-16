# Background Workers

Periodic maintenance jobs (expiry, cleanup) that keep balances and tables
consistent. They run **only on the elected leader instance**, so with multiple
replicas the jobs never run concurrently.

## Leader election

Coordinated through the `leaders` table (`role_name = "bonus_worker"`).

- Each instance identifies itself as `<hostname>:9091` with a lease TTL of
  **30s**.
- `LeaderService` ticks every **10s** and calls `LeaderRepo.TryAcquireLock`,
  an upsert that succeeds only if the current lease is stale
  (`updated_at < NOW() - ttl`) **or** already owned by this instance. A
  successful call by the current leader also refreshes the lease, so ticking
  every 10s against a 30s TTL gives three heartbeats of slack.
- On **gaining** leadership → `startWorkers`; on **losing** it (or on
  shutdown/context-cancel) → `stopWorkers`, which cancels the worker context and
  waits for all goroutines (`sync.WaitGroup`).

Each worker runs its job once immediately, then on its own `time.Ticker`. Job
errors are logged and do not stop the loop; the loop exits on context cancel.

## Registered workers

Registered in `cmd/api/main.go`. Intervals come from env vars (seconds).

| Worker               | Env var                         | Default | Job → service method    | Effect                                                              |
|----------------------|---------------------------------|---------|-------------------------|---------------------------------------------------------------------|
| `TTLWorker`          | `WORKER_TTL_INTERVAL`           | 5s      | `CancelAllExpiredHolds` | Cancel expired active holds; return held points to `available`      |
| `BatchExpiryWorker`  | `WORKER_BATCH_EXPIRY_INTERVAL`  | 5s      | `ExpireAllBatches`      | Zero out expired batches; deduct lost points from `available`       |
| `BatchCleanupWorker` | `WORKER_BATCH_CLEANUP_INTERVAL` | 60s     | `DeleteExpiredBatches`  | Delete spent (`remaining = 0`) batches expired > 30 days ago        |
| `HoldCleanupWorker`  | `WORKER_HOLD_CLEANUP_INTERVAL`  | 60s     | `CleanupOldHolds`       | Delete settled holds expired > 30 days ago (`hold_batches` cascade) |

## Job details

### TTLWorker → `CancelAllExpiredHolds`

Opens a transaction and calls `HoldRepo.CancelAllExpiredHolds` — a single
atomic statement that, for every `active` hold past `expires_at`: restores
`available` / decreases `held`, writes a `cancel` ledger entry, and sets
`status = 'cancelled'`. The leading `FOR UPDATE` CTE serializes it against
concurrent confirm/cancel on the same hold.

### BatchExpiryWorker → `ExpireAllBatches`

Same shape via `BatchRepo.ExpireAllBatches`: for every batch with
`remaining > 0` past `expires_at`, deducts the lost points from `available`,
writes an `expiry` ledger entry, and sets `remaining = 0`.

### BatchCleanupWorker → `DeleteExpiredBatches`

Calls `BatchRepo.DeleteExpiredZeroBatches(30)` — purges fully-spent batches that
expired more than 30 days ago. `ON DELETE RESTRICT` on `hold_batches` prevents
removing a batch still referenced by a hold.

### HoldCleanupWorker → `CleanupOldHolds`

Calls `HoldRepo.DeleteOldHolds(30)` — deletes `confirmed`/`cancelled` holds that
expired more than 30 days ago; linked `hold_batches` rows are removed via
`ON DELETE CASCADE`.

## Notes / follow-ups

- The two expiry sweeps lock **all** matching rows in one transaction. Under a
  large backlog (e.g. after a long worker outage) this can hold many locks at
  once; batching with `FOR UPDATE SKIP LOCKED` is a possible future improvement.
- Cleanup retention (30 days) is currently hard-coded in the worker wrappers,
  not env-configurable.
  