CREATE TABLE IF NOT EXISTS batches
(
    id         BIGSERIAL,
    user_id    UUID        NOT NULL,
    amount     BIGINT      NOT NULL CHECK (amount > 0),
    remaining  BIGINT      NOT NULL CHECK (remaining >= 0),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, id)
);

-- ORDER BY expires_at for FEFO (GET /balance, POST /hold)
CREATE INDEX IF NOT EXISTS idx_batches_user_expires ON batches (user_id, expires_at) WHERE remaining > 0;

CREATE INDEX IF NOT EXISTS idx_batches_expired ON batches (expires_at) WHERE remaining > 0;