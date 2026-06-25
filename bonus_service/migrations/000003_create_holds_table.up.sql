CREATE TABLE IF NOT EXISTS holds
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    UUID        NOT NULL,
    order_id   UUID        NOT NULL,
    amount     BIGINT      NOT NULL CHECK (amount > 0),
    status     VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'confirmed', 'cancelled')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, order_id)
);

-- TTL worker: find stale active holds
CREATE INDEX IF NOT EXISTS idx_holds_expires_status ON holds (expires_at, status) WHERE status = 'active';

-- POST /confirm, POST /cancel
CREATE INDEX IF NOT EXISTS idx_holds_order_id ON holds (order_id);

-- GET /balance (list user's active holds)
CREATE INDEX IF NOT EXISTS idx_holds_user_status ON holds (user_id, status);
