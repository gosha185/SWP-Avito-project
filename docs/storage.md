# Storage Layer Documentation

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

- `error` — `ErrBatchNotFound` if batch doesn't exist

**Important:** This method does NOT check for expiration or any other conditions. The caller must ensure the batch is
still valid (e.g., not expired).

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

- `error` — database error

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

**Purpose:** Get all active holds with `expires_at < NOW()`.

**When to use:** TTL worker — finds holds that need automatic cancellation.

**Parameters:**

- `ctx context.Context` — request context
- `tx *sql.Tx` — open transaction (required)

**Returns:**

- `[]models.Hold` — list of expired holds (locked with `FOR UPDATE`)
- `error` — database error

**Important:** Rows are locked with `FOR UPDATE` to prevent concurrent processing. This method only finds holds — cancellation logic (restoring balances, batches, ledger) must be implemented in the service layer.

## HoldBatchRepo

Methods for managing many-to-many relationships between holds and batches (`hold_batches` table).

### CreateHoldBatch

**Purpose:** Link a hold to a batch with the amount taken from that batch.

**When to use:** During Hold operation — after selecting FEFO batches and creating the hold.

**Parameters:**

- `ctx context.Context`
- `tx *sql.Tx`
- `holdID int64`
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

### isDuplicateKeyError

**Purpose:** Detects if an error is a PostgreSQL unique violation (error code 23505).

**When to use:** Used internally by `Insert` to identify duplicate key errors and return `ErrLedgerDuplicate`.

**Parameters:**

- `err error` — the error returned from a database operation

**Returns:**

- `bool` — `true` if the error is a unique violation, `false` otherwise

**Note:** Private helper function. Not exported.

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
