# Migrations Documentation

This document explains the database schema and indexes used in the bonus service.

---

## Tables Overview

| Table          | Purpose                                            |
|----------------|----------------------------------------------------|
| `balances`     | Current user balance (available + held)            |
| `batches`      | Point batches with expiration dates (FEFO support) |
| `holds`        | Two-phase spending: hold → confirm / cancel        |
| `hold_batches` | Many-to-many link between holds and batches        |
| `ledger`       | Immutable audit log (append-only)                  |
| `leaders`      | Leader-election lock for background workers        |

---

## balances

Stores the current balance for each user.

| Column       | Type          | Description                          |
|--------------|---------------|--------------------------------------|
| `user_id`    | `UUID`        | Primary key, user identifier         |
| `available`  | `BIGINT`      | Points available for use (>= 0)      |
| `held`       | `BIGINT`      | Points locked in active holds (>= 0) |
| `updated_at` | `TIMESTAMPTZ` | Last update timestamp                |

---

## batches

Stores point batches with expiration dates. Supports FEFO (First-Expiring-First-Out) spending.

**Constraints:**

- `PRIMARY KEY (user_id, id)` — composite primary key for future sharding support

| Column       | Type          | Description                                     |
|--------------|---------------|-------------------------------------------------|
| `id`         | `BIGSERIAL`   | Batch identifier (unique within user partition) |
| `user_id`    | `UUID`        | User who owns these points                      |
| `amount`     | `BIGINT`      | Original batch size (> 0)                       |
| `remaining`  | `BIGINT`      | Points still available in this batch (>= 0)     |
| `expires_at` | `TIMESTAMPTZ` | Points expire after this date                   |
| `created_at` | `TIMESTAMPTZ` | Creation timestamp                              |

**Indexes:**

| Index                      | Columns                                     | Purpose                                                  |
|----------------------------|---------------------------------------------|----------------------------------------------------------|
| `idx_batches_user_expires` | `(user_id, expires_at) WHERE remaining > 0` | FEFO lookups: ORDER BY expires_at for /hold and /balance |
| `idx_batches_expired`      | `(expires_at) WHERE remaining > 0`          | Batch-expiry worker: find batches that have expired      |

---

## holds

Stores active and historical holds for two-phase spending.

| Column       | Type          | Description                                                  |
|--------------|---------------|--------------------------------------------------------------|
| `id`         | `BIGSERIAL`   | Primary key, auto-incrementing hold ID                       |
| `user_id`    | `UUID`        | User who owns this hold                                      |
| `order_id`   | `UUID`        | Order identifier from external service (UNIQUE with user_id) |
| `amount`     | `BIGINT`      | Amount held (> 0)                                            |
| `status`     | `VARCHAR(20)` | `active` / `confirmed` / `cancelled` (CHECK constraint)      |
| `expires_at` | `TIMESTAMPTZ` | Expiry time, set by caller as `now + hours`                  |
| `created_at` | `TIMESTAMPTZ` | Creation timestamp                                           |

**Constraints:**

- `UNIQUE(user_id, order_id)` — prevents duplicate holds for the same order

**Indexes:**

| Index                      | Columns                                        | Purpose                                 |
|----------------------------|------------------------------------------------|-----------------------------------------|
| `idx_holds_expires_status` | `(expires_at, status) WHERE status = 'active'` | TTL worker: find stale active holds     |
| `idx_holds_order_id`       | `(order_id)`                                   | POST /confirm, POST /cancel lookups     |
| `idx_holds_user_status`    | `(user_id, status)`                            | GET /balance (list user's active holds) |

---

## hold_batches

Many-to-many link between holds and batches. Stores which batches were locked by a hold.

| Column          | Type     | Description                                               |
|-----------------|----------|-----------------------------------------------------------|
| `hold_id`       | `BIGINT` | References holds(id) ON DELETE CASCADE                    |
| `batch_user_id` | `UUID`   | User who owns the batch (part of composite FK to batches) |
| `batch_id`      | `BIGINT` | Batch id (part of composite FK to batches)                |
| `amount`        | `BIGINT` | Amount from this batch used in the hold (> 0)             |

**Constraints:**

- `PRIMARY KEY (hold_id, batch_user_id, batch_id)` — composite primary key for future sharding support
- `FOREIGN KEY (batch_user_id, batch_id) REFERENCES batches (user_id, id) ON DELETE RESTRICT`

**Indexes:**

| Index                       | Columns                     | Purpose                                           |
|-----------------------------|-----------------------------|---------------------------------------------------|
| `idx_hold_batches_batch_id` | `(batch_user_id, batch_id)` | Batch cleanup: check if batch is still referenced |

---

## ledger

Immutable audit log. Append-only, never updated or deleted.

**Constraints:**

- `PRIMARY KEY (user_id, id)` — composite primary key for future sharding support

| Column           | Type           | Description                                                             |
|------------------|----------------|-------------------------------------------------------------------------|
| `id`             | `BIGSERIAL`    | Ledger entry identifier (unique within user partition)                  |
| `user_id`        | `UUID`         | User affected by this operation                                         |
| `operation_type` | `VARCHAR(20)`  | `accrual` / `hold` / `confirm` / `cancel` / `expiry` (CHECK constraint) |
| `amount`         | `BIGINT`       | Absolute amount (> 0)                                                   |
| `batch_id`       | `BIGINT`       | Batch ID (NULL for operations without batch)                            |
| `external_key`   | `VARCHAR(255)` | Idempotency key from caller (UNIQUE)                                    |
| `created_at`     | `TIMESTAMPTZ`  | Entry timestamp                                                         |
| `metadata`       | `JSONB`        | Additional context (order_id, hold_id, etc.)                            |

**Indexes:**

| Index                     | Columns                 | Purpose                                            |
|---------------------------|-------------------------|----------------------------------------------------|
| `idx_ledger_external_key` | `(external_key)` UNIQUE | Idempotency: prevent duplicate operations          |
| `idx_ledger_user_created` | `(user_id, created_at)` | Transaction history: GET /balance/:user_id/history |

---

## leaders

Distributed lock for background workers (leader election). One row per role;
the holder is whoever last wrote `leader_id` within the TTL window. See
[workers.md](workers.md).

| Column       | Type           | Description                                         |
|--------------|----------------|-----------------------------------------------------|
| `role_name`  | `VARCHAR(255)` | Primary key, worker role (e.g. `bonus_worker`)      |
| `leader_id`  | `VARCHAR(255)` | Current holder instance id (`<host>:9091`)          |
| `updated_at` | `TIMESTAMPTZ`  | Lease heartbeat — lock is stale once older than TTL |

---

## Index Usage Summary

| Index                       | Table          | Used In                            |
|-----------------------------|----------------|------------------------------------|
| `idx_batches_user_expires`  | `batches`      | `POST /hold`, `GET /balance`       |
| `idx_batches_expired`       | `batches`      | Batch-expiry / cleanup workers     |
| `idx_holds_expires_status`  | `holds`        | TTL background worker              |
| `idx_holds_order_id`        | `holds`        | `POST /confirm`, `POST /cancel`    |
| `idx_holds_user_status`     | `holds`        | `GET /balance` (list active holds) |
| `idx_hold_batches_batch_id` | `hold_batches` | Batch cleanup, audit               |
| `idx_ledger_external_key`   | `ledger`       | All write operations (idempotency) |
| `idx_ledger_user_created`   | `ledger`       | `GET /balance/:user_id/history`    |

---

## Operation Type Reference

| Operation | Description                     | Used In                    |
|-----------|---------------------------------|----------------------------|
| `accrual` | Points added to user balance    | `POST /accrue`             |
| `hold`    | Points locked for an order      | `POST /hold`               |
| `confirm` | Hold confirmed, points spent    | `POST /confirm`            |
| `cancel`  | Hold cancelled, points released | `POST /cancel`, TTL worker |
| `expiry`  | Points expired automatically    | TTL / batch-expiry workers |