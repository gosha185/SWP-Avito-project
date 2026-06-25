CREATE TABLE IF NOT EXISTS hold_batches
(
    hold_id       BIGINT NOT NULL REFERENCES holds (id) ON DELETE CASCADE,
    batch_user_id UUID   NOT NULL,
    batch_id      BIGINT NOT NULL,
    amount        BIGINT NOT NULL CHECK (amount > 0),
    PRIMARY KEY (hold_id, batch_user_id, batch_id),
    FOREIGN KEY (batch_user_id, batch_id) REFERENCES batches (user_id, id) ON DELETE RESTRICT
);

-- Batch cleanup: check if batch is still referenced by any hold
CREATE INDEX IF NOT EXISTS idx_hold_batches_batch_id ON hold_batches (batch_user_id, batch_id);
