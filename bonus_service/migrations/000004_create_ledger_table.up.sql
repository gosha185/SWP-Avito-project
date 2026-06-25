CREATE TABLE IF NOT EXISTS ledger
(
    id             BIGSERIAL,
    user_id        UUID         NOT NULL,
    operation_type VARCHAR(20)  NOT NULL
        CHECK (
            operation_type IN (
                               'accrual',
                               'hold',
                               'confirm',
                               'cancel',
                               'expiry'
                )
            ),
    amount         BIGINT       NOT NULL CHECK (amount > 0),
    batch_id       BIGINT,
    external_key   VARCHAR(255) NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    metadata       JSONB,
    PRIMARY KEY (user_id, id)
);

-- Idempotency: prevent duplicate operations
CREATE UNIQUE INDEX IF NOT EXISTS idx_ledger_external_key ON ledger (external_key);

-- Transaction history: GET /balance/:user_id/history
CREATE INDEX IF NOT EXISTS idx_ledger_user_created ON ledger (user_id, created_at);
