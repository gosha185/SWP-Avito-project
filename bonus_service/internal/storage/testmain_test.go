package storage_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const fullSchema = `
CREATE TABLE IF NOT EXISTS balances (
    user_id    UUID PRIMARY KEY,
    available  BIGINT      NOT NULL DEFAULT 0 CHECK (available >= 0),
    held       BIGINT      NOT NULL DEFAULT 0 CHECK (held >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS batches (
    id         BIGSERIAL,
    user_id    UUID        NOT NULL,
    amount     BIGINT      NOT NULL CHECK (amount > 0),
    remaining  BIGINT      NOT NULL CHECK (remaining >= 0),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, id)
);

CREATE INDEX IF NOT EXISTS idx_batches_user_expires ON batches (user_id, expires_at) WHERE remaining > 0;
CREATE INDEX IF NOT EXISTS idx_batches_expired ON batches (expires_at) WHERE remaining > 0;

CREATE TABLE IF NOT EXISTS holds (
    id         BIGSERIAL PRIMARY KEY,
    user_id    UUID        NOT NULL,
    order_id   UUID        NOT NULL,
    amount     BIGINT      NOT NULL CHECK (amount > 0),
    status     VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'confirmed', 'cancelled')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, order_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_holds_user_id_id ON holds (user_id, id);

CREATE TABLE IF NOT EXISTS hold_batches (
    hold_id       BIGINT NOT NULL REFERENCES holds (id) ON DELETE CASCADE,
    batch_user_id UUID   NOT NULL,
    batch_id      BIGINT NOT NULL,
    amount        BIGINT NOT NULL CHECK (amount > 0),
    PRIMARY KEY (hold_id, batch_user_id, batch_id),
    FOREIGN KEY (batch_user_id, batch_id) REFERENCES batches (user_id, id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_holds_expires_status ON holds (expires_at, status) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_holds_order_id ON holds (order_id);
CREATE INDEX IF NOT EXISTS idx_holds_user_status ON holds (user_id, status);
CREATE INDEX IF NOT EXISTS idx_hold_batches_batch_id ON hold_batches (batch_user_id, batch_id);

CREATE TABLE IF NOT EXISTS ledger (
    id             BIGSERIAL,
    user_id        UUID         NOT NULL,
    operation_type VARCHAR(20)  NOT NULL CHECK (operation_type IN ('accrual','hold','confirm','cancel','expiry')),
    amount         BIGINT       NOT NULL CHECK (amount > 0),
    batch_id       BIGINT,
    external_key   VARCHAR(255) NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    metadata       JSONB,
    PRIMARY KEY (user_id, id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ledger_external_key ON ledger (external_key);
CREATE INDEX IF NOT EXISTS idx_ledger_user_created ON ledger (user_id, created_at);
`

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(ctx) })

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err = db.PingContext(ctx); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	if _, err = db.Exec(fullSchema); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	return db
}
