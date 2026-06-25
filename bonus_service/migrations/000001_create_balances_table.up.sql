CREATE TABLE IF NOT EXISTS balances
(
    user_id    UUID PRIMARY KEY,
    available  BIGINT      NOT NULL DEFAULT 0 CHECK (available >= 0),
    held       BIGINT      NOT NULL DEFAULT 0 CHECK (held >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
