# Storage Layer Documentation

The `internal/storage` package implements the repository pattern: one repository
per table, each wrapping `*sql.DB`. It is the only layer that talks SQL.

## Conventions

- Every method takes `context.Context` as the first argument.
- **Read** methods run on the pool (`*sql.DB`) directly.
- **Write** methods that must be atomic with others take a `*sql.Tx` — the
  service layer owns the transaction and decides commit/rollback.
- Row locks use `SELECT … FOR UPDATE`.
- Inserts use `RETURNING id, created_at` to populate the model in place.
- Expected outcomes return the sentinel errors below; unexpected failures return
  wrapped errors (`%w`).

## Sentinel Errors (`errors.go`)

| Error                    | Meaning                                           |
|--------------------------|---------------------------------------------------|
| `ErrBalanceNotFound`     | balance row missing                               |
| `ErrBalanceNotUpdated`   | `UPDATE balances` affected 0 rows                 |
| `ErrBatchNotFound`       | batch missing or `remaining < amount` on decrease |
| `ErrHoldNotFound`        | hold missing or not `active`                      |
| `ErrLedgerDuplicate`     | `external_key` already used (idempotency hit)     |
| `ErrIncorrectInput`      | non-positive amount/duration in validation        |
| `ErrInsufficientBalance` | not enough `available` to cover the operation     |
| `ErrOrderAlreadyHeld`    | order already has a hold (`UNIQUE` violation)     |

## BalanceRepo

### GetBalance

**Purpose:** Read current balance without locking.

**When to use:** GET /balance endpoints (read-only). Transaction not required.

**Parameters:**

- `ctx context.Context` — request context
- `userID uuid.UUID` — user identifier

**Returns:**

- `*models.Balance` — balance record (zero balance if user not found)
- `error` — database error

### GetBalanceForUpdate

**Purpose:** Read balance with row lock (`SELECT FOR UPDATE`). Prevents concurrent modifications.

**When to use:** Before any write operation (accrue, hold, confirm, cancel).

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `userID uuid.UUID` — user identifier

**Returns:**

- `*models.Balance` — balance record
- `error` — database error

**Important:** Creates a zero balance record if user does not exist.

### UpdateBalance

**Purpose:** Update balance by deltas (available ± amount, held ± amount).

**When to use:** After validation checks, inside the same transaction.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `userID uuid.UUID` — user identifier
- `deltaAvailable int64` — change for available (positive = add, negative = subtract)
- `deltaHeld int64` — change for held (positive = add, negative = subtract)

**Returns:**

- `error` — database error; `ErrBalanceNotUpdated` if no rows affected

### UpdateBalancesForExpiredBatches

**Purpose:** Decrease `available` for all users whose batches have expired (`remaining > 0`, `expires_at < NOW()`).

**When to use:** Batch-expiry worker.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)

**Returns:**

- `int64` — number of updated rows
- `error` — database error

**⚠️ Superseded:** the balance update is now part of the atomic `BatchRepo.ExpireAllBatches`. This method is no longer
called and is kept until a cleanup PR.

### UpdateBalancesForExpiredHolds

**Purpose:** Restore `available` and decrease `held` for all users whose holds have expired (`status = 'active'`,
`expires_at < NOW()`).

**When to use:** TTL worker.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)

**Returns:**

- `int64` — number of updated rows
- `error` — database error

**⚠️ Superseded:** the balance update is now part of the atomic `HoldRepo.CancelAllExpiredHolds`. Kept until a cleanup
PR.

## BatchRepo

### GetExpiringBatches

**Purpose:** Get all batches with remaining > 0, sorted by `expires_at` (FEFO order). No locking.

**When to use:** Read-only balance breakdowns.

**Parameters:**

- `ctx context.Context` — request context
- `userID uuid.UUID` — user identifier

**Returns:**

- `[]models.BonusBatch` — list of batches sorted by expires_at (oldest first)
- `error` — database error

---

### GetExpiringBatchesForUpdate

**Purpose:** Get all batches with remaining > 0, sorted by `expires_at` (FEFO order) and lock them using
`SELECT FOR UPDATE`.

**When to use:** `POST /hold` before reserving points from batches.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `userID uuid.UUID` — user identifier

**Returns:**

- `[]models.BonusBatch` — list of locked batches sorted by expires_at (oldest first)
- `error` — database error

**Important:** Must be called inside a transaction.

---

### GetExpiringSum

**Purpose:** Calculate total points that will expire in the next N days.

**When to use:** GET /balance?days=N.

**Parameters:**

- `ctx context.Context` — request context
- `userID uuid.UUID` — user identifier
- `days int` — number of days to look ahead

**Returns:**

- `int64` — sum of remaining points expiring within the next N days
- `error` — database error

**Note:** Optional feature. Can be omitted if only aggregated sum is needed.

---

### GetExpiringBreakdown

**Purpose:** Return a detailed breakdown of points that will expire in the next N days.

**When to use:** GET /balance?days=N&breakdown=true.

**Parameters:**

- `ctx context.Context` — request context
- `userID uuid.UUID` — user identifier
- `days int` — number of days to look ahead

**Returns:**

- `[]models.ExpiryBreakdown` — grouped expiration information (`days_left`, `amount`)
- `error` — database error

---

### CreateBatch

**Purpose:** Create a new batch of points.

**When to use:** POST /accrue.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `batch *models.BonusBatch` — batch data (must have UserID, Amount, Remaining, ExpiresAt)

**Returns:**

- `error` — database error

**Transaction:** Required.

**Note:** `batch.ID` and `batch.CreatedAt` are populated by the database via `RETURNING`.

---

### IncreaseBatchRemaining

**Purpose:** Return points back to a batch.

**When to use:** Cancel hold (points returned to original batches).

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `batchID int64` — batch identifier
- `amount int64` — points to return

**Returns:**

- `bool` — `true` if the batch was updated; `false` if it does not exist or has already expired
- `error` — database error

**Important:** Only restores points to non-expired batches (`expires_at > NOW()`); expired batches are not topped up (
points are burned instead).

---

### DecreaseBatchRemaining

**Purpose:** Spend points from a single batch.

**When to use:** Confirm operation (final spend).

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `batchID int64` — batch identifier
- `amount int64` — points to spend

**Returns:**

- `error` — `ErrBatchNotFound` if batch doesn't exist or has insufficient remaining

**Important:** The method checks `remaining >= amount` before updating.

---

### GetAllBatches

**Purpose:** Return all user batches, including fully spent batches.

**When to use:** Admin/debug endpoints.

**Parameters:**

- `ctx context.Context` — request context
- `userID uuid.UUID` — user identifier

**Returns:**

- `[]models.BonusBatch` — all batches ordered by creation time (newest first)
- `error` — database error

---

### SetBatchRemaining

**Purpose:** Set batch remaining amount directly.

**When to use:** Admin tools, migrations, recovery procedures.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `batchID int64` — batch identifier
- `newRemaining int64` — new remaining value

**Returns:**

- `error` — `ErrBatchNotFound` if batch does not exist

**Important:** Bypasses normal spending validation and should be used with caution.

---

### GetExpiredBatches

**Purpose:** Get all batches with `remaining > 0` and `expires_at < NOW()`. No locking.

**When to use:** Historically used by the batch-expiry worker to list batches to expire.

**Parameters:**

- `ctx context.Context` — request context

**Returns:**

- `[]models.BonusBatch` — expired batches (oldest first)
- `error` — database error

**⚠️ Superseded:** the expiry sweep is now done atomically by `ExpireAllBatches`; this method is no longer called.

---

### DeleteExpiredZeroBatches

**Purpose:** Delete fully-spent (`remaining = 0`) batches that expired more than `daysOld` days ago.

**When to use:** Batch-cleanup worker.

**Parameters:**

- `ctx context.Context` — request context
- `daysOld int` — retention window in days

**Returns:**

- `int64` — number of deleted rows
- `error` — database error

**Important:** `ON DELETE RESTRICT` on `hold_batches` prevents removing a batch still referenced by a hold.

---

### ExpireAllBatches

**Purpose:** Atomically expire all overdue batches in a single statement — deduct the lost points from `available`,
write `expiry` ledger entries, and set `remaining = 0`.

**When to use:** Batch-expiry worker.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)

**Returns:**

- `int64` — number of expired batches
- `error` — database error

**Important:** A leading CTE locks the target rows with `FOR UPDATE` before any dependent write, so the
balance/ledger/remaining updates all act on the same snapshot and cannot race with a concurrent hold spending from the
same batch.

## HoldRepo

### GetHoldByID

**Purpose:** Get a hold by its ID.

**When to use:** Read-only hold lookups.

**Parameters:**

- `ctx context.Context` — request context
- `holdID int64` — hold identifier

**Returns:**

- `*models.Hold` — hold record
- `error` — `ErrHoldNotFound` if not found

### GetHoldByOrderID

**Purpose:** Find active hold by order_id with row lock (`SELECT FOR UPDATE`).

**When to use:** POST /confirm, POST /cancel.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `userID uuid.UUID` — user identifier
- `orderID uuid.UUID` — order identifier

**Returns:**

- `*models.Hold` — hold record (status = active)
- `error` — `ErrHoldNotFound` if no active hold found

### GetActiveHoldsByUser

**Purpose:** Get all active holds for a user.

**When to use:** GET /balance (show held points).

**Parameters:**

- `ctx context.Context` — request context
- `userID uuid.UUID` — user identifier

**Returns:**

- `[]models.Hold` — list of active holds (newest first)
- `error` — database error

### CreateHold

**Purpose:** Create a new hold record.

**When to use:** POST /hold.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `hold *models.Hold` — hold data (must have UserID, OrderID, Amount, Status, ExpiresAt)

**Returns:**

- `error` — database error; `ErrOrderAlreadyHeld` if the order already has a hold (`UNIQUE` violation)

**Note:** `hold.ID` and `hold.CreatedAt` are populated by the database via `RETURNING`.

### UpdateHoldStatus

**Purpose:** Change hold status to confirmed or cancelled.

**When to use:** POST /confirm (status → confirmed), POST /cancel (status → cancelled).

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `holdID int64` — hold identifier
- `status string` — new status ("confirmed" or "cancelled")

**Returns:**

- `error` — `ErrHoldNotFound` if hold not found or not active

### GetExpiredHolds

**Purpose:** Get all active holds with `expires_at < NOW()`. No locking.

**When to use:** Historically used by the TTL worker to list holds to cancel.

**Parameters:**

- `ctx context.Context` — request context

**Returns:**

- `[]models.Hold` — list of expired active holds (oldest first)
- `error` — database error

**⚠️ Superseded:** cancellation is now done atomically by `CancelAllExpiredHolds`; this read-only method is no longer
called.

---

### DeleteOldHolds

**Purpose:** Delete `confirmed`/`cancelled` holds that expired more than `daysOld` days ago.

**When to use:** Hold-cleanup worker.

**Parameters:**

- `ctx context.Context` — request context
- `daysOld int` — retention window in days

**Returns:**

- `int64` — number of deleted rows
- `error` — database error

**Note:** Linked `hold_batches` rows are removed via `ON DELETE CASCADE`.

---

### CancelAllExpiredHolds

**Purpose:** Atomically cancel all expired active holds in a single statement — restore `available` / decrease `held`,
write `cancel` ledger entries, and set `status = 'cancelled'`.

**When to use:** TTL worker.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)

**Returns:**

- `int64` — number of cancelled holds
- `error` — database error

**Important:** A leading CTE locks the target rows with `FOR UPDATE` before any dependent write, so the operation is
serialized against a concurrent confirm/cancel on the same hold.

## HoldBatchRepo

Methods for managing many-to-many relationships between holds and batches (`hold_batches` table).

### CreateHoldBatch

**Purpose:** Link a hold to a batch with the amount taken from that batch.

**When to use:** During Hold operation — after selecting FEFO batches and creating the hold.

**Parameters:**

- `ctx context.Context`
- `tx *sql.Tx`
- `holdID int64`
- `batchUserID uuid.UUID`
- `batchID int64`
- `amount int64`

**Returns:**

- `error`

---

### GetHoldBatchesByHoldID

**Purpose:** Get all hold-batch links for a hold.

**When to use:** Cancel operation, TTL cancellation, audit.

**Parameters:**

- `ctx context.Context`
- `tx *sql.Tx`
- `holdID int64`

**Returns:**

- `[]models.HoldBatch`
- `error`

## LedgerRepo

### Insert

**Purpose:** Append an immutable audit record to the ledger.

**When to use:** After every write operation (accrue, hold, confirm, cancel, expiry).

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `entry *models.LedgerEntry` — audit entry (must have UserID, OperationType, Amount, ExternalKey)

**Returns:**

- `error` — `ErrLedgerDuplicate` if external_key already exists

**Important:** The unique index on `external_key` ensures idempotency.

### GetByExternalKey

**Purpose:** Check if an operation with this external_key already exists. Used for idempotency checks.

**When to use:** Before processing any write request (accrue, hold, confirm, cancel) to detect duplicate calls.

**Parameters:**

- `ctx context.Context` — request context
- `externalKey string` — idempotency key provided by the caller

**Returns:**

- `*models.LedgerEntry` — existing ledger entry, or `nil` if not found
- `error` — database error

### GetByExternalKeyTx

**Purpose:** Same idempotency check as `GetByExternalKey`, but with `SELECT FOR UPDATE` to lock the row and prevent race
conditions.

**When to use:** Inside a transaction, when the idempotency check must be serialized with the writing.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)
- `externalKey string` — idempotency key provided by the caller

**Returns:**

- `*models.LedgerEntry` — existing ledger entry, or `nil` if not found
- `error` — database error

### GetHistory

**Purpose:** Get paginated transaction history for a user.

**When to use:** GET /balance/:user_id/history endpoint.

**Parameters:**

- `ctx context.Context` — request context
- `userID uuid.UUID` — user identifier
- `limit int` — maximum number of entries to return (page size)
- `offset int` — number of entries to skip (pagination offset)

**Returns:**

- `[]models.LedgerEntry` — list of ledger entries (newest first)
- `error` — database error

### GetHistoryCount

**Purpose:** Get total number of transactions for a user. Used for pagination calculations.

**When to use:** Before calling `GetHistory` to determine total pages or to show total count in the response.

**Parameters:**

- `ctx context.Context` — request context
- `userID uuid.UUID` — user identifier

**Returns:**

- `int` — total count of ledger entries for the user
- `error` — database error

**Note:** Used together with `GetHistory` to implement pagination (total pages = ceil(count / limit)).

### GetLedgerByOrderID

**Purpose:** Find all operations related to a specific order.

**When to use:** Admin/debug endpoints, order-level audit.

**Parameters:**

- `ctx context.Context` — request context
- `orderID uuid.UUID` — order identifier

**Returns:**

- `[]models.LedgerEntry` — list of entries (newest first)
- `error` — database error

**Note:** Uses JSONB search in `metadata` field.

### InsertExpiryEntries

**Purpose:** Insert `expiry` ledger entries for all expired batches (`remaining > 0`, `expires_at < NOW()`).

**When to use:** Batch-expiry worker.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)

**Returns:**

- `int64` — number of inserted rows
- `error` — database error

**⚠️ Superseded:** ledger insertion is now part of the atomic `BatchRepo.ExpireAllBatches`. Kept until a cleanup PR.

### InsertCancelEntries

**Purpose:** Insert `cancel` ledger entries for all expired active holds.

**When to use:** TTL worker.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)

**Returns:**

- `int64` — number of inserted rows
- `error` — database error

**⚠️ Superseded:** ledger insertion is now part of the atomic `HoldRepo.CancelAllExpiredHolds`. Kept until a cleanup PR.

### isDuplicateKeyError

**Purpose:** Detects if an error is a PostgreSQL unique violation (error code 23505).

**When to use:** Used internally by `Insert` to identify duplicate key errors and return `ErrLedgerDuplicate`.

**Parameters:**

- `err error` — the error returned from a database operation

**Returns:**

- `bool` — `true` if the error is a unique violation, `false` otherwise

**Note:** Private helper function. Not exported.

## LeaderRepo

Coordinates background workers so they run on a single instance at a time
(`leaders` table). See [workers.md](workers.md).

### TryAcquireLock

**Purpose:** Attempt to become the leader for a role. Upserts into `leaders`, succeeding only if the current lease is
stale (`updated_at < NOW() - ttl`) or already owned by this instance; a successful call by the current holder also
renews the lease.

**Parameters:**

- `ctx context.Context` — request context
- `roleName string` — worker role (e.g. `bonus_worker`)
- `instanceID string` — this instance's id (`<host>:9091`)
- `ttlSeconds int` — lease TTL in seconds

**Returns:**

- `bool` — `true` if this instance now holds the lease
- `error` — database error

### RenewLock

**Purpose:** Refresh `updated_at` for the current leader.

**Parameters:**

- `ctx context.Context` — request context
- `roleName string` — worker role
- `instanceID string` — this instance's id

**Returns:**

- `bool` — `true` if the lease was renewed
- `error` — error if the leader record was not found

**Note:** Defined for completeness; lease renewal is currently handled by `TryAcquireLock`.

## DB

### NewDB

**Purpose:** Creates a new PostgreSQL connection pool with sane defaults.

**Parameters:**

- `dsn string` — PostgreSQL connection string (defaults to local development if empty)

**Returns:**

- `*sql.DB` — connection pool
- `error` — database error

**Connection pool defaults:**

- `MaxOpenConns: 25`
- `MaxIdleConns: 25`
- `ConnMaxLifetime: 5 minutes`